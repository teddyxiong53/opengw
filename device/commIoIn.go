/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-13 13:42:09
@FilePath: /goAdapter-Raw/device/commIoIn.go
*/
package device

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type IoInInterfaceParam struct {
	Name string   `json:"Name"`
	FD   *os.File `json:"-"`
}

type CommunicationIoInTemplate struct {
	Name  string              `json:"Name"`  //接口名称
	Type  string              `json:"Type"`  //接口类型,比如serial,IoIn,udp,http
	Param *IoInInterfaceParam `json:"Param"` //接口参数
	err   error               `json:"-"`
}

var _ CommunicationInterface = (*CommunicationIoInTemplate)(nil)

func (c *CommunicationIoInTemplate) Error() error {
	return c.err
}

func (c *CommunicationIoInTemplate) Unique() string {
	return fmt.Sprintf("type:%s ioin:%s", c.Type, c.Param.Name)
}
func (c *CommunicationIoInTemplate) Open() error {

	fd, err := os.OpenFile(c.Param.Name, os.O_RDWR, 0777)
	if err != nil {
		c.err = err
		return err
	}
	c.Param.FD = fd

	return nil
}

func (c *CommunicationIoInTemplate) Close() error {

	if c.Param.FD == nil {
		return errors.New("tcp client conn is not initialized")
	}
	return c.Param.FD.Close()
}

func (c *CommunicationIoInTemplate) Write(data []byte) (i int, err error) {

	if c.Param.FD == nil {
		err = errors.New("tcp client conn is not initialized")
		return
	}
	return c.Param.FD.Write(data)
}

func (c *CommunicationIoInTemplate) Read(data []byte) (i int, err error) {

	if c.Param.FD == nil {
		err = errors.New("tcp client conn is not initialized")
		return
	}

	return c.Param.FD.Read(data)

}

func (c *CommunicationIoInTemplate) GetName() string {
	return c.Name
}

func (c *CommunicationIoInTemplate) GetType() string {
	return c.Type
}
func (c *CommunicationIoInTemplate) GetParam() interface{} {
	return c.Param
}

func (c *CommunicationIoInTemplate) GetTimeOut() string {
	return ""
}

func (c *CommunicationIoInTemplate) GetInterval() string {
	return ""
}
