package config

import (
	"encoding/json"
	"fmt"
	"goAdapter/setting"
	"log"
	"os"
	"path/filepath"
)

var exeCurDir    string


/**************获取配置信息************************/
func GetConf(){
	exeCurDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println(exeCurDir)

	//if serialParaRead(&serialParamList) == true{
	//	fmt.Println("read serialParam",serialParamList)
	//}

	if SerialParaRead() == true{

	}

	if NetworkParaRead() == true{
		for _,v := range setting.NetworkParamList.NetworkParam{
			log.Printf("networkParam %s,%+v\n",v.Name,v)
		}
	}
}


func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func SerialParaRead() bool{

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

		setting.SerialInterface.SerialParam = make([]setting.SerialParamTemplate,0)
		err = json.Unmarshal(data[:dataCnt],&setting.SerialInterface)
		if err != nil{
			fmt.Println("serialpara unmarshal err",err)

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

		setting.SerialInterface.SerialParam = make([]setting.SerialParamTemplate,0)

		setting.SerialInterface.SerialParam = append(setting.SerialInterface.SerialParam,setting.SerialParamTemplate{
			ID       : "1",
			Name     : "/dev/ttyUSB0",
			BaudRate : "9600",
			DataBits : "8",
			StopBits : "1",
			Parity   : "N",
			Timeout  : "1000",
			Interval : "1000"})

		SerialParaWrite()

		return false
	}
}


func SerialParaWrite(){

	fileDir := exeCurDir + "/selfpara/serialpara.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("open serialpara.json err",err)
		return
	}
	defer fp.Close()

	sJson,_ := json.Marshal(setting.SerialInterface)
	fmt.Println(string(sJson))

	_, err = fp.Write(sJson)
	if err != nil {
		fmt.Println("write serialpara.json err",err)
	}
	fmt.Println("write serialpara.json sucess")
}

func NetworkParaRead() bool{

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
	}else{
		fmt.Println("networkpara.json is not exist")


		os.MkdirAll(exeCurDir + "/selfpara", os.ModePerm)
		fp, err := os.Create(fileDir)
		if err != nil{
			fmt.Println("create networkpara.json err",err)
			return false
		}
		defer fp.Close()

		setting.NetworkParamList.NetworkParam = append(setting.NetworkParamList.NetworkParam,setting.NetworkParamTemplate{
			ID        : "1",
			Name      : "eth0",
			DHCP      : "1",
			IP        : "192.168.4.156",
			Netmask   : "255.255.255.0",
			Broadcast : "192.168.4.255"})
		setting.NetworkParamList.NetworkParam = append(setting.NetworkParamList.NetworkParam,setting.NetworkParamTemplate{
			ID        : "2",
			Name      : "eth1",
			DHCP      : "1",
			IP        : "192.168.4.156",
			Netmask   : "255.255.255.0",
			Broadcast : "192.168.4.255"})
		NetworkParaWrite()

		return true
	}
}

func NetworkParaWrite(){

	fileDir := exeCurDir + "/selfpara/networkpara.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("open networkpara.json err",err)
	}
	defer fp.Close()

	sJson,_ := json.Marshal(setting.NetworkParamList)
	fmt.Println(string(sJson))

	_, err = fp.Write(sJson)
	if err != nil {
		fmt.Println("write networkpara.json err",err)
	}
	fp.Sync()
}

