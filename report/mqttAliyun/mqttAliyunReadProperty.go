package mqttAliyun

import (
	"encoding/json"
	"goAdapter/device"
	"strings"
)

func ReportServiceAliyunProcessGetSubDeviceProperty(r *ReportServiceParamAliyunTemplate, message MQTTAliyunMessageTemplate,
	gw MQTTAliyunRegisterTemplate, cmdName string) {

	addrArray := strings.Split(message.Params["Addr"].(string), ",")
	for _, v := range addrArray {
		for _, n := range r.NodeList {
			if v == n.Param.DeviceName {
				cmd := device.CommunicationCmdTemplate{}
				cmd.CollInterfaceName = "coll1"
				cmd.DeviceName = n.Addr
				cmd.FunName = "GetRealVariables"
				paramStr, _ := json.Marshal(message.Params)
				cmd.FunPara = string(paramStr)

				//if len(device.CommunicationManage) > 0 {
				//	if device.CommunicationManage[0].CommunicationManageAddEmergency(cmd) == true {
				//		payload := MQTTAliyunThingServiceAckTemplate{
				//			Identifier: cmdName,
				//			ID:         message.ID,
				//			Code:       200,
				//			Data:       make(map[string]interface{}),
				//		}
				//		MQTTAliyunThingServiceAck(r.GWParam.MQTTClient, gw, payload)
				//	} else {
				//		payload := MQTTAliyunThingServiceAckTemplate{
				//			Identifier: cmdName,
				//			ID:         message.ID,
				//			Code:       1000,
				//			Data:       make(map[string]interface{}),
				//		}
				//		MQTTAliyunThingServiceAck(r.GWParam.MQTTClient, gw, payload)
				//	}
				//}
			}
		}
	}
}
