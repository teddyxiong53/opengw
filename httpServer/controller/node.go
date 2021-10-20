package controller

import (
	"encoding/json"
	"fmt"
	"goAdapter/httpServer/model"
	"net/http"
	"strings"

	"goAdapter/device"

	"github.com/gin-gonic/gin"
)

func AddNode(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}
	nodeInfo := &struct {
		InterfaceName string `json:"CollInterfaceName"`
		DAddr         string `json:"Addr"`
		DType         string `json:"Type"`
		DName         string `json:"Name"`
	}{}

	err = json.Unmarshal(data, nodeInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}

	v := device.CollectInterfaceMap.Get(nodeInfo.InterfaceName)
	if v == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("interface %s is not exist", nodeInfo.InterfaceName),
		})
		return
	}

	for _, node := range v.DeviceNodes {
		if node.Addr == nodeInfo.DAddr {
			context.JSON(200, model.Response{
				Code:    "1",
				Message: fmt.Sprintf("设备地址与已存在设备【%s】冲突!",node.Name),
			})
			return
		}
		if node.Name == nodeInfo.DName {
			context.JSON(200, model.Response{
				Code:    "1",
				Message: fmt.Sprintf("设备 %s 名称冲突!", nodeInfo.DName),
			})
			return
		}
	}

	err = v.AddDeviceNode(nodeInfo.DName, nodeInfo.DType, nodeInfo.DAddr)
	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}

	//设备增加了设置变化状态
	device.CollectInterfaceMap.Lock()
	device.CollectInterfaceMap.Changed = true
	device.CollectInterfaceMap.Unlock()

	context.JSON(200, struct {
		Code    string
		Message string
	}{
		Code:    "0",
		Message: fmt.Sprintf("add node %s of interface %s success", nodeInfo.InterfaceName, nodeInfo.InterfaceName),
	})

}

func ModifyNode(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}
	nodeInfo := &struct {
		InterfaceName string `json:"CollInterfaceName"`
		Name          string `json:"Name"`
		DType         string `json:"Type"`
		Addr          string `json:"Addr"`
	}{}

	err = json.Unmarshal(data, nodeInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}
	v := device.CollectInterfaceMap.Get(nodeInfo.InterfaceName)
	if v == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("interface %s is not exists!", nodeInfo.InterfaceName),
		})
		return
	}

	for _, node := range v.DeviceNodes {
		if node.Name == nodeInfo.Name {
			node.Type = nodeInfo.DType
			node.Addr = nodeInfo.Addr
			device.CollectInterfaceMap.Lock()
			device.CollectInterfaceMap.Changed = true
			device.CollectInterfaceMap.Unlock()
			context.JSON(200, struct {
				Code    string
				Message string
			}{
				Code:    "0",
				Message: fmt.Sprintf("modify node %s of interface %s success", nodeInfo.InterfaceName, nodeInfo.InterfaceName),
			})
		}
	}

}

func ModifyNodes(context *gin.Context) {

	data, err := context.GetRawData()
	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}

	nodeInfo := &struct {
		InterfaceName string   `json:"CollInterfaceName"`
		DType         string   `json:"Type"`
		Name          []string `json:"Name"`
	}{}

	err = json.Unmarshal(data, nodeInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}
	v := device.CollectInterfaceMap.Get(nodeInfo.InterfaceName)
	if v == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("interfacename %s is not exists", nodeInfo.InterfaceName),
		})
		return
	}

	for _, name := range nodeInfo.Name {
		for i, node := range v.DeviceNodes {
			if node.Name == name {
				v.DeviceNodes[i].Type = nodeInfo.DType
			}
		}
	}
	device.CollectInterfaceMap.Lock()
	device.CollectInterfaceMap.Changed = true
	device.CollectInterfaceMap.Unlock()
	context.JSON(200, struct {
		Code    string
		Message string
	}{
		Code: "0",
	})
}

func GetNode(context *gin.Context) {

	sName := context.Query("CollInterfaceName")
	sAddr := context.Query("Addr")

	aParam := &struct {
		Code    string
		Message string
		Data    *device.DeviceNodeTemplate
	}{}
	v := device.CollectInterfaceMap.Get(sName)
	if v == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("collect interface %s is not exists!", sName),
		})
		return
	}
	aParam.Data = v.GetDeviceNode(sAddr)
	aParam.Code = "0"
	context.JSON(200, aParam)
}

func DeleteNode(context *gin.Context) {

	data, err := context.GetRawData()
	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}
	nodeInfo := &struct {
		InterfaceName string   `json:"CollInterfaceName"`
		DName         []string `json:"Name"`
	}{
		InterfaceName: "",
		DName:         make([]string, 0),
	}

	err = json.Unmarshal(data, nodeInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}

	v := device.CollectInterfaceMap.Get(nodeInfo.InterfaceName)
	if v == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("interface %s is not exists!", nodeInfo.InterfaceName),
		})
		return
	}

	for _, name := range nodeInfo.DName {
		v.DeleteDeviceNode(name)
	}
	device.CollectInterfaceMap.Lock()
	device.CollectInterfaceMap.Changed = true
	device.CollectInterfaceMap.Unlock()

	context.JSON(200, struct {
		Code    string
		Message string
	}{
		Code:    "0",
		Message: fmt.Sprintf("delete nodes 【%s】 success", strings.Join(nodeInfo.DName, ",")),
	})

}

/**
从缓存中获取设备变量
*/
func GetNodeVariableFromCache(context *gin.Context) {

	type VariableTemplate struct {
		Index     int         `json:"index"` // 变量偏移量
		Name      string      `json:"name"`  // 变量名
		Label     string      `json:"lable"` // 变量标签
		Value     interface{} `json:"value"` // 变量值
		Explain   interface{} `json:"explain"`
		TimeStamp string      `json:"timestamp"` // 变量时间戳
		Type      string      `json:"type"`      // 变量类型
	}

	sName := context.Query("CollInterfaceName")
	sAddr := context.Query("Addr")

	i := device.CollectInterfaceMap.Get(sName)
	if i == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("interface %s is not exists!", sName),
		})
		return
	}
	node := i.GetDeviceNode(sAddr)
	if node == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("no suce device addr %s of interface %s", sAddr, sName),
		})
		return
	}
	vData := make([]VariableTemplate, 0, len(node.Properties))

	vt := VariableTemplate{}
	for _, v := range node.Properties {
		vt.Name = v.Name
		vt.Label = v.Explain
		switch v.Type {
		case device.PropertyTypeUInt32:
			vt.Type = "uint32"
		case device.PropertyTypeInt32:
			vt.Type = "int32"
		case device.PropertyTypeDouble:
			vt.Type = "double"
		case device.PropertyTypeString:
			vt.Type = "string"
		}

		if len(v.Value) > 0 {
			last := v.Value[len(v.Value)-1]
			vt.Value = last.Value
			vt.Explain = last.Explain
			vt.TimeStamp = last.TimeStamp
			vt.Index = last.Index
		}
		vData = append(vData, vt)
	}
	context.JSON(200, &struct {
		Code    string
		Message string
		Data    []VariableTemplate
	}{
		Code: "0",
		Data: vData,
	})
}

func GetNodeHistoryVariable(context *gin.Context) {

	sName := context.Query("CollInterfaceName")
	sAddr := context.Query("Addr")
	sVariable := context.Query("VariableName")

	aParam := &struct {
		Code    string
		Message string
		Data    []model.DeviceTSLPropertyValueTemplate
	}{}

	i := device.CollectInterfaceMap.Get(sName)
	if i == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("interface %s is not exists!", sName),
		})
		return
	}
	node := i.GetDeviceNode(sAddr)
	if node == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("no suce device addr %s of interface %s", sAddr, sName),
		})
		return
	}
	for _, v := range node.Properties {
		if v.Name == sVariable {
			aParam.Code = "0"
			aParam.Data = v.Value
		}
	}
	context.JSON(200, aParam)
}

func GetNodeReadVariable(context *gin.Context) {

	type VariableTemplate struct {
		Index     int         `json:"index"` // 变量偏移量
		Name      string      `json:"name"`  // 变量名
		Label     string      `json:"lable"` // 变量标签
		Value     interface{} `json:"value"` // 变量值
		Explain   interface{} `json:"explain"`
		TimeStamp string      `json:"timestamp"` // 变量时间戳
		Type      string      `json:"type"`      // 变量类型
	}

	sName := context.Query("CollInterfaceName")
	sAddr := context.Query("Addr")

	i := device.CollectInterfaceMap.Get(sName)
	if i == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("interface %s is not exists!", sName),
		})
		return
	}
	node := i.GetDeviceNode(sAddr)
	if node == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("no suce device addr %s of interface %s", sAddr, sName),
		})
		return
	}

	manager := device.CollectInterfaceMap.Get(sName)
	if manager == nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("comm manager is not initialized with interface %s", sName),
		})
		return
	}
	var cmd device.CommunicationCmdTemplate
	cmd.CollInterfaceName = sName
	cmd.DeviceName = node.Name
	cmd.FunName = device.GETREAL
	cmdRX := manager.CommunicationManager.CommunicationManageAddEmergency(cmd)

	if cmdRX.Err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: cmdRX.Err.Error(),
		})
		return
	}
	var variable VariableTemplate
	var vs = make([]VariableTemplate, 0, len(node.Properties))
	for _, v := range node.Properties {
		variable.Name = v.Name
		variable.Label = v.Explain
		variable.Label = v.Explain
		switch v.Type {
		case device.PropertyTypeUInt32:
			variable.Type = "uint32"
		case device.PropertyTypeInt32:
			variable.Type = "int32"
		case device.PropertyTypeDouble:
			variable.Type = "double"
		case device.PropertyTypeString:
			variable.Type = "string"
		}
		if len(v.Value) > 0 {
			last := v.Value[len(v.Value)-1]
			variable.Value = last.Value
			variable.Explain = last.Explain
			variable.TimeStamp = last.TimeStamp
			variable.Index = last.Index
		}
		vs = append(vs, variable)
	}
	context.JSON(200, &struct {
		Code    string
		Message string
		Data    []VariableTemplate
	}{
		Code: "0",
		Data: vs,
	})
}

/**
  从设备中获取设备变量
*/
func GetNodeVariableFromDevice(context *gin.Context) {

	// sName := context.Query("interfaceName")
	// sAddr := context.Query("addr")
	//
	// aParam := &struct {
	//	Code    string
	//	Message string
	//	Data    []api.VariableTemplate
	// }{}
	//
	//
	// for _,v := range device.CollectInterfaceMap {
	//	if v.CollInterfaceName == nodeInfo.InterfaceName {
	//
	//	}
	// }
	//
	//		iID, _ := strconv.Atoi(sID)
	//		for k, v := range device.CollectInterfaceMap[iID].DeviceNodeMap {
	//			if v.Addr == sAddr {
	//
	//				cmd := device.CommunicationCmd{}
	//				cmd.InterfaceID = device.InterFaceID0
	//				cmd.DeviceAddr = v.Addr
	//				cmd.FunName = "GenerateGetRealVariables"
	//				if device.CommunicationManageAddEmergency(cmd) == true {
	//					aParam.Code = "0"
	//					aParam.Message = ""
	//					aParam.Data = device.CollectInterfaceMap[iID].DeviceNodeMap[k].VariableMap
	//				} else {
	//					aParam.Code = "1"
	//					aParam.Message = ""
	//					aParam.Data = device.CollectInterfaceMap[iID].DeviceNodeMap[k].VariableMap
	//
	//				}
	//				sJson, _ := json.Marshal(aParam)
	//				context.String(http.StatusOK, string(sJson))
	//				return
	//			}
	//		}
	//		aParam.Code = "1"
	//		aParam.Message = "node is noexist"
	//		sJson, _ := json.Marshal(aParam)
	//		context.String(http.StatusOK, string(sJson))
}
