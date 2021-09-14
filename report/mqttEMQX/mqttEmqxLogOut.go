/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 15:17:27
@FilePath: /goAdapter-Raw/report/mqttEMQX/mqttEmqxLogOut.go
*/
package mqttEmqx

import (
	"encoding/json"
	"goAdapter/pkg/mylog"
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
	LogoutTopic := "/sys/thing/event/logout/post/" + param.Param.ClientID

	sJson, _ := json.Marshal(nodeLogout)
	if len(nodeLogout.Params) > 0 {

		mylog.Logger.Debugf("node publish LogoutMsg: %s", sJson)
		mylog.Logger.Infof("node publish topic: %s", LogoutTopic)

		if param.MQTTClient != nil {
			token := param.MQTTClient.Publish(LogoutTopic, 0, false, sJson)
			token.Wait()
		}
	}

	return MsgID
}

func (r *ReportServiceParamEmqxTemplate) NodeLogOut(name []string) bool {

	nodeMap := make([]string, 0)
	status := false

	mylog.Logger.Debugf("nodeLogoutName %v", name)
	for _, d := range name {
		for _, v := range r.NodeList {
			if d == v.Name {
				nodeMap = append(nodeMap, v.Param.ClientID)

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
