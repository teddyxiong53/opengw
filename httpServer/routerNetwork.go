package httpServer

import (
	"fmt"
	"net/http"

	"goAdapter/setting"

	"github.com/gin-gonic/gin"
)

func apiSetNetwork(context *gin.Context) {
	networkParam := setting.NetworkParamTemplate{}

	err := context.ShouldBindJSON(&networkParam)
	if err != nil {
		fmt.Println("rNetworkParam json unMarshall err,", err)
		context.JSON(http.StatusOK, Response{
			Code:    "1",
			Message: "",
			Data:    "",
		})
		return
	}

	setting.NetworkParamList.SetNetworkParam(networkParam)
	context.JSON(http.StatusOK, Response{
		Code:    "0",
		Message: "",
		Data:    "",
	})
}

func apiGetNetwork(context *gin.Context) {
	context.JSON(http.StatusOK, ResponseData{
		Code: "0",
		Data: *setting.NetworkParamList,
	})
}

func apiGetNetworkLinkState(context *gin.Context) {
	// context.JSON(http.StatusOK, ResponseData{
	// 	Code: "0",
	// 	Data: setting.NetworkLinkState,
	// })
}
