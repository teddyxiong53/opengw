package httpServer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func routerCommon() http.Handler {
	router := gin.Default()
	//router := gin.New()

	systemRouter := router.Group("/api/v1/system")
	{
		systemRouter.POST("/reboot",apiSystemReboot)

		systemRouter.GET("/status",apiGetSystemStatus)

		systemRouter.GET("/loginParam",apiSystemLoginParam)
	}

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
	}

	return router
}