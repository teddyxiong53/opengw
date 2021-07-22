package mqttAliyun

import (
	"goAdapter/setting"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTAliyunReceiveFrameTemplate struct {
	Topic   string
	Payload []byte
}

type MQTTAliyunLogInDataTemplate struct {
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

type MQTTAliyunLogInAckTemplate struct {
	ID      string                        `json:"id"`
	Code    int32                         `json:"code"`
	Message string                        `json:"message"`
	Data    []MQTTAliyunLogInDataTemplate `json:"data"`
}

type MQTTAliyunLogOutDataTemplate struct {
	Code       int32  `json:"code"`
	Message    string `json:"message"`
	ProductKey string `json:"productKey"`
	DeviceName string `json:"deviceName"`
}

type MQTTAliyunLogOutAckTemplate struct {
	ID      string                         `json:"id"`
	Code    int32                          `json:"code"`
	Message string                         `json:"message"`
	Data    []MQTTAliyunLogOutDataTemplate `json:"data"`
}

type MQTTAliyunInvokeThingsServiceRequestTemplate struct {
	ID      string                 `json:"id"`
	Version string                 `json:"version"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

type MQTTAliyunInvokeThingsServiceAckTemplate struct {
	ID   string                 `json:"id"`
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
}

type MQTTAliyunReportPropertyAckTemplate struct {
	Code    int32  `json:"code"`
	Data    string `json:"-"`
	ID      string `json:"id"`
	Message string `json:"message"`
	Method  string `json:"method"`
	Version string `json:"version"`
}

//发送数据回调函数
func ReceiveMessageHandler(client MQTT.Client, msg MQTT.Message) {

	for k, v := range ReportServiceParamListAliyun.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTAliyunReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}
			setting.Logger.Debugf("Recv TOPIC: %s", receiveFrame.Topic)
			setting.Logger.Debugf("Recv MSG: %v", receiveFrame.Payload)
			ReportServiceParamListAliyun.ServiceList[k].ReceiveFrameChan <- receiveFrame
		}
	}
}
