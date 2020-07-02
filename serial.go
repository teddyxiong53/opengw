package main

import (
	"encoding/json"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"os"
	"strconv"
)

type SerialParamTemplate struct{
	ID       string     `json:"ID"`
	Name     string     `json:"Name"`
	BaudRate string 	`json:"BaudRate"`
	DataBits string		`json:"DataBits"`		// Data bits: 5, 6, 7 or 8 (default 8)
	StopBits string		`json:"StopBits"`		// Stop bits: 1 or 2 (default 1)
	Parity 	 string     `json:"Parity"`			// Parity: N - None, E - Even, O - Odd (default E),(The use of no parity requires 2 stop bits.)
	Timeout  string     `json:"Timeout"`		//通信超时
	Interval string		`json:"Interval"`		//通信间隔
}

type SerialInterface struct{
	SerialParam 		[]SerialParamTemplate	`json:"SerialParam"`
	SerialPort          []*serial.Port			`json:"-"`
	SerialStatus        []bool                  `json:"-"`
}

var serialInterface SerialInterface


func newSerialInterface(param SerialParamTemplate) (bool,*serial.Port){

	serialBaud,_ := strconv.Atoi(param.BaudRate)

	var serialParity serial.Parity
	switch param.Parity {
		case "N":
			serialParity = serial.ParityNone
		case "O":
			serialParity = serial.ParityOdd
		case "E":
			serialParity = serial.ParityEven
	}
	var serialStop serial.StopBits
	switch param.StopBits {
	case "1":
		serialStop = serial.Stop1
	case "1.5":
		serialStop = serial.Stop1Half
	case "2":
		serialStop = serial.Stop2
	}

	serialConfig := &serial.Config{
		Name: param.Name,
		Baud: serialBaud,
		Parity:serialParity,
		StopBits: serialStop,
	}

	serial, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Println(err)
		return false,nil
	}

	return true,serial
}

func serialParaRead() bool{

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

		err = json.Unmarshal(data[:dataCnt],&serialInterface)
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

		return false
	}
}


func serialParaWrite(){

	fileDir := exeCurDir + "/selfpara/serialpara.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println("open serialpara.json err",err)
		return
	}
	defer fp.Close()

	sJson,_ := json.Marshal(serialInterface)
	fmt.Println(string(sJson))

	_, err = fp.Write(sJson)
	if err != nil {
		fmt.Println("write serialpara.json err",err)
	}
	fmt.Println("write serialpara.json sucess")
}

func SerialInterfaceInit(){

	//读取串口配置信息
	if serialParaRead() == false{

		serialInterface.SerialParam = make([]SerialParamTemplate,0)

		serialInterface.SerialParam = append(serialInterface.SerialParam,SerialParamTemplate{
			ID       : "1",
			Name     : "/dev/ttyUSB0",
			BaudRate : "9600",
			DataBits : "8",
			StopBits : "1",
			Parity   : "N",
			Timeout  : "1000",
			Interval : "1000"})
		serialInterface.SerialParam = append(serialInterface.SerialParam,SerialParamTemplate{
			ID       : "2",
			Name     : "/dev/ttyUSB1",
			BaudRate : "9600",
			DataBits : "8",
			StopBits : "1",
			Parity   : "N",
			Timeout  : "1000",
			Interval : "1000"})

		serialParaWrite()
	}


	//打开串口
	serialInterface.SerialPort = make([]*serial.Port,0)
	serialInterface.SerialStatus = make([]bool,0)
	//for _,v := range serialInterface.SerialParam{
	//
	//	status,port := newSerialInterface(v)
	//	serialInterface.SerialStatus = append(serialInterface.SerialStatus,status)
	//	if status == true{
	//		serialInterface.SerialPort = append(serialInterface.SerialPort,port)
	//	}
	//}
}

func serialOpen(index int){

	if index > len(serialInterface.SerialParam){
		log.Println("serial index is noexist")
		return
	}

	status,port := newSerialInterface(serialInterface.SerialParam[index])
	serialInterface.SerialStatus = append(serialInterface.SerialStatus,status)
	if status == true{
		serialInterface.SerialPort = append(serialInterface.SerialPort,port)
	}
}

