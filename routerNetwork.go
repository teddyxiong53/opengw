package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func apiSetNetwork(context *gin.Context){

	bodyBuf := make([]byte,1024)
	n,_ := context.Request.Body.Read(bodyBuf)

	fmt.Println(string(bodyBuf[:n]))

	//获取写寄存器的参数
	rNetworkParam := &struct{
		ID   string         `json:ID`
		Name string         `json:"Name"`
		DHCP string         `json:"DHCP"`
		IP string           `json:"IP"`
		Netmask string      `json:"Netmask"`
		Broadcast string    `json:"Broadcast"`
		MAC string          `json:"MAC"`
	}{}

	err := json.Unmarshal(bodyBuf[:n],rNetworkParam)
	if err != nil {
		fmt.Println("rNetworkParam json unMarshall err,",err)

		aParam := struct{
			Code string			`json:"Code"`
			Message string		`json:"Message"`
			Data string			`json:"Data"`
		}{
			Code:"1",
			Message:"",
			Data:"",
		}
		sJson,_ := json.Marshal(aParam)

		context.String(http.StatusOK,string(sJson))
		return
	}

	if (*rNetworkParam).ID == "1"{
		networkParamList.NetworkParam[0] = *rNetworkParam
		setNetworkParam("1",*rNetworkParam)
	}else if (*rNetworkParam).ID == "2"{
		networkParamList.NetworkParam[1] = *rNetworkParam
		setNetworkParam("2",*rNetworkParam)
	}

	networkParaWrite(networkParamList)


	aParam := struct{
		Code string			`json:"Code"`
		Message string		`json:"Message"`
		Data string			`json:"Data"`
	}{
		Code:"1",
		Message:"",
		Data:"",
	}
	sJson,_ := json.Marshal(aParam)

	context.String(http.StatusOK,string(sJson))

}

func apiGetNetwork(context *gin.Context){

	aParam := struct{
		Code string
		Message string
		Data NetworkParamList
	}{Code:"0"}

	aParam.Data = getNetworkParam()

	sJson,_ := json.Marshal(aParam)

	context.String(http.StatusOK,string(sJson))
}

func apiGetNetworkLinkState(context *gin.Context){
	aParam := struct{
		Code string
		Message string
		Data NetworkLinkState
	}{Code:"0"}

	aParam.Data = networkLinkState

	sJson,_ := json.Marshal(aParam)

	context.String(http.StatusOK,string(sJson))
}
