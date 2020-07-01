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

	DeviceInterfaceMap[interfaceInfo.InterfaceID] = NewDeviceInterface(interfaceInfo.InterfaceID,
		interfaceInfo.PollPeriod,
		interfaceInfo.OfflinePeriod,0)

	WriteDeviceInterfaceManageToJson()

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

	DeviceInterfaceMap[interfaceInfo.InterfaceID].ModifyDeviceInterface(interfaceInfo.PollPeriod,
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
		Data    DeviceInterfaceTemplate
	}{}

	iID,_ := strconv.Atoi(sID)

	if iID < len(DeviceInterfaceMap){
		aParam.Code = "0"
		aParam.Message = ""
		aParam.Data = *DeviceInterfaceMap[iID]
	}else{
		aParam.Code = "1"
		aParam.Message = "interface is noexist"
		aParam.Data = DeviceInterfaceTemplate{}
	}

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}


func apiGetAllInterfaceInfo(context *gin.Context){

	aParam := &struct{
		Code 	string
		Message string
		Data    [MaxDeviceInterfaceManage]*DeviceInterfaceTemplate
	}{}

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = DeviceInterfaceMap

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
	status,aParam.Message = DeviceInterfaceMap[nodeInfo.InterfaceID].AddDeviceNode(nodeInfo.DAddr,nodeInfo.DType)
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

	DeviceInterfaceMap[nodeInfo.InterfaceID].DeleteDeviceNode(nodeInfo.DAddr,nodeInfo.DType)

	WriteDeviceInterfaceManageToJson()

	aParam.Code = "0"
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}
