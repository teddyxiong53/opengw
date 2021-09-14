package httpServer

import (
	"goAdapter/httpServer/controller"
	"goAdapter/httpServer/middleware"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	router.Static("/static", exeCurDir+"/webroot/static")
	router.Static("/plugin", exeCurDir+"/plugin/")
	router.Static("/layui", exeCurDir+"/webroot/layui")

	router.StaticFile("/", exeCurDir+"/webroot/index.html")
	router.StaticFile("/favicon.ico", exeCurDir+"/webroot/favicon.ico")
	router.StaticFile("/serverConfig.json", exeCurDir+"/webroot/serverConfig.json")

	loginRouter := router.Group("/api/v1/system/")
	{
		loginRouter.POST("/login", controller.Login)
	}

	router.Use(middleware.JWTAuth())
	{
		systemRouter := router.Group("/api/v1/system")
		{
			systemRouter.POST("/reboot", controller.SystemReboot)

			systemRouter.GET("/status", controller.GetSystemStatus)

			systemRouter.GET("/backup", controller.BackupFiles)

			systemRouter.POST("/recover", controller.RecoverFiles)

			systemRouter.POST("/update", controller.SystemUpdate)

			systemRouter.GET("/loginParam", controller.SystemLoginParam)

			systemRouter.GET("/MemUseList", controller.SystemMemoryUseList)

			systemRouter.GET("/DiskUseList", controller.SystemDiskUseList)

			systemRouter.GET("/DeviceOnlineList", controller.SystemDeviceOnlineList)

			systemRouter.GET("/DevicePacketLossList", controller.SystemDevicePacketLossList)

			systemRouter.POST("/systemRTC", controller.SystemSetSystemRTC)
		}

		logRouter := router.Group("/api/v1/log")
		{
			logRouter.GET("/param", controller.GetLogParam)

			logRouter.POST("/param", controller.SetLogParam)

			logRouter.GET("/filesInfo", controller.GetLogFilesInfo)

			logRouter.DELETE("/files", controller.DeleteLogFile)

			logRouter.GET("/file", controller.GetLogFile)
		}

		ntpRouter := router.Group("/api/v1/system/ntp")
		{
			ntpRouter.POST("/hostAddr", controller.SystemSetNTPHost)

			ntpRouter.GET("/hostAddr", controller.SystemGetNTPHost)
		}

		networkRouter := router.Group("/api/v1/network")
		{
			networkRouter.POST("/param", controller.AddNetwork)

			networkRouter.PUT("/param", controller.ModifyNetwork)

			networkRouter.DELETE("/param", controller.DeleteNetwork)

			networkRouter.GET("/param", controller.GetNetwork)

			networkRouter.GET("/linkstate", controller.GetNetworkLinkState)

		}

		serialRouter := router.Group("/api/v1/serial")
		{

			serialRouter.GET("/param", controller.GetSerial)
		}

		deviceRouter := router.Group("/api/v1/device")
		{
			//获取所有接口信息
			deviceRouter.GET("/allInterface", controller.GetAllInterfaceInfo)
			//采集接口
			collInterfaceGroup := deviceRouter.Group("/interface")
			{
				//增加采集接口
				collInterfaceGroup.POST("", controller.AddInterface)

				//修改采集接口
				collInterfaceGroup.PUT("", controller.ModifyInterface)

				//删除采集接口
				collInterfaceGroup.DELETE("", controller.DeleteInterface)

				//获取接口信息
				collInterfaceGroup.GET("", controller.GetInterfaceInfo)

			}

			// 节点
			nodeGroup := deviceRouter.Group("/node")
			{ //增加节点
				nodeGroup.POST("", controller.AddNode)
				//修改单个节点
				nodeGroup.PUT("", controller.ModifyNode)
				//查看节点
				nodeGroup.GET("", controller.GetNode)
				//删除节点
				nodeGroup.DELETE("", controller.DeleteNode)

			}

			//查看节点变量
			deviceRouter.GET("/nodeVariable", controller.GetNodeVariableFromCache)

			//查看节点历史变量
			deviceRouter.GET("/nodeHistoryVariable", controller.GetNodeHistoryVariable)

			//查看节点变量实时值
			deviceRouter.GET("/nodeRealVariable", controller.GetNodeReadVariable)

			//模板
			tmpGroup := deviceRouter.Group("/template")
			{
				//增加设备模板
				tmpGroup.POST("", controller.AddTemplate)

				//获取设备模板
				tmpGroup.GET("", controller.GetTemplate)
			}

			commInterfaceGroup := deviceRouter.Group("/commInterface")
			{
				//获取通信接口
				commInterfaceGroup.GET("", controller.GetCommInterface)

				//增加通信接口
				commInterfaceGroup.POST("", controller.AddCommInterface)

				//修改通信接口
				commInterfaceGroup.PUT("", controller.ModifyCommInterface)

				//删除通信接口
				commInterfaceGroup.DELETE("", controller.DeleteCommInterface)
			}

			//调用设备服务
			//	deviceRouter.POST("/service", apiInvokeService)
		}

		toolRouter := router.Group("/api/v1/tool")
		{
			//获取通信报文
			toolRouter.POST("/commMessage", controller.GetCommMessage)
		}

		pluginRouter := router.Group("/api/v1/update")
		{
			pluginRouter.POST("/plugin", controller.UpdatePlugin)
		}

		ReportRouter := router.Group("/api/v1/report")
		{
			ReportRouter.POST("/param", controller.SetReportGWParam)

			ReportRouter.GET("/param", controller.GetReportGWParam)

			ReportRouter.DELETE("/param", controller.DeleteReportGWParam)

			ReportRouter.POST("/node/param", controller.SetReportNodeWParam)

			ReportRouter.POST("/nodes/param", controller.BatchAddReportNodeParam)

			ReportRouter.GET("/node/param", controller.GetReportNodeWParam)

			ReportRouter.DELETE("/node/param", controller.DeleteReportNodeWParam)
		}
	}

	return router
}
