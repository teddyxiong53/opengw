package httpServer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goAdapter/setting"
	"net/http"
	"os/exec"
)

func apiSystemReboot(context *gin.Context){

	aParam := struct{
		Code string			`json:"Code"`
		Message string		`json:"Message"`
		Data string			`json:"Data"`
	}{
		Code:"1",
		Message:"",
		Data:"",
	}

	aParam.Code = "0"

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))

	setting.SystemReboot()
}

func apiGetSystemStatus(context *gin.Context){


	setting.GetMemState()
	setting.GetDiskState()
	setting.GetRunTime()

	aParam := struct{
		Code string
		Message string
		Data setting.SystemStateTemplate
	}{"0","",setting.SystemState}

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiSystemLoginParam(context *gin.Context){


	aParam := struct{
		Code string
		Message string
		Data LoginResult `json:"Data"`
	}{"0","",loginResult}

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

// 定义登陆逻辑
// model.LoginReq中定义了登陆的请求体(name,passwd)
func apiLogin(c *gin.Context) {
	var loginReq LoginReq
	if c.BindJSON(&loginReq) == nil {
		// 登陆逻辑校验(查库，验证用户是否存在以及登陆信息是否正确)
		isPass, user, err := LoginCheck(loginReq)
		// 验证通过后为该次请求生成token
		if isPass {
			generateToken(c, user)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"Code"		: "-1",
				"Message"	: "验证失败" + err.Error(),
				"Data"		: "",
			})
			return
		}

	}else {
		c.JSON(http.StatusOK, gin.H{
			"Code"		: "-1",
			"Message"	: "用户数据解析失败",
			"Data"		: "",
		})
		return
	}
}

func apiSystemMemoryUseList(context *gin.Context){

	aParam := struct{
		Code string
		Message string
		Data setting.DataStreamTemplate
	}{"0","",*setting.MemoryDataStream}

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiSystemDiskUseList(context *gin.Context){

	aParam := struct{
		Code string
		Message string
		Data setting.DataStreamTemplate
	}{"0","",*setting.DiskDataStream}

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiSystemDeviceOnlineList(context *gin.Context){

	aParam := struct{
		Code string
		Message string
		Data setting.DataStreamTemplate
	}{"0","",*setting.DeviceOnlineDataStream}

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiSystemDevicePacketLossList(context *gin.Context){

	aParam := struct{
		Code string
		Message string
		Data setting.DataStreamTemplate
	}{"0","",*setting.DevicePacketLossDataStream}

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiSystemSetSystemRTC(context *gin.Context){

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

	rRTC := &struct{
		systemRTC  string
	}{}
	err := json.Unmarshal(bodyBuf[:n],rRTC)
	if err != nil {
		fmt.Println("rRTC json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	cmd := exec.Command("date","-s",rRTC.systemRTC)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Start()

	aParam.Code = "0"

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}