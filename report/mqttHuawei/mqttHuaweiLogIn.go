package mqttHuawei

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"goAdapter/setting"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHuaweiRegisterTemplate struct {
	RemoteIP     string
	RemotePort   string
	DeviceID   string `json:"DeviceName"`
	DeviceSecret string `json:"DeviceSecret"`
}

type MQTTHuaweiNodeRegisterTemplate struct {
	DeviceID   string `json:"DeviceName"`
	DeviceSecret string `json:"DeviceSecret"`
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
	setting.Logger.Debugf("clientID %s",clientID)
	opts.SetClientID(clientID)
	opts.SetUsername(param.DeviceID)
	setting.Logger.Debugf("DeviceSecret %s",param.DeviceSecret)
	//passWord := hmacSha256(param.DeviceSecret, timeStamp())
	passWord := hmacSha256(param.DeviceSecret, timeStapmStatic)
	setting.Logger.Debugf("passWord %s",passWord)
	opts.SetPassword(passWord)
	opts.SetKeepAlive(250 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler)
	opts.SetAutoReconnect(false)

	// create and start a client using the above ClientOptions
	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		setting.Logger.Errorf("Connect Huawei IoT Cloud fail,%s", token.Error())
		return false, nil
	}
	setting.Logger.Info("Connect Huawei IoT Cloud Sucess")


	//subTopic := ""
	//属性上报回应
	//subTopic = "$oc/devices" + param.DeviceID + "/sys/messages/down"
	//MQTTHuaweiSubscribeTopic(mqttClient, subTopic)
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

	return true, mqttClient
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
		DeviceID:   r.GWParam.Param.DeviceID,
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

	/*
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

	 */

	return MsgID
}

func (r *ReportServiceParamHuaweiTemplate) NodeLogin(name []string) bool {

	status := false
	/*
	nodeList := make([]MQTTHuaweiNodeRegisterTemplate, 0)
	nodeParam := MQTTHuaweiNodeRegisterTemplate{}

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

	 */

	return status
}
