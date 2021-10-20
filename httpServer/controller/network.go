package controller

import (
	"fmt"
	"goAdapter/httpServer/model"
	"goAdapter/pkg/mylog"
	"goAdapter/pkg/network"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddNetwork(context *gin.Context) {
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
		mylog.Logger.Warnf("rNetworkParam json unMarshall err,", err)
		context.JSON(http.StatusOK, model.Response{
			Code: "1",
		})
		return
	}

	param := network.NetworkParamTemplate{
		Name:    networkParam.Name,
		DHCP:    networkParam.DHCP,
		IP:      networkParam.IP,
		Netmask: networkParam.Netmask,
		Gateway: networkParam.Gateway,
	}

	err = network.NetworkParamList.AddNetworkParam(param)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
	} else {
		context.JSON(http.StatusOK, model.Response{
			Code: "0",
		})
	}
}

func ModifyNetwork(context *gin.Context) {
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
		context.JSON(http.StatusOK, model.Response{
			Code: "1",
		})
		return
	}

	param := network.NetworkParamTemplate{
		Name:    networkParam.Name,
		DHCP:    networkParam.DHCP,
		IP:      networkParam.IP,
		Netmask: networkParam.DHCP,
		Gateway: networkParam.Gateway,
	}

	err = network.NetworkParamList.ModifyNetworkParam(param)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, model.Response{
		Code:    "0",
		Message: "",
	})
}

func DeleteNetwork(context *gin.Context) {
	type NetworkSetParamTemplate struct {
		Name string `json:"Name"`
	}

	networkParam := NetworkSetParamTemplate{}

	err := context.ShouldBindJSON(&networkParam)
	if err != nil {
		fmt.Println("rNetworkParam json unMarshall err,", err)
		context.JSON(http.StatusOK, model.Response{
			Code: "1",
		})
		return
	}

	status, _ := network.NetworkParamList.DeleteNetworkParam(networkParam.Name)
	if status {
		context.JSON(http.StatusOK, model.Response{
			Code: "0",
		})
	} else {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "网络名称不存在",
		})
	}

}

func GetNetwork(context *gin.Context) {
	data, err := network.ParseNetworks()
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("get networks error:%v", err),
		})
		return
	}
	var r = make([]*network.NetworkParamTemplate, 0)
	for _, i := range data {
		var matched bool
		for _, j := range network.NetworkParamList.NetworkParam {
			if i.Name == j.Name {
				//以json文件为主
				r = append(r, j)
				matched = true
				if err := j.CmdSetStaticIP(); err != nil {
					mylog.ZAPS.Errorf("set network card %s error:%v", j.Name, err)
				}
				break
			}
		}
		if !matched {
			r = append(r, i)
		}
	}

	list := network.NetworkParamListTemplate{
		NetworkParam: r,
	}
	context.JSON(http.StatusOK, model.Response{
		Code: "0",
		Data: list,
	})
}

func GetNetworkLinkState(context *gin.Context) {
	// context.JSON(http.StatusOK, ResponseData{
	// 	Code: "0",
	// 	Data: setting.NetworkLinkState,
	// })
}
