package mqttHuawei

import (
	"encoding/json"
	"goAdapter/setting"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func MQTTHuaweiNodeLogOut(client MQTT.Client, gw ReportServiceGWParamHuaweiTemplate, node []MQTTHuaweiNodeRegisterTemplate) int {

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
		logOutTopic := "$oc/devices/" + gw.Param.DeviceID + "/sys/events/up"

		setting.Logger.Debugf("node publish logOutMsg: %s", sJson)
		setting.Logger.Infof("node publish topic: %s", logOutTopic)
		if client != nil {
			token := client.Publish(logOutTopic, 0, false, sJson)
			token.Wait()
		}
	}

	return MsgID
}

func (r *ReportServiceParamHuaweiTemplate) NodeLogOut(name []string) bool {

	status := false

	nodeList := make([]MQTTHuaweiNodeRegisterTemplate, 0)
	nodeParam := MQTTHuaweiNodeRegisterTemplate{}

	setting.Logger.Debugf("nodeLogOutName %v", name)
	for _, d := range name {
		for k, v := range r.NodeList {
			if d == v.Name {
				nodeParam.DeviceID = v.Param.DeviceID
				r.NodeList[k].CommStatus = "offLine"
				nodeParam.Status = "OFFLINE"
				nodeList = append(nodeList, nodeParam)
				MQTTHuaweiNodeLogOut(r.GWParam.MQTTClient, r.GWParam, nodeList)
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
