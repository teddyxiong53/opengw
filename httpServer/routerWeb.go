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
		loginRouter.POST("/login",apiLogin)
	}

	router.GET("/", func(context *gin.Context) {
		context.File("webroot/index.html")
	})

	router.Static("/static","webroot/static")

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
			systemRouter.POST("/reboot",apiSystemReboot)

			systemRouter.GET("/status",apiGetSystemStatus)

			systemRouter.GET("/loginParam",apiSystemLoginParam)

			systemRouter.GET("/MemUseList",apiSystemMemoryUseList)

			systemRouter.GET("/DiskUseList",apiSystemDiskUseList)

			systemRouter.GET("/DeviceOnlineList",apiSystemDeviceOnlineList)

			systemRouter.GET("/DevicePacketLossList",apiSystemDevicePacketLossList)

			systemRouter.POST("/systemRTC",apiSystemSetSystemRTC)
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
		//serialRouter := router.Group("/api/v1/serial")
		{
			serialRouter.POST("",apiSetSerial)

			serialRouter.GET("",apiGetSerial)
		}

		deviceRouter := router.Group("/api/v1/device")
		{
			//增加接口
			deviceRouter.POST("/interface",apiAddInterface)

			//修改接口
			deviceRouter.PUT("/interface",apiModifyInterface)

			//获取接口信息
			deviceRouter.GET("/interface",apiGetInterfaceInfo)

			//获取所有接口信息
			deviceRouter.GET("/allInterface",apiGetAllInterfaceInfo)

			//增加节点
			deviceRouter.POST("/node",apiAddNode)

			//修改节点
			deviceRouter.PUT("/node",apiModifyNode)

			//查看节点
			deviceRouter.GET("/node",apiGetNode)

			//删除节点
			deviceRouter.DELETE("/node",apiDeleteNode)

			//增加设备模板
			deviceRouter.POST("/template",apiAddTemplate)

			//获取设备模板
			deviceRouter.GET("/template",apiGetTemplate)
		}

		remoteRouter := router.Group("/api/v1/remote")
		{
			remoteRouter.GET("/param",apiGetRemotePlatformParam)
		}
	}

	return router
}


