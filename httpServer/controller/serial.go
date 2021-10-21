/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 19:24:41
@FilePath: /goAdapter-Raw/httpServer/controller/serial.go
*/
package controller

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"goAdapter/device"
	"goAdapter/httpServer/model"
	"goAdapter/pkg/serial"

	"github.com/gin-gonic/gin"
)

func GetSerial(context *gin.Context) {
	type SerialPortNameTemplate struct {
		Name string `json:"Name"`
	}

	data := make([]SerialPortNameTemplate, 0)

	SerialPortName := SerialPortNameTemplate{}
	for _, v := range serial.SerialPortNameTemplateMap.Name {
		SerialPortName.Name = v
		data = append(data, SerialPortName)
	}

	context.JSON(http.StatusOK, model.Response{
		Code:    "0",
		Message: "",
		Data:    data,
	})
}

func SendDirectDataToCollInterface(context *gin.Context) {

	data, err := context.GetRawData()
	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}
	serviceInfo := struct {
		CollInterfaceName string
		DirectData        string
	}{}
	err = json.Unmarshal(data, &serviceInfo)
	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("unmarshal error:%v", err),
		})
		return
	}

	coll := device.CollectInterfaceMap.Get(serviceInfo.CollInterfaceName)
	if coll == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("no such collect named %s", serviceInfo.CollInterfaceName),
		})
		return
	}

	//去掉字符串中的空格
	serviceInfo.DirectData = strings.ReplaceAll(serviceInfo.DirectData, " ", "")
	if len(serviceInfo.DirectData)%2 != 0 {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "输入数字不是偶数位",
		})
		return
	}

	reqData, err := hex.DecodeString(serviceInfo.DirectData)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("decode hex string error:%v", err),
			Data:    serviceInfo.DirectData,
		})
		return
	}
	req := device.CommunicationDirectDataReqTemplate{
		CollInterfaceName: serviceInfo.CollInterfaceName,
		Data:              reqData,
	}
	coll.CommunicationManager.DirectDataRequestChan <- req
	select {
	case ack := <-coll.CommunicationManager.DirectDataAckChan:
		if ack.Error == nil {
			context.JSON(200, model.Response{
				Code: "0",
				Data: fmt.Sprintf("% X", ack.Data),
			})
		} else {
			context.JSON(200, model.Response{
				Code:    "1",
				Message: ack.Error.Error(),
			})
		}
	case <-time.After(time.Second * 2):
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("等待响应超过2秒,已放弃"),
		})
		return
	}

}
