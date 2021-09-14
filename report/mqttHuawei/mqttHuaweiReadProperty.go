/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 19:13:56
@FilePath: /goAdapter-Raw/report/mqttHuawei/mqttHuaweiReadProperty.go
*/
package mqttHuawei

import (
	"encoding/json"
	"goAdapter/device"
	"goAdapter/pkg/mylog"
)

func MQTTHuaweiGetPropertiesAck(r *ReportServiceParamHuaweiTemplate, service MQTTHuaweiServiceTemplate) {

	type MQTTHuaweiDeviceServiceTemplate struct {
		Services []MQTTHuaweiServiceTemplate `json:"services"`
	}

	deviceService := MQTTHuaweiDeviceServiceTemplate{
		Services: make([]MQTTHuaweiServiceTemplate, 0),
	}

	deviceService.Services = append(deviceService.Services, service)

	sJson, _ := json.Marshal(deviceService)
	mylog.Logger.Debugf("thingServiceAck post msg: %s", sJson)

	serviceTopic := "$oc/devices/" + r.GWParam.Param.DeviceID + "/sys/properties/get/response/" + service.ServiceID
	mylog.Logger.Infof("thingServiceAck post topic: %s", serviceTopic)

	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(serviceTopic, 0, false, sJson)
		token.Wait()
	}

}

func ReportServiceHuaweiProcessGetProperties(r *ReportServiceParamHuaweiTemplate, request MQTTHuaweiGetPropertiesRequestTemplate) {

	x := 0
	for k, v := range r.NodeList {
		if v.Param.DeviceID == request.ObjectDeviceID {
			x = k
			break
		}
	}
	var name string
	for k, v := range device.CollectInterfaceMap {
		if v.CollInterfaceName == r.NodeList[x].CollInterfaceName {
			name = k
			break
		}
	}
	i := 0
	for k, v := range device.CollectInterfaceMap[name].DeviceNodes {
		if v.Name == r.NodeList[x].Name {
			i = k
			break
		}
	}

	cmd := device.CommunicationCmdTemplate{}
	cmd.CollInterfaceName = device.CollectInterfaceMap[name].CollInterfaceName
	cmd.DeviceName = device.CollectInterfaceMap[name].DeviceNodes[i].Name
	cmd.FunName = "GetRealVariables"
	cmd.FunPara = ""

	cmdRX := device.CommunicationManage.ManagerTemp[name].CommunicationManageAddEmergency(cmd)
	if cmdRX.Err == nil {
		mylog.Logger.Debugf("GetRealVariables ok")
		service := MQTTHuaweiServiceTemplate{}
		for _, v := range device.CollectInterfaceMap[name].DeviceNodes[i].VariableMap {
			if v.Name == request.ServiceID {
				if len(v.Values) >= 1 {
					index := len(v.Values) - 1
					service := MQTTHuaweiServiceTemplate{}
					service.ServiceID = v.Name
					service.Properties.Value = v.Values[index].Value
				}
			}
		}
		MQTTHuaweiGetPropertiesAck(r, service)
	}
}
