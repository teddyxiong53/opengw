package mqttHuawei

import (
	"goAdapter/setting"

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

//发送数据回调函数
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListHuawei.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTHuaweiReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}
			setting.Logger.Debugf("Recv TOPIC: %s\n", receiveFrame.Topic)
			setting.Logger.Debugf("Recv MSG: %s\n", receiveFrame.Payload)
			ReportServiceParamListHuawei.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
