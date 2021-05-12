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
				methodParam := strings.Split(reqFrame.Method, ".")
				if len(methodParam) != 3 {
					continue
				}
				name := reqFrame.Params["Name"].(string)
				for _, n := range r.NodeList {
					if name == n.Param.DeviceName {
						cmd := device.CommunicationCmdTemplate{}
						cmd.CollInterfaceName = n.CollInterfaceName
						cmd.DeviceName = n.Name
						cmd.FunName = methodParam[2]
						paramStr, _ := json.Marshal(reqFrame.Params)
						cmd.FunPara = string(paramStr)

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
