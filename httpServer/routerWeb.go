package httpServer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RouterWeb() http.Handler {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//router := gin.New()

	loginRouter := router.Group("/api/v1/system/")
	{
		loginRouter.POST("/login", apiLogin)
	}

	router.GET("/", func(context *gin.Context) {
		context.File("webroot/index.html")
	})

	router.Static("/static", "webroot/static")

	router.GET("/favicon.ico", func(context *gin.Context) {
		context.File("webroot/favicon.ico")
	})

	router.GET("/serverConfig.json", func(context *gin.Context) {
		context.File("webroot/serverConfig.json")
	})

	router.Use(JWTAuth())
	{
		systemRouter := router.Group("/api/v1/system")
		{
			systemRouter.POST("/reboot", apiSystemReboot)

			systemRouter.GET("/status", apiGetSystemStatus)

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

		//userRouter := router.Group("/api/v1/modbus")
		//{
		//	userRouter.POST("/cmd03",apiReadHoldReg)
		//
		//	userRouter.POST("/cmd10",apiWriteMultiReg)
		//}

		networkRouter := router.Group("/api/v1/network")
		{
			networkRouter.POST("/param", apiSetNetwork)

			networkRouter.GET("/param", apiGetNetwork)

			networkRouter.GET("/linkstate", apiGetNetworkLinkState)
		}

		networkDHCPRouter := router.Group("/api/v1/network/dhcp")
		{
			networkDHCPRouter.POST("", apiSetNetwork)

			networkDHCPRouter.GET("", apiGetNetwork)
		}

		serialRouter := router.Group("/api/v1/serial/param")
		{

			serialRouter.GET("", apiGetSerial)
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

			//增加串口通信接口
			deviceRouter.POST("/commSerialInterface", apiAddCommSerialInterface)

			//修改串口通信接口
			deviceRouter.PUT("/commSerialInterface", apiModifyCommSerialInterface)

			//删除串口通信接口
			deviceRouter.DELETE("/commSerialInterface", apiDeleteCommSerialInterface)
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

		remoteRouter := router.Group("/api/v1/remote")
		{
			remoteRouter.GET("/param", apiGetRemotePlatformParam)
		}
	}

	return router
}
