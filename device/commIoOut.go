/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-13 13:42:34
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
	Name  string               `json:"Name"`  //接口名称
	Type  string               `json:"Type"`  //接口类型,比如serial,IoOut,udp,http
	Param *IoOutInterfaceParam `json:"Param"` //接口参数
	err   error                `json:"-"`
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
