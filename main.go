package main

import (
	"github.com/robfig/cron"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"time"
)


var (
	g errgroup.Group
)

func main(){

	//记录起始时间
	setTimeStart()

	log.Println("HRx-WDT300 V0.0.1")


	MemoryDataStream = newDataStreamTemplate("内存使用率")
	DiskDataStream   = newDataStreamTemplate("硬盘使用率")
	DeviceOnlineDataStream = newDataStreamTemplate("设备在线率")
	DevicePacketLossDataStream = newDataStreamTemplate("通信丢包率")

	/**************获取配置文件***********************/
	getConf()

	/**************串口初始化***********************/
	SerialInterfaceInit()

	/**************网口初始化***********************/
	for _,v := range networkParamList.NetworkParam{
		log.Println("set network ",v.Name)
		setNetworkParam(v.ID,v)
	}
	networkParamList = getNetworkParam()

	/**************变量模板初始化****************/
	CommunicationManageInit()
	//deviceManageStart()

	DeviceNodeManageInit()

	/**************目标平台初始化****************/
	RemotePlatformInit()

	/**************创建定时获取网络状态的任务***********************/
	// 定义一个cron运行器
	cronGetNetStatus := cron.New()
	// 定时5秒，每5秒执行print5
	cronGetNetStatus.AddFunc("*/5 * * * * *", getNetworkStatus)

	// 定时60秒
	cronGetNetStatus.AddFunc("*/60 * * * * *", CommunicationManagePoll)

	// 定时60秒
	//cronGetNetStatus.AddFunc("*/10 * * * * *", CommunicationManageAddEmergencyTest)

	// 定时60s
	cronGetNetStatus.AddFunc("*/60 * * * * *", func (){
		//threadModuleParam.threadModuleReadNetStatus()
	})

	// 定时60秒
	cronGetNetStatus.AddFunc("*/60 * * * * *",CollectSystemParam)

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
		Handler:      routerWeb(),
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
