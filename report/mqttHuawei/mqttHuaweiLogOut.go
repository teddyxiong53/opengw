package mqttHuawei

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func MQTTHuaweiNodeLoginOut(client MQTT.Client, gw MQTTHuaweiRegisterTemplate, node []MQTTHuaweiNodeRegisterTemplate) int {

	/*
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
		setting.Logger.Infof("node publish logOutMsg: %s\n", sJson)
		setting.Logger.Debugf("node publish topic: %s\n", loginOutTopic)

		token := client.Publish(loginOutTopic, 0, false, sJson)
		token.Wait()
	}

	 */

	return MsgID
}

func (r *ReportServiceParamHuaweiTemplate) NodeLogOut(name []string) bool {

	status := false
	/*
	nodeList := make([]MQTTHuaweiNodeRegisterTemplate, 0)
	nodeParam := MQTTHuaweiNodeRegisterTemplate{}

	for _, d := range name {
		for k, v := range r.NodeList {
			if d == v.Name {
				if v.ReportStatus == "offLine" {
					setting.Logger.Infof("service:%s,%s is already offLine", r.GWParam.ServiceName, v.Name)
				} else {
					nodeParam.DeviceSecret = v.Param.DeviceSecret
					nodeParam.DeviceName = v.Param.DeviceName
					nodeParam.ProductKey = v.Param.ProductKey

					nodeList = append(nodeList, nodeParam)
					r.NodeList[k].CommStatus = "offLine"

					mqttHuaweiRegister := MQTTHuaweiRegisterTemplate{
						RemoteIP:     r.GWParam.IP,
						RemotePort:   r.GWParam.Port,
						ProductKey:   r.GWParam.Param.ProductKey,
						DeviceName:   r.GWParam.Param.DeviceName,
						DeviceSecret: r.GWParam.Param.DeviceSecret,
					}
					MQTTHuaweiNodeLoginOut(r.GWParam.MQTTClient, mqttHuaweiRegister, nodeList)
					select {
					case frame := <-r.ReceiveLogOutAckFrameChan:
						{
							if frame.Code == 200 {

							}
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
	}

	 */
	return status
}
