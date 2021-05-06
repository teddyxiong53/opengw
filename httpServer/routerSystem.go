package httpServer

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"goAdapter/setting"

	"github.com/gin-gonic/gin"
)

func apiSystemReboot(context *gin.Context) {
	context.JSON(http.StatusOK, Response{
		Code:    "0",
		Message: "",
		Data:    "",
	})

	setting.SystemReboot()
}

func apiGetSystemStatus(context *gin.Context) {

	setting.GetMemState()
	setting.GetDiskState()
	setting.GetRunTime()
	context.JSON(http.StatusOK, ResponseData{
		"0",
		"",
		setting.SystemState,
	})
}

func apiSystemLoginParam(context *gin.Context) {
	context.JSON(http.StatusOK, ResponseData{
		"0",
		"",
		loginResult,
	})
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
				"Code":    "-1",
				"Message": "验证失败" + err.Error(),
				"Data":    "",
			})
			return
		}

	} else {
		c.JSON(http.StatusOK, gin.H{
			"Code":    "-1",
			"Message": "用户数据解析失败",
			"Data":    "",
		})
		return
	}
}

func apiSystemMemoryUseList(context *gin.Context) {
	context.JSON(http.StatusOK, ResponseData{
		"0",
		"",
		*setting.MemoryDataStream,
	})
}

func apiSystemDiskUseList(context *gin.Context) {
	context.JSON(http.StatusOK, ResponseData{
		"0",
		"",
		*setting.DiskDataStream,
	})
}

func apiSystemDeviceOnlineList(context *gin.Context) {
	context.JSON(http.StatusOK, ResponseData{
		"0",
		"",
		*setting.DeviceOnlineDataStream,
	})
}

func apiSystemDevicePacketLossList(context *gin.Context) {
	context.JSON(http.StatusOK, ResponseData{
		"0",
		"",
		*setting.DevicePacketLossDataStream,
	})
}

func apiSystemSetSystemRTC(context *gin.Context) {
	rRTC := &struct {
		SystemRTC string `json:"systemRTC"`
	}{}
	err := context.ShouldBindJSON(rRTC)
	if err != nil {
		fmt.Println("rRTC json unMarshall err,", err)
		context.JSON(http.StatusOK, Response{
			Code:    "1",
			Message: "json unMarshall err",
			Data:    "",
		})
		return
	}
	setting.Logger.Debugf("systemRTC %v", rRTC)
	setting.SystemSetRTC(rRTC.SystemRTC)
	context.JSON(http.StatusOK, Response{
		Code:    "0",
		Message: "",
		Data:    "",
	})
}

func apiSystemSetNTPHost(context *gin.Context) {

	rNTPHostAddr := setting.NTPHostAddrTemplate{}

	err := context.ShouldBindJSON(&rNTPHostAddr)
	if err != nil {
		fmt.Println("rNTPHostAddr json unMarshall err,", err)

		context.JSON(http.StatusOK, Response{
			Code:    "1",
			Message: "json unMarshall err",
			Data:    "",
		})
		return
	}

	setting.NTPHostAddr = rNTPHostAddr
	setting.WriteNTPHostAddrToJson()
	context.JSON(http.StatusOK, Response{
		Code:    "0",
		Message: "",
		Data:    "",
	})
}

func apiSystemGetNTPHost(context *gin.Context) {
	context.JSON(http.StatusOK, ResponseData{
		Code:    "0",
		Message: "",
		Data:    setting.NTPHostAddr,
	})
}

func apiBackupFiles(context *gin.Context) {

	status, name := setting.BackupFiles()
	if status == true {
		//返回文件流
		context.Writer.Header().Add("Content-Disposition",
			fmt.Sprintf("attachment;filename=%s", filepath.Base(name)))
		context.File(name) //返回文件路径，自动调用http.ServeFile方法

	} else {
		context.JSON(http.StatusOK, ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
	}
}

func apiRecoverFiles(context *gin.Context) {

	// 获取文件头
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusOK, ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
		return
	}

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileName := exeCurDir + "/selfpara/" + file.Filename

	//保存文件到服务器本地
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		context.JSON(http.StatusOK, ResponseData{
			Code:    "1",
			Message: "Save File Error",
			Data:    "",
		})

		return
	}

	//恢复
	status := setting.RecoverFiles(file.Filename)
	if status == true {
		context.JSON(http.StatusOK, ResponseData{
			Code:    "0",
			Message: "",
			Data:    "",
		})
	} else {
		context.JSON(http.StatusOK, ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
	}
}

func apiSystemUpdate(context *gin.Context) {

	// 获取文件头
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusOK, ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
		return
	}

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileName := exeCurDir + "/config/" + file.Filename

	//保存文件到服务器本地
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		context.JSON(http.StatusOK, ResponseData{
			Code:    "1",
			Message: "Save File Error",
			Data:    "",
		})
		return
	}

	//升级文件解析
	status := setting.Update(file.Filename)
	if status == true {
		context.JSON(http.StatusOK, ResponseData{
			Code:    "0",
			Message: "",
			Data:    "",
		})
		setting.SystemReboot()
	}
}
