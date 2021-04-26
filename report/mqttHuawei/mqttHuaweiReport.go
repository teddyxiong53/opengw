package mqttHuawei

import (
	"encoding/json"
	"goAdapter/setting"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTHuaweiPropertyValueTemplate struct {
	Value interface{} `json:"value"`
}

type MQTTHuaweiPropertyParamsTemplate struct {
	ServiceID  string                          `json:"service_id"`
	Properties MQTTHuaweiPropertyValueTemplate `json:"properties"`
}

type MQTTHuaweiNodeValueTemplate struct {
	DeviceID string                             `json:"device_id"`
	Services []MQTTHuaweiPropertyParamsTemplate `json:"services"`
}

type MQTTHuaweiReportPropertyTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
}

func MQTTHuaweiGWPropertyPost(client MQTT.Client, gw MQTTHuaweiRegisterTemplate, services []MQTTHuaweiPropertyParamsTemplate) int {

	type MQTTPropertyPayloadTemplate struct {
		Services []MQTTHuaweiPropertyParamsTemplate `json:"services"`
	}

	PropertyPayload := MQTTPropertyPayloadTemplate{
		Services: services,
	}

	MsgID++

	sJson, _ := json.Marshal(PropertyPayload)

	propertyPostTopic := "$oc/devices/" + gw.DeviceID + "/sys/properties/report"

	setting.Logger.Infof("gw property post topic: %s", propertyPostTopic)
	setting.Logger.Debugf("gw property post msg: %s", sJson)
	if client != nil {
		token := client.Publish(propertyPostTopic, 0, false, sJson)
		token.Wait()
	}

	return MsgID
}

func (r *ReportServiceParamHuaweiTemplate) GWPropertyPost() {

	services := make([]MQTTHuaweiPropertyParamsTemplate, 0)

	propertyValue := MQTTHuaweiPropertyValueTemplate{}
	propertyParams := MQTTHuaweiPropertyParamsTemplate{}

	//
	propertyParams.ServiceID = "MemTotal"
	propertyValue.Value = setting.SystemState.MemTotal
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "MemUse"
	propertyValue.Value = setting.SystemState.MemUse
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "DiskTotal"
	propertyValue.Value = setting.SystemState.DiskTotal
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "DiskUse"
	propertyValue.Value = setting.SystemState.DiskUse
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "Name"
	propertyValue.Value = setting.SystemState.Name
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "SN"
	propertyValue.Value = setting.SystemState.SN
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "HardVer"
	propertyValue.Value = setting.SystemState.HardVer
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "SoftVer"
	propertyValue.Value = setting.SystemState.SoftVer
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "SystemRTC"
	propertyValue.Value = setting.SystemState.SystemRTC
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "RunTime"
	propertyValue.Value = setting.SystemState.RunTime
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "DeviceOnline"
	propertyValue.Value = setting.SystemState.DeviceOnline
	propertyParams.Properties = propertyValue
	services = append(services, propertyParams)

	propertyParams.ServiceID = "DevicePacketLoss"
	propertyValue.Value = setting.SystemState.DevicePacketLoss
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
	setting.Logger.Debugf("service %s,gw ReportErrCnt %d", r.GWParam.Param.DeviceID, r.GWParam.ReportErrCnt)
	//清空接收缓存
	for i := 0; i < len(r.ReceiveReportPropertyAckFrameChan); i++ {
		<-r.ReceiveReportPropertyAckFrameChan
	}
	MQTTHuaweiGWPropertyPost(r.GWParam.MQTTClient, mqttHuaweiRegister, services)

	select {
	case frame := <-r.ReceiveReportPropertyAckFrameChan:
		{
			setting.Logger.Debugf("frameCode %v", frame.Code)
			if frame.Code == 200 {
				r.GWParam.ReportErrCnt--
				setting.Logger.Debugf("%s,MQTTHuaweiGWPropertyPost OK", r.GWParam.ServiceName)
			} else {
				setting.Logger.Debugf("%s,MQTTHuaweiGWPropertyPost Err", r.GWParam.ServiceName)
			}
		}
	case <-time.After(time.Millisecond * 2000):
		{
			setting.Logger.Debugf("%s,MQTTHuaweiGWPropertyPost Err", r.GWParam.ServiceName)
		}
	}
}

func MQTTHuaweiNodePropertyPost(client MQTT.Client, gw MQTTHuaweiRegisterTemplate, nodeMap []MQTTHuaweiNodeValueTemplate) int {

	/*
		type MQTTPropertyValueTemplate struct {
			Value interface{} `json:"value"`
		}

		type MQTTNodeIdentityTemplate struct {
			ProductKey string `json:"productKey"`
			DeviceName string `json:"deviceName"`
		}

		type MQTTNodePropertyParamsTemplate struct {
			Identity   MQTTNodeIdentityTemplate             `json:"identity"`
			Properties map[string]MQTTPropertyValueTemplate `json:"properties"`
		}

		type MQTTNodesPropertyParamsTemplate struct {
			SubDevices []MQTTNodePropertyParamsTemplate `json:"subDevices"`
		}

		type MQTTPropertyPayloadTemplate struct {
			ID      string                          `json:"id"`
			Version string                          `json:"version"`
			Params  MQTTNodesPropertyParamsTemplate `json:"params"`
			Method  string                          `json:"method"`
		}

		MQTTNodesPropertyParams := MQTTNodesPropertyParamsTemplate{
			SubDevices: make([]MQTTNodePropertyParamsTemplate, 0),
		}

		for _, d := range nodeMap {
			MQTTNodePropertyParams := MQTTNodePropertyParamsTemplate{
				Properties: make(map[string]MQTTPropertyValueTemplate, 0),
			}
			MQTTNodePropertyParams.Identity.DeviceName = d.DeviceName
			MQTTNodePropertyParams.Identity.ProductKey = d.ProductKey

			for _, v := range d.ValueMap {
				MQTTPropertyValue := MQTTPropertyValueTemplate{}
				MQTTPropertyValue.Value = v.Value
				MQTTNodePropertyParams.Properties[v.Name] = MQTTPropertyValue
			}

			MQTTNodesPropertyParams.SubDevices = append(MQTTNodesPropertyParams.SubDevices, MQTTNodePropertyParams)
		}

		PropertyPayload := MQTTPropertyPayloadTemplate{
			ID:      strconv.Itoa(MsgID),
			Params:  MQTTNodesPropertyParams,
			Version: "1.0",
			Method:  "thing.event.property.pack.post",
		}
		MsgID++

		sJson, _ := json.Marshal(PropertyPayload)

		//propertyPostTopic := "/sys/" + MQTTAliyunGWParam.GWParam.ProductKey + "/" + MQTTAliyunGWParam.GWParam.DeviceName + "/thing/event/property/post"
		propertyPostTopic := "/sys/" + gw.ProductKey + "/" + gw.DeviceName + "/thing/event/property/pack/post"
		setting.Logger.Infof("node property post topic: %s\n", propertyPostTopic)
		setting.Logger.Debugf("node property post msg: %s\n", sJson)
		if client != nil {
			token := client.Publish(propertyPostTopic, 0, false, sJson)
			token.Wait()
		}

	*/

	return MsgID
}

func (r *ReportServiceParamHuaweiTemplate) AllNodePropertyPost() {

	/*
				//上报故障计数值先加，收到正确回应后清0
				for i := 0; i < len(r.NodeList); i++ {
					r.NodeList[i].ReportErrCnt++
				}

				pageCnt := len(r.NodeList) / 20 //单包最大发送20个设备
				if len(r.NodeList)%20 != 0 {
					pageCnt += 1
				}
				//log.Printf("pageCnt %v\n", pageCnt)
				for pageIndex := 0; pageIndex < pageCnt; pageIndex++ {
					//log.Printf("pageIndex %v\n", pageIndex)
					if pageIndex != (pageCnt - 1) {
						NodeValueMap := make([]MQTTHuaweiNodeValueTemplate, 0)
						valueMap := make([]MQTTHuaweiValueTemplate, 0)

						node := r.NodeList[20*pageIndex : 20*pageIndex+20]
						//log.Printf("nodeList %v\n", node)
						for _, n := range node {
							for _, c := range device.CollectInterfaceMap {
								if c.CollInterfaceName == n.CollInterfaceName {
									for _, d := range c.DeviceNodeMap {
										if d.Name == n.Name {
											for _, v := range d.VariableMap {
												if len(v.Value) >= 1 {
													index := len(v.Value) - 1
													mqttHuaweiValue := MQTTHuaweiValueTemplate{}
													propertyParams.ServiceID = v.Name
													propertyValue.Value = v.Value[index].Value
													propertyParams.Properties = propertyValue
		services = append(services, propertyParams)
												}
											}
											NodeValue := MQTTHuaweiNodeValueTemplate{}
											NodeValue.ValueMap = valueMap
											NodeValue.ProductKey = n.Param.ProductKey
											NodeValue.DeviceName = n.Param.DeviceName
											NodeValueMap = append(NodeValueMap, NodeValue)
										}
									}
								}
							}
						}

						mqttHuaweiRegister := MQTTHuaweiRegisterTemplate{
							RemoteIP:     r.GWParam.IP,
							RemotePort:   r.GWParam.Port,
							ProductKey:   r.GWParam.Param.ProductKey,
							DeviceName:   r.GWParam.Param.DeviceName,
							DeviceSecret: r.GWParam.Param.DeviceSecret,
						}
						//清空接收缓存
						for i := 0; i < len(r.ReceiveReportPropertyAckFrameChan); i++ {
							<-r.ReceiveReportPropertyAckFrameChan
						}
						MQTTHuaweiNodePropertyPost(r.GWParam.MQTTClient, mqttHuaweiRegister, NodeValueMap)
						select {
						case frame := <-r.ReceiveReportPropertyAckFrameChan:
							{
								if frame.Code == 200 {
									setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost OK", r.GWParam.ServiceName)
								} else {
									setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost Err", r.GWParam.ServiceName)
								}
							}
						case <-time.After(time.Millisecond * 2000):
							{
								setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost Err", r.GWParam.ServiceName)
							}
						}
					} else { //最后一页
						NodeValueMap := make([]MQTTHuaweiNodeValueTemplate, 0)
						valueMap := make([]MQTTHuaweiValueTemplate, 0)
						node := r.NodeList[20*pageIndex : len(r.NodeList)]
						//log.Printf("nodeList %v\n", node)
						for _, n := range node {
							for _, c := range device.CollectInterfaceMap {
								if c.CollInterfaceName == n.CollInterfaceName {
									for _, d := range c.DeviceNodeMap {
										if d.Name == n.Name {
											for _, v := range d.VariableMap {
												if len(v.Value) >= 1 {
													index := len(v.Value) - 1
													mqttHuaweiValue := MQTTHuaweiValueTemplate{}
													propertyParams.ServiceID = v.Name
													propertyValue.Value = v.Value[index].Value
													propertyParams.Properties = propertyValue
		services = append(services, propertyParams)
												}
											}
											NodeValue := MQTTHuaweiNodeValueTemplate{}
											NodeValue.ValueMap = valueMap
											NodeValue.ProductKey = n.Param.ProductKey
											NodeValue.DeviceName = n.Param.DeviceName
											NodeValueMap = append(NodeValueMap, NodeValue)
										}
									}
								}
							}
						}

						mqttHuaweiRegister := MQTTHuaweiRegisterTemplate{
							RemoteIP:     r.GWParam.IP,
							RemotePort:   r.GWParam.Port,
							ProductKey:   r.GWParam.Param.ProductKey,
							DeviceName:   r.GWParam.Param.DeviceName,
							DeviceSecret: r.GWParam.Param.DeviceSecret,
						}
						//清空接收缓存
						for i := 0; i < len(r.ReceiveReportPropertyAckFrameChan); i++ {
							<-r.ReceiveReportPropertyAckFrameChan
						}
						MQTTHuaweiNodePropertyPost(r.GWParam.MQTTClient, mqttHuaweiRegister, NodeValueMap)

						select {
						case frame := <-r.ReceiveReportPropertyAckFrameChan:
							{
								if frame.Code == 200 {
									setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost OK", r.GWParam.ServiceName)
								} else {
									setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost Err", r.GWParam.ServiceName)
								}
							}
						case <-time.After(time.Millisecond * 2000):
							{
								setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost Err", r.GWParam.ServiceName)
							}
						}
					}
				}

	*/
}

//指定设备上传属性
func (r *ReportServiceParamHuaweiTemplate) NodePropertyPost(name []string) {

	/*
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
		//log.Printf("pageCnt %v\n", pageCnt)
		for pageIndex := 0; pageIndex < pageCnt; pageIndex++ {
			//log.Printf("pageIndex %v\n", pageIndex)
			if pageIndex != (pageCnt - 1) {
				NodeValueMap := make([]MQTTHuaweiNodeValueTemplate, 0)
				valueMap := make([]MQTTHuaweiValueTemplate, 0)

				node := nodeList[20*pageIndex : 20*pageIndex+20]
				//log.Printf("nodeList %v\n", node)
				for _, n := range node {
					for _, c := range device.CollectInterfaceMap {
						if c.CollInterfaceName == n.CollInterfaceName {
							for _, d := range c.DeviceNodeMap {
								if d.Name == n.Name {
									for _, v := range d.VariableMap {
										if len(v.Value) >= 1 {
											index := len(v.Value) - 1
											mqttHuaweiValue := MQTTHuaweiValueTemplate{}
											propertyParams.ServiceID = v.Name
											propertyValue.Value = v.Value[index].Value
											propertyParams.Properties = propertyValue
											services = append(services, propertyParams)
										}
									}
									NodeValue := MQTTHuaweiNodeValueTemplate{}
									NodeValue.ValueMap = valueMap
									NodeValue.ProductKey = n.Param.ProductKey
									NodeValue.DeviceName = n.Param.DeviceName
									NodeValueMap = append(NodeValueMap, NodeValue)
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

				MQTTHuaweiNodePropertyPost(r.GWParam.MQTTClient, mqttHuaweiRegister, NodeValueMap)
				select {
				case frame := <-r.ReceiveReportPropertyAckFrameChan:
					{
						if frame.Code == 200 {
							setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost OK", r.GWParam.ServiceName)
						} else {
							setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost Err", r.GWParam.ServiceName)
						}
					}
				case <-time.After(time.Millisecond * 2000):
					{
						setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost Err", r.GWParam.ServiceName)
					}
				}
			} else { //最后一页
				NodeValueMap := make([]MQTTHuaweiNodeValueTemplate, 0)
				valueMap := make([]MQTTHuaweiValueTemplate, 0)
				node := nodeList[20*pageIndex : len(nodeList)]
				//log.Printf("nodeList %v\n", node)
				for _, n := range node {
					for _, c := range device.CollectInterfaceMap {
						if c.CollInterfaceName == n.CollInterfaceName {
							for _, d := range c.DeviceNodeMap {
								if d.Name == n.Name {
									for _, v := range d.VariableMap {
										if len(v.Value) >= 1 {
											index := len(v.Value) - 1
											mqttHuaweiValue := MQTTHuaweiValueTemplate{}
											propertyParams.ServiceID = v.Name
											propertyValue.Value = v.Value[index].Value
											propertyParams.Properties = propertyValue
											services = append(services, propertyParams)
										}
									}
									NodeValue := MQTTHuaweiNodeValueTemplate{}
									NodeValue.ValueMap = valueMap
									NodeValue.ProductKey = n.Param.ProductKey
									NodeValue.DeviceName = n.Param.DeviceName
									NodeValueMap = append(NodeValueMap, NodeValue)
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
				//setting.Logger.Debugf("NodeValueMap %v", NodeValueMap)
				MQTTHuaweiNodePropertyPost(r.GWParam.MQTTClient, mqttHuaweiRegister, NodeValueMap)

				select {
				case frame := <-r.ReceiveReportPropertyAckFrameChan:
					{
						if frame.Code == 200 {
							setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost OK", r.GWParam.ServiceName)
						} else {
							setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost Err", r.GWParam.ServiceName)
						}
					}
				case <-time.After(time.Millisecond * 2000):
					{
						setting.Logger.Debugf("%s,MQTTHuaweiNodePropertyPost Err", r.GWParam.ServiceName)
					}
				}
			}
		}

	*/
}
