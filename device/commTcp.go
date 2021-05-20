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

type TcpClientInterfaceParam struct {
	IP       string `json:"IP"`
	Port     string `json:"Port"`
	Timeout  string `json:"Timeout"`  //通信超时
	Interval string `json:"Interval"` //通信间隔
}

type CommunicationTcpClientTemplate struct {
	Name  string                  `json:"Name"`  //接口名称
	Type  string                  `json:"Type"`  //接口类型,比如serial,TcpClient,udp,http
	Param TcpClientInterfaceParam `json:"Param"` //接口参数
	Conn  net.Conn                `json:"-"`     //通信句柄
}

var CommunicationTcpClientMap = make([]*CommunicationTcpClientTemplate, 0)

func (c *CommunicationTcpClientTemplate) Open() bool {
	conn, err := net.DialTimeout("TcpClient", c.Param.IP+":"+c.Param.Port, 2*time.Second)
	if err != nil {
		//setting.Logger.Errorf("%s,TcpClient open err,%v", c.Name, err)
		return false
	} else {
		setting.Logger.Debugf("%s,TcpClient open ok", c.Name)
	}
	c.Conn = conn
	return true
}

func (c *CommunicationTcpClientTemplate) Close() bool {
	if c.Conn != nil {
		err := c.Conn.Close()
		if err != nil {
			return false
		}
	}
	return true
}

func (c *CommunicationTcpClientTemplate) WriteData(data []byte) int {

	if c.Conn != nil {
		cnt, err := c.Conn.Write(data)
		if err != nil {
			setting.Logger.Errorf("%s,TcpClient write err,%v", c.Name, err)
			err = c.Conn.Close()
			if err != nil {
				setting.Logger.Errorf("%s,TcpClient close err,%v", c.Name, err)
			}
			c.Open()
			return 0
		}
		return cnt
	}
	return 0
}

func (c *CommunicationTcpClientTemplate) ReadData(data []byte) int {

	if c.Conn != nil {
		cnt, err := c.Conn.Read(data)
		//setting.Logger.Debugf("%s,TcpClient read data cnt %v", c.Name, cnt)
		if err != nil {
			//setting.Logger.Errorf("%s,TcpClient read err,%v", c.Name, err)
			return 0
		}
		return cnt
	}
	return 0
}

func (c *CommunicationTcpClientTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationTcpClientTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationTcpClientTemplate) GetInterval() string {
	return c.Param.Interval
}

func ReadCommTcpClientInterfaceListFromJson() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/commTcpClientInterface.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open commTcpClientInterface.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt], &CommunicationTcpClientMap)
		if err != nil {
			log.Println("commTcpClientInterface unmarshal err", err)
			return false
		}
		return true
	} else {
		log.Println("commTcpClientInterface.json is not exist")

		return false
	}
}

func WriteCommTcpClientInterfaceListToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/commTcpClientInterface.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open commTcpClientInterface.json err", err)
		return
	}
	defer fp.Close()

	sJson, _ := json.Marshal(CommunicationTcpClientMap)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write commTcpClientInterface.json err", err)
	}
	setting.Logger.Infof("write commTcpClientInterface.json sucess")
}
