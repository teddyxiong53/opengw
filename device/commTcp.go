package device

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type TcpInterfaceParam struct{
	Name     	string     	`json:"Name"`
	IP 			string 		`json:"IP"`
	Port 		string		`json:"Port"`
	Timeout  	string     	`json:"Timeout"`		//通信超时
	Interval 	string		`json:"Interval"`		//通信间隔
}

type CommunicationTcpTemplate struct{
	CommunicationTemplate
	Param   TcpInterfaceParam     					`json:"Param"`			//接口参数
}

var CommunicationTcpMap = make([]CommunicationTcpTemplate,0)

func (c *CommunicationTcpTemplate)Open() bool{

	return true
}

func (c *CommunicationTcpTemplate)Close() bool{

	return true
}

func (c *CommunicationTcpTemplate)WriteData(data []byte) int{

	return 0
}

func (c *CommunicationTcpTemplate)ReadData(data []byte) int{

	return 0
}

func ReadCommTcpInterfaceListFromJson() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/commTcpInterface.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open commTcpInterface.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt], &CommunicationTcpMap)
		if err != nil {
			log.Println("commTcpInterface unmarshal err", err)
			return false
		}
		return true
	} else {
		log.Println("commTcpInterface.json is not exist")

		return false
	}
}

func WriteCommTcpInterfaceListToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/commTcpInterface.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open commTcpInterface.json err", err)
		return
	}
	defer fp.Close()

	sJson, _ := json.Marshal(CommunicationTcpMap)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write commTcpInterface.json err", err)
	}
	log.Println("write commTcpInterface.json sucess")
}