package mqttEmqx

import (
	"encoding/json"
	"goAdapter/device"
	"goAdapter/pkg/mylog"
)

type MQTTEmqxInvokeServiceAckParamTemplate struct {
	ClientID  string `json:"clientID"`
	CmdName   string `json:"cmdName"`
	CmdStatus int    `json:"cmdStatus"`
}

type MQTTEmqxInvokeServiceAckTemplate struct {
	ID      string                                  `json:"id"`
	Version string                                  `json:"version"`
	Code    int                                     `json:"Code"`
	Params  []MQTTEmqxInvokeServiceAckParamTemplate `json:"params"`
}

type MQTTEmqxInvokeServiceRequestParamTemplate struct {
	ClientID  string                 `json:"clientID"`
	CmdName   string                 `json:"cmdName"`
	CmdParams map[string]interface{} `json:"cmdParams"`
}

type MQTTEmqxInvokeServiceRequestTemplate struct {
	ID      string                                      `json:"id"`
	Version string                                      `json:"version"`
	Ack     int                                         `json:"ack"`
	Params  []MQTTEmqxInvokeServiceRequestParamTemplate `json:"params"`
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxInvokeServiceAck(reqFrame MQTTEmqxInvokeServiceRequestTemplate, code int, ackParams []MQTTEmqxInvokeServiceAckParamTemplate) {

	ackFrame := MQTTEmqxInvokeServiceAckTemplate{
		ID:      reqFrame.ID,
		Version: reqFrame.Version,
		Code:    code,
		Params:  ackParams,
	}

	sJson, _ := json.Marshal(ackFrame)
	serviceInvokeTopic := "/sys/thing/event/service/invoke_reply/" + r.GWParam.Param.ClientID

	mylog.Logger.Infof("service invoke_reply topic: %s", serviceInvokeTopic)
	mylog.Logger.Debugf("service invoke_reply: %v", string(sJson))
	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(serviceInvokeTopic, 0, false, sJson)
		token.Wait()
	}
}

func (r *ReportServiceParamEmqxTemplate) ReportServiceEmqxProcessInvokeService(reqFrame MQTTEmqxInvokeServiceRequestTemplate) {

	ReadStatus := false

	ackParams := make([]MQTTEmqxInvokeServiceAckParamTemplate, 0)

	for _, v := range reqFrame.Params {
		for _, node := range r.NodeList {
			if v.ClientID == node.Param.ClientID {
				//从上报节点中找到相应节点
				tmps := device.CollectInterfaceMap.GetAll()
				for _, coll := range tmps {
					if coll.CollInterfaceName == node.CollInterfaceName {
						for _, n := range coll.DeviceNodes {
							if n.Name == node.Name {
								//从采集服务中找到相应节点
								cmd := device.CommunicationCmdTemplate{}
								cmd.CollInterfaceName = node.CollInterfaceName
								cmd.DeviceName = node.Name
								cmd.FunName = device.LUAFUNC(v.CmdName)
								paramStr, _ := json.Marshal(v.CmdParams)
								cmd.FunPara = string(paramStr)
								ackParam := MQTTEmqxInvokeServiceAckParamTemplate{
									ClientID: node.Param.ClientID,
									CmdName:  v.CmdName,
								}
								//从采集队列中找到

								ackData := coll.CommunicationManager.CommunicationManageAddEmergency(cmd)
								if ackData.Err == nil {
									ReadStatus = true
									ackParam.CmdStatus = 0
								} else {
									ReadStatus = false
									ackParam.CmdStatus = 1
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
		r.ReportServiceEmqxInvokeServiceAck(reqFrame, 0, ackParams)
	} else {
		r.ReportServiceEmqxInvokeServiceAck(reqFrame, 1, ackParams)
	}
}
