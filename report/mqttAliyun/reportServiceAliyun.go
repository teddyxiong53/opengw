package mqttAliyun

import (
	"encoding/json"
	"fmt"
	"goAdapter/device"
	"goAdapter/setting"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/robfig/cron"
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

//阿里云上报服务参数，网关参数，节点参数
type ReportServiceParamAliyunTemplate struct {
	GWParam                             ReportServiceGWParamAliyunTemplate
	NodeList                            []ReportServiceNodeParamAliyunTemplate
	ReceiveFrameChan                    chan MQTTAliyunReceiveFrameTemplate               `json:"-"`
	LogInRequestFrameChan               chan []string                                     `json:"-"` //上线
	ReceiveLogInAckFrameChan            chan MQTTAliyunLogInAckTemplate                   `json:"-"`
	LogOutRequestFrameChan              chan []string                                     `json:"-"`
	ReceiveLogOutAckFrameChan           chan MQTTAliyunLogOutAckTemplate                  `json:"-"`
	ReportPropertyRequestFrameChan      chan MQTTAliyunReportPropertyTemplate             `json:"-"`
	ReceiveReportPropertyAckFrameChan   chan MQTTAliyunReportPropertyAckTemplate          `json:"-"`
	InvokeThingsServiceRequestFrameChan chan MQTTAliyunInvokeThingsServiceRequestTemplate `json:"-"`
	InvokeThingsServiceAckFrameChan     chan MQTTAliyunInvokeThingsServiceAckTemplate     `json:"-"`
}

type ReportServiceParamListAliyunTemplate struct {
	ServiceList []*ReportServiceParamAliyunTemplate
}

//实例化上报服务
var ReportServiceParamListAliyun = &ReportServiceParamListAliyunTemplate{
	ServiceList: make([]*ReportServiceParamAliyunTemplate, 0),
}

func init() {

	ReportServiceParamListAliyun.ReadParamFromJson()

	//初始化
	for _, v := range ReportServiceParamListAliyun.ServiceList {
		v.ReceiveFrameChan = make(chan MQTTAliyunReceiveFrameTemplate, 100)
		v.LogInRequestFrameChan = make(chan []string, 0)
		v.ReceiveLogInAckFrameChan = make(chan MQTTAliyunLogInAckTemplate, 5)
		v.LogOutRequestFrameChan = make(chan []string, 0)
		v.ReceiveLogOutAckFrameChan = make(chan MQTTAliyunLogOutAckTemplate, 5)
		v.ReportPropertyRequestFrameChan = make(chan MQTTAliyunReportPropertyTemplate, 50)
		v.ReceiveReportPropertyAckFrameChan = make(chan MQTTAliyunReportPropertyAckTemplate, 50)
		v.InvokeThingsServiceRequestFrameChan = make(chan MQTTAliyunInvokeThingsServiceRequestTemplate, 50)
		v.InvokeThingsServiceAckFrameChan = make(chan MQTTAliyunInvokeThingsServiceAckTemplate, 50)

		go ReportServiceAliyunPoll(v)
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
			r.NodeList = append(r.NodeList[:k], r.NodeList[k+1:]...)
			ReportServiceParamListAliyun.WriteParamToJson()
			return index
		}
	}
	return index
}

func (r *ReportServiceParamAliyunTemplate) ProcessUpLinkFrame() {

	for {
		select {
		case reqFrame := <-r.LogInRequestFrameChan:
			{
				r.LogIn(reqFrame)
			}
		case reqFrame := <-r.LogOutRequestFrameChan:
			{
				r.LogOut(reqFrame)
			}
		case reqFrame := <-r.ReportPropertyRequestFrameChan:
			{
				if reqFrame.DeviceType == "gw" {
					r.GWPropertyPost()
				} else if reqFrame.DeviceType == "node" {
					r.NodePropertyPost(reqFrame.DeviceName)
				}
			}
		}
	}
}

func (r *ReportServiceParamAliyunTemplate) ProcessDownLinkFrame() {

	for {
		select {
		case frame := <-r.ReceiveFrameChan:
			{
				setting.Logger.Debugf("Recv TOPIC: %s\n", frame.Topic)
				setting.Logger.Debugf("Recv MSG: %s\n", frame.Payload)

				if strings.Contains(frame.Topic, "/thing/event/property/pack/post_reply") { //网关、子设备上报属性回应

					ackFrame := MQTTAliyunReportPropertyAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.Logger.Errorf("ReportPropertyAck json unmarshal err")
						return
					}
					r.ReceiveReportPropertyAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/combine/batch_login_reply") { //子设备上线回应

					ackFrame := MQTTAliyunLogInAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.Logger.Warningf("LogInAck json unmarshal err")
						return
					}
					r.ReceiveLogInAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/combine/batch_logout_reply") { //子设备下线回应

					ackFrame := MQTTAliyunLogOutAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.Logger.Errorf("LogOutAck json unmarshal err")
						return
					}
					r.ReceiveLogOutAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/thing/service") { //设备服务调用
					serviceFrame := MQTTAliyunInvokeThingsServiceRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &serviceFrame)
					if err != nil {
						setting.Logger.Errorf("serviceFrame json unmarshal err")
						return
					}
					r.InvokeThingsServiceRequestFrameChan <- serviceFrame
				} else if strings.Contains(frame.Topic, "/thing/service/property/set") { //设置属性请求

				}
			}
		}
	}
}

func (r *ReportServiceParamAliyunTemplate) LogIn(nodeName []string) {

	//清空接收chan，避免出现有上次接收的缓存
	for i := 0; i < len(r.ReceiveLogInAckFrameChan); i++ {
		<-r.ReceiveLogInAckFrameChan
	}

	r.NodeLogin(nodeName)
}

func (r *ReportServiceParamAliyunTemplate) LogOut(nodeName []string) {

	//清空接收chan，避免出现有上次接收的缓存
	for i := 0; i < len(r.ReceiveLogOutAckFrameChan); i++ {
		<-r.ReceiveLogOutAckFrameChan
	}

	r.NodeLogOut(nodeName)
}

//查看上报服务中设备通信状态
func (r *ReportServiceParamAliyunTemplate) ReportCommStatusTimeFun() {

	setting.Logger.Infof("service:%s,CheckCommStatus", r.GWParam.ServiceName)
	for k, n := range r.NodeList {
		name := make([]string, 0)
		for _, c := range device.CollectInterfaceMap {
			if c.CollInterfaceName == n.CollInterfaceName {
				for _, d := range c.DeviceNodeMap {
					if n.Name == d.Name {
						//通信状态发生了改变
						if d.CommStatus != n.CommStatus {
							if d.CommStatus == "onLine" {
								setting.Logger.Infof("DeviceOnline %v\n", n.Name)
								name = append(name, n.Name)
								r.LogInRequestFrameChan <- name
							} else if d.CommStatus == "offLine" {
								setting.Logger.Infof("DeviceOffline %v\n", n.Name)
								name = append(name, n.Name)
								r.LogOutRequestFrameChan <- name
							}
							r.NodeList[k].CommStatus = d.CommStatus
						}
					}
				}
			}
		}
	}
}

func (r *ReportServiceParamAliyunTemplate) ReportTimeFun() {

	if r.GWParam.ReportStatus == "onLine" {
		//网关上报
		reportGWProperty := MQTTAliyunReportPropertyTemplate{
			DeviceType: "gw",
		}
		r.ReportPropertyRequestFrameChan <- reportGWProperty

		//全部末端设备上报
		nodeName := make([]string, 0)
		for _, v := range r.NodeList {
			nodeName = append(nodeName, v.Name)
		}
		setting.Logger.Debugf("report Nodes %v", nodeName)
		if len(nodeName) > 0 {
			reportNodeProperty := MQTTAliyunReportPropertyTemplate{
				DeviceType: "node",
				DeviceName: nodeName,
			}
			r.ReportPropertyRequestFrameChan <- reportNodeProperty
		}
	}
}

//查看上报服务中设备是否离线
func (r *ReportServiceParamAliyunTemplate) ReportOfflineTimeFun() {

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

	//每10s查看一下上报节点的通信状态
	reportCommStatusTime := fmt.Sprintf("@every %dm%ds", 10/60, 10%60)
	setting.Logger.Infof("reportServiceAliyun reportCommStatusTime%v", reportCommStatusTime)

	reportTime := fmt.Sprintf("@every %dm%ds", r.GWParam.ReportTime/60, r.GWParam.ReportTime%60)
	setting.Logger.Infof("reportServiceAliyun reportTime%v", reportTime)

	reportOfflineTime := fmt.Sprintf("@every %dm%ds", (3*r.GWParam.ReportTime)/60, (3*r.GWParam.ReportTime)%60)
	setting.Logger.Infof("reportServiceAliyun reportOfflineTime%v", reportOfflineTime)

	_ = cronProcess.AddFunc(reportCommStatusTime, r.ReportCommStatusTimeFun)
	_ = cronProcess.AddFunc(reportOfflineTime, r.ReportOfflineTimeFun)
	_ = cronProcess.AddFunc(reportTime, r.ReportTimeFun)

	cronProcess.Start()
	defer cronProcess.Stop()

	go r.ProcessUpLinkFrame()

	go r.ProcessDownLinkFrame()

	go r.ProcessInvokeThingsService()

	//name := make([]string, 0)
	for {
		switch reportState {
		case 0:
			{
				if r.GWLogin() == true {
					reportState = 1

				} else {
					time.Sleep(5 * time.Second)
				}
			}
		case 1:
			{
				//网关
				if r.GWParam.ReportStatus == "offLine" {
					reportState = 0
					r.GWParam.ReportErrCnt = 0
				}

				//for k, v := range r.NodeList {
				//	commStatus := "offLine"
				//	for _, d := range device.CollectInterfaceMap {
				//		if v.CollInterfaceName == v.CollInterfaceName {
				//			for _, n := range d.DeviceNodeMap {
				//				if v.Name == n.Name {
				//					commStatus = n.CommStatus
				//					break
				//				}
				//			}
				//		}
				//	}
				//	if commStatus != v.CommStatus {
				//		if commStatus == "onLine" {
				//			//节点发生了上线
				//			setting.Logger.Debugf("service %s,node %s onLine", v.ServiceName, v.Name)
				//			name = append(name, v.Name)
				//			r.LogInRequestFrameChan <- name
				//			name = name[0:0]
				//		} else if commStatus == "offLine" {
				//			//节点发生了离线
				//			setting.Logger.Debugf("service %s,node %s onLine", v.ServiceName, v.Name)
				//			name = append(name, v.Name)
				//			r.LogOutRequestFrameChan <- name
				//			name = name[0:0]
				//		}
				//		r.NodeList[k].CommStatus = commStatus
				//	}
				//}

				//节点有属性变化
				//for _, c := range device.CollectInterfaceMap {
				//	for i := 0; i < len(c.PropertyReportChan); i++ {
				//		nodeName := <-c.PropertyReportChan
				//		for _, v := range r.NodeList {
				//			if v.Name == nodeName {
				//				if v.ReportStatus == "offLine" { //当设备上报状态是离线时立马发送设备上线
				//					name = append(name, nodeName)
				//					r.LogInRequestFrameChan <- name
				//					name = name[0:0]
				//				} else {
				//					name = append(name, nodeName)
				//					reportPropertyTemplate := MQTTAliyunReportPropertyTemplate{
				//						DeviceType: "node",
				//						DeviceName: name,
				//					}
				//					r.ReportPropertyRequestFrameChan <- reportPropertyTemplate
				//					name = name[0:0]
				//				}
				//			}
				//		}
				//	}
				//}
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
