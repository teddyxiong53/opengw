package mqttAliyun

import (
	"goAdapter/device"
	"goAdapter/setting"
	"time"
)

func (r *ReportServiceParamAliyunTemplate) GWPropertyPost() {

	valueMap := make([]MQTTAliyunValueTemplate, 0)

	mqttAliyunValue := MQTTAliyunValueTemplate{}

	mqttAliyunValue.Name = "MemTotal"
	mqttAliyunValue.Value = setting.SystemState.MemTotal
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "MemUse"
	mqttAliyunValue.Value = setting.SystemState.MemUse
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "DiskTotal"
	mqttAliyunValue.Value = setting.SystemState.DiskTotal
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "DiskUse"
	mqttAliyunValue.Value = setting.SystemState.DiskUse
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "Name"
	mqttAliyunValue.Value = setting.SystemState.Name
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "SN"
	mqttAliyunValue.Value = setting.SystemState.SN
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "HardVer"
	mqttAliyunValue.Value = setting.SystemState.HardVer
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "SoftVer"
	mqttAliyunValue.Value = setting.SystemState.SoftVer
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "SystemRTC"
	mqttAliyunValue.Value = setting.SystemState.SystemRTC
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "RunTime"
	mqttAliyunValue.Value = setting.SystemState.RunTime
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "DeviceOnline"
	mqttAliyunValue.Value = setting.SystemState.DeviceOnline
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunValue.Name = "DevicePacketLoss"
	mqttAliyunValue.Value = setting.SystemState.DevicePacketLoss
	valueMap = append(valueMap, mqttAliyunValue)

	mqttAliyunRegister := MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	//上报故障先加，收到正确回应后清0
	r.GWParam.ReportErrCnt++
	setting.Logger.Debugf("service %s,gw ReportErrCnt %d", r.GWParam.Param.DeviceName, r.GWParam.ReportErrCnt)

	MQTTAliyunGWPropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, valueMap)
	select {
	case frame := <-r.ReceiveReportNodePropertyAckFrameChan:
		{
			if frame.Code == 200 {

			}
		}
	case <-time.After(time.Millisecond * 2000):
		{

		}
	}

}

func (r *ReportServiceParamAliyunTemplate) AllNodePropertyPost() {

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
			NodeValueMap := make([]MQTTAliyunNodeValueTemplate, 0)
			valueMap := make([]MQTTAliyunValueTemplate, 0)

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
										mqttAliyunValue := MQTTAliyunValueTemplate{}
										mqttAliyunValue.Name = v.Name
										mqttAliyunValue.Value = v.Value[index].Value
										valueMap = append(valueMap, mqttAliyunValue)
									}
								}
								NodeValue := MQTTAliyunNodeValueTemplate{}
								NodeValue.ValueMap = valueMap
								NodeValue.ProductKey = n.Param.ProductKey
								NodeValue.DeviceName = n.Param.DeviceName
								NodeValueMap = append(NodeValueMap, NodeValue)
							}
						}
					}
				}
			}

			mqttAliyunRegister := MQTTAliyunRegisterTemplate{
				RemoteIP:     r.GWParam.IP,
				RemotePort:   r.GWParam.Port,
				ProductKey:   r.GWParam.Param.ProductKey,
				DeviceName:   r.GWParam.Param.DeviceName,
				DeviceSecret: r.GWParam.Param.DeviceSecret,
			}

			MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap)
			select {
			case frame := <-r.ReceiveReportNodePropertyAckFrameChan:
				{
					if frame.Code == 200 {

					}
				}
			case <-time.After(time.Millisecond * 1000):
				{

				}
			}
		} else { //最后一页
			NodeValueMap := make([]MQTTAliyunNodeValueTemplate, 0)
			valueMap := make([]MQTTAliyunValueTemplate, 0)
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
										mqttAliyunValue := MQTTAliyunValueTemplate{}
										mqttAliyunValue.Name = v.Name
										mqttAliyunValue.Value = v.Value[index].Value
										valueMap = append(valueMap, mqttAliyunValue)
									}
								}
								NodeValue := MQTTAliyunNodeValueTemplate{}
								NodeValue.ValueMap = valueMap
								NodeValue.ProductKey = n.Param.ProductKey
								NodeValue.DeviceName = n.Param.DeviceName
								NodeValueMap = append(NodeValueMap, NodeValue)
							}
						}
					}
				}
			}

			mqttAliyunRegister := MQTTAliyunRegisterTemplate{
				RemoteIP:     r.GWParam.IP,
				RemotePort:   r.GWParam.Port,
				ProductKey:   r.GWParam.Param.ProductKey,
				DeviceName:   r.GWParam.Param.DeviceName,
				DeviceSecret: r.GWParam.Param.DeviceSecret,
			}
			//setting.Logger.Debugf("NodeValueMap %v", NodeValueMap)
			MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap)
			select {
			case frame := <-r.ReceiveReportNodePropertyAckFrameChan:
				{
					if frame.Code == 200 {

					}
				}
			case <-time.After(time.Millisecond * 1000):
				{

				}
			}
		}
	}
}

//指定设备上传属性
func (r *ReportServiceParamAliyunTemplate) NodePropertyPost(name []string) {

	nodeList := make([]ReportServiceNodeParamAliyunTemplate, 0)
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
			NodeValueMap := make([]MQTTAliyunNodeValueTemplate, 0)
			valueMap := make([]MQTTAliyunValueTemplate, 0)

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
										mqttAliyunValue := MQTTAliyunValueTemplate{}
										mqttAliyunValue.Name = v.Name
										mqttAliyunValue.Value = v.Value[index].Value
										valueMap = append(valueMap, mqttAliyunValue)
									}
								}
								NodeValue := MQTTAliyunNodeValueTemplate{}
								NodeValue.ValueMap = valueMap
								NodeValue.ProductKey = n.Param.ProductKey
								NodeValue.DeviceName = n.Param.DeviceName
								NodeValueMap = append(NodeValueMap, NodeValue)
							}
						}
					}
				}
			}

			mqttAliyunRegister := MQTTAliyunRegisterTemplate{
				RemoteIP:     r.GWParam.IP,
				RemotePort:   r.GWParam.Port,
				ProductKey:   r.GWParam.Param.ProductKey,
				DeviceName:   r.GWParam.Param.DeviceName,
				DeviceSecret: r.GWParam.Param.DeviceSecret,
			}

			MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap)
			select {
			case frame := <-r.ReceiveReportNodePropertyAckFrameChan:
				{
					if frame.Code == 200 {

					}
				}
			case <-time.After(time.Millisecond * 1000):
				{

				}
			}
		} else { //最后一页
			NodeValueMap := make([]MQTTAliyunNodeValueTemplate, 0)
			valueMap := make([]MQTTAliyunValueTemplate, 0)
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
										mqttAliyunValue := MQTTAliyunValueTemplate{}
										mqttAliyunValue.Name = v.Name
										mqttAliyunValue.Value = v.Value[index].Value
										valueMap = append(valueMap, mqttAliyunValue)
									}
								}
								NodeValue := MQTTAliyunNodeValueTemplate{}
								NodeValue.ValueMap = valueMap
								NodeValue.ProductKey = n.Param.ProductKey
								NodeValue.DeviceName = n.Param.DeviceName
								NodeValueMap = append(NodeValueMap, NodeValue)
							}
						}
					}
				}
			}

			mqttAliyunRegister := MQTTAliyunRegisterTemplate{
				RemoteIP:     r.GWParam.IP,
				RemotePort:   r.GWParam.Port,
				ProductKey:   r.GWParam.Param.ProductKey,
				DeviceName:   r.GWParam.Param.DeviceName,
				DeviceSecret: r.GWParam.Param.DeviceSecret,
			}
			//setting.Logger.Debugf("NodeValueMap %v", NodeValueMap)
			MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap)
			select {
			case frame := <-r.ReceiveReportNodePropertyAckFrameChan:
				{
					if frame.Code == 200 {

					}
				}
			case <-time.After(time.Millisecond * 1000):
				{

				}
			}
		}
	}
}
