/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-10-08 09:39:19
@FilePath: /goAdapter-Raw/device/commIoOut.go
*/
package device

import (
	"errors"
	"fmt"
	"os"
)

type IoOutInterfaceParam struct {
	Name string   `json:"Name"`
	FD   *os.File `json:"-"`
}

type CommunicationIoOutTemplate struct {
	Name     string               `json:"Name"`  //接口名称
	Type     string               `json:"Type"`  //接口类型,比如serial,IoOut,udp,http
	Param    *IoOutInterfaceParam `json:"Param"` //接口参数
	err      error                `json:"-"`
	Bindings []string             `json:"Bindings"`
}

var _ CommunicationInterface = (*CommunicationIoOutTemplate)(nil)

func (c *CommunicationIoOutTemplate) Error() error {
	return c.err
}
func (c *CommunicationIoOutTemplate) Unique() string {
	return fmt.Sprintf("type:%s ioout:%s", c.Type, c.Param.Name)
}
func (c *CommunicationIoOutTemplate) Open() error {

	fd, err := os.OpenFile(c.Param.Name, os.O_RDWR, 0666)
	if err != nil {
		c.err = err
		return err
	}
	c.Param.FD = fd

	return nil
}

func (c *CommunicationIoOutTemplate) Close() error {

	if c.Param.FD == nil {
		return errors.New("tcp client conn is not initialized")
	}
	return c.Param.FD.Close()
}

func (c *CommunicationIoOutTemplate) Write(data []byte) (i int, err error) {

	if c.Param.FD == nil {
		err = errors.New("tcp client conn is not initialized")
		return
	}
	return c.Param.FD.Write(data)
}

func (c *CommunicationIoOutTemplate) Read(data []byte) (i int, err error) {

	return
}

func (c *CommunicationIoOutTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationIoOutTemplate) GetType() string {
	return c.Type
}

func (c *CommunicationIoOutTemplate) GetParam() interface{} {
	return c.Param
}

func (c *CommunicationIoOutTemplate) GetTimeOut() string {
	return ""
}

func (c *CommunicationIoOutTemplate) GetInterval() string {
	return ""
}

func (c *CommunicationIoOutTemplate) Bind(name string) {
	if c.Bindings == nil {
		c.Bindings = make([]string, 0)
	}
	c.Bindings = append(c.Bindings, name)
}
func (c *CommunicationIoOutTemplate) UnBind(name string) {
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
func (c *CommunicationIoOutTemplate) BindNames() []string {
	if c.Bindings == nil {
		c.Bindings = make([]string, 0)

	}
	return c.Bindings
}
