package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goAdapter/config"
	"goAdapter/setting"
	"net/http"
)

func apiSetSerial(context *gin.Context){

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
	rSerialParam := &setting.SerialParamTemplate{
	}
	err := json.Unmarshal(bodyBuf[:n],rSerialParam)
	if err != nil {
		fmt.Println("rSerialParam json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson,_ := json.Marshal(aParam)

		context.String(http.StatusOK,string(sJson))
		return
	}

	for k,v := range setting.SerialInterface.SerialParam{

		if v.ID == rSerialParam.ID{
			setting.SerialInterface.SerialParam[k] = *rSerialParam
			config.SerialParaWrite()
		}else{
			aParam.Code = "1"
			aParam.Message = "serial ID is not exist"

			sJson,_ := json.Marshal(aParam)
			context.String(http.StatusOK,string(sJson))
			return
		}
	}

	aParam.Code = "0"
	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiGetSerial(context *gin.Context){

	aParam := struct{
		Code string				`json:"Code"`
		Message string			`json:"Message"`
		Data setting.SerialInterfaceTemplate	`json:"Data"`
	}{Code:"0"}

	aParam.Data = setting.SerialInterface

	sJson,_ := json.Marshal(aParam)

	context.String(http.StatusOK,string(sJson))
}