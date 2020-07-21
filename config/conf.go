package config

import (
	"encoding/json"
	"fmt"
	"goAdapter/setting"
	"log"
	"os"
	"path/filepath"
)

var exeCurDir string

/**************获取配置信息************************/
func GetConf() {
	exeCurDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println(exeCurDir)

	if NetworkParaRead() == true {
		for _, v := range setting.NetworkParamList.NetworkParam {
			log.Printf("networkParam %s,%+v\n", v.Name, v)
		}
	}
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func NetworkParaRead() bool {

	fileDir := exeCurDir + "/selfpara/networkpara.json"

	if FileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			fmt.Println("open networkpara.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 500)
		dataCnt, err := fp.Read(data)

		//fmt.Println(string(data[:dataCnt]))

		err = json.Unmarshal(data[:dataCnt], &setting.NetworkParamList)
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

		setting.NetworkParamList.NetworkParam = append(setting.NetworkParamList.NetworkParam, setting.NetworkParamTemplate{
			ID:        "1",
			Name:      "eth0",
			DHCP:      "1",
			IP:        "192.168.4.156",
			Netmask:   "255.255.255.0",
			Broadcast: "192.168.4.255"})
		setting.NetworkParamList.NetworkParam = append(setting.NetworkParamList.NetworkParam, setting.NetworkParamTemplate{
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

	fileDir := exeCurDir + "/selfpara/networkpara.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("open networkpara.json err", err)
	}
	defer fp.Close()

	sJson, _ := json.Marshal(setting.NetworkParamList)
	fmt.Println(string(sJson))

	_, err = fp.Write(sJson)
	if err != nil {
		fmt.Println("write networkpara.json err", err)
	}
	fp.Sync()
}
