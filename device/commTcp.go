/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-10-08 09:38:40
@FilePath: /goAdapter-Raw/device/commTcp.go
*/
package device

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type TcpClientInterfaceParam struct {
	IP       string `json:"IP"`
	Port     string `json:"Port"`
	Timeout  string `json:"Timeout"`  //通信超时
	Interval string `json:"Interval"` //通信间隔
}

type CommunicationTcpClientTemplate struct {
	Name     string                   `json:"Name"`  //接口名称
	Type     string                   `json:"Type"`  //接口类型,比如serial,TcpClient,udp,http
	Param    *TcpClientInterfaceParam `json:"Param"` //接口参数
	Conn     net.Conn                 `json:"-"`     //通信句柄
	err      error                    `json:"-"`     //open串口是否出错
	Bindings []string                 `json:"Bindings"`
}

var _ CommunicationInterface = (*CommunicationTcpClientTemplate)(nil)

func (c *CommunicationTcpClientTemplate) Error() error {
	return c.err
}

func (c *CommunicationTcpClientTemplate) Unique() string {
	return fmt.Sprintf("type:%s, tcpclient:%s", c.Type, c.Param.IP+":"+c.Param.Port)
}

func (c *CommunicationTcpClientTemplate) Open() error {
	conn, err := net.DialTimeout("tcp", c.Param.IP+":"+c.Param.Port, 2*time.Second)
	if err != nil {
		c.err = err
		return err
	}
	c.Conn = conn
	return nil
}

func (c *CommunicationTcpClientTemplate) Close() error {
	if c.Conn == nil {
		return errors.New("tcp client conn is not initialized")
	}
	return c.Conn.Close()
}

func (c *CommunicationTcpClientTemplate) Write(data []byte) (i int, err error) {

	if c.Conn == nil {
		err = errors.New("tcp client conn is not initialized")
		return
	}
	return c.Conn.Write(data)

}

func (c *CommunicationTcpClientTemplate) Read(data []byte) (i int, err error) {

	if c.Conn == nil {
		err = errors.New("tcp client conn is not initialized")
		return
	}
	return c.Conn.Read(data)
}

func (c *CommunicationTcpClientTemplate) GetName() string {
	return c.Name
}
func (c *CommunicationTcpClientTemplate) GetType() string {
	return c.Type
}

func (c *CommunicationTcpClientTemplate) GetParam() interface{} {
	return c.Param
}

func (c *CommunicationTcpClientTemplate) GetTimeOut() string {
	return c.Param.Timeout
}

func (c *CommunicationTcpClientTemplate) GetInterval() string {
	return c.Param.Interval
}

func (c *CommunicationTcpClientTemplate) Bind(name string) {
	if c.Bindings == nil {
		c.Bindings = make([]string, 0)
	}
	c.Bindings = append(c.Bindings, name)
}
func (c *CommunicationTcpClientTemplate) UnBind(name string) {
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
}
func (c *CommunicationTcpClientTemplate) BindNames() []string {
	if c.Bindings == nil {
		c.Bindings = make([]string, 0)

	}
	return c.Bindings
}
