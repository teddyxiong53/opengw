package mqttAliyun

import (
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
	ID      string                       `json:"id"`
	Code    int32                        `json:"code"`
	Message string                       `json:"message"`
	Data    MQTTAliyunLogOutDataTemplate `json:"data"`
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

	for _, v := range ReportServiceParamListAliyun.ServiceList {
		if v.GWParam.MQTTClient == client {
			receiveFrame := MQTTAliyunReceiveFrameTemplate{
				Topic:   msg.Topic(),
				Payload: msg.Payload(),
			}
			v.ReceiveFrameChan <- receiveFrame
		}
	}
}

//func ProcessPropertyPost(r *ReportServiceParamAliyunTemplate) {
//
//	for {
//		select {
//		case postParam := <-r.PropertyPostChan:
//			{
//				setting.Logger.Tracef("service %s,postParam %v,postChanCnt %v", r.GWParam.ServiceName, postParam, len(r.PropertyPostChan))
//				if postParam.DeviceType == 0 { //网关上报
//					r.GWPropertyPost()
//				} else if postParam.DeviceType == 1 { //末端设备上报
//					r.NodePropertyPost(postParam.DeviceName)
//				}
//			}
//		}
//	}
//}
