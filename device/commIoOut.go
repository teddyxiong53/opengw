package device

import (
	"encoding/json"
	"goAdapter/setting"
	"log"
	"os"
	"path/filepath"
)

type IoOutInterfaceParam struct {
	Name string   `json:"Name"`
	FD   *os.File `json:"-"`
}

type CommunicationIoOutTemplate struct {
	Name  string              `json:"Name"`  //接口名称
	Type  string              `json:"Type"`  //接口类型,比如serial,IoOut,udp,http
	Param IoOutInterfaceParam `json:"Param"` //接口参数
}

var CommunicationIoOutMap = make([]*CommunicationIoOutTemplate, 0)

func (c *CommunicationIoOutTemplate) Open() bool {

	fd, err := os.OpenFile(c.Param.Name, os.O_RDWR, 0666)
	if err != nil {
		setting.Logger.Errorf("IoOut open err,%v", err)
		return false
	}
	c.Param.FD = fd

	return true
}

func (c *CommunicationIoOutTemplate) Close() bool {

	if c.Param.FD != nil {
		err := c.Param.FD.Close()
		if err != nil {
			setting.Logger.Errorf("IoOut close err,%v", err)
			return false
		}
	}

	return true
}

func (c *CommunicationIoOutTemplate) WriteData(data []byte) int {

	if c.Param.FD != nil {
		//setting.Logger.Debugf("IoOut write %v", data)
		if len(data) > 0 {
			_, err := c.Param.FD.Write(data)
			if err != nil {
				setting.Logger.Errorf("IoOut write err,%v", err)
			}
		}
		return 0
	}
	return 0
}

func (c *CommunicationIoOutTemplate) ReadData(data []byte) int {

	return 0
}

func (c *CommunicationIoOutTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationIoOutTemplate) GetTimeOut() string {
	return ""
}

func (c *CommunicationIoOutTemplate) GetInterval() string {
	return ""
}

func ReadCommIoOutInterfaceListFromJson() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/commIoOutInterface.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open commIoOutInterface.json err", err)
			return false
		}
		defer func(fp *os.File) {
			err = fp.Close()
			if err != nil {

			}
		}(fp)

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt], &CommunicationIoOutMap)
		if err != nil {
			log.Println("commIoOutInterface unmarshal err", err)
			return false
		}
		return true
	} else {
		log.Println("commIoOutInterface.json is not exist")

		return false
	}
}

func WriteCommIoOutInterfaceListToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/commIoOutInterface.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open commIoOutInterface.json err", err)
		return
	}
	defer func(fp *os.File) {
		err = fp.Close()
		if err != nil {

		}
	}(fp)

	sJson, _ := json.Marshal(CommunicationIoOutMap)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write commIoOutInterface.json err", err)
	}
	setting.Logger.Infof("write commIoOutInterface.json sucess")
}
