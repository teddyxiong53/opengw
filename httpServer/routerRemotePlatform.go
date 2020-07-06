package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goAdapter/setting"
	"net/http"
)

func apiGetRemotePlatformParam(context *gin.Context){

	aParam := struct{
		Code string
		Message string
		Data setting.RemotePlatformTemplate
	}{"0","",*setting.RemotePlatform}

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiSetHTTPParam(context *gin.Context){

	aParam := struct{
		Code string			`json:"Code"`
		Message string		`json:"Message"`
		Data string			`json:"Data"`
	}{
		Code:"1",
		Message:"",
		Data:"",
	}

	bodyBuf := make([]byte,1024)
	n,_ := context.Request.Body.Read(bodyBuf)

	fmt.Println(string(bodyBuf[:n]))

	//获取写寄存器的参数
	rHttpParam := &setting.HttpRemoteTemplate{}
	err := json.Unmarshal(bodyBuf[:n],rHttpParam)
	if err != nil {
		fmt.Println("rHttpParam json unMarshall err,",err)

		aParam.Message = "json unMarshall err"
		sJson,_ := json.Marshal(aParam)

		context.String(http.StatusOK,string(sJson))
		return
	}



	aParam.Code = "0"
	sJson,_ := json.Marshal(aParam)

	context.String(http.StatusOK,string(sJson))
}