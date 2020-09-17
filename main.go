package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/robfig/cron"
	"golang.org/x/sync/errgroup"

	"goAdapter/device"
	"goAdapter/httpServer"
	"goAdapter/setting"
)

var (
	g errgroup.Group
)



func main() {

	/**************获取配置文件***********************/
	setting.GetConf()

	setting.LogerInit(setting.LogLevel,setting.LogSaveToFile,setting.LogFileMaxCnt)

	//记录起始时间
	setting.GetTimeStart()

	setting.Logger.Info("goteway V0.0.1")

	setting.MemoryDataStream 			= setting.NewDataStreamTemplate("内存使用率")
	setting.DiskDataStream 				= setting.NewDataStreamTemplate("硬盘使用率")
	setting.DeviceOnlineDataStream 		= setting.NewDataStreamTemplate("设备在线率")
	setting.DevicePacketLossDataStream 	= setting.NewDataStreamTemplate("通信丢包率")

	/**************网口初始化***********************/
	setting.NetworkParaRead()
	for _, v := range setting.NetworkParamList.NetworkParam {
		if setting.SetNetworkParam(v.ID, v) == true{
			setting.Logger.WithFields(logrus.Fields{
				"networkName": v.Name,
				"status": "true",
			}).Info("set network")
		}else {
			setting.Logger.WithFields(logrus.Fields{
				"networkName": v.Name,
				"status": "false",
			}).Info("set network")
		}
	}
	setting.NetworkParamList = setting.GetNetworkParam()

	/**************变量模板初始化****************/
	device.DeviceNodeManageInit()

	/**************目标平台初始化****************/
	setting.RemotePlatformInit()
	/**************NTP校时初始化****************/
	setting.NTPInit()

	/**************创建定时获取网络状态的任务***********************/
	// 定义一个cron运行器
	cronProcess := cron.New()
	// 定时5秒，每5秒执行print5
	cronProcess.AddFunc("*/5 * * * * *", setting.GetNetworkStatus)

	// 定时
	for _,v := range device.CollectInterfaceMap{
		CommunicationManage := device.NewCommunicationManageTemplate(v)
		//CommunicationManage.CollInterfaceName = v.CollInterfaceName
		str := fmt.Sprintf("@every %dm%ds",v.PollPeriod/60,v.PollPeriod%60)
		setting.Logger.Infof("str %+v",str)

		//cronGetNetStatus.AddFunc("10 */1 * * * *", CommunicationManage.CommunicationManagePoll)
		cronProcess.AddFunc(str, CommunicationManage.CommunicationManagePoll)

		go CommunicationManage.CommunicationManageDel()
	}

	// 定时60秒,定时获取系统信息
	cronProcess.AddFunc("*/60 * * * * *", setting.CollectSystemParam)

	// 每天0点,定时获取NTP服务器的时间，并校时
	cronProcess.AddFunc("0 0 0 * * ?", func(){
		setting.NTPGetTime()
	})

	// 定时60秒,mqtt发布消息
	//cronGetNetStatus.AddFunc("*/30 * * * * *", mqttClient.MqttAppPublish)

	cronProcess.Start()
	defer cronProcess.Stop()

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
		setting.Logger.Fatal(err)
	}
}
