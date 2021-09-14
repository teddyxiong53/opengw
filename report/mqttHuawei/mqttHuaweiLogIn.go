package mqttHuawei

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"goAdapter/pkg/mylog"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHuaweiRegisterTemplate struct {
	RemoteIP     string
	RemotePort   string
	DeviceID     string `json:"DeviceID"`
	DeviceSecret string `json:"DeviceSecret"`
}

type MQTTHuaweiNodeRegisterTemplate struct {
	DeviceID string `json:"device_id"`
	Status   string `json:"status"`
}

var timeStapmStatic string = "2021042609"
var MsgID int = 0

// 时间戳：为设备连接平台时的UTC时间，格式为YYYYMMDDHH，如UTC 时间2018/7/24 17:56:20 则应表示为2018072417。
func timeStamp() string {
	strFormatTime := time.Now().Format("2006-01-02 15:04:05")
	strFormatTime = strings.ReplaceAll(strFormatTime, "-", "")
	strFormatTime = strings.ReplaceAll(strFormatTime, " ", "")
	strFormatTime = strFormatTime[0:10]
	return strFormatTime
}

func assembleClientId(deviceID string) string {
	segments := make([]string, 4)
	segments[0] = deviceID
	segments[1] = "0"
	segments[2] = "0"
	//segments[3] = timeStamp()
	segments[3] = timeStapmStatic

	return strings.Join(segments, "_")
}

func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func MQTTHuaweiGWLogin(param MQTTHuaweiRegisterTemplate, publishHandler MQTT.MessageHandler) (bool, MQTT.Client) {

	opts := MQTT.NewClientOptions()
	opts.AddBroker(param.RemoteIP)

	clientID := assembleClientId(param.DeviceID)
	mylog.Logger.Debugf("clientID %s", clientID)
	opts.SetClientID(clientID)
	opts.SetUsername(param.DeviceID)
	mylog.Logger.Debugf("DeviceSecret %s", param.DeviceSecret)
	//passWord := hmacSha256(param.DeviceSecret, timeStamp())
	passWord := hmacSha256(param.DeviceSecret, timeStapmStatic)
	mylog.Logger.Debugf("passWord %s", passWord)
	opts.SetPassword(passWord)
	opts.SetKeepAlive(250 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler)
	opts.SetAutoReconnect(false)

	// create and start a client using the above ClientOptions
	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		mylog.Logger.Errorf("Connect Huawei IoT Cloud fail %s", token.Error())
		return false, nil
	}
	mylog.Logger.Info("Connect Huawei IoT Cloud Sucess")

	subTopic := ""
	//平台查询属性
	subTopic = "$oc/devices" + param.DeviceID + "/sys/properties/get/#"
	MQTTHuaweiSubscribeTopic(mqttClient, subTopic)
	/*
		//属性设置
		subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/service/property/set"
		MQTTHuaweiSubscribeTopic(mqttClient, subTopic)

		//服务调用(服务不需要主动订阅，平台自动订阅)
		//subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/service/RemoteCmdOpen"
		//MQTTHuaweiSubscribeTopic(mqttClient, subTopic)

		//子设备注册
		subTopic = "/sys/" + param.ProductKey + "/" + param.DeviceName + "/thing/sub/register_reply"
		MQTTHuaweiSubscribeTopic(mqttClient, subTopic)

	*/
	//服务调用
	subTopic = "$oc/devices/" + param.DeviceID + "/sys/commands/#"
	MQTTHuaweiSubscribeTopic(mqttClient, subTopic)

	return true, mqttClient
}

func MQTTHuaweiSubscribeTopic(client MQTT.Client, topic string) {

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		mylog.Logger.Warningf("Subscribe topic %s fail %v", topic, token.Error())
	}
	mylog.Logger.Info("Subscribe topic " + topic + " success")
}

func (r *ReportServiceParamHuaweiTemplate) GWLogin() bool {

	mqttHuaweiRegister := MQTTHuaweiRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		DeviceID:     r.GWParam.Param.DeviceID,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	status := false
	if r.GWParam.MQTTClient != nil {
		r.GWParam.MQTTClient.Disconnect(5)
	}
	status, r.GWParam.MQTTClient = MQTTHuaweiGWLogin(mqttHuaweiRegister, ReceiveMessageHandler)
	if status == true {
		r.GWParam.ReportStatus = "onLine"
	}

	return status
}

func MQTTHuaweiNodeLoginIn(client MQTT.Client, gw ReportServiceGWParamHuaweiTemplate, node []MQTTHuaweiNodeRegisterTemplate) int {

	type NodeStatusesTemplate struct {
		DeviceStatuses []MQTTHuaweiNodeRegisterTemplate `json:"device_statuses"`
	}

	type NodeRegisterTemplate struct {
		ServiceID string               `json:"service_id"`
		EventType string               `json:"event_type"`
		Paras     NodeStatusesTemplate `json:"paras"`
	}

	type NodeServicesTemplate struct {
		Services []NodeRegisterTemplate `json:"services"`
	}

	nodeStatuses := NodeStatusesTemplate{}
	nodeStatuses.DeviceStatuses = node

	nodeRegister := NodeRegisterTemplate{
		ServiceID: "$sub_device_manager",
		EventType: "sub_device_update_status",
		Paras:     nodeStatuses,
	}

	nodeServices := NodeServicesTemplate{
		Services: make([]NodeRegisterTemplate, 0),
	}
	nodeServices.Services = append(nodeServices.Services, nodeRegister)

	sJson, _ := json.Marshal(nodeServices)
	if len(node) > 0 {
		//批量注册
		loginInTopic := "$oc/devices/" + gw.Param.DeviceID + "/sys/events/up"

		mylog.Logger.Debugf("node publish logInMsg: %s", sJson)
		mylog.Logger.Infof("node publish topic: %s", loginInTopic)
		if client != nil {
			token := client.Publish(loginInTopic, 0, false, sJson)
			token.Wait()
		}
	}

	return MsgID
}

func (r *ReportServiceParamHuaweiTemplate) NodeLogin(name []string) bool {

	status := false

	nodeList := make([]MQTTHuaweiNodeRegisterTemplate, 0)
	nodeParam := MQTTHuaweiNodeRegisterTemplate{}

	mylog.Logger.Debugf("nodeLoginName %v", name)
	for _, d := range name {
		for k, v := range r.NodeList {
			if d == v.Name {
				nodeParam.DeviceID = v.Param.DeviceID
				r.NodeList[k].CommStatus = "onLine"
				nodeParam.Status = "ONLINE"
				nodeList = append(nodeList, nodeParam)
				MQTTHuaweiNodeLoginIn(r.GWParam.MQTTClient, r.GWParam, nodeList)
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
