package device

import (
	"fmt"
	"goAdapter/config"
	"io"
	"strconv"

	s2 "github.com/jacobsa/go-serial/serial"
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
	Name     string                `json:"Name"`  //接口名称
	Type     string                `json:"Type"`  //接口类型,比如serial,tcp,udp,http
	Param    *SerialInterfaceParam `json:"Param"` //接口参数
	Port     io.ReadWriteCloser    `json:"-"`     //通信句柄
	err      error                 `json:"-"`
	Bindings []string              `json:"Bindings"`
}

var _ CommunicationInterface = (*CommunicationSerialTemplate)(nil)

func (c *CommunicationSerialTemplate) Open() (err error) {

	serialParam := c.Param
	serialBaud, err := strconv.Atoi(serialParam.BaudRate)
	if err != nil {
		return err
	}

	var serialParity s2.ParityMode
	switch serialParam.Parity {
	case "N":
		serialParity = s2.PARITY_NONE
	case "O":
		serialParity = s2.PARITY_ODD
	case "E":
		serialParity = s2.PARITY_EVEN
	default:
		return fmt.Errorf("serial parity not valid:%s", serialParam.Parity)
	}

	var serialStop uint = 1
	switch serialParam.StopBits {
	case "1":
		serialStop = 1
	case "1.5":
		serialStop = 1
	case "2":
		serialStop = 2
	}

	var databit int
	databit, err = strconv.Atoi(c.Param.DataBits)
	if err != nil {
		return err
	}
	// serialConfig := &serial.Config{
	// 	Name:        serialParam.Name,
	// 	Baud:        serialBaud,
	// 	Parity:      serialParity,
	// 	StopBits:    serialStop,
	// 	ReadTimeout: time.Second * time.Duration(config.Cfg.SerialCfg.ReadTimeOut),
	// }

	serialConfig := s2.OpenOptions{
		PortName:              serialParam.Name,
		BaudRate:              uint(serialBaud),
		ParityMode:            serialParity,
		StopBits:              serialStop,
		DataBits:              uint(databit),
		InterCharacterTimeout: uint(config.Cfg.SerialCfg.ReadTimeOut),
	}
	//c.Port, err = serial.OpenPort(serialConfig)
	c.Port, err = s2.Open(serialConfig)
	if err != nil {
		c.err = err
		return err
	}
	return nil
}

func (c *CommunicationSerialTemplate) Close() error {
	if c.Port == nil {
		return fmt.Errorf("port close error:port %s is not initialized", c.Param.Name)
	}
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

//通讯超时
func (c *CommunicationSerialTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

//命令间隔
func (c *CommunicationSerialTemplate) GetInterval() string {
	return c.Param.Interval
}

func (c *CommunicationSerialTemplate) Error() error {
	return c.err
}

func (c *CommunicationSerialTemplate) Unique() string {
	return fmt.Sprintf("type:%s serial:%s", c.Type, c.Param.Name)
}

func (c *CommunicationSerialTemplate) Bind(name string) {
	if c.Bindings == nil {
		c.Bindings = make([]string, 0)
	}
	var alreadyExists bool
	for _, v := range c.Bindings {
		if v == name {
			alreadyExists = true
			break
		}
	}
	if !alreadyExists {
		c.Bindings = append(c.Bindings, name)
	}
	fmt.Printf("after bind:%v\n", c.Bindings)
}

func (c *CommunicationSerialTemplate) UnBind(name string) {
	if c.Bindings == nil {
		c.Bindings = make([]string, 0)
		return
	}
	var index int
	for k, v := range c.Bindings {
		if v == name {
			index = k
		}
	}
	c.Bindings = append(c.Bindings[:index], c.Bindings[index+1:]...)
	fmt.Printf("after unbind :%v\n", c.Bindings)

}

func (c *CommunicationSerialTemplate) BindNames() []string {
	if c.Bindings == nil {
		c.Bindings = make([]string, 0)

	}
	fmt.Printf("get bindings:%v\n", c.Bindings)

	return c.Bindings
}
