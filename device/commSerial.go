package device

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tarm/serial"
)

type SerialInterfaceParam struct {
	Name     string `json:"Name"`
	BaudRate string `json:"BaudRate"`
	DataBits string `json:"DataBits"` //数据位: 5, 6, 7 or 8 (default 8)
	StopBits string `json:"StopBits"` //停止位: 1 or 2 (default 1)
	Parity   string `json:"Parity"`   //校验: N - None, E - Even, O - Odd (default E),(The use of no parity requires 2 stop bits.)
	Timeout  string `json:"Timeout"`  //通信超时
	Interval string `json:"Interval"` //通信间隔
}

type CommunicationSerialTemplate struct {
	Name  string                `json:"Name"`  //接口名称
	Type  string                `json:"Type"`  //接口类型,比如serial,tcp,udp,http
	Param *SerialInterfaceParam `json:"Param"` //接口参数
	Port  *serial.Port          `json:"-"`     //通信句柄
	err   error                 `json:"-"`
}

var _ CommunicationInterface = (*CommunicationSerialTemplate)(nil)

func (c *CommunicationSerialTemplate) Open() error {

	serialParam := c.Param
	serialBaud, err := strconv.Atoi(serialParam.BaudRate)
	if err != nil {
		return err
	}

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
		Name:        serialParam.Name,
		Baud:        serialBaud,
		Parity:      serialParity,
		StopBits:    serialStop,
		ReadTimeout: time.Millisecond * 10,
	}

	c.Port, err = serial.OpenPort(serialConfig)
	if err != nil {
		c.err = err
		return err
	}
	return nil
}

func (c *CommunicationSerialTemplate) Close() error {
	return c.Port.Close()
}

func (c *CommunicationSerialTemplate) Write(data []byte) (i int, err error) {
	if c.Port == nil {
		err = fmt.Errorf("port %s is not initialized", c.Param.Name)
		return
	}

	return c.Port.Write(data)
}

func (c *CommunicationSerialTemplate) Read(data []byte) (i int, err error) {

	if c.Port == nil {
		err = fmt.Errorf("port %s is not initialized", c.Param.Name)
		return
	}

	return c.Port.Read(data)
}

func (c *CommunicationSerialTemplate) GetName() string {
	return c.Name
}
func (c *CommunicationSerialTemplate) GetType() string {
	return c.Type
}

func (c *CommunicationSerialTemplate) GetParam() interface{} {
	return c.Param
}

func (c *CommunicationSerialTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationSerialTemplate) GetInterval() string {
	return c.Param.Interval
}

func (c *CommunicationSerialTemplate) Error() error {
	return c.err
}

func (c *CommunicationSerialTemplate) Unique() string {
	return fmt.Sprintf("type:%s serial:%s", c.Type, c.Param.Name)
}
