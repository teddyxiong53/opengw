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
	CommStatus string
	GWParam    ReportServiceGWParamAliyunTemplate
	NodeList   []ReportServiceNodeParamAliyunTemplate
}

type ReportServiceParamListAliyunTemplate struct {
	ServiceList []*ReportServiceParamAliyunTemplate
}

var ReportServiceParamListAliyun = &ReportServiceParamListAliyunTemplate{
	ServiceList: make([]*ReportServiceParamAliyunTemplate, 0),
}

func init() {

	ReportServiceParamListAliyun.ReadParamFromJson()

	for _, v := range ReportServiceParamListAliyun.ServiceList {

		go ReportServiceAliyunPoll(v)
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
	//log.Printf("TOPIC: %s\n", msg.Topic())
	//log.Printf("MSG: %s\n", msg.Payload())

	for _, v := range ReportServiceParamListAliyun.ServiceList {
		//log.Printf("GW %v\n", v.GWParam)
		if v.GWParam.MQTTClient == client {
			//message := ReportServiceAliyunMessageTemplate{
			//	Topic:   msg.Topic(),
			//	Payload: msg.Payload(),
			//}
			//v.GWParam.MessageChan <- message
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

	_, r.GWParam.MQTTClient = mqttClient.MQTTAliyunGWLogin(mqttAliyunRegister, GWPublishHandler)

	return true
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

	return true
}

func (r *ReportServiceParamAliyunTemplate) GWPropertyPost() {

	valueMap := make([]mqttClient.MQTTAliyunValueTemplate, 0)

	for _, n := range r.NodeList {
		for _, c := range device.CollectInterfaceMap {
			if c.CollInterfaceName == n.CollInterfaceName {
				for _, d := range c.DeviceNodeMap {
					////无线主站中TD200的地址为0
					//if d.Addr == "0" {
					//	for _, v := range d.VariableMap {
					//		if v.Name == "Chan" {
					//			if len(v.Value) > 1 {
					//				index := len(v.Value) - 1
					//				mqttAliyunValue := mqttClient.MQTTAliyunValueTemplate{}
					//				mqttAliyunValue.Name = v.Name
					//				mqttAliyunValue.Value = v.Value[index].Value
					//				valueMap = append(valueMap, mqttAliyunValue)
					//			}
					//		} else if v.Name == "SystemID" {
					//			if len(v.Value) > 1 {
					//				index := len(v.Value) - 1
					//				mqttAliyunValue := mqttClient.MQTTAliyunValueTemplate{}
					//				mqttAliyunValue.Name = v.Name
					//				mqttAliyunValue.Value = v.Value[index].Value
					//				valueMap = append(valueMap, mqttAliyunValue)
					//			}
					//		}
					//	}
					//}
					if d.Addr == "0" {
						mqttAliyunValue := mqttClient.MQTTAliyunValueTemplate{}
						mqttAliyunValue.Name = "Chan"
						mqttAliyunValue.Value = 12
						valueMap = append(valueMap, mqttAliyunValue)
					}
				}
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

	mqttClient.MQTTAliyunGWPropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, valueMap)
	//}
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

	//if len(valueMap) > 0 {
	mqttAliyunRegister := mqttClient.MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	mqttClient.MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap)
	//}
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
	//}
}

func ReportServiceAliyunPoll(r *ReportServiceParamAliyunTemplate) {

	// 定义一个cron运行器
	cronProcess := cron.New()

	str := fmt.Sprintf("@every %dm%ds", r.GWParam.ReportTime/60, r.GWParam.ReportTime%60)
	setting.Logger.Infof("reportServiceAliyun %+v", str)

	//cronProcess.AddFunc(str, r.GWPropertyPost)
	//cronProcess.AddFunc(str, r.AllNodePropertyPost)

	cronProcess.Start()
	defer cronProcess.Stop()

	addr := make([]string, 0)

	r.GWLogin()
	//r.NodeLogin()

	for {

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

		time.Sleep(100 * time.Millisecond)
	}
}

func ReportServiceAliyunProcessMessage(r *ReportServiceParamAliyunTemplate, topic string, payload []byte) {

	log.Printf("TOPIC: %s\n", topic)
	log.Printf("MSG: %s\n", payload)

	//属性设置
	if strings.Contains(topic, "/thing/service/property/set") {

		property := mqttClient.MQTTAliyunPropertySetTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			log.Printf("/thing/service/property/set json unmarshal err")
			return
		}
		log.Printf("param %v\n", property.Params)

		splitTopic := strings.Split(topic, "/")
		log.Printf("msg %v\n", splitTopic)
		if len(splitTopic) > 2 {
			deviceName := splitTopic[2]
			//判断网关
			if r.GWParam.Param.DeviceName == deviceName {

			} else {
				//判断设备
				for _, v := range r.NodeList {
					if v.Param.DeviceName == deviceName {
						cmd := device.CommunicationCmdTemplate{}
						cmd.CollInterfaceName = "coll1"
						cmd.DeviceAddr = v.Addr
						cmd.FunName = "SetRemoteCmdAdjust"
						//cmd.FunPara = string(bodyBuf[:n])

						if len(device.CommunicationManage) > 0 {
							if device.CommunicationManage[0].CommunicationManageAddEmergency(cmd) == true {

							}
						}
					}
				}
			}
		}
	}
}
