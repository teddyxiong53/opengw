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
				if reqFrame.Params["Name"] == nil {
					continue
				}
				name := reqFrame.Params["Name"].(string)
				for _, n := range r.NodeList {
					if name == n.Param.DeviceName {
						for _, v := range device.CommunicationManage {
							if v.CollInterface.CollInterfaceName == n.CollInterfaceName {
								cmd := device.CommunicationCmdTemplate{}
								cmd.CollInterfaceName = n.CollInterfaceName
								cmd.DeviceName = n.Name
								cmd.FunName = methodParam[2]
								paramStr, _ := json.Marshal(reqFrame.Params)
								cmd.FunPara = string(paramStr)

								ack := MQTTAliyunInvokeThingsServiceAckTemplate{}
								ackData := v.CommunicationManageAddEmergency(cmd)
								if ackData.Status {
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
}
