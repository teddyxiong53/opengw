package mqttHuawei

import (
	"encoding/json"
	"goAdapter/device"
	"goAdapter/pkg/mylog"
	"goAdapter/pkg/system"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHuaweiPropertyValueTemplate struct {
	Value interface{} `json:"value"`
}

type MQTTHuaweiServiceTemplate struct {
	ServiceID  string                          `json:"service_id"`
	Properties MQTTHuaweiPropertyValueTemplate `json:"properties"`
}

type MQTTHuaweiDeviceServiceTemplate struct {
	DeviceID string                      `json:"device_id"`
	Services []MQTTHuaweiServiceTemplate `json:"services"`
}

type MQTTHuaweiReportPropertyTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
}

func MQTTHuaweiGWPropertyPost(client MQTT.Client, gw MQTTHuaweiRegisterTemplate, services []MQTTHuaweiServiceTemplate) int {

	type MQTTPropertyPayloadTemplate struct {
		Services []MQTTHuaweiServiceTemplate `json:"services"`
	}

	PropertyPayload := MQTTPropertyPayloadTemplate{
		Services: services,
	}

	MsgID = 0

	sJson, _ := json.Marshal(PropertyPayload)

	propertyPostTopic := "$oc/devices/" + gw.DeviceID + "/sys/properties/report"

	mylog.Logger.Infof("huawei gw property post topic: %s", propertyPostTopic)
	mylog.Logger.Debugf("huawei gw property post msg: %v", sJson)
	if client != nil {
		token := client.Publish(propertyPostTopic, 1, false, sJson)
		if token.WaitTimeout(2*time.Second) == true {
			MsgID = 0
			mylog.Logger.Debugf("huawei gw property post msg sucess")
		} else {
			MsgID = 1
			mylog.Logger.Debugf("huawei gw property post msg fail")
		}
	}

	return MsgID
}

func (r *ReportServiceParamHuaweiTemplate) GWPropertyPost() {

	services := make([]MQTTHuaweiServiceTemplate, 0)

	propertyValue := MQTTHuaweiPropertyValueTemplate{}
	propertyParams := MQTTHuaweiServiceTemplate{}

	//
	propertyParams.ServiceID = "MemTotal"
	propertyValue.Value = system.SystemState.MemTotal
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "MemUse"
	propertyValue.Value = system.SystemState.MemUse
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "DiskTotal"
	propertyValue.Value = system.SystemState.DiskTotal
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "DiskUse"
	propertyValue.Value = system.SystemState.DiskUse
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "Name"
	propertyValue.Value = system.SystemState.Name
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "SN"
	propertyValue.Value = system.SystemState.SN
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "HardVer"
	propertyValue.Value = system.SystemState.HardVer
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "SoftVer"
	propertyValue.Value = system.SystemState.SoftVer
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "SystemRTC"
	propertyValue.Value = system.SystemState.SystemRTC
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "RunTime"
	propertyValue.Value = system.SystemState.RunTime
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "DeviceOnline"
	propertyValue.Value = system.SystemState.DeviceOnline
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "DevicePacketLoss"
	propertyValue.Value = system.SystemState.DevicePacketLoss
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	mqttHuaweiRegister := MQTTHuaweiRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		DeviceID:     r.GWParam.Param.DeviceID,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	//上报故障先加，收到正确回应后清0
	r.GWParam.ReportErrCnt++
	mylog.Logger.Debugf("service %s gw ReportErrCnt %d", r.GWParam.ServiceName, r.GWParam.ReportErrCnt)
	//清空接收缓存
	for i := 0; i < len(r.ReceiveReportPropertyAckFrameChan); i++ {
		<-r.ReceiveReportPropertyAckFrameChan
	}
	if MQTTHuaweiGWPropertyPost(r.GWParam.MQTTClient, mqttHuaweiRegister, services) == 0 {
		r.GWParam.ReportErrCnt--
		mylog.Logger.Debugf("%s MQTTHuaweiGWPropertyPost OK", r.GWParam.ServiceName)
	} else {
		mylog.Logger.Debugf("%s MQTTHuaweiGWPropertyPost Err", r.GWParam.ServiceName)
	}
}

func MQTTHuaweiNodePropertyPost(client MQTT.Client, gw MQTTHuaweiRegisterTemplate, deviceServiceMap []MQTTHuaweiDeviceServiceTemplate) int {

	type MQTTHuaweiDevicesServiceTemplate struct {
		Devices []MQTTHuaweiDeviceServiceTemplate `json:"devices"`
	}

	DevicesService := MQTTHuaweiDevicesServiceTemplate{
		Devices: deviceServiceMap,
	}

	sJson, _ := json.Marshal(DevicesService)

	propertyPostTopic := "$oc/devices/" + gw.DeviceID + "/sys/gateway/sub_devices/properties/report"
	mylog.Logger.Infof("node property post topic: %v", propertyPostTopic)
	mylog.Logger.Debugf("node property post msg: %v", sJson)

	MsgID = 0
	if client != nil {
		token := client.Publish(propertyPostTopic, 1, false, sJson)
		if token.WaitTimeout(2*time.Second) == true {
			MsgID = 0
		} else {
			MsgID = 1
		}
	}

	return MsgID
}

func (r *ReportServiceParamHuaweiTemplate) NodePropertyPost(name []string) {

	nodeList := make([]ReportServiceNodeParamHuaweiTemplate, 0)
	for _, n := range name {
		for k, v := range r.NodeList {
			if n == v.Name {
				nodeList = append(nodeList, v)
				//上报故障计数值先加，收到正确回应后清0
				r.NodeList[k].ReportErrCnt++
			}
		}
	}

	pageCnt := len(nodeList) / 20 //单包最大发送20个设备
	if len(nodeList)%20 != 0 {
		pageCnt += 1
	}
	for pageIndex := 0; pageIndex < pageCnt; pageIndex++ {
		if pageIndex != (pageCnt - 1) {
			node := nodeList[20*pageIndex : 20*pageIndex+20]
			DeviceServiceMap := make([]MQTTHuaweiDeviceServiceTemplate, 0)
			for _, n := range node {
				tmps := device.CollectInterfaceMap.GetAll()
				for _, c := range tmps {
					if c.CollInterfaceName == n.CollInterfaceName {
						for _, d := range c.DeviceNodes {
							if d.Name == n.Name {
								ServiceMap := make([]MQTTHuaweiServiceTemplate, 0)
								for _, v := range d.Properties {
									if len(v.Value) >= 1 {
										index := len(v.Value) - 1
										service := MQTTHuaweiServiceTemplate{}
										service.ServiceID = v.Name
										service.Properties.Value = v.Value[index].Value
										ServiceMap = append(ServiceMap, service)
									}
								}
								deviceService := MQTTHuaweiDeviceServiceTemplate{}
								deviceService.DeviceID = n.Param.DeviceID
								deviceService.Services = ServiceMap
								DeviceServiceMap = append(DeviceServiceMap, deviceService)
							}
						}
					}
				}
			}

			mqttHuaweiRegister := MQTTHuaweiRegisterTemplate{
				RemoteIP:     r.GWParam.IP,
				RemotePort:   r.GWParam.Port,
				DeviceID:     r.GWParam.Param.DeviceID,
				DeviceSecret: r.GWParam.Param.DeviceSecret,
			}

			if MQTTHuaweiNodePropertyPost(r.GWParam.MQTTClient, mqttHuaweiRegister, DeviceServiceMap) == 0 {
				mylog.Logger.Debugf("%s MQTTHuaweiNodePropertyPost OK", r.GWParam.ServiceName)
				for _, n := range node {
					for k, v := range r.NodeList {
						if n.Name == v.Name {
							//上报故障计数值先加，收到正确回应后清0
							r.NodeList[k].ReportErrCnt--
						}
					}
				}
			} else {
				mylog.Logger.Debugf("%s MQTTHuaweiNodePropertyPost Err", r.GWParam.ServiceName)
			}
		} else { //最后一页
			node := nodeList[20*pageIndex:]
			DeviceServiceMap := make([]MQTTHuaweiDeviceServiceTemplate, 0)
			for _, n := range node {
				tmps := device.CollectInterfaceMap.GetAll()
				for _, c := range tmps {
					if c.CollInterfaceName == n.CollInterfaceName {
						for _, d := range c.DeviceNodes {
							if d.Name == n.Name {
								ServiceMap := make([]MQTTHuaweiServiceTemplate, 0)
								for _, v := range d.Properties {
									if len(v.Value) >= 1 {
										index := len(v.Value) - 1
										service := MQTTHuaweiServiceTemplate{}
										service.ServiceID = v.Name
										service.Properties.Value = v.Value[index].Value
										ServiceMap = append(ServiceMap, service)
									}
								}
								deviceService := MQTTHuaweiDeviceServiceTemplate{}
								deviceService.DeviceID = n.Param.DeviceID
								deviceService.Services = ServiceMap
								DeviceServiceMap = append(DeviceServiceMap, deviceService)
							}
						}
					}
				}
			}

			mqttHuaweiRegister := MQTTHuaweiRegisterTemplate{
				RemoteIP:     r.GWParam.IP,
				RemotePort:   r.GWParam.Port,
				DeviceID:     r.GWParam.Param.DeviceID,
				DeviceSecret: r.GWParam.Param.DeviceSecret,
			}
			if MQTTHuaweiNodePropertyPost(r.GWParam.MQTTClient, mqttHuaweiRegister, DeviceServiceMap) == 0 {
				mylog.Logger.Debugf("%s MQTTHuaweiNodePropertyPost OK", r.GWParam.ServiceName)
				for _, n := range node {
					for k, v := range r.NodeList {
						if n.Name == v.Name {
							//上报故障计数值先加，收到正确回应后清0
							r.NodeList[k].ReportErrCnt--
						}
					}
				}
			} else {
				mylog.Logger.Debugf("%s MQTTHuaweiNodePropertyPost Err", r.GWParam.ServiceName)
			}
		}
	}
}
