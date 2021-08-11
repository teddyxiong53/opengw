package mqttEmqx

import (
	"encoding/json"
	"fmt"
	"goAdapter/device"
	"goAdapter/setting"
	"time"
)

type MQTTEmqxReadPropertyRequestParamPropertyTemplate struct {
	Name string `json:"name"`
}

type MQTTEmqxReadPropertyRequestParamTemplate struct {
	ClientID   string                                             `json:"clientID"`
	Properties []MQTTEmqxReadPropertyRequestParamPropertyTemplate `json:"properties"`
}

type MQTTEmqxReadPropertyRequestTemplate struct {
	ID      string                                     `json:"id"`
	Version string                                     `json:"version"`
	Ack     int                                        `json:"ack"`
	Params  []MQTTEmqxReadPropertyRequestParamTemplate `json:"params"`
}

type MQTTEmqxReadPropertyAckParamPropertyTemplate struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

type MQTTEmqxReadPropertyAckParamTemplate struct {
	ClientID   string                                         `json:"clientID"`
	Properties []MQTTEmqxReadPropertyAckParamPropertyTemplate `json:"properties"`
}

type MQTTEmqxReadPropertyAckTemplate struct {
	ID      string                                 `json:"id"`
	Version string                                 `json:"version"`
	Code    int                                    `json:"code"`
	Params  []MQTTEmqxReadPropertyAckParamTemplate `json:"params"`
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxReadPropertyAck(reqFrame MQTTEmqxReadPropertyRequestTemplate, code int, ackParams []MQTTEmqxReadPropertyAckParamTemplate) {

	ackFrame := MQTTEmqxReadPropertyAckTemplate{
		ID:      reqFrame.ID,
		Version: reqFrame.Version,
		Code:    code,
		Params:  ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	propertyPostTopic := "/sys/thing/event/property/get_reply/" + r.GWParam.Param.ClientID

	setting.Logger.Infof("property get_reply topic: %s", propertyPostTopic)
	setting.Logger.Debugf("property get_reply: %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson)
		token.Wait()
	}
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxProcessReadProperty(reqFrame MQTTEmqxReadPropertyRequestTemplate) {

	ReadStatus := false

	ackParams := make([]MQTTEmqxReadPropertyAckParamTemplate, 0)

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
								cmd.FunName = "GetRealVariables"
								nameMap := make([]string, 0)
								for _, pro := range v.Properties {
									nameMap = append(nameMap, pro.Name)
								}
								paramStr, _ := json.Marshal(nameMap)
								cmd.FunPara = string(paramStr)
								ackParam := MQTTEmqxReadPropertyAckParamTemplate{
									ClientID: node.Param.ClientID,
								}
								property := MQTTEmqxReadPropertyAckParamPropertyTemplate{}
								timeStamp := time.Now().Unix()
								//从采集队列中找到
								for _, comm := range device.CommunicationManage {
									if comm.CollInterface == coll {
										ackData := comm.CommunicationManageAddEmergency(cmd)
										if ackData.Status {
											ReadStatus = true
											for _, p := range v.Properties {
												for _, variable := range n.VariableMap {
													if p.Name == variable.Name {
														if len(variable.Value) >= 1 {
															index := len(variable.Value) - 1
															property.Name = variable.Name
															property.Timestamp = timeStamp
															switch t := variable.Value[index].Value.(type) {
															case uint8, uint16, int16, uint32, uint64:
																property.Value = fmt.Sprintf("%d", variable.Value[index].Value)
															case string:
																property.Value = variable.Value[index].Value.(string)
															default:
																setting.Logger.Debugf("valueType %T", t)
															}
															ackParam.Properties = append(ackParam.Properties, property)
														}
													}
												}
											}
										} else {
											ReadStatus = false
											for _, p := range v.Properties {
												property.Name = p.Name
												property.Value = -1
												ackParam.Properties = append(ackParam.Properties, property)
											}
										}
									}
								}
								ackParams = append(ackParams, ackParam)
							}
						}
					}
				}
			}
		}
	}

	if ReadStatus == true {
		r.ReportServiceEmqxReadPropertyAck(reqFrame, 0, ackParams)
	} else {
		r.ReportServiceEmqxReadPropertyAck(reqFrame, 1, ackParams)
	}
}
