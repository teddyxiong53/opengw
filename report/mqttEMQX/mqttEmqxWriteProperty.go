package mqttEmqx

import (
	"encoding/json"
	"goAdapter/device"
	"goAdapter/setting"
)

type MQTTEmqxWritePropertyRequestParamPropertyTemplate struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type MQTTEmqxWritePropertyRequestParamTemplate struct {
	ClientID   string                                              `json:"clientID"`
	Properties []MQTTEmqxWritePropertyRequestParamPropertyTemplate `json:"properties"`
}

type MQTTEmqxWritePropertyRequestTemplate struct {
	ID      string                                      `json:"id"`
	Version string                                      `json:"version"`
	Ack     int                                         `json:"ack"`
	Params  []MQTTEmqxWritePropertyRequestParamTemplate `json:"params"`
}

type MQTTEmqxWritePropertyAckTemplate struct {
	ID      string                                      `json:"id"`
	Version string                                      `json:"version"`
	Code    int                                         `json:"code"`
	Params  []MQTTEmqxWritePropertyRequestParamTemplate `json:"params"`
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxWritePropertyAck(reqFrame MQTTEmqxWritePropertyRequestTemplate, code int, ackParams []MQTTEmqxWritePropertyRequestParamTemplate) {

	ackFrame := MQTTEmqxWritePropertyAckTemplate{
		ID:      reqFrame.ID,
		Version: reqFrame.Version,
		Code:    code,
		Params:  ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	propertyPostTopic := "/sys/thing/event/property/set_reply/" + r.GWParam.Param.ClientID

	setting.Logger.Infof("property set_reply topic: %s", propertyPostTopic)
	setting.Logger.Debugf("property set_reply: %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson)
		token.Wait()
	}
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxProcessWriteProperty(reqFrame MQTTEmqxWritePropertyRequestTemplate) {

	writeStatus := false

	ackParams := make([]MQTTEmqxWritePropertyRequestParamTemplate, 0)

	for _, v := range reqFrame.Params {
		for _, node := range r.NodeList {
			if v.ClientID == node.Param.ClientID {
				//从上报节点中找到相应节点
				for _, coll := range device.CollectInterfaceMap {
					if coll.CollInterfaceName == node.CollInterfaceName {
						for _, n := range coll.DeviceNodeMap {
							if n.Name == node.Name {
								//从采集服务中找到相应节点
								cmd := device.CommunicationCmdTemplate{}
								cmd.CollInterfaceName = node.CollInterfaceName
								cmd.DeviceName = node.Name
								cmd.FunName = "SetVariables"
								valueMap := make(map[string]interface{})
								for _, pro := range v.Properties {
									valueMap[pro.Name] = pro.Value
								}
								paramStr, _ := json.Marshal(valueMap)
								cmd.FunPara = string(paramStr)
								param := MQTTEmqxWritePropertyRequestParamTemplate{
									ClientID: node.Param.ClientID,
								}
								property := MQTTEmqxWritePropertyRequestParamPropertyTemplate{}
								for _, comm := range device.CommunicationManage {
									if comm.CollInterface == coll {
										ackData := comm.CommunicationManageAddEmergency(cmd)
										if ackData.Status {
											writeStatus = true
											for _, p := range v.Properties {
												property.Name = p.Name
												property.Value = 0
												param.Properties = append(param.Properties, property)
											}
										} else {
											writeStatus = false
											for _, p := range v.Properties {
												property.Name = p.Name
												property.Value = 1
												param.Properties = append(param.Properties, property)
											}
										}
									}
								}
								ackParams = append(ackParams, param)
							}
						}
					}
				}
			}
		}
	}

	if writeStatus == true {
		r.ReportServiceEmqxWritePropertyAck(reqFrame, 0, ackParams)
	} else {
		r.ReportServiceEmqxWritePropertyAck(reqFrame, 1, ackParams)
	}
}
