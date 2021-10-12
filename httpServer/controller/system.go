package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"goAdapter/device"
	"goAdapter/httpServer/middleware"
	"goAdapter/httpServer/model"
	"goAdapter/pkg/backup"
	"goAdapter/pkg/mylog"
	"goAdapter/pkg/ntp"
	"goAdapter/pkg/system"

	"github.com/gin-gonic/gin"
)

func SystemReboot(context *gin.Context) {
	context.JSON(http.StatusOK, model.Response{
		Code: "0",
	})

	system.SystemReboot()
}

func GetSystemStatus(context *gin.Context) {

	system.GetMemState()
	system.GetDiskState()
	system.GetRunTime()
	context.JSON(http.StatusOK, model.Response{
		Code: "0",
		Data: system.SystemState,
	})
}

func SystemLoginParam(context *gin.Context) {
	context.JSON(http.StatusOK, model.Response{
		Code: "0",
		Data: middleware.Result,
	})
}

// 定义登陆逻辑
// model.LoginReq中定义了登陆的请求体(name,passwd)
func Login(c *gin.Context) {
	var loginReq middleware.LoginReq
	if c.BindJSON(&loginReq) == nil {
		// 登陆逻辑校验(查库，验证用户是否存在以及登陆信息是否正确)
		isPass, user, err := middleware.LoginCheck(loginReq)
		// 验证通过后为该次请求生成token
		if isPass {
			middleware.GenerateToken(c, user)
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
		})
		return
	}
}

func SystemMemoryUseList(context *gin.Context) {
	context.JSON(http.StatusOK, model.Response{
		Code: "0",
		Data: *system.MemoryDataStream,
	})
}

func SystemDiskUseList(context *gin.Context) {
	context.JSON(http.StatusOK, model.Response{
		Code: "0",
		Data: *system.DiskDataStream,
	})
}

func SystemDeviceOnlineList(context *gin.Context) {
	context.JSON(http.StatusOK, model.Response{
		Code: "0",
		Data: *system.DeviceOnlineDataStream,
	})
}

func SystemDevicePacketLossList(context *gin.Context) {
	context.JSON(http.StatusOK, model.Response{
		Code: "0",
		Data: *system.DevicePacketLossDataStream,
	})
}

func SystemSetSystemRTC(context *gin.Context) {
	rRTC := &struct {
		SystemRTC string `json:"systemRTC"`
	}{}
	err := context.ShouldBindJSON(rRTC)
	if err != nil {
		fmt.Println("rRTC json unMarshall err,", err)
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "json unMarshall err",
		})
		return
	}
	mylog.Logger.Debugf("systemRTC %v", rRTC)
	system.SystemSetRTC(rRTC.SystemRTC)
	context.JSON(http.StatusOK, model.Response{
		Code: "0",
	})
}

func SystemSetNTPHost(context *gin.Context) {

	rNTPHostAddr := ntp.NTPHostAddrTemplate{}

	err := context.ShouldBindJSON(&rNTPHostAddr)
	if err != nil {
		fmt.Println("rNTPHostAddr json unMarshall err,", err)

		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "json unMarshall err",
			Data:    "",
		})
		return
	}

	ntp.NTPHostAddr = rNTPHostAddr
	ntp.WriteNTPHostAddrToJson()
	context.JSON(http.StatusOK, model.Response{
		Code:    "0",
		Message: "",
		Data:    "",
	})
}

func SystemGetNTPHost(context *gin.Context) {
	context.JSON(http.StatusOK, model.Response{
		Code:    "0",
		Message: "",
		Data:    ntp.NTPHostAddr,
	})
}

// BackupConfigs 备份plugin目录和selfpara目录
func BackupConfigs(context *gin.Context) {

	err := backup.Zip(device.BACKUPZIP, device.SELFPARAPATH, device.PLUGINPATH)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("backup config  error:%v", err),
		})
		return
	}
	//返回文件流
	context.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filepath.Base(device.BACKUPZIP)))
	defer os.Remove(device.BACKUPZIP)
	context.File(device.BACKUPZIP) //返回文件路径，自动调用http.ServeFile方法

}

func RecoverFiles(context *gin.Context) {

	// 获取文件头
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "no form file named file",
		})
		return
	}

	dst := device.BACKUPZIP
	//保存文件到服务器本地
	if err := context.SaveUploadedFile(file, dst); err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("Save File Error:%v", err),
		})

		return
	}

	if err = backup.Unzip(dst, "."); err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("unzip %s error:%v", dst, err),
		})
		return
	}
	context.JSON(http.StatusOK, model.Response{
		Code:    "0",
		Message: fmt.Sprintf("unzip %s success", dst),
	})
}

func SystemUpdate(context *gin.Context) {

	// 获取文件头
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
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
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "Save File Error",
			Data:    "",
		})
		return
	}

	//升级文件解析
	status := backup.Update(file.Filename)
	if status == true {
		context.JSON(http.StatusOK, model.Response{
			Code:    "0",
			Message: "",
			Data:    "",
		})
		system.SystemReboot()
	}
}
