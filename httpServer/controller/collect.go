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
	if v := device.CollectInterfaceMap.Get(interfaceInfo.CollInterfaceName); v != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("采集接口【%s】已经存在", interfaceInfo.CollInterfaceName),
		})
		return
	}
	comm := device.CommunicationInterfaceMap.Get(interfaceInfo.CommInterfaceName)
	if comm == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("未找到已注册的通讯接口【%s】!", interfaceInfo.CommInterfaceName),
		})
		return
	}
	//通信口绑定采集接口
	comm.Bind(interfaceInfo.CollInterfaceName)

	//是否只针对特定采集口 比如串口或者TCP客户端
	if used, collectName := device.CollectInterfaceMap.CommCheck(comm.GetName()); used {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("通讯口【%s】已经被接口【%s】使用!", comm.GetName(), collectName),
		})
		return
	}

	nodeManage, _ := device.NewCollectInterface(&interfaceInfo)
	device.CollectInterfaceMap.Add(nodeManage)
	// 废弃
	// device.CommunicationManage.Collectors <- &device.CollectInterfaceStatus{
	// 	Tmp: nodeManage,
	// 	ACT: device.ADD,
	// }
	context.JSON(200, struct {
		Code    string
		Message string
	}{
		Code:    "0",
		Message: fmt.Sprintf("add interface %s success", interfaceInfo.CollInterfaceName),
	})

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

	old := device.CollectInterfaceMap.Get(interfaceInfo.CollInterfaceName)
	if old == nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("collInterface【%s】is not exist", interfaceInfo.CollInterfaceName),
		})
		return
	}

	//只有关键参数更改才会发布这个消息
	if interfaceInfo.OfflinePeriod != old.OfflinePeriod || interfaceInfo.PollPeriod != old.PollPeriod || interfaceInfo.CommInterfaceName != old.CommInterfaceName {
		device.CollectInterfaceMap.Update(interfaceInfo)
		//废弃
		//先删
		// device.CommunicationManage.Collectors <- &device.CollectInterfaceStatus{
		// 	Tmp: old,
		// 	ACT: device.DELETE,
		// }
		// old.CollInterfaceName = interfaceInfo.CollInterfaceName
		// old.CommInterfaceName = interfaceInfo.CommInterfaceName
		// old.PollPeriod = interfaceInfo.PollPeriod
		// old.OfflinePeriod = interfaceInfo.OfflinePeriod
		// //后增
		// device.CommunicationManage.Collectors <- &device.CollectInterfaceStatus{
		// 	Tmp: old,
		// 	ACT: device.ADD,
		// }
	}

	context.JSON(200, struct {
		Code    string
		Message string
	}{
		Code:    "0",
		Message: fmt.Sprintf("modify interface %s success", interfaceInfo.CollInterfaceName),
	})

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

	if ok := device.CollectInterfaceMap.Delete(interfaceInfo.CollectInterfaceName); !ok {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("key %s is not exist!", interfaceInfo.CollectInterfaceName),
		})
		return
	}
	context.JSON(200, struct {
		Code    string
		Message string
	}{
		Code:    "0",
		Message: fmt.Sprintf("delete interface %s success", interfaceInfo.CollectInterfaceName),
	})
}

//接口详情
func GetInterfaceInfo(context *gin.Context) {

	sName := context.Query("CollInterfaceName")

	aParam := &struct {
		Code    string
		Message string
		Data    interface{}
	}{}

	v := device.CollectInterfaceMap.Get(sName)
	if v == nil {
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
	tmps := device.CollectInterfaceMap.GetAll()
	for _, v := range tmps {
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
