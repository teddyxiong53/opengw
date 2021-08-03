package mqttHuawei

import (
	"encoding/json"
	"goAdapter/device"
	"goAdapter/setting"
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
	setting.Logger.Debugf("thingServiceAck post msg: %s", sJson)

	serviceTopic := "$oc/devices/" + r.GWParam.Param.DeviceID + "/sys/properties/get/response/" + service.ServiceID
	setting.Logger.Infof("thingServiceAck post topic: %s", serviceTopic)

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
	y := 0
	for k, v := range device.CollectInterfaceMap {
		if v.CollInterfaceName == r.NodeList[x].CollInterfaceName {
			y = k
			break
		}
	}
	i := 0
	for k, v := range device.CollectInterfaceMap[y].DeviceNodeMap {
		if v.Name == r.NodeList[x].Name {
			i = k
			break
		}
	}

	cmd := device.CommunicationCmdTemplate{}
	cmd.CollInterfaceName = device.CollectInterfaceMap[y].CollInterfaceName
	cmd.DeviceName = device.CollectInterfaceMap[y].DeviceNodeMap[i].Name
	cmd.FunName = "GetRealVariables"
	cmd.FunPara = ""

	cmdRX := device.CommunicationManage[y].CommunicationManageAddEmergency(cmd)
	if cmdRX.Status == true {
		setting.Logger.Debugf("GetRealVariables ok")
		service := MQTTHuaweiServiceTemplate{}
		for _, v := range device.CollectInterfaceMap[y].DeviceNodeMap[i].VariableMap {
			if v.Name == request.ServiceID {
				if len(v.Value) >= 1 {
					index := len(v.Value) - 1
					service := MQTTHuaweiServiceTemplate{}
					service.ServiceID = v.Name
					service.Properties.Value = v.Value[index].Value
				}
			}
		}
		MQTTHuaweiGetPropertiesAck(r, service)
	}
}
