package setting

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/safchain/ethtool"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type NetworkParamTemplate struct{
	ID   string         `json:ID`
	Name string         `json:"Name"`
	DHCP string         `json:"DHCP"`
	IP string           `json:"IP"`
	Netmask string      `json:"Netmask"`
	Broadcast string    `json:"Broadcast"`
	MAC string          `json:"MAC"`
}

type NetworkParamListTemplate struct{
	NetworkParam []NetworkParamTemplate
}

type NetworkLinkStateTemplate struct{
	State [2]uint32
}

var NetworkParamList NetworkParamListTemplate
var NetworkLinkState NetworkLinkStateTemplate


func cmdSetDHCP(id int){

	cmd := exec.Command("udhcpc","-i",NetworkParamList.NetworkParam[id].Name,"5")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()	//执行到此会阻塞5s

	str := out.String()

	log.Println(str)
}

func cmdSetStaticIP(id int){

	strNetMask   := "netmask " + NetworkParamList.NetworkParam[id].Netmask
	cmd := exec.Command("ifconfig",
		NetworkParamList.NetworkParam[id].Name,
		NetworkParamList.NetworkParam[id].IP,
		strNetMask)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Start()	//执行到此,直接往后执行

	cmd2 := exec.Command("/sbin/route","add","default","gw",NetworkParamList.NetworkParam[id].Broadcast)
	cmd2.Stdout = &out
	cmd2.Start()	//执行到此,直接往后执行
}

func GetNetworkStatus(){

	for _,v := range NetworkParamList.NetworkParam{
		GetLinkState(v.ID)
	}
}

func GetLinkState(id string){

	ethHandle, _ := ethtool.NewEthtool()
	defer ethHandle.Close()

	var state uint32

	if id == "1"{
		state, _ = ethHandle.LinkState(NetworkParamList.NetworkParam[0].Name)
		NetworkLinkState.State[0] = state
	}else if id == "2"{
		state, _ = ethHandle.LinkState(NetworkParamList.NetworkParam[1].Name)
		NetworkLinkState.State[1] = state
	}
}


//获取当前网络参数
func GetNetworkParam() NetworkParamListTemplate{

	for k,v := range NetworkParamList.NetworkParam{

		GetLinkState(v.ID)
		ethInfo,err := GetNetInformation(v.Name)
		//log.Printf("ethInfo %+v\n",ethInfo)
		if err == nil{
			NetworkParamList.NetworkParam[k].IP = ethInfo.IP
			NetworkParamList.NetworkParam[k].Netmask = ethInfo.Mask
			NetworkParamList.NetworkParam[k].Broadcast = ethInfo.GatewayIP
			NetworkParamList.NetworkParam[k].MAC = strings.ToUpper(ethInfo.Mac)
		}
	}

	return NetworkParamList
}

//设置网络参数
func SetNetworkParam(id string,param NetworkParamTemplate) bool{

	GetLinkState(id)

	if id == "1"{
		if NetworkLinkState.State[0] == 0{
			//log.Printf("setNetworkParam %s err\n",id)
			return false
		}

		NetworkParamList.NetworkParam[0] = param

		if NetworkParamList.NetworkParam[0].DHCP == "1"{
			//开启动态IP
			cmdSetDHCP(0)
		}else if NetworkParamList.NetworkParam[0].DHCP == "0"{
			//开启静态IP
			cmdSetStaticIP(0)
		}
	}else if id == "2"{
		if NetworkLinkState.State[1] == 0{
			//log.Printf("setNetworkParam %s err\n",id)
			return false
		}

		NetworkParamList.NetworkParam[1] = param

		if NetworkParamList.NetworkParam[1].DHCP == "1"{
			//开启动态IP
			cmdSetDHCP(1)
		}else if NetworkParamList.NetworkParam[1].DHCP == "0"{
			//开启静态IP
			cmdSetStaticIP(1)
		}
	}

	return true
}

func findNetCard(name string) (string, error) {
	if runtime.GOOS == "linux" {
		inters, err := net.Interfaces()
		if err != nil {
			return "", err
		}
		for _, v := range inters {
			if v.Name == name {
				return name, nil
			}
		}
	}
	return "", fmt.Errorf("not support GOOS(%s) and GOARCH(%s)",
		runtime.GOOS, runtime.GOARCH)
}

// HardwareAddr get mac address, if failed,it is empty
func HardwareAddr(name string) (net.HardwareAddr, error) {
	netCard, err := findNetCard(name)
	if err != nil {
		return net.HardwareAddr{}, err
	}
	inter, err := net.InterfaceByName(netCard)
	if err != nil {
		return net.HardwareAddr{}, err
	}
	return inter.HardwareAddr, err
}

// NetInformation 网络信息
type NetInformation struct {
	InterName    string // 网卡名
	HardwareAddr net.HardwareAddr
	Mac          string
	IP           string
	Mask         string
	GatewayIP    string
}

// 通过网卡获得 MAC IP IPMask GatewayIP
func GetNetInformation(netName string) (NetInformation, error) {
	info := NetInformation{}

	card, err := findNetCard(netName)
	if err != nil {
		return info, err
	}
	info.InterName = card

	inter, err := net.InterfaceByName(card)
	if err != nil {
		return info, err
	}

	info.HardwareAddr = inter.HardwareAddr
	info.Mac = hex.EncodeToString(inter.HardwareAddr)

	address, err := inter.Addrs()
	if err != nil {
		return info, err
	}
	for _, addr := range address {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				info.IP = ipnet.IP.String()
				info.Mask = net.IP(ipnet.Mask).String()
			}
		}
	}
	//获取网关Ip
	out, err := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("route -n | grep %s | grep UG | awk '{print $2}'", card)).Output()
	if err != nil {
		return info, err
	}
	info.GatewayIP = strings.Trim(string(out), "\n")
	return info, nil
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func NetworkParaRead() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/networkpara.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			fmt.Println("open networkpara.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 500)
		dataCnt, err := fp.Read(data)

		//fmt.Println(string(data[:dataCnt]))

		err = json.Unmarshal(data[:dataCnt], &NetworkParamList)
		if err != nil {
			fmt.Println("networkpara unmarshal err", err)

			return false
		}
		return true
	} else {
		fmt.Println("networkpara.json is not exist")

		os.MkdirAll(exeCurDir+"/selfpara", os.ModePerm)
		fp, err := os.Create(fileDir)
		if err != nil {
			fmt.Println("create networkpara.json err", err)
			return false
		}
		defer fp.Close()

		NetworkParamList.NetworkParam = append(NetworkParamList.NetworkParam, NetworkParamTemplate{
			ID:        "1",
			Name:      "eth0",
			DHCP:      "1",
			IP:        "192.168.4.156",
			Netmask:   "255.255.255.0",
			Broadcast: "192.168.4.255"})
		NetworkParamList.NetworkParam = append(NetworkParamList.NetworkParam, NetworkParamTemplate{
			ID:        "2",
			Name:      "eth1",
			DHCP:      "1",
			IP:        "192.168.4.156",
			Netmask:   "255.255.255.0",
			Broadcast: "192.168.4.255"})
		NetworkParaWrite()

		return true
	}
}

func NetworkParaWrite() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/networkpara.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("open networkpara.json err", err)
	}
	defer fp.Close()

	sJson, _ := json.Marshal(NetworkParamList)
	fmt.Println(string(sJson))

	_, err = fp.Write(sJson)
	if err != nil {
		fmt.Println("write networkpara.json err", err)
	}
	fp.Sync()
}