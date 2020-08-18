package main

import (
	"fmt"
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

	//log输出行号和ms
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	//记录起始时间
	setting.GetTimeStart()

	log.Println("goteway V0.0.1")

	setting.MemoryDataStream 			= setting.NewDataStreamTemplate("内存使用率")
	setting.DiskDataStream 				= setting.NewDataStreamTemplate("硬盘使用率")
	setting.DeviceOnlineDataStream 		= setting.NewDataStreamTemplate("设备在线率")
	setting.DevicePacketLossDataStream 	= setting.NewDataStreamTemplate("通信丢包率")

	/**************获取配置文件***********************/
	setting.GetConf()

	/**************网口初始化***********************/
	setting.NetworkParaRead()
	for _, v := range setting.NetworkParamList.NetworkParam {
		log.Println("set network ", v.Name)
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
		log.Printf("str %+v\n",str)

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
