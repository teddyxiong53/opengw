package report

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/robfig/cron"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"goAdapter/device"
	mqttClient "goAdapter/mqttClient/mqttAliyun"
	"goAdapter/setting"
)

type ReportServiceNodeParamAliyunTemplate struct {
	ServiceName       string
	CollInterfaceName string
	Name              string
	Addr              string
	CommStatus        string
	ReportErrCnt      int `json:"-"`
	ReportStatus      string
	Protocol          string
	Param             struct {
		ProductKey   string
		DeviceName   string
		DeviceSecret string
	}
}

type ReportServiceAliyunMessageTemplate struct {
	Topic   string
	Payload []byte
}

type ReportServiceGWParamAliyunTemplate struct {
	ServiceName string
	IP          string
	Port        string
	ReportTime  int
	Protocol    string
	Param       struct {
		ProductKey   string
		DeviceName   string
		DeviceSecret string
	}
	MQTTClient MQTT.Client `json:"-"`
}

type ReportServiceParamAliyunTemplate struct {
	CommStatus  string
	GWParam     ReportServiceGWParamAliyunTemplate
	NodeList    []ReportServiceNodeParamAliyunTemplate
	MessageChan chan ReportServiceMessageAliyunTemplate `json:"-"`
}

type ReportServiceParamListAliyunTemplate struct {
	ServiceList []*ReportServiceParamAliyunTemplate
}

type ReportServiceMessageAliyunTemplate struct {
	Topic   string
	Payload []byte
}

var ReportServiceParamListAliyun = &ReportServiceParamListAliyunTemplate{
	ServiceList: make([]*ReportServiceParamAliyunTemplate, 0),
}

func init() {

	ReportServiceParamListAliyun.ReadParamFromJson()

	//初始化chan
	for _, v := range ReportServiceParamListAliyun.ServiceList {
		v.MessageChan = make(chan ReportServiceMessageAliyunTemplate, 10)

		//go ReportServiceAliyunPoll(v)
	}
}

func (s *ReportServiceParamListAliyunTemplate) ReadParamFromJson() bool {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/reportServiceParamListAliyun.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open reportServiceParamListAliyun.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt], s)
		if err != nil {
			log.Println("reportServiceParamListAliyun unmarshal err", err)
			return false
		}
		log.Println("read reportServiceParamListAliyun.json success")

		return true
	} else {
		log.Println("reportServiceParamListAliyun.json is not exist")

		return false
	}
}

func (s *ReportServiceParamListAliyunTemplate) WriteParamToJson() {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/reportServiceParamListAliyun.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open reportServiceParamListAliyun.json err", err)
		return
	}
	defer fp.Close()

	sJson, _ := json.Marshal(*s)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write reportServiceParamListAliyun.json err", err)
	}
	log.Println("write reportServiceParamListAliyun.json success")
}

func (s *ReportServiceParamListAliyunTemplate) AddReportService(param ReportServiceGWParamAliyunTemplate) {

	for k, v := range s.ServiceList {
		//存在相同的，表示修改;不存在表示增加
		if v.GWParam.ServiceName == param.ServiceName {

			s.ServiceList[k].GWParam = param
			s.WriteParamToJson()
			return
		}
	}

	ReportServiceParam := &ReportServiceParamAliyunTemplate{
		GWParam:    param,
		CommStatus: "offLine",
	}
	s.ServiceList = append(s.ServiceList, ReportServiceParam)

	s.WriteParamToJson()
}

func (s *ReportServiceParamListAliyunTemplate) DeleteReportService(serviceName string) {

	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == serviceName {

			s.ServiceList = append(s.ServiceList[:k], s.ServiceList[k+1:]...)
			s.WriteParamToJson()
			return
		}
	}
}

func (r *ReportServiceParamAliyunTemplate) AddReportNode(param ReportServiceNodeParamAliyunTemplate) {

	param.CommStatus = "offLine"
	param.ReportStatus = "offLine"

	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Addr == param.Addr {
			r.NodeList[k] = param
			ReportServiceParamListAliyun.WriteParamToJson()
			return
		}
	}

	//节点不存在则新建
	r.NodeList = append(r.NodeList, param)
	ReportServiceParamListAliyun.WriteParamToJson()

	log.Printf("param %v\n", ReportServiceParamListAliyun)
}

func (r *ReportServiceParamAliyunTemplate) DeleteReportNode(addr string) {

	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Addr == addr {
			r.NodeList = append(r.NodeList[:k], r.NodeList[k+1:]...)
			ReportServiceParamListAliyun.WriteParamToJson()
			return
		}
	}
}

func GWPublishHandler(client MQTT.Client, msg MQTT.Message) {

	for _, v := range ReportServiceParamListAliyun.ServiceList {
		if v.GWParam.MQTTClient == client {
			go ReportServiceAliyunProcessMessage(v, msg.Topic(), msg.Payload())
		}
	}
}

func (r *ReportServiceParamAliyunTemplate) GWLogin() bool {

	mqttAliyunRegister := mqttClient.MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	status := false
	status, r.GWParam.MQTTClient = mqttClient.MQTTAliyunGWLogin(mqttAliyunRegister, GWPublishHandler)

	return status
}

func (r *ReportServiceParamAliyunTemplate) NodeLogin(addr []string) bool {

	nodeList := make([]mqttClient.MQTTAliyunNodeRegisterTemplate, 0)
	nodeParam := mqttClient.MQTTAliyunNodeRegisterTemplate{}

	for _, d := range addr {
		for k, v := range r.NodeList {
			if d == v.Addr {
				nodeParam.DeviceSecret = v.Param.DeviceSecret
				nodeParam.DeviceName = v.Param.DeviceName
				nodeParam.ProductKey = v.Param.ProductKey
				nodeList = append(nodeList, nodeParam)
				r.NodeList[k].CommStatus = "onLine"
			}
		}
	}

	mqttAliyunRegister := mqttClient.MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}
	mqttClient.MQTTAliyunNodeLoginIn(r.GWParam.MQTTClient, mqttAliyunRegister, nodeList)

	timerOut := time.NewTimer(500 * time.Millisecond)
	select {
	case ackMessage := <-r.MessageChan:
		if strings.Contains(ackMessage.Topic, "/combine/batch_login_reply") {
			log.Printf("Node combine/login ok")
		} else {
			log.Printf("Node combine/login err")
		}
	case <-timerOut.C:
		timerOut.Stop()
		log.Printf("Node combine/login err")

	}
	return true
}

func (r *ReportServiceParamAliyunTemplate) NodeLogOut(addr []string) bool {

	nodeList := make([]mqttClient.MQTTAliyunNodeRegisterTemplate, 0)
	nodeParam := mqttClient.MQTTAliyunNodeRegisterTemplate{}

	for _, d := range addr {
		for k, v := range r.NodeList {
			if d == v.Addr {
				nodeParam.DeviceSecret = v.Param.DeviceSecret
				nodeParam.DeviceName = v.Param.DeviceName
				nodeParam.ProductKey = v.Param.ProductKey

				nodeList = append(nodeList, nodeParam)
				r.NodeList[k].CommStatus = "offLine"
			}
		}
	}

	mqttAliyunRegister := mqttClient.MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}
	mqttClient.MQTTAliyunNodeLoginOut(r.GWParam.MQTTClient, mqttAliyunRegister, nodeList)

	timerOut := time.NewTimer(500 * time.Millisecond)
	select {
	case ackMessage := <-r.MessageChan:
		if strings.Contains(ackMessage.Topic, "/combine/batch_logout_reply") {
			log.Printf("Node combine/logout ok")
		} else {
			log.Printf("Node combine/logout err")
		}
	case <-timerOut.C:
		timerOut.Stop()
		log.Printf("Node combine/logout err")
	}

	return true
}

func (r *ReportServiceParamAliyunTemplate) GWPropertyPost() {

	valueMap := make([]mqttClient.MQTTAliyunValueTemplate, 0)

	mqttAliyunValue := mqttClient.MQTTAliyunValueTemplate{}

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

	mqttAliyunRegister := mqttClient.MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	mqttClient.MQTTAliyunGWPropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, valueMap)

	timerOut := time.NewTimer(500 * time.Millisecond)
	select {
	case ackMessage := <-r.MessageChan:
		if strings.Contains(ackMessage.Topic, "/thing/event/property/pack/post_reply") {
			log.Printf("gw property post ok")
		} else {
			log.Printf("gw property post err")
		}
	case <-timerOut.C:
		timerOut.Stop()
		log.Printf("gw property post err")
	}

}

func (r *ReportServiceParamAliyunTemplate) AllNodePropertyPost() {

	NodeValueMap := make([]mqttClient.MQTTAliyunNodeValueTemplate, 0)
	valueMap := make([]mqttClient.MQTTAliyunValueTemplate, 0)

	for _, n := range r.NodeList {
		for _, c := range device.CollectInterfaceMap {
			if c.CollInterfaceName == n.CollInterfaceName {
				for _, d := range c.DeviceNodeMap {
					if d.Addr == n.Addr {
						for _, v := range d.VariableMap {
							if len(v.Value) > 1 {
								index := len(v.Value) - 1
								mqttAliyunValue := mqttClient.MQTTAliyunValueTemplate{}
								mqttAliyunValue.Name = v.Name
								mqttAliyunValue.Value = v.Value[index].Value
								valueMap = append(valueMap, mqttAliyunValue)
							}
						}
						NodeValue := mqttClient.MQTTAliyunNodeValueTemplate{}
						NodeValue.ValueMap = valueMap
						NodeValue.ProductKey = n.Param.ProductKey
						NodeValue.DeviceName = n.Param.DeviceName
						NodeValueMap = append(NodeValueMap, NodeValue)
					}
				}
			}
		}
	}

	mqttAliyunRegister := mqttClient.MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	mqttClient.MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap)

	timerOut := time.NewTimer(500 * time.Millisecond)
	select {
	case ackMessage := <-r.MessageChan:
		if strings.Contains(ackMessage.Topic, "/thing/event/property/pack/post_reply") {
			log.Printf("node property post ok")
			for k, _ := range r.NodeList {
				r.NodeList[k].ReportErrCnt = 0
				r.NodeList[k].ReportStatus = "onLine"
			}
		} else {
			log.Printf("node property post err")
			for k, _ := range r.NodeList {
				r.NodeList[k].ReportErrCnt++
				if r.NodeList[k].ReportErrCnt >= 3 {
					r.NodeList[k].ReportErrCnt = 0
					r.NodeList[k].ReportStatus = "offLine"
				}
			}
		}
	case <-timerOut.C:
		timerOut.Stop()
		log.Printf("node property post err")
	}
}

//指定设备上传属性
func (r *ReportServiceParamAliyunTemplate) NodePropertyPost(addr []string) {

	NodeValueMap := make([]mqttClient.MQTTAliyunNodeValueTemplate, 0)
	valueMap := make([]mqttClient.MQTTAliyunValueTemplate, 0)

	for _, a := range addr {
		for _, n := range r.NodeList {
			if a == n.Addr {
				deviceName := n.Param.DeviceName
				productKey := n.Param.ProductKey
				for _, c := range device.CollectInterfaceMap {
					if n.CollInterfaceName == c.CollInterfaceName {
						for _, d := range c.DeviceNodeMap {
							if a == d.Addr {
								for _, v := range d.VariableMap {
									if len(v.Value) > 1 {
										index := len(v.Value) - 1
										mqttAliyunValue := mqttClient.MQTTAliyunValueTemplate{}
										mqttAliyunValue.Name = v.Name
										mqttAliyunValue.Value = v.Value[index].Value
										valueMap = append(valueMap, mqttAliyunValue)
									}
								}
							}
						}
					}
				}
				NodeValue := mqttClient.MQTTAliyunNodeValueTemplate{}
				NodeValue.ValueMap = valueMap
				NodeValue.ProductKey = productKey
				NodeValue.DeviceName = deviceName
				NodeValueMap = append(NodeValueMap, NodeValue)
			}
		}
	}

	//if len(valueMap) > 0 {
	mqttAliyunRegister := mqttClient.MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	mqttClient.MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap)

	timerOut := time.NewTimer(500 * time.Millisecond)
	select {
	case ackMessage := <-r.MessageChan:
		if strings.Contains(ackMessage.Topic, "/thing/event/property/pack/post_reply") {
			log.Printf("node property post ok")
			for _, a := range addr {
				for k, n := range r.NodeList {
					if a == n.Addr {
						r.NodeList[k].ReportErrCnt = 0
						r.NodeList[k].ReportStatus = "onLine"
					}
				}
			}
		} else {
			log.Printf("node property post err")
			for _, a := range addr {
				for k, n := range r.NodeList {
					if a == n.Addr {
						r.NodeList[k].ReportErrCnt++
						if r.NodeList[k].ReportErrCnt >= 3 {
							r.NodeList[k].ReportErrCnt = 0
							r.NodeList[k].ReportStatus = "offLine"
						}
					}
				}
			}
		}
	case <-timerOut.C:
		timerOut.Stop()
		log.Printf("node property post err")
	}
}

func ReportServiceAliyunPoll(r *ReportServiceParamAliyunTemplate) {

	reportState := 0

	// 定义一个cron运行器
	cronProcess := cron.New()

	str := fmt.Sprintf("@every %dm%ds", r.GWParam.ReportTime/60, r.GWParam.ReportTime%60)
	setting.Logger.Infof("reportServiceAliyun %+v", str)

	cronProcess.Start()
	defer cronProcess.Stop()

	addr := make([]string, 0)

	for {
		switch reportState {
		case 0:
			{
				if r.GWLogin() == true {
					reportState = 1

					cronProcess.AddFunc(str, r.GWPropertyPost)
					cronProcess.AddFunc(str, r.AllNodePropertyPost)
				} else {
					time.Sleep(5 * time.Second)
				}
			}
		case 1:
			{
				//节点发生了上线
				for _, c := range device.CollectInterfaceMap {
					for i := 0; i < len(c.OnlineReportChan); i++ {
						addr = append(addr, <-c.OnlineReportChan)
					}
				}
				if len(addr) > 0 {
					log.Printf("DeviceOnline %v\n", addr)
					r.NodeLogin(addr)
					addr = addr[0:0]
				}

				//节点发生了离线
				for _, c := range device.CollectInterfaceMap {
					for i := 0; i < len(c.OfflineReportChan); i++ {
						addr = append(addr, <-c.OfflineReportChan)
					}
				}
				if len(addr) > 0 {
					log.Printf("DeviceOffline %v\n", addr)
					r.NodeLogOut(addr)
					addr = addr[0:0]
				}

				//节点有属性变化
				for _, c := range device.CollectInterfaceMap {
					for i := 0; i < len(c.PropertyReportChan); i++ {
						addr = append(addr, <-c.PropertyReportChan)
					}
				}
				if len(addr) > 0 {
					log.Printf("DevicePropertyChanged %v\n", addr)
					r.NodePropertyPost(addr)
					addr = addr[0:0]
				}
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func ReportServiceAliyunProcessRemoteCmd(r *ReportServiceParamAliyunTemplate, message mqttClient.MQTTAliyunMessageTemplate,
	gw mqttClient.MQTTAliyunRegisterTemplate, cmdName string) {

	addrArray := strings.Split(message.Params["Addr"].(string), ",")
	for _, v := range addrArray {
		cmd := device.CommunicationCmdTemplate{}
		cmd.CollInterfaceName = "coll1"
		cmd.DeviceAddr = v
		cmd.FunName = cmdName
		paramStr, _ := json.Marshal(message.Params)
		cmd.FunPara = string(paramStr)

		if len(device.CommunicationManage) > 0 {
			if device.CommunicationManage[0].CommunicationManageAddEmergency(cmd) == true {
				payload := mqttClient.MQTTAliyunThingServiceAckTemplate{
					Identifier: cmdName,
					ID:         message.ID,
					Code:       200,
					Data:       make(map[string]interface{}),
				}
				mqttClient.MQTTAliyunThingServiceAck(r.GWParam.MQTTClient, gw, payload)
			} else {
				payload := mqttClient.MQTTAliyunThingServiceAckTemplate{
					Identifier: cmdName,
					ID:         message.ID,
					Code:       1000,
					Data:       make(map[string]interface{}),
				}
				mqttClient.MQTTAliyunThingServiceAck(r.GWParam.MQTTClient, gw, payload)
			}
		}
	}
}

func ReportServiceAliyunProcessGetSubDeviceProperty(r *ReportServiceParamAliyunTemplate, message mqttClient.MQTTAliyunMessageTemplate,
	gw mqttClient.MQTTAliyunRegisterTemplate, cmdName string) {

	addrArray := strings.Split(message.Params["Addr"].(string), ",")
	for _, v := range addrArray {
		cmd := device.CommunicationCmdTemplate{}
		cmd.CollInterfaceName = "coll1"
		cmd.DeviceAddr = v
		cmd.FunName = cmdName
		paramStr, _ := json.Marshal(message.Params)
		cmd.FunPara = string(paramStr)

		if len(device.CommunicationManage) > 0 {
			if device.CommunicationManage[0].CommunicationManageAddEmergency(cmd) == true {
				payload := mqttClient.MQTTAliyunThingServiceAckTemplate{
					Identifier: cmdName,
					ID:         message.ID,
					Code:       200,
					Data:       make(map[string]interface{}),
				}
				mqttClient.MQTTAliyunThingServiceAck(r.GWParam.MQTTClient, gw, payload)
			} else {
				payload := mqttClient.MQTTAliyunThingServiceAckTemplate{
					Identifier: cmdName,
					ID:         message.ID,
					Code:       1000,
					Data:       make(map[string]interface{}),
				}
				mqttClient.MQTTAliyunThingServiceAck(r.GWParam.MQTTClient, gw, payload)
			}
		}
	}
}

func ReportServiceAliyunProcessMessage(r *ReportServiceParamAliyunTemplate, topic string, payload []byte) {

	log.Printf("Recv TOPIC: %s\n", topic)
	log.Printf("Recv MSG: %s\n", payload)

	message := ReportServiceMessageAliyunTemplate{
		Topic:   topic,
		Payload: payload,
	}

	mqttAliyunRegister := mqttClient.MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	if strings.Contains(topic, "/thing/event/property/pack/post_reply") { //上报属性回应

		property := mqttClient.MQTTAliyunPropertyPostAckTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			log.Printf("PropertyPostAck json unmarshal err")
			return
		}
		log.Printf("code %v\n", property.Code)

		r.MessageChan <- message
	} else if strings.Contains(topic, "/combine/batch_login_reply") { //子设备上线回应

		type MQTTAliyunLogInAckTemplate struct {
			ID      string `json:"id"`
			Code    int32  `json:"code"`
			Message string `json:"message"`
			data    string `json:"data"`
		}

		property := MQTTAliyunLogInAckTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			log.Printf("LogInAck json unmarshal err")
			return
		}
		log.Printf("code %v\n", property.Code)

		r.MessageChan <- message
	} else if strings.Contains(topic, "/combine/batch_logout_reply") { //子设备下线回应

		type MQTTAliyunLogOutAckTemplate struct {
			ID      string `json:"id"`
			Code    int32  `json:"code"`
			Message string `json:"message"`
			data    string `json:"data"`
		}

		property := MQTTAliyunLogOutAckTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			log.Printf("LogOutAck json unmarshal err")
			return
		}
		log.Printf("code %v\n", property.Code)

		r.MessageChan <- message
	} else if strings.Contains(topic, "/thing/service/property/set") { //设置属性请求

		cmd := device.CommunicationCmdTemplate{}
		cmd.CollInterfaceName = "coll1"
		//cmd.DeviceAddr = property["Addr"]
		cmd.FunName = "SetRemoteCmdAdjust"
		//cmd.FunPara = string(bodyBuf[:n])

		if len(device.CommunicationManage) > 0 {
			if device.CommunicationManage[0].CommunicationManageAddEmergency(cmd) == true {

			}
		}
	} else if strings.Contains(topic, "/thing/service/SetRemoteCmdOpen") {

		property := mqttClient.MQTTAliyunMessageTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			log.Printf("processMessage json unmarshal err")
			return
		}
		log.Printf("param %v\n", property.Params)
		ReportServiceAliyunProcessRemoteCmd(r, property, mqttAliyunRegister, "SetRemoteCmdOpen")
	} else if strings.Contains(topic, "/thing/service/SetRemoteCmdClose") {

		property := mqttClient.MQTTAliyunMessageTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			log.Printf("processMessage json unmarshal err")
			return
		}
		log.Printf("param %v\n", property.Params)
		ReportServiceAliyunProcessRemoteCmd(r, property, mqttAliyunRegister, "SetRemoteCmdClose")
	} else if strings.Contains(topic, "/thing/service/SetRemoteCmdAdjust") {

		property := mqttClient.MQTTAliyunMessageTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			log.Printf("processMessage json unmarshal err")
			return
		}
		log.Printf("param %v\n", property.Params)
		ReportServiceAliyunProcessRemoteCmd(r, property, mqttAliyunRegister, "SetRemoteCmdAdjust")
	} else if strings.Contains(topic, "/thing/service/GetSubDeviceProperty") { //读取子设备的属性

		property := mqttClient.MQTTAliyunMessageTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			log.Printf("processMessage json unmarshal err")
			return
		}
		log.Printf("param %v\n", property.Params)
		ReportServiceAliyunProcessGetSubDeviceProperty(r, property, mqttAliyunRegister, "GetDeviceRealVariables")
	}
}
