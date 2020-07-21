package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goAdapter/api"
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
	status,aParam.Message = device.DeviceInterfaceMap[nodeInfo.InterfaceID].AddDeviceNode(nodeInfo.DType,nodeInfo.DAddr)
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
		Data    []api.VariableTemplate
	}{}

	iID,_ := strconv.Atoi(sID)
	for k,v := range device.DeviceInterfaceMap[iID].DeviceNodeMap{
		if v.Addr == sAddr{
			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = device.DeviceInterfaceMap[iID].DeviceNodeMap[k].VariableMap
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

	typeInfo := &struct{
		TemplateName string					`json:"templateName"`		//模板名称
		TemplateType string					`json:"templateType"`		//模板型号
		TemplateMessage string              `json:"templateMessage"`	//备注信息
	}{}

	err := json.Unmarshal(bodyBuf[:n],typeInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	index := len(device.DeviceNodeTypeMap.DeviceNodeType)
	template := device.DeviceNodeTypeTemplate{
		TemplateName:typeInfo.TemplateName,
		TemplateType:typeInfo.TemplateType,
		TemplateID: index,
		TemplateMessage:typeInfo.TemplateMessage,
	}

	device.DeviceNodeTypeMap.DeviceNodeType = append(device.DeviceNodeTypeMap.DeviceNodeType,template)

	device.WriteDeviceNodeTypeMapToJson()

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

	if len(device.DeviceNodeTypeMap.DeviceNodeType) > 0{
		aParam.Code = "0"
		aParam.Message = ""
		aParam.Data = device.DeviceNodeTypeMap.DeviceNodeType
	}else{
		aParam.Code = "1"
		aParam.Message = "nodeTypeCnt is 0"
		aParam.Data = device.DeviceNodeTypeMap.DeviceNodeType
	}

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

/**
	从缓存中获取设备变量
 */
func apiGetNodeVariableFromCache(context *gin.Context){

	sID := context.Query("interfaceID")
	sAddr := context.Query("addr")

	aParam := &struct{
		Code 	string
		Message string
		Data    []api.VariableTemplate
	}{}

	iID,_ := strconv.Atoi(sID)
	for k,v := range device.DeviceInterfaceMap[iID].DeviceNodeMap{
		if v.Addr == sAddr{
			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = device.DeviceInterfaceMap[iID].DeviceNodeMap[k].VariableMap
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

/**
	从设备中获取设备变量
*/
func apiGetNodeVariableFromDevice(context *gin.Context){

	sID := context.Query("interfaceID")
	sAddr := context.Query("addr")

	aParam := &struct{
		Code 	string
		Message string
		Data    []api.VariableTemplate
	}{}

	iID,_ := strconv.Atoi(sID)
	for k,v := range device.DeviceInterfaceMap[iID].DeviceNodeMap{
		if v.Addr == sAddr{

			cmd := device.CommunicationCmd{}
			cmd.InterfaceID = device.InterFaceID0
			cmd.DeviceAddr = v.Addr
			cmd.FunName = "GenerateGetRealVariables"
			if device.CommunicationManageAddEmergency(cmd) == true{
				aParam.Code = "0"
				aParam.Message = ""
				aParam.Data = device.DeviceInterfaceMap[iID].DeviceNodeMap[k].VariableMap
			}else{
				aParam.Code = "1"
				aParam.Message = ""
				aParam.Data = device.DeviceInterfaceMap[iID].DeviceNodeMap[k].VariableMap

			}
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

func apiAddCommInterface(context *gin.Context){

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
	//fmt.Println(string(bodyBuf[:n]))

	interfaceInfo := struct{
		Name 	string							`json:"Name"`			//接口名称
		Type    string          				`json:"Type"`			//接口类型,比如serial,tcp,udp,http
		Param   *json.RawMessage                `json:"Param"`
	}{}

	err := json.Unmarshal(bodyBuf[:n],&interfaceInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	//log.Printf("info %+v\n",interfaceInfo)
	//log.Printf("type %+v\n",reflect.TypeOf(interfaceInfo.Param))
	//switch t:= interfaceInfo.Param.(type){
	//case device.SerialInterfaceParam:
	//	log.Printf("param %+v\n",t)
	//	device.CommInterfaceList.AddCommInterface(t.Name,t.Type,t.Param)
	//default:
	//	aParam.Code = "1"
	//	aParam.Message = "param is noexist"
	//	aParam.Data = ""
	//	sJson,_ := json.Marshal(aParam)
	//	context.String(http.StatusOK,string(sJson))
	//}



    var msg json.RawMessage
	switch interfaceInfo.Type{
	case "serial":
		serial := &device.SerialInterfaceParam{}
		err := json.Unmarshal(msg,serial)
		if err != nil{
			log.Println("CommunicationSerialInterface json unMarshall err,",err)
			break
		}
		log.Printf("type %+v\n",serial)
		//device.CommInterfaceList.AddCommInterface(serial.Name,serial.Type,serial.Param)
	case "tcp":
	}

	//device.WriteCommInterfaceListToJson()

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiGetCommInterface(context *gin.Context){

	type CommunicationInterfaceTemplate struct{
		Name 	string										`json:"Name"`			//接口名称
		Type    string          							`json:"Type"`			//接口类型,比如serial,tcp,udp,http
		Param   interface{} 							    `json:"Param"`			//接口参数
	}

	type CommunicationInterfaceManageTemplate struct{
		InterfaceCnt int
		InterfaceMap []CommunicationInterfaceTemplate
	}

	aParam := &struct{
		Code 	string
		Message string
		Data    CommunicationInterfaceManageTemplate
	}{}

	CommunicationInterfaceManage := CommunicationInterfaceManageTemplate{
		InterfaceCnt: 0,
		InterfaceMap: make([]CommunicationInterfaceTemplate,0),
	}

	aParam.Code = "0"
	aParam.Message = ""
	for _,v := range device.CommunicationSerialInterfaceList.SerialInterfaceMap{

		CommunicationInterface := CommunicationInterfaceTemplate{
			Name: v.Name,
			Type: v.Type,
			Param: v.Param,
		}
		CommunicationInterfaceManage.InterfaceCnt++
		CommunicationInterfaceManage.InterfaceMap = append(CommunicationInterfaceManage.InterfaceMap,
														CommunicationInterface)
	}
	aParam.Data = CommunicationInterfaceManage

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

func apiAddCommSerialInterface(context *gin.Context){

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
	//fmt.Println(string(bodyBuf[:n]))

	interfaceInfo := struct{
		Name 	string							`json:"Name"`			//接口名称
		Type    string          				`json:"Type"`			//接口类型,比如serial,tcp,udp,http
		Param   device.SerialInterfaceParam     `json:"Param"`
	}{}

	err := json.Unmarshal(bodyBuf[:n],&interfaceInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	SerialInterface := device.CommunicationSerialInterface{
		Param:interfaceInfo.Param,
		CommunicationTemplate:device.CommunicationTemplate{
			Name:interfaceInfo.Name,
			Type:interfaceInfo.Type,
		},
	}

	device.CommunicationSerialInterfaceList.SerialInterfaceMap = append(device.CommunicationSerialInterfaceList.SerialInterfaceMap,
															SerialInterface)
	device.WriteCommSerialInterfaceListToJson()

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiModifyCommSerialInterface(context *gin.Context){

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

	interfaceInfo := struct{
		Name 	string							`json:"Name"`			//接口名称
		Type    string          				`json:"Type"`			//接口类型,比如serial,tcp,udp,http
		Param   device.SerialInterfaceParam     `json:"Param"`
	}{}

	err := json.Unmarshal(bodyBuf[:n],&interfaceInfo)
	if err != nil {
		fmt.Println("CommSerialInterface json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	for k,v := range device.CommunicationSerialInterfaceList.SerialInterfaceMap{
		//判断名称是否一致
		if v.Name == interfaceInfo.Name{
			device.CommunicationSerialInterfaceList.SerialInterfaceMap[k].Type = interfaceInfo.Type
			device.CommunicationSerialInterfaceList.SerialInterfaceMap[k].Param = interfaceInfo.Param
			device.WriteCommSerialInterfaceListToJson()

			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = ""
			sJson,_ := json.Marshal(aParam)
			context.String(http.StatusOK,string(sJson))
			return
		}
	}

	aParam.Code = "1"
	aParam.Message = "addr is not exist"
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}

func apiDeleteCommSerialInterface(context *gin.Context){

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

	interfaceInfo := struct{
		Name 	string							`json:"Name"`			//接口名称
		Type    string          				`json:"Type"`			//接口类型,比如serial,tcp,udp,http
	}{}

	err := json.Unmarshal(bodyBuf[:n],&interfaceInfo)
	if err != nil {
		fmt.Println("CommSerialInterface json unMarshall err,",err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson,_ := json.Marshal(aParam)
		context.String(http.StatusOK,string(sJson))
		return
	}

	for k,v := range device.CommunicationSerialInterfaceList.SerialInterfaceMap{
		//判断名称是否一致
		if v.Name == interfaceInfo.Name{

			device.CommunicationSerialInterfaceList.SerialInterfaceMap = append(device.CommunicationSerialInterfaceList.SerialInterfaceMap[:k],
																		device.CommunicationSerialInterfaceList.SerialInterfaceMap[k+1:]...)
			device.WriteCommSerialInterfaceListToJson()

			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = ""
			sJson,_ := json.Marshal(aParam)
			context.String(http.StatusOK,string(sJson))
			return
		}
	}

	aParam.Code = "1"
	aParam.Message = "addr is not exist"
	aParam.Data = ""

	sJson,_ := json.Marshal(aParam)
	context.String(http.StatusOK,string(sJson))
}
