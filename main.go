package main

import (
	"fmt"
	"goAdapter/device"
	"goAdapter/httpServer"
	"goAdapter/report"
	"goAdapter/setting"

	"github.com/robfig/cron"
)

func main() {

	/**************获取配置文件***********************/
	setting.GetConf()

	setting.LogerInit()

	//记录起始时间
	setting.GetTimeStart()

	setting.Logger.Debugf("goteway V0.0.1")

	setting.MemoryDataStream = setting.NewDataStreamTemplate("内存使用率")
	setting.DiskDataStream = setting.NewDataStreamTemplate("硬盘使用率")
	setting.DeviceOnlineDataStream = setting.NewDataStreamTemplate("设备在线率")
	setting.DevicePacketLossDataStream = setting.NewDataStreamTemplate("通信丢包率")

	/**************网口初始化***********************/
	setting.NetworkParaRead()
	setting.NetworkParamList.GetNetworkParam()

	/**************变量模板初始化****************/
	device.DeviceNodeManageInit()

	/**************NTP校时初始化****************/
	setting.NTPInit()

	/**************创建定时获取网络状态的任务***********************/
	// 定义一个cron运行器
	cronProcess := cron.New()
	// 定时5秒，每5秒执行print5
	_ = cronProcess.AddFunc("*/5 * * * * *", setting.NetworkParamList.GetNetworkParam)

	// 定时
	for k, v := range device.CollectInterfaceMap {
		device.CommunicationManage = append(device.CommunicationManage, device.NewCommunicationManageTemplate(v))
		//CommunicationManage.CollInterfaceName = v.CollInterfaceName
		str := fmt.Sprintf("@every %dm%ds", v.PollPeriod/60, v.PollPeriod%60)
		setting.Logger.Infof("str %+v", str)

		_ = cronProcess.AddFunc(str, device.CommunicationManage[k].CommunicationManagePoll)

		go device.CommunicationManage[k].CommunicationManageDel()
	}

	// 定时60秒,定时获取系统信息
	_ = cronProcess.AddFunc("*/60 * * * * *", setting.CollectSystemParam)

	// 每天0点,定时获取NTP服务器的时间，并校时
	_ = cronProcess.AddFunc("0 0 0 * * ?", func() {
		setting.NTPGetTime()
	})

	// 定时60秒,mqtt发布消息
	//cronGetNetStatus.AddFunc("*/30 * * * * *", mqttClient.MqttAppPublish)

	cronProcess.Start()
	defer cronProcess.Stop()

	report.ReportServiceInit()

	for _, v := range device.CommunicationManage {
		v.CommunicationManagePoll()
	}

	httpServer.RouterWeb()
}
