/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-15 15:48:28
@LastEditors: WalkMiao
@LastEditTime: 2021-10-07 23:07:02
@FilePath: /goAdapter-Raw/main.go
*/
package main

import (
	"context"
	"goAdapter/config"
	"goAdapter/device"
	"goAdapter/httpServer"
	"goAdapter/initialize"
	"goAdapter/pkg/mylog"
	"goAdapter/pkg/ntp"
	"goAdapter/pkg/system"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"
)

func main() {

	/**************初始化配置以及日志***********************/
	initialize.Init()
	mylog.Logger.Debugf("%s %s", system.SystemState.Name, system.SystemState.SoftVer)

	/**************订阅主题************************/
	quitChan := make(chan struct{}, 1)
	device.SubScribeCollect("collect.*", quitChan)
	device.SubScribeComunication("comm.*", quitChan)
	device.SubScribeTSL("property.*", quitChan)

	/**************变量模板初始化****************/
	if err := device.NodeManageInit(); err != nil {
		mylog.ZAP.Error("初始化模板和接口失败", zap.Error(err))
		return
	}

	/**************创建定时获取网络状态的任务***********************/
	schedule := gocron.NewScheduler()
	// 定时60秒,定时获取系统信息
	schedule.Every(60).Seconds().Do(system.CollectSystemParam)
	// 每天0点,定时获取NTP服务器的时间，并校时
	schedule.Every(1).Day().At("00:00").Do(ntp.NTPGetTime)
	schedule.Every(1).Hour().Do(device.WriteAllCfg)
	// 定时60秒,mqtt发布消息
	//cronGetNetStatus.AddFunc("*/30 * * * * *", mqttClient.MqttAppPublish)

	schedule.Start()
	defer schedule.Clear()

	router := httpServer.Router()
	server := http.Server{
		Addr:    ":" + config.Cfg.ServerCfg.Port,
		Handler: router,
	}
	sigChan := make(chan os.Signal, 1)

	go func() {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
		for {
			select {
			case <-sigChan:
				if err := device.WriteAllCfg(); err != nil {
					mylog.ZAP.Error("保存配置错误", zap.Error(err))
				}
				if err := server.Shutdown(context.Background()); err != nil {
					log.Println(color.RedString("shutdown server error:%v", err))
				}
				quitChan <- struct{}{}
				return

			}
		}

	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		mylog.ZAPS.Debugf("server listen and serve error:%v", err)
		return
	}

	mylog.ZAPS.Debug("服务器正常退出....")
}
