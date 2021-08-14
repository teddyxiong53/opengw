package mqttEmqx

import (
	"encoding/json"
	"goAdapter/setting"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTNodeLoginParamTemplate struct {
	ClientID  string `json:"clientID"`
	Timestamp int64  `json:"timestamp"`
}

type MQTTNodeLoginTemplate struct {
	ID      string                       `json:"id"`
	Version string                       `json:"version"`
	Params  []MQTTNodeLoginParamTemplate `json:"params"`
}

var MsgID int = 0

func MQTTEmqxGWLogin(param ReportServiceGWParamEmqxTemplate, publishHandler MQTT.MessageHandler) (bool, MQTT.Client) {

	opts := MQTT.NewClientOptions().AddBroker(param.IP + ":" + param.Port)

	opts.SetClientID(param.Param.ClientID)
	opts.SetUsername(param.Param.UserName)
	//hs256 := sha256.New()
	//hs256.Write([]byte("zhsHrx123456@"))
	//password := hs256.Sum(nil)
	//param.Password = string(hex.EncodeToString(password))
	//setting.Logger.Debugf("Emqx password %v", param.Password)
	//param.Password = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InpocyIsImNsaWVudGlkIjoiIiwiaWF0IjoxNjI2NzAwMzE1fQ.EICw6uVoP-_X2iKcdkmBTJevFspm7Nz9ipHjJpr8eHg"
	opts.SetPassword(param.Param.Password)
	opts.SetKeepAlive(60 * 2 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler)
	opts.SetAutoReconnect(false)

	// create and start a client using the above ClientOptions
	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		setting.Logger.Errorf("Connect Emqx IoT Cloud fail %s", token.Error())
		return false, nil
	}
	setting.Logger.Info("Connect Emqx IoT Cloud Sucess")

	subTopic := ""
	//子设备上线回应
	subTopic = "/sys/thing/event/login/post_reply/" + param.Param.ClientID
	MQTTEmqxSubscribeTopic(mqttClient, subTopic)

	//子设备下线回应
	subTopic = "/sys/thing/event/logout/post_reply/" + param.Param.ClientID
	MQTTEmqxSubscribeTopic(mqttClient, subTopic)

	//属性设置上报回应
	subTopic = "/sys/thing/event/property/post_reply/" + param.Param.ClientID
	MQTTEmqxSubscribeTopic(mqttClient, subTopic)

	//订阅属性下发请求
	subTopic = "/sys/thing/event/property/set/" + param.Param.ClientID
	MQTTEmqxSubscribeTopic(mqttClient, subTopic)

	//订阅属性查询请求
	subTopic = "/sys/thing/event/property/get/" + param.Param.ClientID
	MQTTEmqxSubscribeTopic(mqttClient, subTopic)

	//订阅服务调用请求
	subTopic = "/sys/thing/event/service/invoke/" + param.Param.ClientID
	MQTTEmqxSubscribeTopic(mqttClient, subTopic)

	return true, mqttClient

}

func MQTTEmqxSubscribeTopic(client MQTT.Client, topic string) {

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		setting.Logger.Warningf("Subscribe topic %s fail %v", topic, token.Error())
	}
	setting.Logger.Info("Subscribe topic " + topic + " success")
}

func (r *ReportServiceParamEmqxTemplate) GWLogin() bool {
	status := false
	status, r.GWParam.MQTTClient = MQTTEmqxGWLogin(r.GWParam, ReceiveMessageHandler)
	if status == true {
		r.GWParam.ReportStatus = "onLine"
	}

	return status
}

func MQTTEmqxNodeLoginIn(param ReportServiceGWParamEmqxTemplate, nodeMap []string) int {

	nodeLogin := MQTTNodeLoginTemplate{
		ID:      strconv.Itoa(MsgID),
		Version: "V1.0",
	}
	MsgID++

	for _, v := range nodeMap {
		nodeLoginParam := MQTTNodeLoginParamTemplate{
			ClientID:  v,
			Timestamp: time.Now().Unix(),
		}
		nodeLogin.Params = append(nodeLogin.Params, nodeLoginParam)
	}

	//批量注册
	loginTopic := "/sys/thing/event/login/post/" + param.Param.ClientID

	sJson, _ := json.Marshal(nodeLogin)
	if len(nodeLogin.Params) > 0 {

		setting.Logger.Debugf("node publish logInMsg: %s", sJson)
		setting.Logger.Infof("node publish topic: %s", loginTopic)

		if param.MQTTClient != nil {
			token := param.MQTTClient.Publish(loginTopic, 0, false, sJson)
			token.Wait()
		}
	}

	return MsgID
}

func (r *ReportServiceParamEmqxTemplate) NodeLogIn(name []string) bool {

	nodeMap := make([]string, 0)
	status := false

	setting.Logger.Debugf("nodeLoginName %v", name)
	for _, d := range name {
		for _, v := range r.NodeList {
			if d == v.Name {
				nodeMap = append(nodeMap, v.Param.ClientID)

				MQTTEmqxNodeLoginIn(r.GWParam, nodeMap)
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
