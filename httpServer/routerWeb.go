package httpServer

import (
	"goAdapter/setting"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func RouterWeb() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//router := gin.New()

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	router.Static("/static", exeCurDir+"/webroot/static")
	router.Static("/plugin", exeCurDir+"/plugin/")

	router.StaticFile("/", exeCurDir+"/webroot/index.html")
	router.StaticFile("/favicon.ico", exeCurDir+"/webroot/favicon.ico")
	router.StaticFile("/serverConfig.json", exeCurDir+"/webroot/serverConfig.json")

	loginRouter := router.Group("/api/v1/system/")
	{
		loginRouter.POST("/login", apiLogin)
	}

	router.Use(JWTAuth())
	{
		systemRouter := router.Group("/api/v1/system")
		{
			systemRouter.POST("/reboot", apiSystemReboot)

			systemRouter.GET("/status", apiGetSystemStatus)

			systemRouter.GET("/backup", apiBackupFiles)

			systemRouter.POST("/recover", apiRecoverFiles)

			systemRouter.POST("/update", apiSystemUpdate)

			systemRouter.GET("/loginParam", apiSystemLoginParam)

			systemRouter.GET("/MemUseList", apiSystemMemoryUseList)

			systemRouter.GET("/DiskUseList", apiSystemDiskUseList)

			systemRouter.GET("/DeviceOnlineList", apiSystemDeviceOnlineList)

			systemRouter.GET("/DevicePacketLossList", apiSystemDevicePacketLossList)

			systemRouter.POST("/systemRTC", apiSystemSetSystemRTC)
		}

		ntpRouter := router.Group("/api/v1/system/ntp")
		{
			ntpRouter.POST("/hostAddr", apiSystemSetNTPHost)

			ntpRouter.GET("/hostAddr", apiSystemGetNTPHost)
		}

		networkRouter := router.Group("/api/v1/network")
		{
			networkRouter.POST("/param", apiAddNetwork)

			networkRouter.PUT("/param", apiModifyNetwork)

			networkRouter.DELETE("/param", apiDeleteNetwork)

			networkRouter.GET("/param", apiGetNetwork)

			networkRouter.GET("/linkstate", apiGetNetworkLinkState)

		}

		serialRouter := router.Group("/api/v1/serial")
		{

			serialRouter.GET("/param", apiGetSerial)
		}

		deviceRouter := router.Group("/api/v1/device")
		{
			//增加采集接口
			deviceRouter.POST("/interface", apiAddInterface)

			//修改采集接口
			deviceRouter.PUT("/interface", apiModifyInterface)

			//删除采集接口
			deviceRouter.DELETE("/interface", apiDeleteInterface)

			//获取接口信息
			deviceRouter.GET("/interface", apiGetInterfaceInfo)

			//获取所有接口信息
			deviceRouter.GET("/allInterface", apiGetAllInterfaceInfo)

			//增加节点
			deviceRouter.POST("/node", apiAddNode)

			//修改单个节点
			deviceRouter.PUT("/node", apiModifyNode)

			//修改多个节点
			deviceRouter.PUT("/nodes", apiModifyNodes)

			//查看节点
			deviceRouter.GET("/node", apiGetNode)

			//查看节点变量
			deviceRouter.GET("/nodeVariable", apiGetNodeVariableFromCache)

			//查看节点历史变量
			deviceRouter.GET("/nodeHistoryVariable", apiGetNodeHistoryVariableFromCache)

			//删除节点
			deviceRouter.DELETE("/node", apiDeleteNode)

			//增加设备模板
			deviceRouter.POST("/template", apiAddTemplate)

			//获取设备模板
			deviceRouter.GET("/template", apiGetTemplate)

			//获取通信接口
			deviceRouter.GET("/commInterface", apiGetCommInterface)

			//增加通信接口
			deviceRouter.POST("/commInterface", apiAddCommInterface)

			//修改通信接口
			deviceRouter.PUT("/commInterface", apiModifyCommInterface)

			//删除通信接口
			deviceRouter.DELETE("/commInterface", apiDeleteCommInterface)

			//增加串口通信接口
			deviceRouter.POST("/commSerialInterface", apiAddCommSerialInterface)

			//修改串口通信接口
			deviceRouter.PUT("/commSerialInterface", apiModifyCommSerialInterface)

			//删除串口通信接口
			deviceRouter.DELETE("/commSerialInterface", apiDeleteCommSerialInterface)

			//调用设备服务
			deviceRouter.POST("/service", apiInvokeService)
		}

		toolRouter := router.Group("/api/v1/tool")
		{
			//获取通信报文
			toolRouter.POST("/commMessage", apiGetCommMessage)
		}

		pluginRouter := router.Group("/api/v1/update")
		{
			pluginRouter.POST("/plugin", apiUpdatePlugin)
		}

		ReportRouter := router.Group("/api/v1/report")
		{
			ReportRouter.POST("/param", apiSetReportGWParam)

			ReportRouter.GET("/param", apiGetReportGWParam)

			ReportRouter.DELETE("/param", apiDeleteReportGWParam)

			ReportRouter.POST("/node/param", apiSetReportNodeWParam)

			ReportRouter.POST("/nodes/param", apiBatchAddReportNodeParam)

			ReportRouter.GET("/node/param", apiGetReportNodeWParam)

			ReportRouter.DELETE("/node/param", apiDeleteReportNodeWParam)
		}
	}

	if err := router.Run(":8080"); err != nil {
		setting.Logger.Errorf("gin run err,%v", err)
	}
}
