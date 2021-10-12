/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-29 07:58:04
@LastEditors: WalkMiao
@LastEditTime: 2021-10-06 10:46:44
@FilePath: /goAdapter-Raw/httpServer/controller/deviceTSL.go
*/
package controller

import (
	"encoding/json"
	"fmt"
	"goAdapter/device"
	"goAdapter/httpServer/model"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leandro-lugaresi/hub"
)

func AddDeviceTSL(context *gin.Context) {

	data, err := context.GetRawData()
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("gin context get raw data error:%v", err)})
		return
	}

	var tslInfo model.DeviceTSLTemplate

	err = json.Unmarshal(data, &tslInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("unmarshal tslinfo error:%v", err)})
		return
	}

	if ok := device.DeviceTSLMap.AddTSL(&tslInfo); !ok {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: "add tslinfo failed"})
		return
	}
	//驻留内存中程序退出后defer执行
	// if err = device.WriteToJsonFile(device.DEVICETSLJSON); err != nil {
	// 	context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("write to %s file error:%v", device.DEVICETSLJSON, err)})
	// 	return
	// }
	context.JSON(http.StatusOK, model.Response{Code: "0", Message: fmt.Sprintf("add tsl %s success", tslInfo.Name)})

}

func DeleteDeviceTSL(context *gin.Context) {

	data, err := context.GetRawData()
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("gin context get raw data error:%v", err)})
		return
	}

	var tslInfo model.DeviceTSLTemplate
	err = json.Unmarshal(data, &tslInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("unmarshal tslinfo error:%v", err)})
		return
	}

	if ok := device.DeviceTSLMap.DeleteTSL(&tslInfo); !ok {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: "delete tslinfo failed"})
		return
	}
	context.JSON(http.StatusOK, model.Response{Code: "0", Message: fmt.Sprintf("delete tsl %s success", tslInfo.Name)})

}

func ModifyDeviceTSL(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("gin context get raw data error:%v", err)})
		return
	}

	var tslInfo model.DeviceTSLTemplate

	err = json.Unmarshal(data, &tslInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("unmarshal tslinfo error:%v", err)})
		return
	}
	tmp := device.DeviceTSLMap.Get(tslInfo.Name)
	if tmp == nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("no such tsl template:%s", tslInfo.Name)})
		return
	}
	tmp.Explain = tslInfo.Explain
	context.JSON(http.StatusOK, model.Response{Code: "0", Message: fmt.Sprintf("modify tsl template %s success", tslInfo.Name)})
}

func GetDeviceTSL(context *gin.Context) {
	temps := device.DeviceTSLMap.GetAll()
	type TSLInfoTemplate struct {
		Name    string
		Explain string
		Plugin  string
	}
	resp := struct {
		Code    string            `json:"Code"`
		Message string            `json:"Message"`
		Data    []TSLInfoTemplate `json:"Data"`
	}{
		Code:    "0",
		Message: "get tsl templates ok",
		Data:    make([]TSLInfoTemplate, 0, len(temps)),
	}

	for _, v := range temps {
		item := TSLInfoTemplate{
			Name:    v.Name,
			Explain: v.Explain,
			Plugin:  v.Plugin,
		}
		resp.Data = append(resp.Data, item)
	}

	context.JSON(http.StatusOK, resp)
}

func GetDeviceTSLContents(context *gin.Context) {

	tslName := context.Query("TSLName")
	tmp := device.DeviceTSLMap.Get(tslName)
	if tmp == nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("tsl %s is not exists!", tslName)})
		return

	}
	var p = struct {
		Properties []model.DeviceTSLPropertyTemplate
		Services   []model.DeviceTSLServiceTempalte
	}{}
	if tmp.Properties == nil {
		tmp.Properties = make([]model.DeviceTSLPropertyTemplate, 0, 100)
	}
	if tmp.Services == nil {
		tmp.Services = make([]model.DeviceTSLServiceTempalte, 0, 100)
	}

	p.Properties = tmp.Properties
	p.Services = tmp.Services
	resp := model.Response{
		Code: "0",
		Data: p,
	}

	context.JSON(http.StatusOK, resp)
}

func AddDeviceTSLProperty(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("gin context get raw data error:%v", err)})
		return
	}

	tslInfo := &struct {
		TSLName  string
		Property model.DeviceTSLPropertyTemplate
	}{}
	err = json.Unmarshal(data, &tslInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("post tsl property error:%v", err)})
		return
	}
	tmp := device.DeviceTSLMap.Get(tslInfo.TSLName)
	if tmp == nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("tsl %s is not exists!", tslInfo.TSLName)})
		return
	}

	for _, v := range tmp.Properties {
		if v.Name == tslInfo.Property.Name {
			context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("tsl [%s] property [%s] is exists!", tslInfo.TSLName, tslInfo.Property.Name)})
			return
		}
	}
	AddProperty(tmp, tslInfo.Property)
	context.JSON(http.StatusOK, model.Response{Code: "0"})
}

func ModifyDeviceTSLProperty(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("gin context get raw data error:%v", err)})
		return
	}
	tslInfo := &struct {
		TSLName  string
		Property model.DeviceTSLPropertyTemplate
	}{}
	err = json.Unmarshal(data, &tslInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("post tsl property error:%v", err)})
		return
	}
	tmp := device.DeviceTSLMap.Get(tslInfo.TSLName)
	if tmp == nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("tsl %s is not exists!", tslInfo.TSLName)})
		return
	}
	ModifyProperty(tmp, tslInfo.Property)
	context.JSON(http.StatusOK, model.Response{Code: "0"})
}

func DeleteDeviceTSLProperties(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("gin context get raw data error:%v", err)})
		return
	}
	tslInfo := &struct {
		TSLName    string
		Properties []model.DeviceTSLPropertyTemplate
	}{}
	err = json.Unmarshal(data, &tslInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("post tsl property error:%v", err)})
		return
	}
	tmp := device.DeviceTSLMap.Get(tslInfo.TSLName)
	DeleteProperty(tmp, tslInfo.Properties)
	context.JSON(http.StatusOK, model.Response{Code: "0"})
}

func ImportDeviceTSLContents(context *gin.Context) {

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
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("formfile error:%v", err)})
		return
	}

	f, err := file.Open()
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("open temp file  error:%v", err)})
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("read temp file  error:%v", err)})
		return
	}
	var bak = struct {
		Properties []model.DeviceTSLPropertyTemplate
		Services   []model.DeviceTSLServiceTempalte
	}{}
	if err = json.Unmarshal(data, &bak); err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("unmarshal %s  error:%v", file.Filename, err)})
		return
	}
	device.DeviceTSLMap.SetChanged(true)
	if bak.Properties != nil {
		tmp.Properties = bak.Properties

		//发布消息
		device.DeviceTSLMap.Publish(device.PropertySync, hub.Fields{
			"properties": bak.Properties,
			"plugin":     tmp.Plugin,
		})

	}
	if bak.Services != nil {
		tmp.Services = bak.Services
	}
	context.JSON(http.StatusOK, model.Response{Code: "0", Message: fmt.Sprintf("import %s success", file.Filename)})

}

func ExportDeviceTSLContents(context *gin.Context) {

	tslName := context.Query("TSLName")
	tmp := device.DeviceTSLMap.Get(tslName)
	if tmp == nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("tsl %s is not exists!", tslName)})
		return

	}
	var bak = struct {
		Properties []model.DeviceTSLPropertyTemplate
		Services   []model.DeviceTSLServiceTempalte
	}{
		Properties: tmp.Properties,
		Services:   tmp.Services,
	}
	data, err := json.Marshal(bak)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{Code: "1", Message: fmt.Sprintf("marshaling tsl %s  property and services error:%v", tslName, err)})
		return
	}
	filename := fmt.Sprintf("TSL-%s.bak", tslName)
	context.Writer.WriteHeader(http.StatusOK)

	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	context.Header("Content-Type", "application/octet-stream")
	context.Header("Accept-Length", fmt.Sprintf("%d", len(data)))
	context.Header("Content-Type", "application/download")
	context.Writer.Write(data)
}

// 属性操作
func AddProperty(tslTemplate *model.DeviceTSLTemplate, propery model.DeviceTSLPropertyTemplate) {
	device.DeviceTSLMap.Lock()
	tslTemplate.Properties = append(tslTemplate.Properties, propery)
	device.DeviceTSLMap.SetChanged(true)
	device.DeviceTSLMap.Publish(device.PropertyAdd, hub.Fields{
		"property": propery,
		"plugin":   tslTemplate.Plugin,
		"name":     propery.Name,
	})
	device.DeviceTSLMap.Unlock()
}

func ModifyProperty(tslTemplate *model.DeviceTSLTemplate, property model.DeviceTSLPropertyTemplate) {
	device.DeviceTSLMap.Lock()
	var index = -1
	for i := 0; i < len(tslTemplate.Properties); i++ {
		if tslTemplate.Properties[i].Name == property.Name {
			index = i
			tslTemplate.Properties[i] = property
		}
	}
	if index != -1 {
		device.DeviceTSLMap.Publish(device.PropertyUpdate, hub.Fields{
			"property": property,
			"plugin":   tslTemplate.Plugin,
			"name":     property.Name,
		})
	}

	device.DeviceTSLMap.Unlock()
}

func DeleteProperty(tslTemplate *model.DeviceTSLTemplate, properties []model.DeviceTSLPropertyTemplate) error {
	if len(properties) != 1 {
		return fmt.Errorf("deleted properties is not 1")
	}
	device.DeviceTSLMap.Lock()
	var i = -1
LOOP:
	for index, v := range tslTemplate.Properties {
		for _, pName := range properties {
			if v.Name == pName.Name {
				i = index
				break LOOP
			}
		}
	}
	if i != -1 {
		tslTemplate.Properties = append(tslTemplate.Properties[:i], tslTemplate.Properties[i+1:]...)
		device.DeviceTSLMap.Publish(device.PropertyDelete, hub.Fields{
			"property": properties[0],
			"plugin":   tslTemplate.Plugin,
			"name":     properties[0].Name,
		})
		device.DeviceTSLMap.SetChanged(true)
	}
	device.DeviceTSLMap.Unlock()
	return nil
}
