package setting

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type NetworkNameListTemplate struct {
	Name []string `json:"Name"`
}

type NetworkParamTemplate struct {
	Name         string           `json:"Name"` // 网卡名
	DHCP         string           `json:"DHCP"`
	IP           string           `json:"IP"`
	Netmask      string           `json:"Netmask"`
	Broadcast    string           `json:"Broadcast"`
	Gateway      string           `json:"Gateway"`
	MAC          string           `json:"MAC"`
	LinkStatus   uint32           `json:"-"`
	NetFlags     net.Flags        `json:"-"`
	HardwareAddr net.HardwareAddr `json:"-"`
}

type NetworkParamListTemplate struct {
	NetworkParam []*NetworkParamTemplate
}

var NetworkNameList = NetworkNameListTemplate{}
var NetworkParamList = &NetworkParamListTemplate{
	NetworkParam: make([]*NetworkParamTemplate, 0),
}

func init() {

}

func (n *NetworkParamListTemplate) AddNetworkParam(param NetworkParamTemplate) error {

	for _, v := range n.NetworkParam {
		if v.Name == param.Name {
			return errors.New("网络名已经存在")
		}
	}
	n.NetworkParam = append(n.NetworkParam, &param)
	NetworkParaWrite()
	return nil
}

func (n *NetworkParamTemplate) GetNetworkStatus() {

	//if runtime.GOOS == "linux" {
	//	ethHandle, _ := ethtool.NewEthtool()
	//	defer ethHandle.Close()
	//
	//	n.LinkStatus, _ = ethHandle.LinkState(n.Name)
	//	Logger.Debugf("%v LinkStatus %v", n.Name, n.LinkStatus)
	//}
}

//获取当前网络参数
func (n *NetworkParamListTemplate) GetNetworkParam() {

	for _, v := range n.NetworkParam {
		ethInfo, err := GetNetInformation(v.Name)
		if err != nil {
			Logger.Errorf("getNetInfor err,%v\n", err)
			break
		}
		v.IP = ethInfo.IP
		v.Netmask = ethInfo.Netmask
		v.Broadcast = ethInfo.Gateway
		v.MAC = strings.ToUpper(ethInfo.MAC)
		//Logger.Debugf("%v netFlags %v", v.Name, ethInfo.NetFlags)
		if runtime.GOOS == "linux" {
			v.GetNetworkStatus()
		}
	}
}

//设置网络参数
func (n *NetworkParamListTemplate) ModifyNetworkParam(param NetworkParamTemplate) {

	for k, v := range n.NetworkParam {
		if v.Name == param.Name {
			n.NetworkParam[k].DHCP = param.DHCP
			n.NetworkParam[k].IP = param.IP
			n.NetworkParam[k].Netmask = param.Netmask
			n.NetworkParam[k].Gateway = param.Gateway

			NetworkParaWrite()
		}
	}

}

//删除网络参数
func (n *NetworkParamListTemplate) DeleteNetworkParam(name string) (bool, string) {

	for k, v := range n.NetworkParam {
		if v.Name == name {
			n.NetworkParam = append(n.NetworkParam[:k], n.NetworkParam[k+1:]...)
			NetworkParaWrite()
			return true, ""
		}
	}

	return false, "name is not exist"
}

func (n *NetworkParamTemplate) CmdSetDHCP() {

	//cmd := exec.Command("udhcpc","-i",n.Name,"5")
	cmd := exec.Command("udhcpc", "-i", n.Name)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run() //执行到此会阻塞5s

	str := out.String()

	Logger.Debugf(str)
}

func (n *NetworkParamTemplate) CmdSetStaticIP() {

	strNetMask := "netmask " + n.Netmask
	cmd := exec.Command("ifconfig",
		n.Name,
		n.IP,
		strNetMask)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Start() //执行到此,直接往后执行

	cmd2 := exec.Command("/sbin/route", "add", "default", "gw", n.Broadcast)
	cmd2.Stdout = &out
	cmd2.Start() //执行到此,直接往后执行
}

func findNetCard(name string) (string, error) {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
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

// 通过网卡获得 MAC IP IPMask GatewayIP
func GetNetInformation(netName string) (NetworkParamTemplate, error) {
	info := NetworkParamTemplate{}

	card, err := findNetCard(netName)
	if err != nil {
		return info, err
	}
	info.Name = card

	inter, err := net.InterfaceByName(card)
	if err != nil {
		return info, err
	}
	Logger.Tracef("inter %v\n", *inter)
	info.HardwareAddr = inter.HardwareAddr
	info.MAC = hex.EncodeToString(inter.HardwareAddr)
	info.NetFlags = inter.Flags

	address, err := inter.Addrs()
	if err != nil {
		return info, err
	}
	for _, addr := range address {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				info.IP = ipnet.IP.String()
				info.Netmask = net.IP(ipnet.Mask).String()
			}
		}
	}
	//获取网关Ip
	out, err := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("route -n | grep %s | grep UG | awk '{print $2}'", card)).Output()
	if err != nil {
		return info, err
	}
	info.Gateway = strings.Trim(string(out), "\n")
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
			Logger.Errorf("open networkpara.json err,%v", err)
			return false
		}
		defer fp.Close()
		Logger.Infof("open networkpara.json ok")
		data := make([]byte, 500)
		dataCnt, err := fp.Read(data)

		//fmt.Println(string(data[:dataCnt]))

		err = json.Unmarshal(data[:dataCnt], &NetworkParamList)
		if err != nil {
			Logger.Errorf("networkpara unmarshal err,%v", err)

			return false
		}

		for _, v := range NetworkParamList.NetworkParam {
			if v.DHCP == "1" {
				v.CmdSetDHCP()
			} else if v.DHCP == "0" {
				v.CmdSetStaticIP()
			}
		}

		return true
	} else {
		Logger.Errorf("networkpara.json is not exist")

		return true
	}
}

func NetworkParaWrite() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/networkpara.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		Logger.Warnf("open networkpara.json err,%v", err)
	}
	defer fp.Close()

	sJson, _ := json.Marshal(NetworkParamList)
	Logger.Debugf(string(sJson))

	_, err = fp.Write(sJson)
	if err != nil {
		Logger.Warnf("write networkpara.json err,%v", err)
	}
	Logger.Debugf("write networkpara.json ok")
	fp.Sync()
}
