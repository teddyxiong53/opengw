package device

import (
	"encoding/json"
	"goAdapter/setting"
	"io"
	"log"
	"os"
	"path/filepath"
)

type IoInInterfaceParam struct {
	Name string   `json:"Name"`
	FD   *os.File `json:"-"`
}

type CommunicationIoInTemplate struct {
	Name  string             `json:"Name"`  //接口名称
	Type  string             `json:"Type"`  //接口类型,比如serial,IoIn,udp,http
	Param IoInInterfaceParam `json:"Param"` //接口参数
}

var CommunicationIoInMap = make([]*CommunicationIoInTemplate, 0)

func (c *CommunicationIoInTemplate) Open() bool {

	fd, err := os.OpenFile(c.Param.Name, os.O_RDWR, 0777)
	if err != nil {
		setting.Logger.Errorf("IoIn open err,%v", err)
		return false
	}
	c.Param.FD = fd

	return true
}

func (c *CommunicationIoInTemplate) Close() bool {

	if c.Param.FD != nil {
		err := c.Param.FD.Close()
		if err != nil {
			setting.Logger.Errorf("IoIn close err,%v", err)
			return false
		}
	}

	return true
}

func (c *CommunicationIoInTemplate) WriteData(data []byte) int {

	if c.Param.FD != nil {
		//setting.Logger.Debugf("IoIn write %v", data)
		if len(data) > 0 {
			_, err := c.Param.FD.Write(data)
			if err != nil {
				setting.Logger.Errorf("IoIn write err,%v", err)
			}
		}
		return 0
	}
	return 0
}

func (c *CommunicationIoInTemplate) ReadData(data []byte) int {

	if c.Param.FD != nil {
		_, err := c.Param.FD.Seek(0, 0)
		if err != nil {
			setting.Logger.Errorf("IoIn seek err,%v", err)
		}
		cnt, err := c.Param.FD.Read(data)
		if err != nil {
			if err != io.EOF {
				setting.Logger.Errorf("IoIn read err,%v", err)
			}
		}
		//if cnt > 0 {
		//	setting.Logger.Errorf("IoIn read data,%v", data[:cnt])
		//}
		return cnt
	}
	return 0
}

func (c *CommunicationIoInTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationIoInTemplate) GetTimeOut() string {
	return ""
}

func (c *CommunicationIoInTemplate) GetInterval() string {
	return ""
}

func ReadCommIoInInterfaceListFromJson() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/commIoInInterface.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open commIoInInterface.json err", err)
			return false
		}
		defer func(fp *os.File) {
			err = fp.Close()
			if err != nil {

			}
		}(fp)

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt], &CommunicationIoInMap)
		if err != nil {
			log.Println("commIoInInterface unmarshal err", err)
			return false
		}
		return true
	} else {
		log.Println("commIoInInterface.json is not exist")

		return false
	}
}

func WriteCommIoInInterfaceListToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/commIoInInterface.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open commIoInInterface.json err", err)
		return
	}
	defer func(fp *os.File) {
		err = fp.Close()
		if err != nil {

		}
	}(fp)

	sJson, _ := json.Marshal(CommunicationIoInMap)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write commIoInInterface.json err", err)
	}
	setting.Logger.Infof("write commIoInInterface.json sucess")
}
