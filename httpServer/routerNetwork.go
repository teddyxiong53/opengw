package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goAdapter/setting"
	"net/http"
)

func apiSetNetwork(context *gin.Context){

	bodyBuf := make([]byte,1024)
	n,_ := context.Request.Body.Read(bodyBuf)

	fmt.Println(string(bodyBuf[:n]))

	networkParam := setting.NetworkParamTemplate{}

	err := json.Unmarshal(bodyBuf[:n],&networkParam)
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

	setting.NetworkParamList.SetNetworkParam(networkParam)

	aParam := struct{
		Code string			`json:"Code"`
		Message string		`json:"Message"`
		Data string			`json:"Data"`
	}{
		Code:"0",
		Message:"",
		Data:"",
	}
	sJson,_ := json.Marshal(aParam)

	context.String(http.StatusOK,string(sJson))

}

func apiGetNetwork(context *gin.Context){

	aParam := &struct{
		Code string
		Message string
		Data setting.NetworkParamListTemplate
	}{}

	aParam.Code = "0"
	aParam.Data = *setting.NetworkParamList

	sJson,_ := json.Marshal(aParam)

	context.String(http.StatusOK,string(sJson))
}

func apiGetNetworkLinkState(context *gin.Context){
	//aParam := struct{
	//	Code string
	//	Message string
	//	Data setting.NetworkLinkStateTemplate
	//}{Code:"0"}
	//
	//aParam.Data = setting.NetworkLinkState
	//
	//sJson,_ := json.Marshal(aParam)
	//
	//context.String(http.StatusOK,string(sJson))
}
