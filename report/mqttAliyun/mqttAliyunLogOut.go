package mqttAliyun

import (
	"goAdapter/setting"
	"time"
)

func (r *ReportServiceParamAliyunTemplate) NodeLogOut(name []string) bool {

	status := false

	nodeList := make([]MQTTAliyunNodeRegisterTemplate, 0)
	nodeParam := MQTTAliyunNodeRegisterTemplate{}

	for _, d := range name {
		for _, v := range r.NodeList {
			if d == v.Name {
				if v.ReportStatus == "offLine" {
					setting.Logger.Infof("service:%s,%s is already offLine", r.GWParam.ServiceName, v.Name)
				} else {
					nodeParam.DeviceSecret = v.Param.DeviceSecret
					nodeParam.DeviceName = v.Param.DeviceName
					nodeParam.ProductKey = v.Param.ProductKey

					nodeList = append(nodeList, nodeParam)
					//r.NodeList[k].CommStatus = "offLine"

					mqttAliyunRegister := MQTTAliyunRegisterTemplate{
						RemoteIP:     r.GWParam.IP,
						RemotePort:   r.GWParam.Port,
						ProductKey:   r.GWParam.Param.ProductKey,
						DeviceName:   r.GWParam.Param.DeviceName,
						DeviceSecret: r.GWParam.Param.DeviceSecret,
					}
					MQTTAliyunNodeLoginOut(r.GWParam.MQTTClient, mqttAliyunRegister, nodeList)
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
	return status
}
