package mqttEmqx

import (
	"encoding/json"
	"goAdapter/device"
	"goAdapter/pkg/mylog"
	"goAdapter/pkg/system"
	"strconv"
	"time"
)

type MQTTEmqxReportPropertyTemplate struct {
	DeviceType string //设备类型，"gw" "node"
	DeviceName []string
}

type MQTTEmqxPropertyPostParamPropertyTemplate struct {
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
	TimeStamp int64       `json:"timestamp"`
}

type MQTTEmqxPropertyPostParamTemplate struct {
	ClientID   string                                      `json:"clientID"`
	Properties []MQTTEmqxPropertyPostParamPropertyTemplate `json:"properties"`
}

type MQTTEmqxPropertyPostTemplate struct {
	ID      string                              `json:"id"`
	Version string                              `json:"version"`
	Ack     int                                 `json:"ack"`
	Params  []MQTTEmqxPropertyPostParamTemplate `json:"params"`
}

func MQTTEmqxPropertyPost(gwParam ReportServiceGWParamEmqxTemplate, propertyParam []MQTTEmqxPropertyPostParamTemplate) int {

	propertyPost := MQTTEmqxPropertyPostTemplate{
		ID:      strconv.Itoa(MsgID),
		Version: "V1.0",
		Ack:     1,
		Params:  propertyParam,
	}
	MsgID++

	sJson, _ := json.Marshal(propertyPost)
	propertyPostTopic := "/sys/thing/event/property/post/" + gwParam.Param.ClientID

	mylog.Logger.Infof("property post topic: %s", propertyPostTopic)
	mylog.Logger.Debugf("property post msg: %v", string(sJson))
	if gwParam.MQTTClient != nil {
		token := gwParam.MQTTClient.Publish(propertyPostTopic, 0, false, sJson)
		token.Wait()
	}

	return MsgID
}

func (r *ReportServiceParamEmqxTemplate) GWPropertyPost() {

	propertyMap := make([]MQTTEmqxPropertyPostParamPropertyTemplate, 0)

	property := MQTTEmqxPropertyPostParamPropertyTemplate{}

	timeStamp := time.Now().Unix()

	property.Name = "MemTotal"
	property.Value = system.SystemState.MemTotal
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "MemUse"
	property.Value = system.SystemState.MemUse
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "DiskTotal"
	property.Value = system.SystemState.DiskTotal
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "DiskUse"
	property.Value = system.SystemState.DiskUse
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "Name"
	property.Value = system.SystemState.Name
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "SN"
	property.Value = system.SystemState.SN
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "HardVer"
	property.Value = system.SystemState.HardVer
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "SoftVer"
	property.Value = system.SystemState.SoftVer
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "SystemRTC"
	property.Value = system.SystemState.SystemRTC
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "RunTime"
	property.Value = system.SystemState.RunTime
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "DeviceOnline"
	property.Value = system.SystemState.DeviceOnline
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	property.Name = "DevicePacketLoss"
	property.Value = system.SystemState.DevicePacketLoss
	property.TimeStamp = timeStamp
	propertyMap = append(propertyMap, property)

	//上报故障先加，收到正确回应后清0
	r.GWParam.ReportErrCnt++
	mylog.Logger.Debugf("service %s gw ReportErrCnt %d", r.GWParam.Param.ClientID, r.GWParam.ReportErrCnt)
	//清空接收缓存
	for i := 0; i < len(r.ReceiveReportPropertyAckFrameChan); i++ {
		<-r.ReceiveReportPropertyAckFrameChan
	}

	propertyPostParam := MQTTEmqxPropertyPostParamTemplate{
		ClientID:   r.GWParam.Param.ClientID,
		Properties: propertyMap,
	}

	propertyPostParamMap := make([]MQTTEmqxPropertyPostParamTemplate, 0)
	propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)
	MQTTEmqxPropertyPost(r.GWParam, propertyPostParamMap)

	select {
	case frame := <-r.ReceiveReportPropertyAckFrameChan:
		{
			mylog.Logger.Debugf("frameCode %v", frame.Code)
			if frame.Code == 200 {
				r.GWParam.ReportErrCnt--
				mylog.Logger.Debugf("%s MQTTEmqxGWPropertyPost OK", r.GWParam.ServiceName)
			} else {
				mylog.Logger.Debugf("%s MQTTEmqxGWPropertyPost Err", r.GWParam.ServiceName)
			}
		}
	case <-time.After(time.Millisecond * 2000):
		{
			mylog.Logger.Debugf("%s MQTTEmqxGWPropertyPost Err", r.GWParam.ServiceName)
		}
	}
}

//指定设备上传属性
func (r *ReportServiceParamEmqxTemplate) NodePropertyPost(name []string) {

	propertyPostParamMap := make([]MQTTEmqxPropertyPostParamTemplate, 0)
	for _, n := range name {
		for k, v := range r.NodeList {
			if n == v.Name {
				//上报故障计数值先加，收到正确回应后清0
				r.NodeList[k].ReportErrCnt++
				propertyPostParam := MQTTEmqxPropertyPostParamTemplate{
					ClientID: v.Param.ClientID,
				}
				timeStamp := time.Now().Unix()
				tmps := device.CollectInterfaceMap.GetAll()
				for _, c := range tmps {
					if c.CollInterfaceName == v.CollInterfaceName {
						for _, d := range c.DeviceNodes {
							if d.Name == v.Name {
								for _, v := range d.VariableMap {
									if len(v.Values) >= 1 {
										index := len(v.Values) - 1
										property := MQTTEmqxPropertyPostParamPropertyTemplate{}
										property.Name = v.Name
										property.Value = v.Values[index].Value
										property.TimeStamp = timeStamp
										propertyPostParam.Properties = append(propertyPostParam.Properties, property)
									}
								}
							}
						}
					}
				}
				propertyPostParamMap = append(propertyPostParamMap, propertyPostParam)
			}
		}
	}

	mylog.Logger.Debugf("propertyPostParamMap %v", propertyPostParamMap)

	pageCnt := len(propertyPostParamMap) / 20 //单包最大发送20个设备
	if len(propertyPostParamMap)%20 != 0 {
		pageCnt += 1
	}

	for pageIndex := 0; pageIndex < pageCnt; pageIndex++ {
		if pageIndex != (pageCnt - 1) {
			MQTTEmqxPropertyPost(r.GWParam, propertyPostParamMap[pageIndex:pageIndex+20])
		} else { //最后一页
			MQTTEmqxPropertyPost(r.GWParam, propertyPostParamMap[pageIndex+20*(pageCnt-1):])
		}
		select {
		case frame := <-r.ReceiveReportPropertyAckFrameChan:
			{
				if frame.Code == 200 {
					mylog.Logger.Debugf("%s MQTTEmqxNodePropertyPost OK", r.GWParam.ServiceName)
				} else {
					mylog.Logger.Debugf("%s MQTTEmqxNodePropertyPost Err", r.GWParam.ServiceName)
				}
			}
		case <-time.After(time.Millisecond * 2000):
			{
				mylog.Logger.Debugf("%s MQTTEmqxNodePropertyPost Err", r.GWParam.ServiceName)
			}
		}
	}
}

//func (r *ReportServiceParamEmqxTemplate) NodePropertyPost(name []string) {
//
//	nodeList := make([]ReportServiceNodeParamEmqxTemplate, 0)
//	for _, n := range name {
//		for k, v := range r.NodeList {
//			if n == v.Name {
//				nodeList = append(nodeList, v)
//				//上报故障计数值先加，收到正确回应后清0
//				r.NodeList[k].ReportErrCnt++
//			}
//		}
//	}
//
//	pageCnt := len(nodeList) / 20 //单包最大发送20个设备
//	if len(nodeList)%20 != 0 {
//		pageCnt += 1
//	}
//
//	for pageIndex := 0; pageIndex < pageCnt; pageIndex++ {
//		if pageIndex != (pageCnt - 1) {
//			propertyPost := MQTTPropertyPostTemplate{
//				ID: strconv.Itoa(MsgID),
//				Version: "V1.0",
//				Ack: 1,
//			}
//			node := nodeList[20*pageIndex : 20*pageIndex+20]
//			for _, n := range node {
//				for _, c := range device.CollectInterfaceMap {
//					if c.CollInterfaceName == n.CollInterfaceName {
//						for _, d := range c.DeviceNodeMap {
//							if d.Name == n.Name {
//								propertyPostParam := MQTTPropertyPostParamTemplate{
//									ClientID: n
//								}
//								for _, v := range d.VariableMap {
//									if len(v.Value) >= 1 {
//										index := len(v.Value) - 1
//										property := propertyTemplate{}
//										property.Name = v.Name
//										property.Value = v.Value[index].Value
//										valueMap = append(valueMap, property)
//									}
//								}
//								NodeValue := MQTTEmqxNodeValueTemplate{}
//								NodeValue.ValueMap = valueMap
//								NodeValue.ProductKey = n.Param.ProductKey
//								NodeValue.DeviceName = n.Param.DeviceName
//								NodeValueMap = append(NodeValueMap, NodeValue)
//							}
//						}
//					}
//				}
//			}
//
//
//			MsgID++
//			MQTTEmqxPropertyPost(r.GWParam, NodeValueMap)
//			select {
//			case frame := <-r.ReceiveReportPropertyAckFrameChan:
//				{
//					if frame.Code == 200 {
//						mylog.Logger.Debugf("%s MQTTEmqxNodePropertyPost OK", r.GWParam.ServiceName)
//					} else {
//						mylog.Logger.Debugf("%s MQTTEmqxNodePropertyPost Err", r.GWParam.ServiceName)
//					}
//				}
//			case <-time.After(time.Millisecond * 2000):
//				{
//					mylog.Logger.Debugf("%s MQTTEmqxNodePropertyPost Err", r.GWParam.ServiceName)
//				}
//			}
//		} else { //最后一页
//			NodeValueMap := make([]MQTTEmqxNodeValueTemplate, 0)
//			valueMap := make([]propertyTemplate, 0)
//			node := nodeList[20*pageIndex : len(nodeList)]
//			//log.Printf("nodeList %v\n", node)
//			for _, n := range node {
//				for _, c := range device.CollectInterfaceMap {
//					if c.CollInterfaceName == n.CollInterfaceName {
//						for _, d := range c.DeviceNodeMap {
//							if d.Name == n.Name {
//								for _, v := range d.VariableMap {
//									if len(v.Value) >= 1 {
//										index := len(v.Value) - 1
//										property := propertyTemplate{}
//										property.Name = v.Name
//										property.Value = v.Value[index].Value
//										valueMap = append(valueMap, property)
//									}
//								}
//								NodeValue := MQTTEmqxNodeValueTemplate{}
//								NodeValue.ValueMap = valueMap
//								NodeValue.ProductKey = n.Param.ProductKey
//								NodeValue.DeviceName = n.Param.DeviceName
//								NodeValueMap = append(NodeValueMap, NodeValue)
//							}
//						}
//					}
//				}
//			}
//
//			mqttEmqxRegister := MQTTEmqxRegisterTemplate{
//				RemoteIP:     r.GWParam.IP,
//				RemotePort:   r.GWParam.Port,
//				ProductKey:   r.GWParam.Param.ProductKey,
//				DeviceName:   r.GWParam.Param.DeviceName,
//				DeviceSecret: r.GWParam.Param.DeviceSecret,
//			}
//			//mylog.Logger.Debugf("NodeValueMap %v", NodeValueMap)
//			MQTTEmqxNodePropertyPost(r.GWParam.MQTTClient, mqttEmqxRegister, NodeValueMap)
//
//			select {
//			case frame := <-r.ReceiveReportPropertyAckFrameChan:
//				{
//					if frame.Code == 200 {
//						mylog.Logger.Debugf("%s MQTTEmqxNodePropertyPost OK", r.GWParam.ServiceName)
//					} else {
//						mylog.Logger.Debugf("%s MQTTEmqxNodePropertyPost Err", r.GWParam.ServiceName)
//					}
//				}
//			case <-time.After(time.Millisecond * 2000):
//				{
//					mylog.Logger.Debugf("%s MQTTEmqxNodePropertyPost Err", r.GWParam.ServiceName)
//				}
//			}
//		}
//	}
//}
