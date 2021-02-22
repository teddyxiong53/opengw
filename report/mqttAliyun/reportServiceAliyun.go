package mqttAliyun

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/robfig/cron"
	"goAdapter/device"
	"goAdapter/setting"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

//阿里云上报节点参数结构体
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

//阿里云上报网关参数结构体
type ReportServiceGWParamAliyunTemplate struct {
	ServiceName  string
	IP           string
	Port         string
	ReportStatus string
	ReportTime   int
	ReportErrCnt int
	Protocol     string
	Param        struct {
		ProductKey   string
		DeviceName   string
		DeviceSecret string
	}
	MQTTClient MQTT.Client `json:"-"`
}

//阿里云上报服务接收报文结构体
type ReportServiceReceiveMessageAliyunTemplate struct {
	Topic string
	Code  int32
	ID    string
}

//阿里云上报服务发送报文结构体
type ReportServiceSendMessageAliyunTemplate struct {
	Topic      string
	DeviceName []string //发送时会存在多个设备在一条报文里
	ID         string
}

type ReportServicePropertyPostAliyunTemplate struct {
	DeviceType int
	DeviceName []string
}

//阿里云上报服务参数，网关参数，节点参数
type ReportServiceParamAliyunTemplate struct {
	GWParam           ReportServiceGWParamAliyunTemplate
	NodeList          []ReportServiceNodeParamAliyunTemplate
	ReceiveMessageMap []ReportServiceReceiveMessageAliyunTemplate  `json:"-"`
	SendMessageMap    []ReportServiceSendMessageAliyunTemplate     `json:"-"`
	PropertyPostChan  chan ReportServicePropertyPostAliyunTemplate `json:"-"`
}

type ReportServiceParamListAliyunTemplate struct {
	ServiceList []*ReportServiceParamAliyunTemplate
}

//实例化上报服务
var ReportServiceParamListAliyun = &ReportServiceParamListAliyunTemplate{
	ServiceList: make([]*ReportServiceParamAliyunTemplate, 0),
}

var lock sync.Mutex

func init() {

	ReportServiceParamListAliyun.ReadParamFromJson()

	//初始化
	for _, v := range ReportServiceParamListAliyun.ServiceList {
		v.ReceiveMessageMap = make([]ReportServiceReceiveMessageAliyunTemplate, 0)
		v.SendMessageMap = make([]ReportServiceSendMessageAliyunTemplate, 0)
		v.PropertyPostChan = make(chan ReportServicePropertyPostAliyunTemplate, 10)

		go ReportServiceAliyunPoll(v)
		go ProcessPropertyPost(v)
	}
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func (s *ReportServiceParamListAliyunTemplate) ReadParamFromJson() bool {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/reportServiceParamListAliyun.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open reportServiceParamListAliyun.json err,", err)
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
		setting.Logger.Info("read reportServiceParamListAliyun.json ok")
		return true
	} else {
		setting.Logger.Warn("reportServiceParamListAliyun.json is not exist")
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
		setting.Logger.Errorf("write reportServiceParamListAliyun.json err", err)
	}
	setting.Logger.Debugf("write reportServiceParamListAliyun.json success")
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
		GWParam: param,
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
	param.ReportErrCnt = 0

	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == param.Name {
			r.NodeList[k] = param
			ReportServiceParamListAliyun.WriteParamToJson()
			return
		}
	}

	//节点不存在则新建
	r.NodeList = append(r.NodeList, param)
	ReportServiceParamListAliyun.WriteParamToJson()

	setting.Logger.Debugf("param %v\n", ReportServiceParamListAliyun)
}

func (r *ReportServiceParamAliyunTemplate) DeleteReportNode(name string) int {

	index := -1
	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == name {
			index = k
			return index
		}
	}

	return index
}

//发送数据回调函数
func GWPublishHandler(client MQTT.Client, msg MQTT.Message) {

	for _, v := range ReportServiceParamListAliyun.ServiceList {
		if v.GWParam.MQTTClient == client {
			go ReportServiceAliyunProcessMessage(v, msg.Topic(), msg.Payload())
		}
	}
}

func (r *ReportServiceParamAliyunTemplate) GWLogin() bool {

	mqttAliyunRegister := MQTTAliyunRegisterTemplate{
		RemoteIP:     r.GWParam.IP,
		RemotePort:   r.GWParam.Port,
		ProductKey:   r.GWParam.Param.ProductKey,
		DeviceName:   r.GWParam.Param.DeviceName,
		DeviceSecret: r.GWParam.Param.DeviceSecret,
	}

	status := false
	status, r.GWParam.MQTTClient = MQTTAliyunGWLogin(mqttAliyunRegister, GWPublishHandler)
	if status == true {
		r.GWParam.ReportStatus = "onLine"
	}

	return status
}

func (r *ReportServiceParamAliyunTemplate) NodeLogin(name []string) bool {

	nodeList := make([]MQTTAliyunNodeRegisterTemplate, 0)
	nodeParam := MQTTAliyunNodeRegisterTemplate{}
	status := false

	setting.Logger.Debugf("nodeLoginName %v", name)
	for _, d := range name {
		for k, v := range r.NodeList {
			if d == v.Name {
				nodeParam.DeviceSecret = v.Param.DeviceSecret
				nodeParam.DeviceName = v.Param.DeviceName
				nodeParam.ProductKey = v.Param.ProductKey
				nodeList = append(nodeList, nodeParam)
				r.NodeList[k].CommStatus = "onLine"

				mqttAliyunRegister := MQTTAliyunRegisterTemplate{
					RemoteIP:     r.GWParam.IP,
					RemotePort:   r.GWParam.Port,
					ProductKey:   r.GWParam.Param.ProductKey,
					DeviceName:   r.GWParam.Param.DeviceName,
					DeviceSecret: r.GWParam.Param.DeviceSecret,
				}
				MsgId := MQTTAliyunNodeLoginIn(r.GWParam.MQTTClient, mqttAliyunRegister, nodeList) - 1
				MsgIdStr := strconv.Itoa(MsgId)
				sendMessage := ReportServiceSendMessageAliyunTemplate{
					ID: MsgIdStr,
				}
				sendMessage.DeviceName = append(sendMessage.DeviceName, v.Param.DeviceName)
				r.SendMessageMap = append(r.SendMessageMap, sendMessage)
				setting.Logger.Debugf("service:%s,sendMessageMapAdd %v", r.GWParam.ServiceName, r.SendMessageMap)
				//超时3s
				time.AfterFunc(5*time.Second, func() {
					for i, s := range r.SendMessageMap {
						if s.ID == MsgIdStr {
							r.SendMessageMap = append(r.SendMessageMap[:i], r.SendMessageMap[i+1:]...)
						}
					}
				})

				status = true
			}
		}
	}

	return status
}

func (r *ReportServiceParamAliyunTemplate) NodeLogOut(name []string) bool {

	nodeList := make([]MQTTAliyunNodeRegisterTemplate, 0)
	nodeParam := MQTTAliyunNodeRegisterTemplate{}

	for _, d := range name {
		for k, v := range r.NodeList {
			if d == v.Name {
				if v.ReportStatus == "offLine" {
					setting.Logger.Infof("service:%s,%s is already offLine", r.GWParam.ServiceName, v.Name)
				} else {
					nodeParam.DeviceSecret = v.Param.DeviceSecret
					nodeParam.DeviceName = v.Param.DeviceName
					nodeParam.ProductKey = v.Param.ProductKey

					nodeList = append(nodeList, nodeParam)
					r.NodeList[k].CommStatus = "offLine"

					mqttAliyunRegister := MQTTAliyunRegisterTemplate{
						RemoteIP:     r.GWParam.IP,
						RemotePort:   r.GWParam.Port,
						ProductKey:   r.GWParam.Param.ProductKey,
						DeviceName:   r.GWParam.Param.DeviceName,
						DeviceSecret: r.GWParam.Param.DeviceSecret,
					}
					MsgId := MQTTAliyunNodeLoginOut(r.GWParam.MQTTClient, mqttAliyunRegister, nodeList) - 1
					MsgIdStr := strconv.Itoa(MsgId)

					sendMessage := ReportServiceSendMessageAliyunTemplate{
						ID: MsgIdStr,
					}
					for _, v := range nodeList {
						sendMessage.DeviceName = append(sendMessage.DeviceName, v.DeviceName)
					}
					r.SendMessageMap = append(r.SendMessageMap, sendMessage)
				}

			}
		}
	}
	return true
}

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

	MsgId := MQTTAliyunGWPropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, valueMap) - 1
	MsgIdStr := strconv.Itoa(MsgId)

	sendMessage := ReportServiceSendMessageAliyunTemplate{
		ID: MsgIdStr,
	}
	sendMessage.DeviceName = append(sendMessage.DeviceName, r.GWParam.Param.DeviceName)
	r.SendMessageMap = append(r.SendMessageMap, sendMessage)
	setting.Logger.Debugf("service:%s,sendMessageMapAdd %v", r.GWParam.ServiceName, r.SendMessageMap)
	//上报故障先加，收到正确回应后清0
	r.GWParam.ReportErrCnt++
	setting.Logger.Debugf("service %s,gw ReportErrCnt %d", r.GWParam.Param.DeviceName, r.GWParam.ReportErrCnt)

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

			MsgId := MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap) - 1
			MsgIdStr := strconv.Itoa(MsgId)

			sendMessage := ReportServiceSendMessageAliyunTemplate{
				ID: MsgIdStr,
			}
			for _, v := range NodeValueMap {
				sendMessage.DeviceName = append(sendMessage.DeviceName, v.DeviceName)
			}
			r.SendMessageMap = append(r.SendMessageMap, sendMessage)
			setting.Logger.Debugf("service:%s,sendMessageMapAdd %v", r.GWParam.ServiceName, r.SendMessageMap)
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
			MsgId := MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap) - 1
			MsgIdStr := strconv.Itoa(MsgId)

			sendMessage := ReportServiceSendMessageAliyunTemplate{
				ID: MsgIdStr,
			}
			for _, v := range NodeValueMap {
				sendMessage.DeviceName = append(sendMessage.DeviceName, v.DeviceName)
			}
			r.SendMessageMap = append(r.SendMessageMap, sendMessage)
			setting.Logger.Debugf("service:%s,sendMessageMapAdd %v", r.GWParam.ServiceName, r.SendMessageMap)
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

			MsgId := MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap) - 1
			MsgIdStr := strconv.Itoa(MsgId)

			sendMessage := ReportServiceSendMessageAliyunTemplate{
				ID: MsgIdStr,
			}
			for _, v := range NodeValueMap {
				sendMessage.DeviceName = append(sendMessage.DeviceName, v.DeviceName)
			}
			r.SendMessageMap = append(r.SendMessageMap, sendMessage)
			setting.Logger.Debugf("service:%s,sendMessageMapAdd %v", r.GWParam.ServiceName, r.SendMessageMap)
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
			MsgId := MQTTAliyunNodePropertyPost(r.GWParam.MQTTClient, mqttAliyunRegister, NodeValueMap) - 1
			MsgIdStr := strconv.Itoa(MsgId)

			sendMessage := ReportServiceSendMessageAliyunTemplate{
				ID: MsgIdStr,
			}
			for _, v := range NodeValueMap {
				sendMessage.DeviceName = append(sendMessage.DeviceName, v.DeviceName)
			}
			r.SendMessageMap = append(r.SendMessageMap, sendMessage)
			setting.Logger.Debugf("service:%s,sendMessageMapAdd %v", r.GWParam.ServiceName, r.SendMessageMap)
		}
	}
}

func (r *ReportServiceParamAliyunTemplate) PropertyPost() {
	//网关上报
	nameMap := make([]string, 0)
	nameMap = append(nameMap, r.GWParam.Param.DeviceName)
	propertyPost := ReportServicePropertyPostAliyunTemplate{
		DeviceName: nameMap,
		DeviceType: 0,
	}
	r.PropertyPostChan <- propertyPost

	//末端设备上报
	nameMap = nameMap[0:0] //清空slice
	for _, v := range r.NodeList {
		nameMap = append(nameMap, v.Name)
	}
	if len(nameMap) > 0 {
		propertyPost = ReportServicePropertyPostAliyunTemplate{
			DeviceName: nameMap,
			DeviceType: 1,
		}
		r.PropertyPostChan <- propertyPost
	}
}

func ProcessPropertyPost(r *ReportServiceParamAliyunTemplate) {

	for {
		select {
		case postParam := <-r.PropertyPostChan:
			{
				setting.Logger.Tracef("service %s,postParam %v,postChanCnt %v", r.GWParam.ServiceName, postParam, len(r.PropertyPostChan))
				if postParam.DeviceType == 0 { //网关上报
					r.GWPropertyPost()
				} else if postParam.DeviceType == 1 { //末端设备上报
					r.NodePropertyPost(postParam.DeviceName)
				}
			}
		}
	}
}

func ReportServiceAliyunProcessGetSubDeviceProperty(r *ReportServiceParamAliyunTemplate, message MQTTAliyunMessageTemplate,
	gw MQTTAliyunRegisterTemplate, cmdName string) {

	addrArray := strings.Split(message.Params["Addr"].(string), ",")
	for _, v := range addrArray {
		for _, n := range r.NodeList {
			if v == n.Param.DeviceName {
				cmd := device.CommunicationCmdTemplate{}
				cmd.CollInterfaceName = "coll1"
				cmd.DeviceName = n.Addr
				cmd.FunName = "GetRealVariables"
				paramStr, _ := json.Marshal(message.Params)
				cmd.FunPara = string(paramStr)

				if len(device.CommunicationManage) > 0 {
					if device.CommunicationManage[0].CommunicationManageAddEmergency(cmd) == true {
						payload := MQTTAliyunThingServiceAckTemplate{
							Identifier: cmdName,
							ID:         message.ID,
							Code:       200,
							Data:       make(map[string]interface{}),
						}
						MQTTAliyunThingServiceAck(r.GWParam.MQTTClient, gw, payload)
					} else {
						payload := MQTTAliyunThingServiceAckTemplate{
							Identifier: cmdName,
							ID:         message.ID,
							Code:       1000,
							Data:       make(map[string]interface{}),
						}
						MQTTAliyunThingServiceAck(r.GWParam.MQTTClient, gw, payload)
					}
				}
			}
		}
	}
}

func ReportServiceAliyunProcessMessage(r *ReportServiceParamAliyunTemplate, topic string, payload []byte) {

	setting.Logger.Debugf("Recv TOPIC: %s\n", topic)
	setting.Logger.Debugf("Recv MSG: %s\n", payload)

	type ReportServiceAliyunMessageTemplate struct {
		Topic   string
		Payload []byte
	}

	message := ReportServiceAliyunMessageTemplate{
		Topic:   topic,
		Payload: payload,
	}

	if strings.Contains(topic, "/thing/event/property/pack/post_reply") { //上报属性回应
		type MQTTAliyunPropertyPostAckTemplate struct {
			Code    int32  `json:"code"`
			Data    string `json:"-"`
			ID      string `json:"id"`
			Message string `json:"message"`
			Method  string `json:"method"`
			Version string `json:"version"`
		}
		property := MQTTAliyunPropertyPostAckTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			setting.Logger.Errorf("PropertyPostAck json unmarshal err")
			return
		}
		setting.Logger.Debugf("code %v\n", property.Code)
		if property.Code == 200 {
			ackMessage := ReportServiceReceiveMessageAliyunTemplate{
				Topic: message.Topic,
				Code:  property.Code,
				ID:    property.ID,
			}
			lock.Lock()
			for k, v := range r.SendMessageMap {
				if v.ID == ackMessage.ID {
					for _, name := range v.DeviceName {
						if name == r.GWParam.Param.DeviceName { //网关设备
							r.GWParam.ReportStatus = "onLine"
							r.GWParam.ReportErrCnt = 0
							setting.Logger.Infof("service:%s,gw online", r.GWParam.ServiceName)
						} else { //末端设备
							for i, n := range r.NodeList {
								if name == n.Param.DeviceName {
									r.NodeList[i].ReportStatus = "onLine"
									r.NodeList[i].ReportErrCnt = 0
									setting.Logger.Infof("service:%s,%s online", r.GWParam.ServiceName, n.Param.DeviceName)
								}
							}
						}
					}
					setting.Logger.Debugf("service:%s,sendMessageMapPre %v", r.GWParam.ServiceName, r.SendMessageMap)
					setting.Logger.Debugf("k:%v,id:%v", k, v.ID)
					r.SendMessageMap = append(r.SendMessageMap[:k], r.SendMessageMap[k+1:]...)
					setting.Logger.Debugf("service:%s,sendMessageMapNow %v", r.GWParam.ServiceName, r.SendMessageMap)
				}
			}
			lock.Unlock()
		}
	} else if strings.Contains(topic, "/combine/batch_login_reply") { //子设备上线回应
		type MQTTAliyunLogInDataTemplate struct {
			ProductKey string `json:"productKey"`
			DeviceName string `json:"deviceName"`
		}

		type MQTTAliyunLogInAckTemplate struct {
			ID      string                        `json:"id"`
			Code    int32                         `json:"code"`
			Message string                        `json:"message"`
			Data    []MQTTAliyunLogInDataTemplate `json:"data"`
		}

		property := MQTTAliyunLogInAckTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			setting.Logger.Warningf("LogInAck json unmarshal err")
			return
		}
		setting.Logger.Infof("code %v\n", property.Code)
		if property.Code == 200 {
			ackMessage := ReportServiceReceiveMessageAliyunTemplate{
				Topic: message.Topic,
				Code:  property.Code,
				ID:    property.ID,
			}
			for k, v := range r.SendMessageMap {
				if v.ID == ackMessage.ID {
					for _, name := range v.DeviceName {
						for i, n := range r.NodeList {
							if name == n.Param.DeviceName {
								r.NodeList[i].ReportStatus = "onLine"
								r.NodeList[i].ReportErrCnt = 0

							}
						}
					}
					r.SendMessageMap = append(r.SendMessageMap[:k], r.SendMessageMap[k+1:]...)
					setting.Logger.Debugf("service:%s,sendMessageMapNow %v", r.GWParam.ServiceName, r.SendMessageMap)
				}
			}
		}
	} else if strings.Contains(topic, "/combine/batch_logout_reply") { //子设备下线回应
		type MQTTAliyunLogOutDataTemplate struct {
			Code       int32  `json:"code"`
			Message    string `json:"message"`
			ProductKey string `json:"productKey"`
			DeviceName string `json:"deviceName"`
		}

		type MQTTAliyunLogOutAckTemplate struct {
			ID      string                       `json:"id"`
			Code    int32                        `json:"code"`
			Message string                       `json:"message"`
			Data    MQTTAliyunLogOutDataTemplate `json:"data"`
		}

		property := MQTTAliyunLogOutAckTemplate{}
		err := json.Unmarshal(payload, &property)
		if err != nil {
			setting.Logger.Errorf("LogOutAck json unmarshal err")
			return
		}
		setting.Logger.Debugf("code %v\n", property.Code)
		if property.Code == 200 {
			//ackMessage := ReportServiceReceiveMessageAliyunTemplate{
			//	Topic: message.Topic,
			//	Code:  property.Code,
			//	ID:    property.ID,
			//}
			//r.ReceiveMessageChan <- ackMessage
		}
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
	}
}

//查看上报服务中设备是否离线
func (r *ReportServiceParamAliyunTemplate) CheckReportServiceOffline() {

	setting.Logger.Infof("service:%s,CheckReportOffline", r.GWParam.ServiceName)
	if r.GWParam.ReportErrCnt >= 3 {
		r.GWParam.ReportStatus = "offLine"
		r.GWParam.ReportErrCnt = 0
		setting.Logger.Warningf("service:%s,gw offline", r.GWParam.ServiceName)
	}

	for k, v := range r.NodeList {
		if v.ReportErrCnt >= 3 {
			r.NodeList[k].ReportStatus = "offLine"
			r.NodeList[k].ReportErrCnt = 0
			setting.Logger.Warningf("service:%s,%s offline", v.ServiceName, v.Name)
		}
	}
}

func ReportServiceAliyunPoll(r *ReportServiceParamAliyunTemplate) {

	reportState := 0

	// 定义一个cron运行器
	cronProcess := cron.New()

	reportTime := fmt.Sprintf("@every %dm%ds", r.GWParam.ReportTime/60, r.GWParam.ReportTime%60)
	setting.Logger.Infof("reportServiceAliyun reportTime%v", reportTime)

	reportOfflineTime := fmt.Sprintf("@every %dm%ds", (3*r.GWParam.ReportTime)/60, (3*r.GWParam.ReportTime)%60)
	setting.Logger.Infof("reportServiceAliyun reportOfflineTime%v", reportOfflineTime)
	_ = cronProcess.AddFunc(reportOfflineTime, r.CheckReportServiceOffline)

	cronProcess.Start()
	defer cronProcess.Stop()

	name := make([]string, 0)

	for {
		switch reportState {
		case 0:
			{
				if r.GWLogin() == true {
					reportState = 1

					_ = cronProcess.AddFunc(reportTime, r.PropertyPost)
					//_ = cronProcess.AddFunc(reportTime, r.GWPropertyPost)
					//_ = cronProcess.AddFunc(reportTime, r.AllNodePropertyPost)
				} else {
					time.Sleep(5 * time.Second)
				}
			}
		case 1:
			{
				//网关
				if r.GWParam.ReportStatus == "offLine" {
					reportState = 0
				}

				//节点发生了上线
				for _, c := range device.CollectInterfaceMap {
					for i := 0; i < len(c.OnlineReportChan); i++ {
						name = append(name, <-c.OnlineReportChan)
					}
				}
				if len(name) > 0 {
					setting.Logger.Infof("DeviceOnline %v\n", name)
					r.NodeLogin(name)
					name = name[0:0]
				}

				//节点发生了离线
				for _, c := range device.CollectInterfaceMap {
					for i := 0; i < len(c.OfflineReportChan); i++ {
						name = append(name, <-c.OfflineReportChan)
					}
				}
				if len(name) > 0 {
					setting.Logger.Infof("DeviceOffline %v\n", name)
					r.NodeLogOut(name)
					name = name[0:0]
				}

				//节点有属性变化
				for _, c := range device.CollectInterfaceMap {
					for i := 0; i < len(c.PropertyReportChan); i++ {
						nodeName := <-c.PropertyReportChan
						for _, v := range r.NodeList {
							if v.Name == nodeName {
								if v.ReportStatus == "offLine" { //当设备上报状态是离线时立马发送设备上线
									name = append(name, nodeName)
									go r.NodeLogin(name)
									name = name[0:0]
								}
							}
						}
					}
				}
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
