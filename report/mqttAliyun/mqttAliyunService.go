package mqttAliyun

import (
	"encoding/json"
	"goAdapter/device"
	"goAdapter/setting"
	"strings"
)

func (r *ReportServiceParamAliyunTemplate) ProcessInvokeThingsService() {

	for {
		select {
		case reqFrame := <-r.InvokeThingsServiceRequestFrameChan:
			{
				setting.Logger.Debugf("reqFrame %v", reqFrame)
				name := reqFrame.Params["Name"].(string)
				for _, n := range r.NodeList {
					if name == n.Param.DeviceName {
						cmd := device.CommunicationCmdTemplate{}
						if strings.Contains(reqFrame.Method, "RelayOpen") {
							cmd.CollInterfaceName = n.CollInterfaceName
							cmd.DeviceName = n.Name
							cmd.FunName = "RelayOpen"
							paramStr, _ := json.Marshal(reqFrame.Params)
							cmd.FunPara = string(paramStr)
						} else if strings.Contains(reqFrame.Method, "RelayClose") {
							cmd.CollInterfaceName = n.CollInterfaceName
							cmd.DeviceName = n.Name
							cmd.FunName = "RelayClose"
							paramStr, _ := json.Marshal(reqFrame.Params)
							cmd.FunPara = string(paramStr)
						}

						ack := MQTTAliyunInvokeThingsServiceAckTemplate{}
						if len(device.CommunicationManage) > 0 {
							if device.CommunicationManage[0].CommunicationManageAddEmergency(cmd) == true {
								ack.ID = reqFrame.ID
								ack.Code = 200
								MQTTAliyunThingServiceAck(r.GWParam.MQTTClient, r.GWParam, ack, cmd.FunName)
							} else {
								ack.ID = reqFrame.ID
								ack.Code = 1000
								MQTTAliyunThingServiceAck(r.GWParam.MQTTClient, r.GWParam, ack, cmd.FunName)
							}
						}
					}
				}
			}
		}
	}
}
