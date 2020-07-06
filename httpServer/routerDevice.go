package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goAdapter/device"
	"log"
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

	device.DeviceInterfaceMap[interfaceInfo.InterfaceID] = device.NewDeviceInterface(interfaceInfo.InterfaceID,
		interfaceInfo.PollPeriod,
		interfaceInfo.OfflinePeriod,0)

	device.WriteDeviceInterfaceManageToJson()

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

	device.DeviceInterfaceMap[interfaceInfo.InterfaceID].ModifyDeviceInterface(interfaceInfo.PollPeriod,
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
		Data    device.DeviceInterfaceTemplate
	}{}

	iID,_ := strconv.Atoi(sID)

	if iID < len(device.DeviceInterfaceMap){
		aParam.Code = "0"
		aParam.Message = ""
		aParam.Data = *device.DeviceInterfaceMap[iID]
	}else{
		aParam.Code = "1"
		aParam.Message = "interface is noexist"
		aParam.Data = device.DeviceInterfaceTemplate{}
	}

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}


func apiGetAllInterfaceInfo(context *gin.Context){

	aParam := &struct{
		Code 	string
		Message string
		Data    [device.MaxDeviceInterfaceManage]*device.DeviceInterfaceTemplate
	}{}

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = device.DeviceInterfaceMap

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
	status,aParam.Message = device.DeviceInterfaceMap[nodeInfo.InterfaceID].AddDeviceNode(nodeInfo.DAddr,nodeInfo.DType)
	if status == true{
		device.WriteDeviceInterfaceManageToJson()

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
	device.WriteDeviceInterfaceManageToJson()

	aParam.Code = "0"
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiGetNode(context *gin.Context){

	sID := context.Query("interfaceID")
	sAddr := context.Query("addr")

	aParam := &struct{
		Code 	string
		Message string
		Data    []device.VariableTemplate
	}{}

	iID,_ := strconv.Atoi(sID)
	for k,v := range device.DeviceInterfaceMap[iID].DeviceNodeAddrMap{
		if v == sAddr{
			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = device.DeviceInterfaceMap[iID].DeviceNodeMap[k].GetDeviceVariablesValue()
			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}
	aParam.Code = "1"
	aParam.Message = "node is noexist"
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
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

	device.DeviceInterfaceMap[nodeInfo.InterfaceID].DeleteDeviceNode(nodeInfo.DAddr,nodeInfo.DType)

	device.WriteDeviceInterfaceManageToJson()

	aParam.Code = "0"
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiAddTemplate(context *gin.Context){

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
		TemplateName string					`json:"templateName"`		//模板名称
		TemplateType string					`json:"templateType"`		//模板型号
		TemplateMessage string              `json:"templateMessage"`	//备注信息
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

	index := len(device.DeviceInterfaceParamMap.DeviceNodeTypeMap)
	template := device.DeviceNodeTypeTemplate{
		TemplateName:interfaceInfo.TemplateName,
		TemplateType:interfaceInfo.TemplateType,
		TemplateID: index,
		TemplateMessage:interfaceInfo.TemplateMessage,
	}

	device.DeviceInterfaceParamMap.DeviceNodeTypeMap = append(device.DeviceInterfaceParamMap.DeviceNodeTypeMap,template)

	device.WriteDeviceInterfaceManageToJson()

	aParam.Code = "0"
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiGetTemplate(context *gin.Context){

	aParam := &struct{
		Code 	string
		Message string
		Data    []device.DeviceNodeTypeTemplate
	}{}

	log.Printf("%+v\n",device.DeviceInterfaceParamMap.DeviceNodeTypeMap)

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = device.DeviceInterfaceParamMap.DeviceNodeTypeMap

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}