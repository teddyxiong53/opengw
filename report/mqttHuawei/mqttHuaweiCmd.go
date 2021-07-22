package mqttHuawei

import (
	"encoding/json"
	"goAdapter/device"
	"goAdapter/setting"
)

type MQTTHuaweiWriteCmdRequestTemplate struct {
	ServiceID      string                 `json:"service_id"`
	ObjectDeviceID string                 `json:"object_device_id"`
	CommandName    string                 `json:"command_name"`
	Paras          map[string]interface{} `json:"paras"`
}

type MQTTHuaweiWriteCmdAckTemplate struct {
	ResultCode   int                    `json:"result_code"`
	ResponseName string                 `json:"response_name"`
	Paras        map[string]interface{} `json:"paras"`
}

func MQTTHuaweiWriteCmdAck(r *ReportServiceParamHuaweiTemplate, requestID string, ack MQTTHuaweiWriteCmdAckTemplate) {

	sJson, _ := json.Marshal(ack)
	setting.Logger.Debugf("writeCmdAck post msg: %s", sJson)

	serviceTopic := "$oc/devices/" + r.GWParam.Param.DeviceID + "/sys/commands/response/request_id=" + requestID
	setting.Logger.Infof("writeCmdAck post topic: %s", serviceTopic)

	if r.GWParam.MQTTClient != nil {
		token := r.GWParam.MQTTClient.Publish(serviceTopic, 0, false, sJson)
		token.Wait()
	}
}

func ReportServiceHuaweiProcessWriteCmd(r *ReportServiceParamHuaweiTemplate, requestID string, request MQTTHuaweiWriteCmdRequestTemplate) {
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
	cmd.FunName = request.CommandName
	paramStr, _ := json.Marshal(request.Paras)
	cmd.FunPara = string(paramStr)

	cmdAck := MQTTHuaweiWriteCmdAckTemplate{}
	if device.CommunicationManage[y].CommunicationManageAddEmergency(cmd) == true {
		setting.Logger.Debugf("WriteCmd ok")
		cmdAck.ResultCode = 0
		cmdAck.ResponseName = request.ServiceID
	} else {
		cmdAck.ResultCode = 1
		cmdAck.ResponseName = request.ServiceID
	}

	MQTTHuaweiWriteCmdAck(r, requestID, cmdAck)
}
