package mqttAliyun

import (
	"encoding/json"
	"goAdapter/pkg/mylog"
	"strconv"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

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

var (
	timeStamp string = "1528018257135"
	MsgID     int    = 0
)

func init() {

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

		mylog.Logger.Debugf("node publish logInMsg: %s", sJson)
		mylog.Logger.Infof("node publish topic: %s", loginInTopic)

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
		mylog.Logger.Infof("node publish logOutMsg: %s", sJson)
		mylog.Logger.Debugf("node publish topic: %s", loginOutTopic)

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

	mylog.Logger.Infof("gw property post topic: %s", propertyPostTopic)
	mylog.Logger.Debugf("gw property post msg: %s", sJson)
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
	mylog.Logger.Infof("node property post topic: %s", propertyPostTopic)
	mylog.Logger.Debugf("node property post msg: %v", PropertyPayload)
	if client != nil {
		token := client.Publish(propertyPostTopic, 0, false, sJson)
		token.Wait()
	}

	return MsgID
}

func MQTTAliyunThingServiceAck(client MQTT.Client, gw ReportServiceGWParamAliyunTemplate, ackMessage MQTTAliyunInvokeThingsServiceAckTemplate, serviceName string) {

	sJson, _ := json.Marshal(ackMessage)
	mylog.Logger.Debugf("thingServiceAck post msg: %s", sJson)

	thingServiceTopic := "/sys/" + gw.Param.ProductKey + "/" + gw.Param.DeviceName +
		"/thing/service/" + serviceName + "_reply"
	mylog.Logger.Infof("thingServiceAck post topic: %s", thingServiceTopic)

	if client != nil {
		token := client.Publish(thingServiceTopic, 0, false, sJson)
		token.Wait()
	}
}
