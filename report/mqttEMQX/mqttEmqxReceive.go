package mqttEmqx

import (
	"goAdapter/setting"

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
			setting.Logger.Debugf("Recv TOPIC: %s", receiveFrame.Topic)
			setting.Logger.Debugf("Recv MSG: %s", receiveFrame.Payload)
			ReportServiceParamListEmqx.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
