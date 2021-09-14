package controller

import (
	"encoding/json"
	"fmt"
	"goAdapter/device"
	"goAdapter/httpServer/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddInterface(context *gin.Context) {
	var interfaceInfo device.CollectInterfaceParamTemplate
	err := context.ShouldBindJSON(&interfaceInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}

	//判断串口是否已经被其他采集接口使用了
	if _, ok := device.CollectInterfaceMap[interfaceInfo.CollInterfaceName]; ok {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("采集接口【%s】已经存在", interfaceInfo.CollInterfaceName),
		})
		return
	}
	comm, ok := device.CommunicationInterfaceMap[interfaceInfo.CommInterfaceName]
	if !ok {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("未找到已注册的通讯接口【%s】!", interfaceInfo.CommInterfaceName),
		})
		return
	}
	for k, v := range device.CollectInterfaceMap {
		//是否只针对特定采集口 比如串口或者TCP客户端
		if v.CommInterfaceName == comm.GetName() {
			context.JSON(200, model.Response{
				Code:    "1",
				Message: fmt.Sprintf("通讯接口【%s】已经被【%s】使用!", comm.GetName(), k),
			})
			return
		}
	}

	nodeManage, _ := device.NewCollectInterface(&interfaceInfo)
	device.CollectInterfaceMap[interfaceInfo.CollInterfaceName] = nodeManage
	device.CommunicationManage.Collectors <- &device.CollectInterfaceStatus{
		Tmp: nodeManage,
		ACT: device.ADD,
	}
	device.WriteJsonErrorHandler(context, device.COLLINTERFACEJSON,
		200, 200, fmt.Sprintf("add interface %s success", interfaceInfo.CollInterfaceName))
}

func ModifyInterface(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}
	var interfaceInfo device.CollectInterfaceParamTemplate
	err = json.Unmarshal(data, &interfaceInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("json unmarshal error:%v", err),
		})
		return
	}

	old, ok := device.CollectInterfaceMap[interfaceInfo.CollInterfaceName]
	if !ok {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "collInterface is not exist",
		})
		return
	}
	//先删
	device.CommunicationManage.Collectors <- &device.CollectInterfaceStatus{
		Tmp: old,
		ACT: device.DELETE,
	}
	old.CollInterfaceName = interfaceInfo.CollInterfaceName
	old.CommInterfaceName = interfaceInfo.CommInterfaceName
	old.PollPeriod = interfaceInfo.PollPeriod
	old.OfflinePeriod = interfaceInfo.OfflinePeriod
	//后增
	device.CommunicationManage.Collectors <- &device.CollectInterfaceStatus{
		Tmp: old,
		ACT: device.ADD,
	}
	device.WriteJsonErrorHandler(context, device.COLLINTERFACEJSON,
		200, 200, fmt.Sprintf("modify interface %s success", interfaceInfo.CollInterfaceName))
}

func DeleteInterface(context *gin.Context) {

	data, err := context.GetRawData()
	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}

	interfaceInfo := &struct {
		CollectInterfaceName string `json:"CollInterfaceName"` // 采集接口名字
	}{}

	err = json.Unmarshal(data, interfaceInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}

	if _, ok := device.CollectInterfaceMap[interfaceInfo.CollectInterfaceName]; !ok {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("key %s is not exist!", interfaceInfo.CollectInterfaceName),
		})
		return
	}
	v, ok := device.CollectInterfaceMap[interfaceInfo.CollectInterfaceName]
	if ok {
		delete(device.CollectInterfaceMap, interfaceInfo.CollectInterfaceName)
		device.CommunicationManage.Collectors <- &device.CollectInterfaceStatus{
			Tmp: v,
			ACT: device.DELETE,
		}
	} else {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("interface %s is not registered", interfaceInfo.CollectInterfaceName),
		})
		return
	}
	device.WriteJsonErrorHandler(context, device.COLLINTERFACEJSON,
		200, 200, fmt.Sprintf("delete interface %s success", interfaceInfo.CollectInterfaceName))
}

//接口详情
func GetInterfaceInfo(context *gin.Context) {

	sName := context.Query("CollInterfaceName")

	aParam := &struct {
		Code    string
		Message string
		Data    interface{}
	}{}

	v, ok := device.CollectInterfaceMap[sName]
	if !ok {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("key %s is not exist", sName),
		})
		return
	}
	if v.DeviceNodes == nil {
		v.DeviceNodes = make([]*device.DeviceNodeTemplate, 0)
	}
	aParam.Code = "0"
	aParam.Data = struct {
		CollInterfaceName   string                       `json:"CollInterfaceName"`   //采集接口
		CommInterfaceName   string                       `json:"CommInterfaceName"`   //通信接口
		PollPeriod          int                          `json:"PollPeriod"`          //采集周期
		OfflinePeriod       int                          `json:"OfflinePeriod"`       //离线超时周期
		DeviceNodeCnt       int                          `json:"DeviceNodeCnt"`       //设备数量
		DeviceNodeOnlineCnt int                          `json:"DeviceNodeOnlineCnt"` //设备在线数量
		DeviceNodeMap       []*device.DeviceNodeTemplate `json:"DeviceNodeMap"`
	}{
		CollInterfaceName:   v.CollInterfaceName,
		CommInterfaceName:   v.CommInterfaceName,
		PollPeriod:          v.PollPeriod,
		OfflinePeriod:       v.OfflinePeriod,
		DeviceNodeCnt:       v.DeviceNodeCnt,
		DeviceNodeOnlineCnt: v.DeviceNodeOnlineCnt,
		DeviceNodeMap:       v.DeviceNodes,
	}

	context.JSON(http.StatusOK, aParam)
}

func GetAllInterfaceInfo(context *gin.Context) {

	type InterfaceParamTemplate struct {
		CollInterfaceName   string `json:"CollInterfaceName"`   // 采集接口
		CommInterfaceName   string `json:"CommInterfaceName"`   // 通信接口
		PollPeriod          int    `json:"PollPeriod"`          // 采集周期
		OfflinePeriod       int    `json:"OfflinePeriod"`       // 离线超时周期
		DeviceNodeCnt       int    `json:"DeviceNodeCnt"`       // 设备数量
		DeviceNodeOnlineCnt int    `json:"DeviceNodeOnlineCnt"` // 设备在线数量
	}

	aParam := &struct {
		Code    string
		Message string
		Data    []InterfaceParamTemplate
	}{}

	aParam.Data = make([]InterfaceParamTemplate, 0)

	aParam.Code = "0"
	aParam.Message = ""
	for _, v := range device.CollectInterfaceMap {
		Param := InterfaceParamTemplate{
			CollInterfaceName:   v.CollInterfaceName,
			CommInterfaceName:   v.CommInterfaceName,
			PollPeriod:          v.PollPeriod,
			OfflinePeriod:       v.OfflinePeriod,
			DeviceNodeCnt:       v.DeviceNodeCnt,
			DeviceNodeOnlineCnt: v.DeviceNodeOnlineCnt,
		}
		aParam.Data = append(aParam.Data, Param)
	}

	context.JSON(http.StatusOK, aParam)
}
