/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 15:19:08
@FilePath: /goAdapter-Raw/report/mqttEMQX/mqttEmqxReceive.go
*/
package mqttEmqx

import (
	"goAdapter/pkg/mylog"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTEmqxReceiveFrameTemplate struct {
	Topic   string
	Payload []byte
}

type MQTTEmqxLogInDataTemplate struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

type MQTTEmqxLogInAckTemplate struct {
	ID      string                      `json:"id"`
	Code    int32                       `json:"code"`
	Message string                      `json:"message"`
	Data    []MQTTEmqxLogInDataTemplate `json:"data"`
}

type MQTTEmqxLogOutDataTemplate struct {
	Code       int32  `json:"code"`
	Message    string `json:"message"`
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

type MQTTEmqxLogOutAckTemplate struct {
	ID      string                       `json:"id"`
	Code    int32                        `json:"code"`
	Message string                       `json:"message"`
	Data    []MQTTEmqxLogOutDataTemplate `json:"data"`
}

type MQTTEmqxReportPropertyAckTemplate struct {
	Code    int32  `json:"code"`
	Data    string `json:"-"`
	ID      string `json:"id"`
	Message string `json:"message"`
	Method  string `json:"method"`
	Version string `json:"version"`
}

//发送数据回调函数
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListEmqx.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTEmqxReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}
			mylog.Logger.Debugf("Recv TOPIC: %s", receiveFrame.Topic)
			mylog.Logger.Debugf("Recv MSG: %s", receiveFrame.Payload)
			ReportServiceParamListEmqx.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
