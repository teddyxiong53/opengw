package mqttAliyun

import (
	"bytes"
	"goAdapter/setting"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTAliyunRegisterTemplate struct {
	RemoteIP     string
	RemotePort   string
	ProductKey   string `json:"ProductKey"`
	DeviceName   string `json:"DeviceName"`
	DeviceSecret string `json:"DeviceSecret"`
}

func MQTTAliyunGWLogin(param MQTTAliyunRegisterTemplate, publishHandler MQTT.MessageHandler) (bool, MQTT.Client) {

	var raw_broker bytes.Buffer

	//MQTT.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)

	raw_broker.WriteString(param.ProductKey)
	raw_broker.WriteString(param.RemoteIP)
	opts := MQTT.NewClientOptions().AddBroker(raw_broker.String())

	auth := MqttClient_CalculateSign(param.ProductKey,
		param.DeviceName,
		param.DeviceSecret, timeStamp)
	opts.SetClientID(auth.mqttClientId)
	opts.SetUsername(auth.username)
	opts.SetPassword(auth.password)
	opts.SetKeepAlive(60 * 2 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler)
	opts.SetAutoReconnect(false)

	// create and start a client using the above ClientOptions
	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		setting.Logger.Errorf("Connect aliyun IoT Cloud fail,%s", token.Error())
		return false, nil
	}
	setting.Logger.Info("Connect aliyun IoT Cloud Sucess")

	subTopic := ""
	//属性上报回应
	subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/event/property/pack/post_reply"
	MQTTAliyunSubscribeTopic(mqttClient, subTopic)

	//属性设置
	subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/service/property/set"
	MQTTAliyunSubscribeTopic(mqttClient, subTopic)

	//服务调用(服务不需要主动订阅，平台自动订阅)
	//subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/service/RemoteCmdOpen"
	//MQTTAliyunSubscribeTopic(mqttClient, subTopic)

	//子设备注册
	subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/sub/register_reply"
	MQTTAliyunSubscribeTopic(mqttClient, subTopic)

	return true, mqttClient

	//MQTTClient_AddTopo()

	//MQTTClient_Register()
}

func MQTTAliyunSubscribeTopic(client MQTT.Client, topic string) {

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		setting.Logger.Warningf("Subscribe topic %s fail,%v", topic, token.Error())
	}
	setting.Logger.Info("Subscribe topic " + topic + " success")
}

func (r *ReportServiceParamAliyunTemplate) GWLogin() bool {

	mqttAliyunRegister := MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	status := false
	status, r.GWParam.MQTTClient = MQTTAliyunGWLogin(mqttAliyunRegister, ReceiveMessageHandler)
	if status == true {
		r.GWParam.ReportStatus = "onLine"
	}

	return status
}

func (r *ReportServiceParamAliyunTemplate) NodeLogin(name []string) bool {

	nodeList := make([]MQTTAliyunNodeRegisterTemplate, 0)
	nodeParam := MQTTAliyunNodeRegisterTemplate{}
	status := false

	setting.Logger.Debugf("nodeLoginName %v", name)
	for _, d := range name {
		for _, v := range r.NodeList {
			if d == v.Name {
				nodeParam.DeviceSecret = v.Param.DeviceSecret
				nodeParam.DeviceName = v.Param.DeviceName
				nodeParam.ProductKey = v.Param.ProductKey
				nodeList = append(nodeList, nodeParam)
				//r.NodeList[k].CommStatus = "onLine"

				mqttAliyunRegister := MQTTAliyunRegisterTemplate{
					RemoteIP:     r.GWParam.IP,
					RemotePort:   r.GWParam.Port,
					ProductKey:   r.GWParam.Param.ProductKey,
					DeviceName:   r.GWParam.Param.DeviceName,
					DeviceSecret: r.GWParam.Param.DeviceSecret,
				}
				MQTTAliyunNodeLoginIn(r.GWParam.MQTTClient, mqttAliyunRegister, nodeList)
				select {
				case frame := <-r.ReceiveLogInAckFrameChan:
					{
						if frame.Code == 200 {
							status = true
						}
					}
				case <-time.After(time.Millisecond * 2000):
					{
						status = false
					}
				}
			}
		}
	}

	return status
}
