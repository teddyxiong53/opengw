package device

import (
	"encoding/json"
	"github.com/tarm/serial"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

type CommunicationSerialTemplate struct{
	CommunicationTemplate
	Param   SerialInterfaceParam     					`json:"Param"`			//接口参数
	Port    *serial.Port								`json:"-"`				//通信句柄
}

var CommunicationSerialMap = make([]CommunicationSerialTemplate,0)

func (c *CommunicationSerialTemplate)Open() bool{

	serialParam := c.Param
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

	serialPort, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Printf("open serial err,%s",err)
		return false
	}else{
		log.Printf("open serial %s ok\n",c.Param.Name)
	}

	c.Port = serialPort
	return true
}

func (c *CommunicationSerialTemplate)Close() bool{

	return true
}

func (c *CommunicationSerialTemplate)WriteData(data []byte) int{

	cnt,_ := c.Port.Write(data)

	return cnt
}

func (c *CommunicationSerialTemplate)ReadData(data []byte) int{

	cnt,_ := c.Port.Read(data)

	return cnt
}

func NewCommunicationSerialTemplate(commName,commType string,param SerialInterfaceParam) *CommunicationSerialTemplate{

	return &CommunicationSerialTemplate{
		Param:param,
		CommunicationTemplate:CommunicationTemplate{
			Name:commName,
			Type:commType,
		},
	}
}

func ReadCommSerialInterfaceListFromJson() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/commSerialInterface.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open commSerialInterface.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		//CommunicationSerialTemplateList.CommunicationSerialMap = make([]CommunicationSerialTemplate,0)

		err = json.Unmarshal(data[:dataCnt], &CommunicationSerialMap)
		if err != nil {
			log.Println("commSerialInterface unmarshal err", err)
			return false
		}
		//log.Printf("CommunicationSerialMap %+v\n",CommunicationSerialTemplateList.CommunicationSerialMap)
		return true
	} else {
		log.Println("commSerialInterface.json is not exist")

		return false
	}
}

func WriteCommSerialInterfaceListToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/commSerialInterface.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open commSerialInterface.json err", err)
		return
	}
	defer fp.Close()

	sJson, _ := json.Marshal(CommunicationSerialMap)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write commSerialInterface.json err", err)
	}
	log.Println("write commSerialInterface.json sucess")
}

