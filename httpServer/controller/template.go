package controller

import (
	"encoding/json"
	"goAdapter/device"
	"goAdapter/httpServer/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddTemplate(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
	}

	typeInfo := &struct {
		TemplateName    string `json:"TemplateName"`    // 模板名称
		TemplateType    string `json:"TemplateType"`    // 模板型号
		TemplateMessage string `json:"TemplateMessage"` // 备注信息
	}{}

	err = json.Unmarshal(data, typeInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}

	template := &device.DeviceTemplate{
		Name:    typeInfo.TemplateName,
		Type:    typeInfo.TemplateType,
		Message: typeInfo.TemplateMessage,
	}

	device.DeviceTemplateMap[template.Type] = template

	context.JSON(http.StatusOK, model.Response{
		Code:    "1",
		Message: "add device template success",
	})
}

func GetTemplate(context *gin.Context) {
	//清空设备模版缓存
	//device.DeviceTemplateMap = make(map[string]*device.DeviceTemplate)

	//获取最新的模版
	// if err := device.ReadPlugins(device.PLUGINPATH); err != nil {
	// 	context.JSON(200, model.Response{
	// 		Code:    "1",
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }

	ts := make([]*device.DeviceTemplate, 0, len(device.DeviceTemplateMap))
	for _, v := range device.DeviceTemplateMap {
		ts = append(ts, v)
	}

	context.JSON(http.StatusOK, &struct {
		Code    string
		Message string
		Data    []*device.DeviceTemplate
	}{
		Code: "0",
		Data: ts,
	})
}
