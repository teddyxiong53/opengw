package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func apiAddInterface(context *gin.Context){

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

	interfaceInfo := &struct{
		InterfaceID  int		`json:"interfaceID"`
		PollPeriod 	int			`json:"pollPeriod"`
		OfflinePeriod int		`json:"offlinePeriod"`
	}{}

	err := json.Unmarshal(bodyBuf[:n],interfaceInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	DeviceNodeManageMap[interfaceInfo.InterfaceID] = NewDeviceNodeManage(interfaceInfo.InterfaceID,
		interfaceInfo.PollPeriod,
		interfaceInfo.OfflinePeriod,0)

	aParam.Code = "0"
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiModifyInterface(context *gin.Context){

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

	interfaceInfo := &struct{
		InterfaceID  int
		PollPeriod int
		OfflinePeriod int
	}{}

	err := json.Unmarshal(bodyBuf[:n],interfaceInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	DeviceNodeManageMap[interfaceInfo.InterfaceID].ModifyDeviceNodeManage(interfaceInfo.PollPeriod,
		interfaceInfo.OfflinePeriod)

	aParam.Code = "0"
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}


func apiGetInterfaceInfo(context *gin.Context){

	sID := context.Query("interfaceID")
	fmt.Println(sID)

	aParam := &struct{
		Code 	string
		Message string
		Data    interface{}
	}{}

	iID,_ := strconv.Atoi(sID)

	nodeManage := struct{
		InterfaceID 		int										`json:"InterfaceID"`			//通信接口
		PollPeriod 			int										`json:"PollPeriod"`				//采集周期
		OfflinePeriod 		int      								`json:"OfflinePeriod"`			//离线超时周期
		DeviceNodeCnt       int                 					`json:"DeviceNodeCnt"`			//设备数量
		DeviceNodeMap       []interface{} 							`json:"DeviceNodeMap"`			//节点链表
		DeviceUseMap       	[50]bool 									`json:"DeviceUseMap"`			//节点使用
	}{
		InterfaceID:DeviceNodeManageMap[iID].InterfaceID,
		PollPeriod:DeviceNodeManageMap[iID].PollPeriod,
		OfflinePeriod:DeviceNodeManageMap[iID].OfflinePeriod,
		DeviceNodeCnt:DeviceNodeManageMap[iID].DeviceNodeCnt,
		DeviceUseMap:DeviceNodeManageMap[iID].DeviceNodeUseMap,
	}

	nodeManage.DeviceNodeMap = make([]interface{},0)
	for k,v := range DeviceNodeManageMap[iID].DeviceNodeUseMap{
		if v == true{
			nodeManage.DeviceNodeMap = append(nodeManage.DeviceNodeMap,DeviceNodeManageMap[iID].DeviceNodeMap[k])
		}
	}

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = nodeManage

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}


func apiGetAllInterfaceInfo(context *gin.Context){

	aParam := &struct{
		Code 	string
		Message string
		Data    [8]*DeviceNodeManage
	}{}

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = DeviceNodeManageMap

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

func apiAddNode(context *gin.Context){

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

	nodeInfo := &struct{
		InterfaceID  	int			`json:"interfaceID"`
		DAddr 			string		`json:"addr"`
		DType 			string		`json:"type"`
	}{}

	err := json.Unmarshal(bodyBuf[:n],nodeInfo)
	if err != nil {
		fmt.Println("nodeInfo json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	var status bool
	status,aParam.Message = DeviceNodeManageMap[nodeInfo.InterfaceID].AddDeviceNode(nodeInfo.DAddr,nodeInfo.DType)
	if status == true{
		WriteDeviceInterfaceManageToJson()

		aParam.Code = "0"
		aParam.Data = ""
	}else{
		aParam.Code = "1"
		aParam.Data = ""
	}

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiModifyNode(context *gin.Context){

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

	nodeInfo := &struct{
		InterfaceID  	int			`json:"interfaceID"`
		DAddr 			string		`json:"addr"`
		DType 			string		`json:"type"`
	}{}

	err := json.Unmarshal(bodyBuf[:n],nodeInfo)
	if err != nil {
		fmt.Println("nodeInfo json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	//DeviceNodeManageMap[nodeInfo.InterfaceID].ModifyDeviceNode(nodeInfo.DAddr,nodeInfo.DType)
	WriteDeviceInterfaceManageToJson()

	aParam.Code = "0"
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiDeleteNode(context *gin.Context){

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

	nodeInfo := &struct{
		InterfaceID  	int			`json:"interfaceID"`
		DAddr 			string		`json:"addr"`
		DType 			string		`json:"type"`
	}{}

	err := json.Unmarshal(bodyBuf[:n],nodeInfo)
	if err != nil {
		fmt.Println("nodeInfo json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	DeviceNodeManageMap[nodeInfo.InterfaceID].DeleteDeviceNode(nodeInfo.DAddr,nodeInfo.DType)

	WriteDeviceInterfaceManageToJson()

	aParam.Code = "0"
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}
