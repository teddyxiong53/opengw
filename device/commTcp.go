package device

import (
	"encoding/json"
	"goAdapter/setting"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

type TcpInterfaceParam struct {
	Name     string `json:"Name"`
	IP       string `json:"IP"`
	Port     string `json:"Port"`
	Timeout  string `json:"Timeout"`  //通信超时
	Interval string `json:"Interval"` //通信间隔
}

type CommunicationTcpTemplate struct {
	CommunicationTemplate
	Param TcpInterfaceParam `json:"Param"` //接口参数
	Conn  net.Conn          `json:"-"`
}

var CommunicationTcpMap = make([]*CommunicationTcpTemplate, 0)

func (c *CommunicationTcpTemplate) Open() bool {
	conn, err := net.DialTimeout("tcp", c.Param.IP+":"+c.Param.Port, 2*time.Second)
	if err != nil {
		//setting.Logger.Errorf("%s,tcp open err,%v", c.Name, err)
		return false
	} else {
		setting.Logger.Debugf("%s,tcp open ok", c.Name)
	}
	c.Conn = conn
	return true
}

func (c *CommunicationTcpTemplate) Close() bool {
	if c.Conn != nil {
		err := c.Conn.Close()
		if err != nil {
			return false
		}
	}
	return true
}

func (c *CommunicationTcpTemplate) WriteData(data []byte) int {

	if c.Conn != nil {
		cnt, err := c.Conn.Write(data)
		if err != nil {
			setting.Logger.Errorf("%s,tcp write err,%v", c.Name, err)
			err = c.Conn.Close()
			if err != nil {
				setting.Logger.Errorf("%s,tcp close err,%v", c.Name, err)
			}
			c.Open()
			return 0
		}
		return cnt
	}
	return 0
}

func (c *CommunicationTcpTemplate) ReadData(data []byte) int {

	if c.Conn != nil {
		cnt, err := c.Conn.Read(data)
		//setting.Logger.Debugf("%s,tcp read data cnt %v", c.Name, cnt)
		if err != nil {
			//setting.Logger.Errorf("%s,tcp read err,%v", c.Name, err)
			return 0
		}
		return cnt
	}
	return 0
}

func (c *CommunicationTcpTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationTcpTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationTcpTemplate) GetInterval() string {
	return c.Param.Interval
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
	setting.Logger.Infof("write commTcpInterface.json sucess")
}
