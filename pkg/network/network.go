package network

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"goAdapter/pkg/mylog"
	"goAdapter/pkg/system"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jackpal/gateway"
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

func (n *NetworkParamListTemplate) AddNetworkParam(param NetworkParamTemplate) error {

	for _, v := range n.NetworkParam {
		if v.Name == param.Name {
			return fmt.Errorf("network card %s already exists", v.Name)
		}
	}
	if err := param.CmdSetStaticIP(); err != nil {
		return err
	}
	n.NetworkParam = append(n.NetworkParam, &param)
	return NetworkParaWrite()
}

//设置网络参数
func (n *NetworkParamListTemplate) ModifyNetworkParam(param NetworkParamTemplate) error {
	if err := param.CmdSetStaticIP(); err != nil {
		return err
	}
	var index = -1
	for i, n := range n.NetworkParam {
		if n.Name == param.Name {
			index = i
		}
	}
	if index != -1 {
		n.NetworkParam[index] = &param
	} else {
		//如果原来没有这个网卡在json中会添加进去
		n.NetworkParam = append(n.NetworkParam, &param)
	}
	return NetworkParaWrite()
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

	mylog.Logger.Debugf(str)
}

func (n *NetworkParamTemplate) CmdSetStaticIP() error {
	fmt.Println("goos:", runtime.GOOS)
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", fmt.Sprintf("netsh interface ip set address %s static %s %s %s", n.Name, n.IP, n.Netmask, n.Gateway))
		msg, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("set ip on %s error:%v", runtime.GOOS, msg)
		}
		return nil

	case "linux":
		strNetMask := "netmask " + n.Netmask
		cmd := exec.Command("ifconfig",
			n.Name,
			n.IP,
			strNetMask)

		msg, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("ifconfig error:%s", msg)
		}

		cmd2 := exec.Command("/sbin/route", "add", "default", "gw", n.Gateway)
		msg, err = cmd2.CombinedOutput()
		if err != nil {
			return fmt.Errorf("add default gw error:%s", msg)
		}
		return nil
	default:
		return fmt.Errorf("%s is not supported", runtime.GOOS)

	}
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
	mylog.Logger.Tracef("inter %v\n", *inter)
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

func NetworkParaRead() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/networkpara.json"

	if system.FileExist(fileDir) {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			mylog.Logger.Errorf("open networkpara.json err,%v", err)
			return false
		}
		defer fp.Close()
		mylog.Logger.Infof("open networkpara.json ok")
		data := make([]byte, 500)
		dataCnt, err := fp.Read(data)

		//fmt.Println(string(data[:dataCnt]))

		err = json.Unmarshal(data[:dataCnt], &NetworkParamList)
		if err != nil {
			mylog.Logger.Errorf("networkpara unmarshal err,%v", err)

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
		mylog.Logger.Errorf("networkpara.json is not exist")

		return true
	}
}

func NetworkParaWrite() error {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/networkpara.json"

	fp, err := os.Create(fileDir)
	if err != nil {
		return fmt.Errorf("create networkpara.json err,%v", err)
	}
	defer fp.Close()

	sJson, _ := json.Marshal(NetworkParamList)
	_, err = fp.Write(sJson)
	if err != nil {
		return fmt.Errorf("write networkpara.json err,%v", err)
	}
	fp.Sync()
	return nil

}

func ParseNetworks() ([]*NetworkParamTemplate, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	networkParams := make([]*NetworkParamTemplate, 0)
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPAddr:
				if v.IP.IsLoopback() || v.IP.To4() == nil {
					continue
				}
				networkParam := NetworkParamTemplate{
					MAC:      iface.HardwareAddr.String(),
					Name:     iface.Name,
					NetFlags: iface.Flags,
					Netmask:  v.IP.DefaultMask().String(),
					IP:       v.IP.String(),
				}
				gw, err := gateway.DiscoverGateway()
				if err == nil {
					networkParam.Gateway = gw.String()
				} else {
					fmt.Println(err)
				}
				mask, err := parseMask(v.IP.DefaultMask())
				if err == nil {
					networkParam.Netmask = mask
				}
				networkParams = append(networkParams, &networkParam)

			case *net.IPNet:
				if v.IP.IsLoopback() || v.IP.To4() == nil {
					continue
				}
				networkParam := NetworkParamTemplate{
					MAC:      iface.HardwareAddr.String(),
					Name:     iface.Name,
					NetFlags: iface.Flags,
					IP:       v.IP.String(),
				}
				gw, err := gateway.DiscoverGateway()
				if err == nil {
					networkParam.Gateway = gw.String()
				} else {
					fmt.Println(err)
				}
				mask, err := parseMask(v.IP.DefaultMask())
				if err == nil {
					networkParam.Netmask = mask
				}
				networkParams = append(networkParams, &networkParam)

			}

		}
	}
	return networkParams, nil
}

func parseMask(mask net.IPMask) (string, error) {
	if len(mask) != 4 {
		return "", fmt.Errorf("not ipv4 mask type")
	}
	return fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3]), nil
}
