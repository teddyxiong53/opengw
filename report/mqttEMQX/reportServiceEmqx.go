package mqttEmqx

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

//上报节点参数结构体
type ReportServiceNodeParamEmqxTemplate struct {
	ServiceName       string
	CollInterfaceName string
	Name              string
	Addr              string
	CommStatus        string
	ReportErrCnt      int `json:"-"`
	ReportStatus      string
	Protocol          string
	Param             struct {
		ClientID string
	}
}

//上报网关参数结构体
type ReportServiceGWParamEmqxTemplate struct {
	ServiceName  string
	IP           string
	Port         string
	ReportStatus string
	ReportTime   int
	ReportErrCnt int
	Protocol     string
	Param        struct {
		UserName string
		Password string
		ClientID string
	}
	MQTTClient MQTT.Client `json:"-"`
}

//上报服务参数，网关参数，节点参数
type ReportServiceParamEmqxTemplate struct {
	GWParam                              ReportServiceGWParamEmqxTemplate
	NodeList                             []ReportServiceNodeParamEmqxTemplate
	ReceiveFrameChan                     chan MQTTEmqxReceiveFrameTemplate         `json:"-"`
	LogInRequestFrameChan                chan []string                             `json:"-"` //上线
	ReceiveLogInAckFrameChan             chan MQTTEmqxLogInAckTemplate             `json:"-"`
	LogOutRequestFrameChan               chan []string                             `json:"-"`
	ReceiveLogOutAckFrameChan            chan MQTTEmqxLogOutAckTemplate            `json:"-"`
	ReportPropertyRequestFrameChan       chan MQTTEmqxReportPropertyTemplate       `json:"-"`
	ReceiveReportPropertyAckFrameChan    chan MQTTEmqxReportPropertyAckTemplate    `json:"-"`
	ReceiveInvokeServiceRequestFrameChan chan MQTTEmqxInvokeServiceRequestTemplate `json:"-"`
	ReceiveInvokeServiceAckFrameChan     chan MQTTEmqxInvokeServiceAckTemplate     `json:"-"`
	ReceiveWritePropertyRequestFrameChan chan MQTTEmqxWritePropertyRequestTemplate `json:"-"`
	ReceiveReadPropertyRequestFrameChan  chan MQTTEmqxReadPropertyRequestTemplate  `json:"-"`
}

type ReportServiceParamListEmqxTemplate struct {
	ServiceList []*ReportServiceParamEmqxTemplate
}

//实例化上报服务
var ReportServiceParamListEmqx = &ReportServiceParamListEmqxTemplate{
	ServiceList: make([]*ReportServiceParamEmqxTemplate, 0),
}

func ReportServiceEmqxInit() {

	ReportServiceParamListEmqx.ReadParamFromJson()

	//初始化
	for _, v := range ReportServiceParamListEmqx.ServiceList {
		v.ReceiveFrameChan = make(chan MQTTEmqxReceiveFrameTemplate, 100)
		v.LogInRequestFrameChan = make(chan []string, 0)
		v.ReceiveLogInAckFrameChan = make(chan MQTTEmqxLogInAckTemplate, 5)
		v.LogOutRequestFrameChan = make(chan []string, 0)
		v.ReceiveLogOutAckFrameChan = make(chan MQTTEmqxLogOutAckTemplate, 5)
		v.ReportPropertyRequestFrameChan = make(chan MQTTEmqxReportPropertyTemplate, 50)
		v.ReceiveReportPropertyAckFrameChan = make(chan MQTTEmqxReportPropertyAckTemplate, 50)
		v.ReceiveInvokeServiceRequestFrameChan = make(chan MQTTEmqxInvokeServiceRequestTemplate, 50)
		v.ReceiveInvokeServiceAckFrameChan = make(chan MQTTEmqxInvokeServiceAckTemplate, 50)
		v.ReceiveWritePropertyRequestFrameChan = make(chan MQTTEmqxWritePropertyRequestTemplate, 50)
		v.ReceiveReadPropertyRequestFrameChan = make(chan MQTTEmqxReadPropertyRequestTemplate, 50)

		go ReportServiceEmqxPoll(v)
	}
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func (s *ReportServiceParamListEmqxTemplate) ReadParamFromJson() bool {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/reportServiceParamListEmqx.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open reportServiceParamListEmqx.json err,", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt], s)
		if err != nil {
			log.Println("reportServiceParamListEmqx unmarshal err", err)
			return false
		}
		setting.Logger.Info("read reportServiceParamListEmqx.json ok")
		return true
	} else {
		setting.Logger.Warn("reportServiceParamListEmqx.json is not exist")
		return false
	}
}

func (s *ReportServiceParamListEmqxTemplate) WriteParamToJson() {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/reportServiceParamListEmqx.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open reportServiceParamListEmqx.json err", err)
		return
	}
	defer fp.Close()

	sJson, _ := json.Marshal(*s)

	_, err = fp.Write(sJson)
	if err != nil {
		setting.Logger.Errorf("write reportServiceParamListEmqx.json err", err)
	}
	setting.Logger.Debugf("write reportServiceParamListEmqx.json success")
}

func (s *ReportServiceParamListEmqxTemplate) AddReportService(param ReportServiceGWParamEmqxTemplate) {

	for k, v := range s.ServiceList {
		//存在相同的，表示修改;不存在表示增加
		if v.GWParam.ServiceName == param.ServiceName {

			s.ServiceList[k].GWParam = param
			s.WriteParamToJson()
			return
		}
	}

	ReportServiceParam := &ReportServiceParamEmqxTemplate{
		GWParam: param,
	}
	s.ServiceList = append(s.ServiceList, ReportServiceParam)

	s.WriteParamToJson()
}

func (s *ReportServiceParamListEmqxTemplate) DeleteReportService(serviceName string) {

	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == serviceName {

			s.ServiceList = append(s.ServiceList[:k], s.ServiceList[k+1:]...)
			s.WriteParamToJson()
			return
		}
	}
}

func (r *ReportServiceParamEmqxTemplate) AddReportNode(param ReportServiceNodeParamEmqxTemplate) {

	param.CommStatus = "offLine"
	param.ReportStatus = "offLine"
	param.ReportErrCnt = 0

	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == param.Name {
			r.NodeList[k] = param
			ReportServiceParamListEmqx.WriteParamToJson()
			return
		}
	}

	//节点不存在则新建
	r.NodeList = append(r.NodeList, param)
	ReportServiceParamListEmqx.WriteParamToJson()

	setting.Logger.Debugf("param %v", ReportServiceParamListEmqx)
}

func (r *ReportServiceParamEmqxTemplate) DeleteReportNode(name string) int {

	index := -1
	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == name {
			index = k
			r.NodeList = append(r.NodeList[:k], r.NodeList[k+1:]...)
			ReportServiceParamListEmqx.WriteParamToJson()
			return index
		}
	}
	return index
}

func (r *ReportServiceParamEmqxTemplate) ProcessUpLinkFrame() {

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
		case reqFrame := <-r.ReceiveWritePropertyRequestFrameChan:
			{
				r.ReportServiceEmqxProcessWriteProperty(reqFrame)
			}
		case reqFrame := <-r.ReceiveReadPropertyRequestFrameChan:
			{
				r.ReportServiceEmqxProcessReadProperty(reqFrame)
			}
		case reqFrame := <-r.ReceiveInvokeServiceRequestFrameChan:
			{
				r.ReportServiceEmqxProcessInvokeService(reqFrame)
			}
		}
	}
}

func (r *ReportServiceParamEmqxTemplate) ProcessDownLinkFrame() {

	for {
		select {
		case frame := <-r.ReceiveFrameChan:
			{
				//setting.Logger.Debugf("Recv TOPIC: %s", frame.Topic)
				//setting.Logger.Debugf("Recv MSG: %v", frame.Payload)

				if strings.Contains(frame.Topic, "/sys/thing/event/property/post_reply") { //网关、子设备上报属性回应

					ackFrame := MQTTEmqxReportPropertyAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.Logger.Errorf("ReportPropertyAck json unmarshal err")
						return
					}
					r.ReceiveReportPropertyAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/event/login/post_reply") { //子设备上线回应

					ackFrame := MQTTEmqxLogInAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.Logger.Warningf("LogInAck json unmarshal err")
						return
					}
					r.ReceiveLogInAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/event/logout/post_reply") { //子设备下线回应

					ackFrame := MQTTEmqxLogOutAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.Logger.Errorf("LogOutAck json unmarshal err")
						return
					}
					r.ReceiveLogOutAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/event/service/invoke") { //设备服务调用
					serviceFrame := MQTTEmqxInvokeServiceRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &serviceFrame)
					if err != nil {
						setting.Logger.Errorf("serviceFrame json unmarshal err")
						return
					}
					r.ReceiveInvokeServiceRequestFrameChan <- serviceFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/event/property/set") { //设置属性请求
					writePropertyFrame := MQTTEmqxWritePropertyRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &writePropertyFrame)
					if err != nil {
						setting.Logger.Errorf("writePropertyFrame json unmarshal err")
						return
					}
					r.ReceiveWritePropertyRequestFrameChan <- writePropertyFrame
				} else if strings.Contains(frame.Topic, "/sys/thing/event/property/get") { //获取属性请求
					readPropertyFrame := MQTTEmqxReadPropertyRequestTemplate{}
					err := json.Unmarshal(frame.Payload, &readPropertyFrame)
					if err != nil {
						setting.Logger.Errorf("readPropertyFrame json unmarshal err")
						return
					}
					r.ReceiveReadPropertyRequestFrameChan <- readPropertyFrame
				}

			}
		}
	}
}

func (r *ReportServiceParamEmqxTemplate) LogIn(nodeName []string) {

	//清空接收chan，避免出现有上次接收的缓存
	for i := 0; i < len(r.ReceiveLogInAckFrameChan); i++ {
		<-r.ReceiveLogInAckFrameChan
	}

	r.NodeLogIn(nodeName)
}

func (r *ReportServiceParamEmqxTemplate) LogOut(nodeName []string) {

	//清空接收chan，避免出现有上次接收的缓存
	for i := 0; i < len(r.ReceiveLogOutAckFrameChan); i++ {
		<-r.ReceiveLogOutAckFrameChan
	}

	r.NodeLogOut(nodeName)
}

//查看上报服务中设备通信状态
func (r *ReportServiceParamEmqxTemplate) ReportCommStatusTimeFun() {

	setting.Logger.Infof("service:%s CheckCommStatus", r.GWParam.ServiceName)
	for k, n := range r.NodeList {
		name := make([]string, 0)
		for _, c := range device.CollectInterfaceMap {
			if c.CollInterfaceName == n.CollInterfaceName {
				for _, d := range c.DeviceNodeMap {
					if n.Name == d.Name {
						//通信状态发生了改变
						if d.CommStatus != n.CommStatus {
							if d.CommStatus == "onLine" {
								setting.Logger.Infof("DeviceOnline %v", n.Name)
								name = append(name, n.Name)
								r.LogInRequestFrameChan <- name
							} else if d.CommStatus == "offLine" {
								setting.Logger.Infof("DeviceOffline %v", n.Name)
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

func (r *ReportServiceParamEmqxTemplate) ReportTimeFun() {

	if r.GWParam.ReportStatus == "onLine" {
		//网关上报
		reportGWProperty := MQTTEmqxReportPropertyTemplate{
			DeviceType: "gw",
		}
		r.ReportPropertyRequestFrameChan <- reportGWProperty

		//全部末端设备上报
		nodeName := make([]string, 0)
		for _, v := range r.NodeList {
			if v.CommStatus == "onLine" {
				nodeName = append(nodeName, v.Name)
			}
		}
		setting.Logger.Debugf("report Nodes %v", nodeName)
		if len(nodeName) > 0 {
			reportNodeProperty := MQTTEmqxReportPropertyTemplate{
				DeviceType: "node",
				DeviceName: nodeName,
			}
			r.ReportPropertyRequestFrameChan <- reportNodeProperty
		}
	}
}

//查看上报服务中设备是否离线
func (r *ReportServiceParamEmqxTemplate) ReportOfflineTimeFun() {

	setting.Logger.Infof("service:%s CheckReportOffline", r.GWParam.ServiceName)
	if r.GWParam.ReportErrCnt >= 3 {
		r.GWParam.ReportStatus = "offLine"
		r.GWParam.ReportErrCnt = 0
		setting.Logger.Warningf("service:%s gw offline", r.GWParam.ServiceName)
	}

	for k, v := range r.NodeList {
		if v.ReportErrCnt >= 3 {
			r.NodeList[k].ReportStatus = "offLine"
			r.NodeList[k].ReportErrCnt = 0
			setting.Logger.Warningf("service:%s %s offline", v.ServiceName, v.Name)
		}
	}
}

func ReportServiceEmqxPoll(r *ReportServiceParamEmqxTemplate) {

	reportState := 0

	// 定义一个cron运行器
	cronProcess := cron.New()

	//每10s查看一下上报节点的通信状态
	reportCommStatusTime := fmt.Sprintf("@every %dm%ds", 10/60, 10%60)
	setting.Logger.Infof("reportServiceEmqx reportCommStatusTime%v", reportCommStatusTime)

	reportTime := fmt.Sprintf("@every %dm%ds", r.GWParam.ReportTime/60, r.GWParam.ReportTime%60)
	setting.Logger.Infof("reportServiceEmqx reportTime%v", reportTime)

	reportOfflineTime := fmt.Sprintf("@every %dm%ds", (3*r.GWParam.ReportTime)/60, (3*r.GWParam.ReportTime)%60)
	setting.Logger.Infof("reportServiceEmqx reportOfflineTime%v", reportOfflineTime)

	_ = cronProcess.AddFunc(reportCommStatusTime, r.ReportCommStatusTimeFun)
	_ = cronProcess.AddFunc(reportOfflineTime, r.ReportOfflineTimeFun)
	_ = cronProcess.AddFunc(reportTime, r.ReportTimeFun)

	cronProcess.Start()
	defer cronProcess.Stop()

	go r.ProcessUpLinkFrame()

	go r.ProcessDownLinkFrame()

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
				//					reportPropertyTemplate := MQTTEmqxReportPropertyTemplate{
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
