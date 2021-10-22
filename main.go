/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-15 15:48:28
@LastEditors: WalkMiao
@LastEditTime: 2021-10-19 12:56:49
@FilePath: /goAdapter-Raw/main.go
*/
package main

import (
	"embed"
	"github.com/fvbock/endless"
	"goAdapter/config"
	"goAdapter/device"
	"goAdapter/httpServer"
	"goAdapter/initialize"
	"goAdapter/pkg/mylog"
	"goAdapter/pkg/ntp"
	"goAdapter/pkg/system"
	"log"
	"net/http"
	"syscall"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"
)

var (
	//go:embed webroot/static
	static embed.FS
	//go:embed webroot/layui
	layui embed.FS
	//go:embed webroot/serialHelper
	serialHelper embed.FS
	//go:embed webroot/index.html
	indexTpl embed.FS
	//go:embed webroot/serverConfig.json
	serverConfig []byte
)

func embedOpts()httpServer.Option{
	return httpServer.WithEmbedFS(
		httpServer.EmbedFS{
		Fs:     static,
		FsPath: "/static",
		SubPath: "webroot/static",

	},
	httpServer.EmbedFS{
			Fs: layui,
			FsPath: "/layui",
			SubPath: "webroot/layui",
	},
	httpServer.EmbedFS{
			Fs: serialHelper,
			FsPath: "/serialHelper",
			SubPath: "webroot/serialHelper",
	},
	httpServer.EmbedFS{
			Fs: indexTpl,
			FsPath: "/",
			SubPath: "webroot/index.html",
			Type: httpServer.HTMLType,
	},
	httpServer.EmbedFS{
			Data: serverConfig,
			Type:    httpServer.FileType,
			FsPath: "/serverConfig.json",
	},
	)
}
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

	router := httpServer.RouterWithOpts(httpServer.WithMode(gin.ReleaseMode), embedOpts())
	addr:= ":" + config.Cfg.ServerCfg.Port
	startServer(addr,router,quitChan)

}

func startServer(addr string,h http.Handler,quitChan chan struct{}){
	server := endless.NewServer(addr, h)
	defer func() {
		if err:=device.WriteAllCfg();err!=nil{
			mylog.ZAP.Error("服务关闭时,保存配置错误", zap.Error(err))
		}
		quitChan<- struct{}{}
	}()
	server.BeforeBegin = func(add string) {
		mylog.ZAPS.Debugf("current pid is %d", syscall.Getpid())
	}

	if err := server.ListenAndServe();err!=nil {
		log.Println(color.RedString("Server ListenAndServe err: %v", err))
	}
	mylog.ZAPS.Debugf("服务器(PID=%d)正常退出....",syscall.Getpid())
}