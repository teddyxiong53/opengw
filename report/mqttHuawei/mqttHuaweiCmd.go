/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 19:11:00
@FilePath: /goAdapter-Raw/report/mqttHuawei/mqttHuaweiCmd.go
*/
package mqttHuawei

import (
	"encoding/json"
	"goAdapter/device"
	"goAdapter/pkg/mylog"
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
	mylog.Logger.Debugf("writeCmdAck post msg: %s", sJson)

	serviceTopic := "$oc/devices/" + r.GWParam.Param.DeviceID + "/sys/commands/response/request_id=" + requestID
	mylog.Logger.Infof("writeCmdAck post topic: %s", serviceTopic)

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
	cmd.FunName = device.LUAFUNC(request.CommandName)
	paramStr, _ := json.Marshal(request.Paras)
	cmd.FunPara = string(paramStr)

	cmdAck := MQTTHuaweiWriteCmdAckTemplate{}
	cmdRX := device.CommunicationManage.ManagerTemp[name].CommunicationManageAddEmergency(cmd)
	if cmdRX.Err == nil {
		mylog.Logger.Debugf("WriteCmd ok")
		cmdAck.ResultCode = 0
		cmdAck.ResponseName = request.ServiceID
	} else {
		cmdAck.ResultCode = 1
		cmdAck.ResponseName = request.ServiceID
	}

	MQTTHuaweiWriteCmdAck(r, requestID, cmdAck)
}
