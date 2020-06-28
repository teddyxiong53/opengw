package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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
	rSerialParam := &SerialParamTemplate{
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

	for k,v := range serialInterface.SerialParam{

		if v.ID == rSerialParam.ID{
			serialInterface.SerialParam[k] = *rSerialParam
			serialParaWrite()
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
		Data SerialInterface	`json:"Data"`
	}{Code:"0"}

	aParam.Data = serialInterface

	sJson,_ := json.Marshal(aParam)

	context.String(http.StatusOK,string(sJson))
}