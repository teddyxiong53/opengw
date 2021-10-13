package controller

import (
	"fmt"
	"goAdapter/device"
	"goAdapter/httpServer/model"
	"goAdapter/pkg/backup"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExportDeviceTSLPlugin(context *gin.Context) {

	tslName := context.Query("TSLName")

	tmp := device.DeviceTSLMap.Get(tslName)
	if tmp == nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("tsl %s is not exists!", tslName)})
		return
	}

	zipPath, err := device.DeviceTSLExportPlugin(tmp.Plugin)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("export plugin %s error:%v!", tmp.Plugin, err),
		})
		return
	}

	//返回文件流
	context.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s", filepath.Base(zipPath)))
	context.File(zipPath) //返回文件路径，自动调用http.ServeFile方法

}

func ImportDeviceTSLPlugin(context *gin.Context) {

	// 获取物模型名称
	tslName := context.PostForm("TSLName")
	tmp := device.DeviceTSLMap.Get(tslName)
	if tmp == nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("tsl %s is not exists!", tslName)})
		return
	}

	// 获取文件头
	file, err := context.FormFile("FileName")
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("form file error:%v", err)})
		return
	}

	zipPath := path.Join(device.PLUGINPATH, file.Filename)
	names := strings.Split(file.Filename, ".")
	if len(names) != 2 {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("invalid tsl template filename:%s", file.Filename)})
		return
	}
	pluginName := names[0]
	if !fileExist(device.PLUGINPATH) {
		os.MkdirAll(device.PLUGINPATH, 0666)
	}
	if fileExist(zipPath) {
		device.DeviceTSLMap.ModifyPlugin(tslName, pluginName)
		context.JSON(http.StatusOK, model.Response{Code: "0", Message: fmt.Sprintf("plugin %s already exists!", pluginName)})
		return
	}
	//保存文件到服务器本地
	if err := context.SaveUploadedFile(file, zipPath); err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("save file error:%v", err)})
		return
	}

	defer os.Remove(zipPath)

	err = backup.Unzip(zipPath, device.PLUGINPATH)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("unzip plugin  %s error:%v!", pluginName, err)})
		return
	}

	if err := device.ReadPlugins(device.PLUGINPATH); err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("roadmap plugin dir  error:%v!", err)})
		return
	}
	if err := device.DeviceTSLMap.ModifyPlugin(tslName, pluginName); err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("modify  plugin  error:%v!", err)})
		return
	}

	context.JSON(http.StatusOK, model.Response{Code: "0", Message: fmt.Sprintf("import plugin  %s success!", pluginName)})

}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
