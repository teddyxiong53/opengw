package mqttHuawei

import (
	"bytes"
	"encoding/json"
	"goAdapter/setting"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHuaweiRegisterTemplate struct {
	RemoteIP     string
	RemotePort   string
	ProductKey   string `json:"ProductKey"`
	DeviceName   string `json:"DeviceName"`
	DeviceSecret string `json:"DeviceSecret"`
}

type MQTTHuaweiNodeRegisterTemplate struct {
	ProductKey   string `json:"ProductKey"`
	DeviceName   string `json:"DeviceName"`
	DeviceSecret string `json:"DeviceSecret"`
}

var timeStamp string = "1528018257135"
var MsgID int = 0

func MQTTHuaweiGWLogin(param MQTTHuaweiRegisterTemplate, publishHandler MQTT.MessageHandler) (bool, MQTT.Client) {

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
		setting.Logger.Errorf("Connect Huawei IoT Cloud fail,%s", token.Error())
		return false, nil
	}
	setting.Logger.Info("Connect Huawei IoT Cloud Sucess")

	subTopic := ""
	//属性上报回应
	subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/event/property/pack/post_reply"
	MQTTHuaweiSubscribeTopic(mqttClient, subTopic)

	//属性设置
	subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/service/property/set"
	MQTTHuaweiSubscribeTopic(mqttClient, subTopic)

	//服务调用(服务不需要主动订阅，平台自动订阅)
	//subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/service/RemoteCmdOpen"
	//MQTTHuaweiSubscribeTopic(mqttClient, subTopic)

	//子设备注册
	subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/sub/register_reply"
	MQTTHuaweiSubscribeTopic(mqttClient, subTopic)

	return true, mqttClient

	//MQTTClient_AddTopo()

	//MQTTClient_Register()
}

func MQTTHuaweiSubscribeTopic(client MQTT.Client, topic string) {

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		setting.Logger.Warningf("Subscribe topic %s fail,%v", topic, token.Error())
	}
	setting.Logger.Info("Subscribe topic " + topic + " success")
}

func (r *ReportServiceParamHuaweiTemplate) GWLogin() bool {

	mqttHuaweiRegister := MQTTHuaweiRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	status := false
	status, r.GWParam.MQTTClient = MQTTHuaweiGWLogin(mqttHuaweiRegister, ReceiveMessageHandler)
	if status == true {
		r.GWParam.ReportStatus = "onLine"
	}

	return status
}

func MQTTHuaweiNodeLoginIn(client MQTT.Client, gw MQTTHuaweiRegisterTemplate, node []MQTTHuaweiNodeRegisterTemplate) int {

	type NodeParamsTemplate struct {
		DeviceName   string `json:"deviceName"`
		ProductKey   string `json:"productKey"`
		Sign         string `json:"sign"`
		SignMethod   string `json:"signMethod"`
		TimeStamp    string `json:"timestamp"`
		ClientID     string `json:"clientId"`
		CleanSession string `json:"cleanSession"`
	}

	type NodeParamsListTemplate struct {
		DeviceList []NodeParamsTemplate `json:"deviceList"`
	}

	type MQTTNodePayloadTemplate struct {
		ID     string                 `json:"id"`
		Params NodeParamsListTemplate `json:"params"`
	}
	//单个注册
	//loginTopic := "/ext/session/" + MQTTAliyunGWParam.GWParam.ProductKey + "/" + MQTTAliyunGWParam.GWParam.DeviceName + "/combine/login"
	//批量注册
	loginInTopic := "/ext/session/" + gw.ProductKey + "/" + gw.DeviceName + "/combine/batch_login"

	NodeParamsList := NodeParamsListTemplate{
		make([]NodeParamsTemplate, 0),
	}

	mqttPayload := MQTTNodePayloadTemplate{
		ID: strconv.Itoa(MsgID),
	}
	MsgID++

	for _, v := range node {
		auth := MqttClient_CalculateSign(v.ProductKey, v.DeviceName, v.DeviceSecret, timeStamp)
		MQTTNodeParams := NodeParamsTemplate{
			DeviceName:   v.DeviceName,
			ProductKey:   v.ProductKey,
			Sign:         auth.password,
			SignMethod:   "hmacSha1",
			TimeStamp:    timeStamp,
			ClientID:     v.ProductKey + "&" + v.DeviceName,
			CleanSession: "true",
		}
		NodeParamsList.DeviceList = append(NodeParamsList.DeviceList, MQTTNodeParams)
	}
	mqttPayload.Params = NodeParamsList
	sJson, _ := json.Marshal(mqttPayload)
	if len(NodeParamsList.DeviceList) > 0 {

		setting.Logger.Debugf("node publish logInMsg: %s\n", sJson)
		setting.Logger.Infof("node publish topic: %s\n", loginInTopic)

		if client != nil {
			token := client.Publish(loginInTopic, 0, false, sJson)
			token.Wait()
		}
	}

	return MsgID
}

func (r *ReportServiceParamHuaweiTemplate) NodeLogin(name []string) bool {

	nodeList := make([]MQTTHuaweiNodeRegisterTemplate, 0)
	nodeParam := MQTTHuaweiNodeRegisterTemplate{}
	status := false

	setting.Logger.Debugf("nodeLoginName %v", name)
	for _, d := range name {
		for k, v := range r.NodeList {
			if d == v.Name {
				nodeParam.DeviceSecret = v.Param.DeviceSecret
				nodeParam.DeviceName = v.Param.DeviceName
				nodeParam.ProductKey = v.Param.ProductKey
				nodeList = append(nodeList, nodeParam)
				r.NodeList[k].CommStatus = "onLine"

				mqttHuaweiRegister := MQTTHuaweiRegisterTemplate{
					RemoteIP:     r.GWParam.IP,
					RemotePort:   r.GWParam.Port,
					ProductKey:   r.GWParam.Param.ProductKey,
					DeviceName:   r.GWParam.Param.DeviceName,
					DeviceSecret: r.GWParam.Param.DeviceSecret,
				}
				MQTTHuaweiNodeLoginIn(r.GWParam.MQTTClient, mqttHuaweiRegister, nodeList)
				select {
				case <-r.ReceiveLogInAckFrameChan:
					{
						status = true
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
