package main

import (
	"fmt"
	//	"log"
	"net/http"
	"time"

	"github.com/robfig/cron"
	"golang.org/x/sync/errgroup"

	"goAdapter/device"
	"goAdapter/httpServer"
	"goAdapter/setting"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

var (
	g errgroup.Group
)

func logInit(){

	//log输出行号和ms
	//log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	// 设置日志格式为json格式　自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	log.SetFormatter(&log.JSONFormatter{})

	path := "/Users/opensource/test/go.log"
	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	`WithMaxAge 和 WithRotationCount二者只能设置一个
	`WithMaxAge` 设置文件清理前的最长保存时间
	`WithRotationCount` 设置文件清理前最多保存的个数
	*/
	// 下面配置日志每隔 60 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationCount(5),
		rotatelogs.WithRotationTime(time.Hour),
	)
	log.SetOutput(writer)

	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	//log.SetOutput(os.Stdout)

	// 设置日志级别为warn以上
	log.SetLevel(log.DebugLevel)
}

func main() {

	logInit()

	//记录起始时间
	setting.GetTimeStart()

	log.Info("goteway V0.0.1")

	setting.MemoryDataStream 			= setting.NewDataStreamTemplate("内存使用率")
	setting.DiskDataStream 				= setting.NewDataStreamTemplate("硬盘使用率")
	setting.DeviceOnlineDataStream 		= setting.NewDataStreamTemplate("设备在线率")
	setting.DevicePacketLossDataStream 	= setting.NewDataStreamTemplate("通信丢包率")

	/**************获取配置文件***********************/
	setting.GetConf()

	/**************网口初始化***********************/
	setting.NetworkParaRead()
	for _, v := range setting.NetworkParamList.NetworkParam {
		//log.Info("set network ", v.Name)

		//log.WithFields(log.Fields{
		//	"networkName": v.Name,
		//}).Info("set network")

		log.Infof("set network %v\n",v.Name)

		setting.SetNetworkParam(v.ID, v)
	}
	setting.NetworkParamList = setting.GetNetworkParam()

	/**************变量模板初始化****************/
	device.DeviceNodeManageInit()

	/**************目标平台初始化****************/
	setting.RemotePlatformInit()

	/**************创建定时获取网络状态的任务***********************/
	// 定义一个cron运行器
	cronGetNetStatus := cron.New()
	// 定时5秒，每5秒执行print5
	cronGetNetStatus.AddFunc("*/5 * * * * *", setting.GetNetworkStatus)

	// 定时
	for _,v := range device.CollectInterfaceMap{
		CommunicationManage := device.NewCommunicationManageTemplate()
		CommunicationManage.CollInterfaceName = v.CollInterfaceName
		str := fmt.Sprintf("@every %dm%ds",v.PollPeriod/60,v.PollPeriod%60)
		log.Infof("str %+v",str)

		//cronGetNetStatus.AddFunc("10 */1 * * * *", CommunicationManage.CommunicationManagePoll)
		cronGetNetStatus.AddFunc(str, CommunicationManage.CommunicationManagePoll)

		go CommunicationManage.CommunicationManageDel()
	}


	// 定时60秒,定时获取系统信息
	cronGetNetStatus.AddFunc("*/60 * * * * *", setting.CollectSystemParam)

	// 定时60秒,mqtt发布消息
	//cronGetNetStatus.AddFunc("*/30 * * * * *", mqttClient.MqttAppPublish)

	cronGetNetStatus.Start()
	defer cronGetNetStatus.Stop()

	//mqttClient.MQTTClient_Init()

	/**************httpserver初始化****************/
	// 默认启动方式，包含 Logger、Recovery 中间件
	serverWeb := &http.Server{
		Addr:         setting.HttpPort,
		Handler:      httpServer.RouterWeb(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return serverWeb.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
