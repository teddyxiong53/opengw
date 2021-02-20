package httpServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goAdapter/setting"
	"net/http"
)

func apiAddNetwork(context *gin.Context) {
	type NetworkSetParamTemplate struct {
		Name    string `json:"Name"`
		DHCP    string `json:"DHCP"`
		IP      string `json:"IP"`
		Netmask string `json:"Netmask"`
		Gateway string `json:"Gateway"`
	}

	networkParam := NetworkSetParamTemplate{}

	err := context.ShouldBindJSON(&networkParam)
	if err != nil {
		setting.Logger.Warnf("rNetworkParam json unMarshall err,", err)
		context.JSON(http.StatusOK, Response{
			Code:    "1",
			Message: "",
			Data:    "",
		})
		return
	}

	param := setting.NetworkParamTemplate{
		Name:    networkParam.Name,
		DHCP:    networkParam.DHCP,
		IP:      networkParam.IP,
		Netmask: networkParam.DHCP,
		Gateway: networkParam.Gateway,
	}

	err = setting.NetworkParamList.AddNetworkParam(param)
	if err != nil {
		context.JSON(http.StatusOK, Response{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, Response{
			Code:    "0",
			Message: "",
			Data:    "",
		})
	}
}

func apiModifyNetwork(context *gin.Context) {
	type NetworkSetParamTemplate struct {
		Name    string `json:"Name"`
		DHCP    string `json:"DHCP"`
		IP      string `json:"IP"`
		Netmask string `json:"Netmask"`
		Gateway string `json:"Gateway"`
	}

	networkParam := NetworkSetParamTemplate{}

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

	param := setting.NetworkParamTemplate{
		Name:    networkParam.Name,
		DHCP:    networkParam.DHCP,
		IP:      networkParam.IP,
		Netmask: networkParam.DHCP,
		Gateway: networkParam.Gateway,
	}

	setting.NetworkParamList.ModifyNetworkParam(param)
	context.JSON(http.StatusOK, Response{
		Code:    "0",
		Message: "",
		Data:    "",
	})
}

func apiDeleteNetwork(context *gin.Context) {
	type NetworkSetParamTemplate struct {
		Name string `json:"Name"`
	}

	networkParam := NetworkSetParamTemplate{}

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

	status, _ := setting.NetworkParamList.DeleteNetworkParam(networkParam.Name)
	if status == true {
		context.JSON(http.StatusOK, Response{
			Code:    "0",
			Message: "",
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, Response{
			Code:    "1",
			Message: "网络名称不存在",
			Data:    "",
		})
	}

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
