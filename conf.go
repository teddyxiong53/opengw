package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var exeCurDir    string


/**************获取配置信息************************/
func getConf(){
	exeCurDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println(exeCurDir)

	//if serialParaRead(&serialParamList) == true{
	//	fmt.Println("read serialParam",serialParamList)
	//}

	if networkParaRead(&networkParamList) == true{
		for _,v := range networkParamList.NetworkParam{
			log.Printf("networkParam %s,%+v\n",v.Name,v)
		}
	}
}


func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

/*
func serialParaRead(param *SerialParamList) bool{

	fileDir := exeCurDir + "/selfpara/serialpara.json"

	if FileExist(fileDir) == true{
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil{
			fmt.Println("open serialpara.json err",err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 500)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt],param)
		if err != nil{
			fmt.Println("serialpara unmarshal err",err)

			serialParamList.SerialParam = append(serialParamList.SerialParam,SerialParam{
				Name     : "/dev/ttyUSB1",
				BaudRate : "9600",
				DataBits : "8",
				StopBits : "1",
				Parity   : "N",
				Timeout  : "1000"})
			serialParaWrite(serialParamList)

			return false
		}
		return true
	}else{
		fmt.Println("/opt/ibox/selfpara/serialpara.json is not exist")

		os.MkdirAll(exeCurDir+"/selfpara", os.ModePerm)
		fileDir = exeCurDir + "/selfpara/serialpara.json"
		fp, err := os.Create(fileDir)
		if err != nil{
			fmt.Println("create serialpara.json err",err)
			return false
		}
		defer fp.Close()

		serialParamList.SerialParam = append(serialParamList.SerialParam,SerialParam{
			ID       : "1",
			Name     : "/dev/ttyUSB1",
			BaudRate : "9600",
			DataBits : "8",
			StopBits : "1",
			Parity   : "N",
			Timeout  : "1000"})
		serialParamList.SerialParam = append(serialParamList.SerialParam,SerialParam{
			ID       : "2",
			Name     : "/dev/ttyUSB2",
			BaudRate : "9600",
			DataBits : "8",
			StopBits : "1",
			Parity   : "N",
			Timeout  : "1000"})
		serialParaWrite(serialParamList)

		return true
	}
}


func serialParaWrite(param SerialParamList){

	fileDir := exeCurDir + "/selfpara/serialpara.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("open serialpara.json err",err)
		return
	}
	defer fp.Close()

	sJson,_ := json.Marshal(param)
	fmt.Println(string(sJson))

	_, err = fp.Write(sJson)
	if err != nil {
		fmt.Println("write serialpara.json err",err)
	}
	fmt.Println("write serialpara.json sucess")
}

 */

func networkParaRead(param *NetworkParamList) bool{

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

		err = json.Unmarshal(data[:dataCnt], param)
		if err != nil {
			fmt.Println("networkpara unmarshal err", err)

			return false
		}
		return true
	}else{
		fmt.Println("networkpara.json is not exist")


		os.MkdirAll(exeCurDir + "/selfpara", os.ModePerm)
		fp, err := os.Create(fileDir)
		if err != nil{
			fmt.Println("create networkpara.json err",err)
			return false
		}
		defer fp.Close()

		networkParamList.NetworkParam = append(networkParamList.NetworkParam,NetworkParam{
			ID        : "1",
			Name      : "eth0",
			DHCP      : "1",
			IP        : "192.168.4.156",
			Netmask   : "255.255.255.0",
			Broadcast : "192.168.4.255"})
		networkParamList.NetworkParam = append(networkParamList.NetworkParam,NetworkParam{
			ID        : "2",
			Name      : "eth1",
			DHCP      : "1",
			IP        : "192.168.4.156",
			Netmask   : "255.255.255.0",
			Broadcast : "192.168.4.255"})
		networkParaWrite(networkParamList)

		return true
	}
}

func networkParaWrite(param NetworkParamList){

	fileDir := exeCurDir + "/selfpara/networkpara.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("open networkpara.json err",err)
	}
	defer fp.Close()

	sJson,_ := json.Marshal(param)
	fmt.Println(string(sJson))

	_, err = fp.Write(sJson)
	if err != nil {
		fmt.Println("write networkpara.json err",err)
	}
	fp.Sync()
}

