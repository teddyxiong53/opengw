package device

import (
	"github.com/tarm/serial"
	"log"
	"reflect"
)

type SerialInterfaceParam struct{
	Name     string 		`json:"Name"`
	BaudRate string 		`json:"BaudRate"`
	DataBits string			`json:"DataBits"`		//数据位: 5, 6, 7 or 8 (default 8)
	StopBits string			`json:"StopBits"`		//停止位: 1 or 2 (default 1)
	Parity 	 string     	`json:"Parity"`			//校验: N - None, E - Even, O - Odd (default E),(The use of no parity requires 2 stop bits.)
	Timeout  string     	`json:"Timeout"`		//通信超时
	Interval string			`json:"Interval"`		//通信间隔
}

type CommunicationSerialInterface struct{
	Name     string 		`json:"Name"`
	BaudRate string 		`json:"BaudRate"`
	DataBits string			`json:"DataBits"`		//数据位: 5, 6, 7 or 8 (default 8)
	StopBits string			`json:"StopBits"`		//停止位: 1 or 2 (default 1)
	Parity 	 string     	`json:"Parity"`			//校验: N - None, E - Even, O - Odd (default E),(The use of no parity requires 2 stop bits.)
	Timeout  string     	`json:"Timeout"`		//通信超时
	Interval string			`json:"Interval"`		//通信间隔
	Port     *serial.Port				`json:"-"`				//通信句柄
}

func (c *CommunicationSerialInterface)Open(param interface{}) bool{


	log.Println("  ",reflect.TypeOf(param))
	//log.Printf("Name is %s\n",serialParam.FieldByName("Name"))

	/*
	serialParam := param.(CommunicationSerialInterface)
	serialBaud,_ := strconv.Atoi(serialParam.BaudRate)

	var serialParity serial.Parity
	switch serialParam.Parity {
	case "N":
		serialParity = serial.ParityNone
	case "O":
		serialParity = serial.ParityOdd
	case "E":
		serialParity = serial.ParityEven
	}

	var serialStop serial.StopBits
	switch serialParam.StopBits {
	case "1":
		serialStop = serial.Stop1
	case "1.5":
		serialStop = serial.Stop1Half
	case "2":
		serialStop = serial.Stop2
	}

	serialConfig := &serial.Config{
		Name: serialParam.Name,
		Baud: serialBaud,
		Parity:serialParity,
		StopBits: serialStop,
		ReadTimeout: time.Millisecond*1,
	}

	serial, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Printf("open serial err,%s",err)
		return false
	}
	serialParam.Port = serial
	 */
	return true
}

func (c *CommunicationSerialInterface)Close() bool{

	return true
}

func (c *CommunicationSerialInterface)WriteData(data []byte) int{

	cnt,_ := c.Port.Write(data)

	return cnt
}

func (c *CommunicationSerialInterface)ReadData(data []byte) int{

	cnt,_ := c.Port.Read(data)

	return cnt
}
