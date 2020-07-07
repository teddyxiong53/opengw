package main

import (
	"goAdapter/config"
	"goAdapter/device"
	"goAdapter/httpServer"
	"goAdapter/setting"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group

)



func main() {

	//记录起始时间
	setting.GetTimeStart()

	log.Println("HRx-WDT300 V0.0.1")


	//pluginTest,_ := p.Lookup("PluginTest")
	//pluginTest.(func())()


	setting.MemoryDataStream = setting.NewDataStreamTemplate("内存使用率")
	setting.DiskDataStream = setting.NewDataStreamTemplate("硬盘使用率")
	setting.DeviceOnlineDataStream = setting.NewDataStreamTemplate("设备在线率")
	setting.DevicePacketLossDataStream = setting.NewDataStreamTemplate("通信丢包率")

	/**************获取配置文件***********************/
	config.GetConf()

	/**************串口初始化***********************/
	setting.SerialInterfaceInit()

	/**************网口初始化***********************/
	for _, v := range setting.NetworkParamList.NetworkParam {
		log.Println("set network ", v.Name)
		setting.SetNetworkParam(v.ID, v)
	}
	setting.NetworkParamList = setting.GetNetworkParam()

	/**************变量模板初始化****************/
	device.CommunicationManageInit()
	//deviceManageStart()

	device.DeviceNodeManageInit()

	/**************目标平台初始化****************/
	setting.RemotePlatformInit()

	/**************创建定时获取网络状态的任务***********************/
	// 定义一个cron运行器
	cronGetNetStatus := cron.New()
	// 定时5秒，每5秒执行print5
	cronGetNetStatus.AddFunc("*/5 * * * * *", setting.GetNetworkStatus)

	// 定时60秒
	cronGetNetStatus.AddFunc("*/30 * * * * *", device.CommunicationManagePoll)

	// 定时60秒
	//cronGetNetStatus.AddFunc("*/10 * * * * *", CommunicationManageAddEmergencyTest)

	// 定时60s
	cronGetNetStatus.AddFunc("*/60 * * * * *", func() {
		//threadModuleParam.threadModuleReadNetStatus()
	})

	// 定时60秒
	cronGetNetStatus.AddFunc("*/60 * * * * *", setting.CollectSystemParam)

	cronGetNetStatus.Start()
	defer cronGetNetStatus.Stop()

	//mqttAppConnect()

	/**************httpserver初始化****************/
	// 默认启动方式，包含 Logger、Recovery 中间件

	//serverCommon := &http.Server{
	//	Addr:         ":8091",
	//	Handler:      routerCommon(),
	//	ReadTimeout:  5 * time.Second,
	//	WriteTimeout: 10 * time.Second,
	//}

	serverWeb := &http.Server{
		Addr:         ":8090",
		Handler:      httpServer.RouterWeb(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//g.Go(func() error {
	//	return serverCommon.ListenAndServe()
	//})

	g.Go(func() error {
		return serverWeb.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
