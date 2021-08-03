package mqttEmqx

import (
	"encoding/json"
	"goAdapter/setting"
	"strconv"
	"time"
)

type MQTTNodeLogoutParamTemplate struct {
	ClientID  string `json:"clientID"`
	Timestamp int64  `json:"timestamp"`
}

type MQTTNodeLogoutTemplate struct {
	ID      string                        `json:"id"`
	Version string                        `json:"version"`
	Params  []MQTTNodeLogoutParamTemplate `json:"params"`
}

func MQTTEmqxNodeLogOut(param ReportServiceGWParamEmqxTemplate, nodeMap []string) int {

	nodeLogout := MQTTNodeLogoutTemplate{
		ID:      strconv.Itoa(MsgID),
		Version: "V1.0",
	}
	MsgID++

	for _, v := range nodeMap {
		nodeLogoutParam := MQTTNodeLogoutParamTemplate{
			ClientID:  v,
			Timestamp: time.Now().Unix(),
		}
		nodeLogout.Params = append(nodeLogout.Params, nodeLogoutParam)
	}

	//批量注册
	LogoutInTopic := "/sys/thing/event/Logout/post/" + param.Param.ClientID

	sJson, _ := json.Marshal(nodeLogout)
	if len(nodeLogout.Params) > 0 {

		setting.Logger.Debugf("node publish LogoutMsg: %s", sJson)
		setting.Logger.Infof("node publish topic: %s", LogoutInTopic)

		if param.MQTTClient != nil {
			token := param.MQTTClient.Publish(LogoutInTopic, 0, false, sJson)
			token.Wait()
		}
	}

	return MsgID
}

func (r *ReportServiceParamEmqxTemplate) NodeLogOut(name []string) bool {

	nodeMap := make([]string, 0)
	status := false

	setting.Logger.Debugf("nodeLogoutName %v", name)
	for _, d := range name {
		for _, v := range r.NodeList {
			if d == v.Name {
				nodeMap = append(nodeMap, v.Name)

				MQTTEmqxNodeLogOut(r.GWParam, nodeMap)
				select {
				case frame := <-r.ReceiveLogOutAckFrameChan:
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
