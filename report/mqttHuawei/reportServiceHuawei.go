package mqttHuawei

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

//华为云上报节点参数结构体
type ReportServiceNodeParamHuaweiTemplate struct {
	ServiceName       string
	CollInterfaceName string
	Name              string
	Addr              string
	CommStatus        string
	ReportErrCnt      int `json:"-"`
	ReportStatus      string
	Protocol          string
	Param             struct {
		DeviceID string
	}
}

//华为云上报网关参数结构体
type ReportServiceGWParamHuaweiTemplate struct {
	ServiceName  string
	IP           string
	Port         string
	ReportStatus string
	ReportTime   int
	ReportErrCnt int
	Protocol     string
	Param        struct {
		DeviceID     string
		DeviceSecret string
	}
	MQTTClient MQTT.Client `json:"-"`
}

//华为云上报服务参数，网关参数，节点参数
type ReportServiceParamHuaweiTemplate struct {
	GWParam                           ReportServiceGWParamHuaweiTemplate
	NodeList                          []ReportServiceNodeParamHuaweiTemplate
	ReceiveFrameChan                  chan MQTTHuaweiReceiveFrameTemplate      `json:"-"`
	LogInRequestFrameChan             chan []string                            `json:"-"` //上线
	ReceiveLogInAckFrameChan          chan MQTTHuaweiLogInAckTemplate          `json:"-"`
	LogOutRequestFrameChan            chan []string                            `json:"-"`
	ReceiveLogOutAckFrameChan         chan MQTTHuaweiLogOutAckTemplate         `json:"-"`
	ReportPropertyRequestFrameChan    chan MQTTHuaweiReportPropertyTemplate    `json:"-"`
	ReceiveReportPropertyAckFrameChan chan MQTTHuaweiReportPropertyAckTemplate `json:"-"`
}

type ReportServiceParamListHuaweiTemplate struct {
	ServiceList []*ReportServiceParamHuaweiTemplate
}

//实例化上报服务
var ReportServiceParamListHuawei = &ReportServiceParamListHuaweiTemplate{
	ServiceList: make([]*ReportServiceParamHuaweiTemplate, 0),
}

func init() {

	ReportServiceParamListHuawei.ReadParamFromJson()

	//初始化
	for _, v := range ReportServiceParamListHuawei.ServiceList {
		v.ReceiveFrameChan = make(chan MQTTHuaweiReceiveFrameTemplate, 100)
		v.LogInRequestFrameChan = make(chan []string, 0)
		v.ReceiveLogInAckFrameChan = make(chan MQTTHuaweiLogInAckTemplate, 5)
		v.LogOutRequestFrameChan = make(chan []string, 0)
		v.ReceiveLogOutAckFrameChan = make(chan MQTTHuaweiLogOutAckTemplate, 5)
		v.ReportPropertyRequestFrameChan = make(chan MQTTHuaweiReportPropertyTemplate, 50)
		v.ReceiveReportPropertyAckFrameChan = make(chan MQTTHuaweiReportPropertyAckTemplate, 50)

		go ReportServiceHuaweiPoll(v)
	}
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func (s *ReportServiceParamListHuaweiTemplate) ReadParamFromJson() bool {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/reportServiceParamListHuawei.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open reportServiceParamListHuawei.json err,", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt], s)
		if err != nil {
			log.Println("reportServiceParamListHuawei unmarshal err", err)
			return false
		}
		setting.Logger.Info("read reportServiceParamListHuawei.json ok")
		return true
	} else {
		setting.Logger.Warn("reportServiceParamListHuawei.json is not exist")
		return false
	}
}

func (s *ReportServiceParamListHuaweiTemplate) WriteParamToJson() {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/reportServiceParamListHuawei.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open reportServiceParamListHuawei.json err", err)
		return
	}
	defer fp.Close()

	sJson, _ := json.Marshal(*s)

	_, err = fp.Write(sJson)
	if err != nil {
		setting.Logger.Errorf("write reportServiceParamListHuawei.json err", err)
	}
	setting.Logger.Debugf("write reportServiceParamListHuawei.json success")
}

func (s *ReportServiceParamListHuaweiTemplate) AddReportService(param ReportServiceGWParamHuaweiTemplate) {

	for k, v := range s.ServiceList {
		//存在相同的，表示修改;不存在表示增加
		if v.GWParam.ServiceName == param.ServiceName {

			s.ServiceList[k].GWParam = param
			s.WriteParamToJson()
			return
		}
	}

	ReportServiceParam := &ReportServiceParamHuaweiTemplate{
		GWParam: param,
	}
	s.ServiceList = append(s.ServiceList, ReportServiceParam)

	s.WriteParamToJson()
}

func (s *ReportServiceParamListHuaweiTemplate) DeleteReportService(serviceName string) {

	for k, v := range s.ServiceList {
		if v.GWParam.ServiceName == serviceName {

			s.ServiceList = append(s.ServiceList[:k], s.ServiceList[k+1:]...)
			s.WriteParamToJson()
			return
		}
	}
}

func (r *ReportServiceParamHuaweiTemplate) AddReportNode(param ReportServiceNodeParamHuaweiTemplate) {

	param.CommStatus = "offLine"
	param.ReportStatus = "offLine"
	param.ReportErrCnt = 0

	//节点存在则进行修改
	for k, v := range r.NodeList {
		//节点已经存在
		if v.Name == param.Name {
			r.NodeList[k] = param
			ReportServiceParamListHuawei.WriteParamToJson()
			return
		}
	}

	//节点不存在则新建
	r.NodeList = append(r.NodeList, param)
	ReportServiceParamListHuawei.WriteParamToJson()

	setting.Logger.Debugf("param %v\n", ReportServiceParamListHuawei)
}

func (r *ReportServiceParamHuaweiTemplate) DeleteReportNode(name string) int {

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

func (r *ReportServiceParamHuaweiTemplate) ProcessUpLinkFrame() {

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

func (r *ReportServiceParamHuaweiTemplate) ProcessDownLinkFrame() {

	for {
		select {
		case frame := <-r.ReceiveFrameChan:
			{
				setting.Logger.Debugf("Recv TOPIC: %s\n", frame.Topic)
				setting.Logger.Debugf("Recv MSG: %s\n", frame.Payload)

				if strings.Contains(frame.Topic, "/thing/event/property/pack/post_reply") { //网关、子设备上报属性回应

					ackFrame := MQTTHuaweiReportPropertyAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.Logger.Errorf("ReportPropertyAck json unmarshal err")
						return
					}
					r.ReceiveReportPropertyAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/combine/batch_login_reply") { //子设备上线回应

					ackFrame := MQTTHuaweiLogInAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.Logger.Warningf("LogInAck json unmarshal err")
						return
					}
					r.ReceiveLogInAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/combine/batch_logout_reply") { //子设备下线回应

					ackFrame := MQTTHuaweiLogOutAckTemplate{}
					err := json.Unmarshal(frame.Payload, &ackFrame)
					if err != nil {
						setting.Logger.Errorf("LogOutAck json unmarshal err")
						return
					}
					r.ReceiveLogOutAckFrameChan <- ackFrame
				} else if strings.Contains(frame.Topic, "/thing/service/property/set") { //设置属性请求

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
		}
	}
}

func (r *ReportServiceParamHuaweiTemplate) LogIn(nodeName []string) {

	//清空接收chan，避免出现有上次接收的缓存
	for i := 0; i < len(r.ReceiveLogInAckFrameChan); i++ {
		<-r.ReceiveLogInAckFrameChan
	}

	r.NodeLogin(nodeName)
}

func (r *ReportServiceParamHuaweiTemplate) LogOut(nodeName []string) {

	//清空接收chan，避免出现有上次接收的缓存
	for i := 0; i < len(r.ReceiveLogOutAckFrameChan); i++ {
		<-r.ReceiveLogOutAckFrameChan
	}

	r.NodeLogOut(nodeName)
}

func (r *ReportServiceParamHuaweiTemplate) ReportTimeOut() {

	if r.GWParam.ReportStatus == "onLine" {
		//网关上报
		reportGWProperty := MQTTHuaweiReportPropertyTemplate{
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
			reportNodeProperty := MQTTHuaweiReportPropertyTemplate{
				DeviceType: "node",
				DeviceName: nodeName,
			}
			r.ReportPropertyRequestFrameChan <- reportNodeProperty
		}
	}
}

//查看上报服务中设备是否离线
func (r *ReportServiceParamHuaweiTemplate) ReportOfflineTime() {

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

func ReportServiceHuaweiPoll(r *ReportServiceParamHuaweiTemplate) {

	reportState := 0

	// 定义一个cron运行器
	cronProcess := cron.New()

	reportTime := fmt.Sprintf("@every %dm%ds", r.GWParam.ReportTime/60, r.GWParam.ReportTime%60)
	setting.Logger.Infof("reportServiceHuawei reportTime%v", reportTime)

	reportOfflineTime := fmt.Sprintf("@every %dm%ds", (3*r.GWParam.ReportTime)/60, (3*r.GWParam.ReportTime)%60)
	setting.Logger.Infof("reportServiceHuawei reportOfflineTime%v", reportOfflineTime)
	_ = cronProcess.AddFunc(reportOfflineTime, r.ReportOfflineTime)
	_ = cronProcess.AddFunc(reportTime, r.ReportTimeOut)

	cronProcess.Start()
	defer cronProcess.Stop()

	go r.ProcessUpLinkFrame()

	go r.ProcessDownLinkFrame()

	name := make([]string, 0)
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

				//节点发生了上线
				for _, c := range device.CollectInterfaceMap {
					for i := 0; i < len(c.OnlineReportChan); i++ {
						name = append(name, <-c.OnlineReportChan)
					}
				}
				if len(name) > 0 {
					setting.Logger.Infof("DeviceOnline %v\n", name)
					r.LogInRequestFrameChan <- name
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
					r.LogOutRequestFrameChan <- name
					name = name[0:0]
				}

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
				//					reportPropertyTemplate := MQTTHuaweiReportPropertyTemplate{
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
