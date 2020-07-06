package setting

import (
	"github.com/tarm/serial"
	"log"
	"strconv"
	"time"
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

type SerialInterfaceTemplate struct{
	SerialParam 		[]SerialParamTemplate	`json:"SerialParam"`
	SerialPort          []*serial.Port			`json:"-"`
	SerialStatus        []bool                  `json:"-"`
}

var SerialInterface SerialInterfaceTemplate


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
		ReadTimeout: time.Millisecond*1,
	}

	serial, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Println(err)
		return false,nil
	}

	return true,serial
}



func SerialInterfaceInit(){


	//打开串口
	SerialInterface.SerialPort = make([]*serial.Port,0)
	SerialInterface.SerialStatus = make([]bool,0)
	//for _,v := range serialInterface.SerialParam{
	//
	//	status,port := newSerialInterface(v)
	//	serialInterface.SerialStatus = append(serialInterface.SerialStatus,status)
	//	if status == true{
	//		serialInterface.SerialPort = append(serialInterface.SerialPort,port)
	//	}
	//}
}

func SerialOpen(index int){

	if index > len(SerialInterface.SerialParam){
		log.Println("serial index is noexist")
		return
	}

	status,port := newSerialInterface(SerialInterface.SerialParam[index])
	SerialInterface.SerialStatus = append(SerialInterface.SerialStatus,status)
	if status == true{
		SerialInterface.SerialPort = append(SerialInterface.SerialPort,port)
	}
}

