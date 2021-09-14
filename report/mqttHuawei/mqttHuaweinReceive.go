/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 19:11:50
@FilePath: /goAdapter-Raw/report/mqttHuawei/mqttHuaweinReceive.go
*/
package mqttHuawei

import (
	"goAdapter/pkg/mylog"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHuaweiReceiveFrameTemplate struct {
	Topic   string
	Payload []byte
}

type MQTTHuaweiLogInDataTemplate struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

type MQTTHuaweiLogInAckTemplate struct {
	ID      string                        `json:"id"`
	Code    int32                         `json:"code"`
	Message string                        `json:"message"`
	Data    []MQTTHuaweiLogInDataTemplate `json:"data"`
}

type MQTTHuaweiLogOutDataTemplate struct {
	Code       int32  `json:"code"`
	Message    string `json:"message"`
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

type MQTTHuaweiLogOutAckTemplate struct {
	ID      string                         `json:"id"`
	Code    int32                          `json:"code"`
	Message string                         `json:"message"`
	Data    []MQTTHuaweiLogOutDataTemplate `json:"data"`
}

type MQTTHuaweiReportPropertyAckTemplate struct {
	Code    int32  `json:"code"`
	Data    string `json:"-"`
	ID      string `json:"id"`
	Message string `json:"message"`
	Method  string `json:"method"`
	Version string `json:"version"`
}

type MQTTHuaweiGetPropertiesRequestTemplate struct {
	ServiceID      string `json:"service_id"`
	ObjectDeviceID string `json:"object_device_id"`
}

//发送数据回调函数
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListHuawei.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTHuaweiReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}

			mylog.Logger.Debugf("Recv TOPIC: %s", receiveFrame.Topic)
			mylog.Logger.Debugf("Recv MSG: %s", receiveFrame.Payload)
			ReportServiceParamListHuawei.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
