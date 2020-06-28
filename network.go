package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/safchain/ethtool"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

type NetworkParam struct{
	ID   string         `json:ID`
	Name string         `json:"Name"`
	DHCP string         `json:"DHCP"`
	IP string           `json:"IP"`
	Netmask string      `json:"Netmask"`
	Broadcast string    `json:"Broadcast"`
	MAC string          `json:"MAC"`
}

type NetworkParamList struct{
	NetworkParam []NetworkParam
}

type NetworkLinkState struct{
	State [2]uint32
}

var networkParamList NetworkParamList
var networkLinkState NetworkLinkState


func cmdSetDHCP(id int){

	cmd := exec.Command("udhcpc","-i",networkParamList.NetworkParam[id].Name,"5")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()	//执行到此会阻塞5s

	str := out.String()

	log.Println(str)
}

func cmdSetStaticIP(id int){

	strNetMask   := "netmask " + networkParamList.NetworkParam[id].Netmask
	cmd := exec.Command("ifconfig",
		networkParamList.NetworkParam[id].Name,
		networkParamList.NetworkParam[id].IP,
		strNetMask)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Start()	//执行到此,直接往后执行

	cmd2 := exec.Command("/sbin/route","add","default","gw",networkParamList.NetworkParam[id].Broadcast)
	cmd2.Stdout = &out
	cmd2.Start()	//执行到此,直接往后执行
}

func getNetworkStatus(){

	for _,v := range networkParamList.NetworkParam{
		getLinkState(v.ID)
	}
}

func getLinkState(id string){

	ethHandle, _ := ethtool.NewEthtool()
	defer ethHandle.Close()

	var state uint32

	if id == "1"{
		state, _ = ethHandle.LinkState(networkParamList.NetworkParam[0].Name)
		networkLinkState.State[0] = state
	}else if id == "2"{
		state, _ = ethHandle.LinkState(networkParamList.NetworkParam[1].Name)
		networkLinkState.State[1] = state
	}
}


//获取当前网络参数
func getNetworkParam() NetworkParamList{

	for k,v := range networkParamList.NetworkParam{

		getLinkState(v.ID)
		ethInfo,err := GetNetInformation(v.Name)
		if err == nil{
			networkParamList.NetworkParam[k].IP = ethInfo.IP
			networkParamList.NetworkParam[k].Netmask = ethInfo.Mask
			networkParamList.NetworkParam[k].Broadcast = ethInfo.GatewayIP
			networkParamList.NetworkParam[k].MAC = ethInfo.Mac
		}
	}

	return networkParamList
}

//设置网络参数
func setNetworkParam(id string,param NetworkParam){

	getLinkState(id)

	if id == "1"{
		if networkLinkState.State[0] == 0{
			log.Printf("setNetworkParam %s err\n",id)
			return
		}

		networkParamList.NetworkParam[0] = param

		if networkParamList.NetworkParam[0].DHCP == "1"{
			//开启动态IP

			cmdSetDHCP(0)
		}else if networkParamList.NetworkParam[0].DHCP == "0"{
			//开启静态IP

			cmdSetStaticIP(0)
		}
	}else if id == "2"{
		if networkLinkState.State[1] == 0{
			log.Printf("setNetworkParam %s err\n",id)
			return
		}

		networkParamList.NetworkParam[1] = param

		if networkParamList.NetworkParam[1].DHCP == "1"{
			//开启动态IP

			cmdSetDHCP(1)
		}else if networkParamList.NetworkParam[1].DHCP == "0"{
			//开启静态IP

			cmdSetStaticIP(1)
		}
	}

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