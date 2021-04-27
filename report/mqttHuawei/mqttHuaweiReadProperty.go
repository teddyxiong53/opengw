package mqttHuawei

import (
	"encoding/json"
	"goAdapter/device"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHuaweiMessageTemplate struct {
	Method  string                 `json:"method"`
	ID      string                 `json:"id"`
	Params  map[string]interface{} `json:"params"`
	Version string                 `json:"version"`
}

type MQTTHuaweiThingServiceAckTemplate struct {
	Identifier string                 `json:"identifier"`
	ID         string                 `json:"id"`
	Code       int                    `json:"code"`
	Data       map[string]interface{} `json:"data"`
}

func MQTTHuaweiThingServiceAck(client MQTT.Client, gw MQTTHuaweiRegisterTemplate, ackMessage MQTTHuaweiThingServiceAckTemplate) {

	/*
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

	*/
}

func ReportServiceHuaweiProcessGetSubDeviceProperty(r *ReportServiceParamHuaweiTemplate, message MQTTHuaweiMessageTemplate,
	gw MQTTHuaweiRegisterTemplate, cmdName string) {

	addrArray := strings.Split(message.Params["Addr"].(string), ",")
	for _, v := range addrArray {
		for _, n := range r.NodeList {
			if v == n.Param.DeviceID {
				cmd := device.CommunicationCmdTemplate{}
				cmd.CollInterfaceName = "coll1"
				cmd.DeviceName = n.Addr
				cmd.FunName = "GetRealVariables"
				paramStr, _ := json.Marshal(message.Params)
				cmd.FunPara = string(paramStr)

				if len(device.CommunicationManage) > 0 {
					if device.CommunicationManage[0].CommunicationManageAddEmergency(cmd) == true {
						payload := MQTTHuaweiThingServiceAckTemplate{
							Identifier: cmdName,
							ID:         message.ID,
							Code:       200,
							Data:       make(map[string]interface{}),
						}
						MQTTHuaweiThingServiceAck(r.GWParam.MQTTClient, gw, payload)
					} else {
						payload := MQTTHuaweiThingServiceAckTemplate{
							Identifier: cmdName,
							ID:         message.ID,
							Code:       1000,
							Data:       make(map[string]interface{}),
						}
						MQTTHuaweiThingServiceAck(r.GWParam.MQTTClient, gw, payload)
					}
				}
			}
		}
	}
}
