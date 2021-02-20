package mqttAliyun

import (
	"bytes"
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"goAdapter/setting"
	"strconv"
	"time"
)

type MQTTAliyunRegisterTemplate struct {
	RemoteIP     string
	RemotePort   string
	ProductKey   string `json:"ProductKey"`
	DeviceName   string `json:"DeviceName"`
	DeviceSecret string `json:"DeviceSecret"`
}

type MQTTAliyunNodeRegisterTemplate struct {
	ProductKey   string `json:"ProductKey"`
	DeviceName   string `json:"DeviceName"`
	DeviceSecret string `json:"DeviceSecret"`
}

type MQTTAliyunValueTemplate struct {
	Value interface{}
	Name  string
}

type MQTTAliyunNodeValueTemplate struct {
	ProductKey string `json:"ProductKey"`
	DeviceName string `json:"DeviceName"`
	ValueMap   []MQTTAliyunValueTemplate
}

//type MQTTAliyunPropertyPostAckTemplate struct {
//	ID   string `json:"id"`
//	Code int32  `json:"code"`
//	Data string `json:"data"`
//}

type MQTTAliyunMessageTemplate struct {
	Method  string                 `json:"method"`
	ID      string                 `json:"id"`
	Params  map[string]interface{} `json:"params"`
	Version string                 `json:"version"`
}

type MQTTAliyunThingServiceAckTemplate struct {
	Identifier string                 `json:"identifier"`
	ID         string                 `json:"id"`
	Code       int                    `json:"code"`
	Data       map[string]interface{} `json:"data"`
}

var (
	timeStamp string = "1528018257135"
	MsgID     int    = 0
)

func init() {

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
		//log.Println(token.Error())
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

	//MQTTAliyunGWParam.NodeLoginIn()

	//MQTTAliyunGWParam.GWPropertyPost()

	//MQTTAliyunGWParam.NodePropertyPost()
}

func MQTTAliyunSubscribeTopic(client MQTT.Client, topic string) {

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		setting.Logger.Warningf("Subscribe topic %s fail,%v", topic, token.Error())
	}
	setting.Logger.Info("Subscribe topic " + topic + " success")
}

func MQTTAliyunNodeLoginIn(client MQTT.Client, gw MQTTAliyunRegisterTemplate, node []MQTTAliyunNodeRegisterTemplate) int {

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

func MQTTAliyunNodeLoginOut(client MQTT.Client, gw MQTTAliyunRegisterTemplate, node []MQTTAliyunNodeRegisterTemplate) int {

	type NodeParamsTemplate struct {
		DeviceName string `json:"deviceName"`
		ProductKey string `json:"productKey"`
	}

	type MQTTNodePayloadTemplate struct {
		ID     string               `json:"id"`
		Params []NodeParamsTemplate `json:"params"`
	}
	//单个注册
	//loginTopic := "/ext/session/" + MQTTAliyunGWParam.GWParam.ProductKey + "/" + MQTTAliyunGWParam.GWParam.DeviceName + "/combine/login"
	//批量注册
	loginOutTopic := "/ext/session/" + gw.ProductKey + "/" + gw.DeviceName + "/combine/batch_logout"

	mqttPayload := MQTTNodePayloadTemplate{
		ID:     strconv.Itoa(MsgID),
		Params: make([]NodeParamsTemplate, 0),
	}
	MsgID++

	for _, v := range node {
		MQTTNodeParams := NodeParamsTemplate{
			DeviceName: v.DeviceName,
			ProductKey: v.ProductKey,
		}
		mqttPayload.Params = append(mqttPayload.Params, MQTTNodeParams)
	}
	sJson, _ := json.Marshal(mqttPayload)
	if len(mqttPayload.Params) > 0 {
		setting.Logger.Infof("node publish logOutMsg: %s\n", sJson)
		setting.Logger.Debugf("node publish topic: %s\n", loginOutTopic)

		token := client.Publish(loginOutTopic, 0, false, sJson)
		token.Wait()
	}

	return MsgID
}

func MQTTAliyunGWPropertyPost(client MQTT.Client, gw MQTTAliyunRegisterTemplate, valueMap []MQTTAliyunValueTemplate) int {

	type MQTTPropertyValueTemplate struct {
		Value interface{} `json:"value"`
	}

	type MQTTPropertyParamsTemplate struct {
		Properties map[string]MQTTPropertyValueTemplate `json:"properties"`
	}

	type MQTTPropertyPayloadTemplate struct {
		ID      string                     `json:"id"`
		Version string                     `json:"version"`
		Params  MQTTPropertyParamsTemplate `json:"params"`
		Method  string                     `json:"method"`
	}

	PropertyParams := MQTTPropertyParamsTemplate{
		Properties: make(map[string]MQTTPropertyValueTemplate, 0),
	}

	PropertyValueTemplate := MQTTPropertyValueTemplate{}
	for _, v := range valueMap {
		PropertyValueTemplate.Value = v.Value
		PropertyParams.Properties[v.Name] = PropertyValueTemplate
	}

	PropertyPayload := MQTTPropertyPayloadTemplate{
		ID:      strconv.Itoa(MsgID),
		Params:  PropertyParams,
		Version: "1.0",
		Method:  "thing.event.property.pack.post",
	}
	MsgID++

	sJson, _ := json.Marshal(PropertyPayload)
	//propertyPostTopic := "/sys/" + MQTTAliyunGWParam.GWParam.ProductKey + "/" + MQTTAliyunGWParam.GWParam.DeviceName + "/thing/event/property/post"
	propertyPostTopic := "/sys/" + gw.ProductKey + "/" + gw.DeviceName + "/thing/event/property/pack/post"

	setting.Logger.Infof("gw property post topic: %s", propertyPostTopic)
	setting.Logger.Debugf("gw property post msg: %s", sJson)
	if client != nil {
		token := client.Publish(propertyPostTopic, 0, false, sJson)
		token.Wait()
	}

	return MsgID
}

func MQTTAliyunNodePropertyPost(client MQTT.Client, gw MQTTAliyunRegisterTemplate, nodeMap []MQTTAliyunNodeValueTemplate) int {

	type MQTTPropertyValueTemplate struct {
		Value interface{} `json:"value"`
	}

	type MQTTNodeIdentityTemplate struct {
		ProductKey string `json:"productKey"`
		DeviceName string `json:"deviceName"`
	}

	type MQTTNodePropertyParamsTemplate struct {
		Identity   MQTTNodeIdentityTemplate             `json:"identity"`
		Properties map[string]MQTTPropertyValueTemplate `json:"properties"`
	}

	type MQTTNodesPropertyParamsTemplate struct {
		SubDevices []MQTTNodePropertyParamsTemplate `json:"subDevices"`
	}

	type MQTTPropertyPayloadTemplate struct {
		ID      string                          `json:"id"`
		Version string                          `json:"version"`
		Params  MQTTNodesPropertyParamsTemplate `json:"params"`
		Method  string                          `json:"method"`
	}

	MQTTNodesPropertyParams := MQTTNodesPropertyParamsTemplate{
		SubDevices: make([]MQTTNodePropertyParamsTemplate, 0),
	}

	for _, d := range nodeMap {
		MQTTNodePropertyParams := MQTTNodePropertyParamsTemplate{
			Properties: make(map[string]MQTTPropertyValueTemplate, 0),
		}
		MQTTNodePropertyParams.Identity.DeviceName = d.DeviceName
		MQTTNodePropertyParams.Identity.ProductKey = d.ProductKey

		for _, v := range d.ValueMap {
			MQTTPropertyValue := MQTTPropertyValueTemplate{}
			MQTTPropertyValue.Value = v.Value
			MQTTNodePropertyParams.Properties[v.Name] = MQTTPropertyValue
		}

		MQTTNodesPropertyParams.SubDevices = append(MQTTNodesPropertyParams.SubDevices, MQTTNodePropertyParams)
	}

	PropertyPayload := MQTTPropertyPayloadTemplate{
		ID:      strconv.Itoa(MsgID),
		Params:  MQTTNodesPropertyParams,
		Version: "1.0",
		Method:  "thing.event.property.pack.post",
	}
	MsgID++

	sJson, _ := json.Marshal(PropertyPayload)

	//propertyPostTopic := "/sys/" + MQTTAliyunGWParam.GWParam.ProductKey + "/" + MQTTAliyunGWParam.GWParam.DeviceName + "/thing/event/property/post"
	propertyPostTopic := "/sys/" + gw.ProductKey + "/" + gw.DeviceName + "/thing/event/property/pack/post"
	setting.Logger.Infof("node property post topic: %s\n", propertyPostTopic)
	setting.Logger.Debugf("node property post msg: %s\n", sJson)
	if client != nil {
		token := client.Publish(propertyPostTopic, 0, false, sJson)
		token.Wait()
	}

	return MsgID
}

func MQTTAliyunThingServiceAck(client MQTT.Client, gw MQTTAliyunRegisterTemplate, ackMessage MQTTAliyunThingServiceAckTemplate) {

	type MQTTThingServicePayloadTemplate struct {
		ID   string                 `json:"id"`
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}

	payload := MQTTThingServicePayloadTemplate{
		ID:   ackMessage.ID,
		Code: ackMessage.Code,
		Data: ackMessage.Data,
	}

	sJson, _ := json.Marshal(payload)
	setting.Logger.Debugf("thingServiceAck post msg: %s\n", sJson)

	thingServiceTopic := "/sys/" + gw.ProductKey + "/" + gw.DeviceName +
		"/thing/service/" + ackMessage.Identifier + "_reply"
	setting.Logger.Infof("thingServiceAck post topic: %s\n", thingServiceTopic)

	if client != nil {
		token := client.Publish(thingServiceTopic, 0, false, sJson)
		token.Wait()
	}
}
