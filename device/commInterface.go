/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 09:55:35
@FilePath: /goAdapter-Raw/device/commInterface.go
*/
package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type CommunicationInterface interface {
	Open() error
	io.ReadWriteCloser //和系统接口保持一致
	GetName() string
	GetType() string
	GetParam() interface{}
	GetTimeOut() string
	GetInterval() string
	Unique() string //串口不能名字一致,网口不能ip和port一致
	Error() error   //打开或者关闭的时候是否产生错误，如果有错误就不执行read这些操作哦了
}

//通信接口Map
var CommunicationInterfaceMap = make(map[string]CommunicationInterface)

func CommInterfaceInit(data []byte) error {
	var temp map[string]map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	for k, i := range temp {

		typ, ok := i["Type"]
		if !ok {
			return errors.New("json file content is not include Type")
		}
		param, ok := i["Param"]
		if !ok {
			return errors.New("json file content is not include Param")
		}
		switch typ {
		case SERIALTYPE:
			sParam := CommunicationSerialTemplate{
				Name: i["Name"].(string),
			}
			sParam.Type = typ.(string)
			var serial SerialInterfaceParam
			data, _ := json.Marshal(param)
			if err := json.Unmarshal(data, &serial); err != nil {
				return err
			}
			sParam.Param = &serial
			CommunicationInterfaceMap[k] = &sParam
		case TCPCLIENTTYPE:
			sParam := CommunicationTcpClientTemplate{
				Name: i["Name"].(string),
			}
			sParam.Type = typ.(string)
			var tcpClient TcpClientInterfaceParam
			data, _ := json.Marshal(param)
			if err := json.Unmarshal(data, &tcpClient); err != nil {
				return err
			}
			sParam.Param = &tcpClient
			CommunicationInterfaceMap[k] = &sParam
		case IOINTYPE:
			sParam := CommunicationIoInTemplate{
				Name: i["Name"].(string),
			}
			sParam.Type = typ.(string)
			var ioIn IoInInterfaceParam
			data, _ := json.Marshal(param)
			if err := json.Unmarshal(data, &ioIn); err != nil {
				return err
			}
			sParam.Param = &ioIn
			CommunicationInterfaceMap[k] = &sParam
		case IOOUTTYPE:
			sParam := CommunicationIoOutTemplate{
				Name: i["Name"].(string),
			}
			sParam.Type = typ.(string)
			var ioOut IoOutInterfaceParam
			data, _ := json.Marshal(param)
			if err := json.Unmarshal(data, &ioOut); err != nil {
				return err
			}
			sParam.Param = &ioOut
			CommunicationInterfaceMap[k] = &sParam
		default:
			return fmt.Errorf("unKnown type of json:%s", typ)
		}
	}
	return nil
}
