package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goAdapter/api"
	"goAdapter/device"
	"log"
	"net/http"
)

func apiAddInterface(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	fmt.Println(string(bodyBuf[:n]))

	interfaceInfo := &struct {
		CollectInterfaceName 	string	`json:"CollInterfaceName"`	 //采集接口名字
		CommInterfaceName 		string  `json:"CommInterfaceName"`   //通信接口名字
		PollPeriod        		int     `json:"PollPeriod"`
		OfflinePeriod     		int     `json:"OfflinePeriod"`
	}{}

	err := json.Unmarshal(bodyBuf[:n], interfaceInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	device.CollectInterfaceMap = append(device.CollectInterfaceMap,device.NewCollectInterface(interfaceInfo.CollectInterfaceName,
		interfaceInfo.CommInterfaceName,
		interfaceInfo.PollPeriod,
		interfaceInfo.OfflinePeriod,0))

	device.WriteCollectInterfaceManageToJson()

	aParam.Code = "0"
	aParam.Data = ""

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiModifyInterface(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	fmt.Println(string(bodyBuf[:n]))

	interfaceInfo := &struct {
		CollectInterfaceName 	string	`json:"CollInterfaceName"`	 	//采集接口名字
		CommInterfaceName 		string  `json:"CommInterfaceName"`       //通信接口名字
		PollPeriod    int
		OfflinePeriod int
	}{}

	err := json.Unmarshal(bodyBuf[:n], interfaceInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	for _,v := range device.CollectInterfaceMap{
		if v.CollInterfaceName == interfaceInfo.CollectInterfaceName{
			v.CommInterfaceName = interfaceInfo.CommInterfaceName
			v.PollPeriod = interfaceInfo.PollPeriod
			v.OfflinePeriod = interfaceInfo.OfflinePeriod

			aParam.Code = "0"
			aParam.Data = ""

			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	aParam.Code = "1"
	aParam.Message = "collInterface is not exist"
	aParam.Data = ""

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiDeleteInterface(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	fmt.Println(string(bodyBuf[:n]))

	interfaceInfo := &struct {
		CollectInterfaceName 	string	`json:"CollInterfaceName"`	 //采集接口名字
	}{}

	err := json.Unmarshal(bodyBuf[:n], interfaceInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	for k,v := range device.CollectInterfaceMap{
		if v.CollInterfaceName == interfaceInfo.CollectInterfaceName{

			device.CollectInterfaceMap = append(device.CollectInterfaceMap[:k],device.CollectInterfaceMap[k+1:]...)

			aParam.Code = "0"
			aParam.Data = ""

			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	aParam.Code = "1"
	aParam.Message = "collInterface is not exist"
	aParam.Data = ""

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiGetInterfaceInfo(context *gin.Context) {

	sName := context.Query("CollInterfaceName")

	aParam := &struct {
		Code    string
		Message string
		Data    device.CollectInterfaceTemplate
	}{}

	for k,v := range device.CollectInterfaceMap{
		if v.CollInterfaceName == sName{

			aParam.Code = "0"
			aParam.Message = ""

			aParam.Data = *device.CollectInterfaceMap[k]

			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	aParam.Code = "1"
	aParam.Message = "interface is not exist"

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiGetAllInterfaceInfo(context *gin.Context) {

	type InterfaceParamTemplate struct{
		CollInterfaceName   string          		`json:"CollInterfaceName"`   	//采集接口
		CommInterfaceName   string					`json:"CommInterfaceName"`   	//通信接口
		PollPeriod          int                   	`json:"PollPeriod"`    			//采集周期
		OfflinePeriod       int                   	`json:"OfflinePeriod"` 			//离线超时周期
		DeviceNodeCnt       int                   	`json:"DeviceNodeCnt"` 			//设备数量
		DeviceNodeOnlineCnt int             		`json:"DeviceNodeOnlineCnt"`	//设备在线数量
	}

	aParam := &struct {
		Code    string
		Message string
		Data    []InterfaceParamTemplate
	}{}

	aParam.Data = make([]InterfaceParamTemplate,0)

	aParam.Code = "0"
	aParam.Message = ""
	for _,v := range device.CollectInterfaceMap{

		Param := InterfaceParamTemplate{
			CollInterfaceName:v.CollInterfaceName,
			CommInterfaceName:v.CommInterfaceName,
			PollPeriod: v.PollPeriod,
			OfflinePeriod: v.OfflinePeriod,
			DeviceNodeCnt: v.DeviceNodeCnt,
			DeviceNodeOnlineCnt: v.DeviceNodeOnlineCnt,
		}
		aParam.Data = append(aParam.Data,Param)
	}

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

func apiAddNode(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	fmt.Println(string(bodyBuf[:n]))

	nodeInfo := &struct {
		InterfaceName string `json:"CollInterfaceName"`
		DAddr         string `json:"Addr"`
		DType         string `json:"Type"`
		DName         string `json:"Name"`
	}{}

	err := json.Unmarshal(bodyBuf[:n], nodeInfo)
	if err != nil {
		fmt.Println("nodeInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	var status bool
	for _,v := range device.CollectInterfaceMap{
		if v.CollInterfaceName == nodeInfo.InterfaceName{

			status,aParam.Message = v.AddDeviceNode(nodeInfo.DName,nodeInfo.DType, nodeInfo.DAddr)
			if status == true {
				device.WriteCollectInterfaceManageToJson()

				aParam.Code = "0"
				aParam.Data = ""
			} else {
				aParam.Code = "1"
				aParam.Data = ""
			}
			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	aParam.Code = "1"
	aParam.Data = ""
	aParam.Message = "interfaceName is not exist"

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiModifyNode(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	fmt.Println(string(bodyBuf[:n]))

	nodeInfo := &struct {
		InterfaceName string    `json:"CollInterfaceName"`
		DAddr         string 	`json:"Addr"`
		DType         string 	`json:"Type"`
		DName         string    `json:"Name"`
	}{}

	err := json.Unmarshal(bodyBuf[:n], nodeInfo)
	if err != nil {
		fmt.Println("nodeInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	var status bool
	for _,v := range device.CollectInterfaceMap{
		if v.CollInterfaceName == nodeInfo.InterfaceName{

			status,aParam.Message = v.AddDeviceNode(nodeInfo.DName,nodeInfo.DType, nodeInfo.DAddr)
			if status == true {
				device.WriteCollectInterfaceManageToJson()

				aParam.Code = "0"
				aParam.Data = ""
			} else {
				aParam.Code = "1"
				aParam.Data = ""
			}
			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	aParam.Code = "0"
	aParam.Data = ""

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiGetNode(context *gin.Context) {

	sName := context.Query("CollInterfaceName")
	sAddr := context.Query("Addr")

	aParam := &struct {
		Code    string
		Message string
		Data    device.DeviceNodeTemplate
	}{}

	for _, v := range device.CollectInterfaceMap {
		if v.CollInterfaceName == sName {
			for _, n := range v.DeviceNodeMap {
				if n.Addr == sAddr {
					aParam.Code = "0"
					aParam.Message = ""
					aParam.Data = *v.GetDeviceNode(sAddr)
					sJson, _ := json.Marshal(aParam)
					context.String(http.StatusOK, string(sJson))
					return
				}
			}
		}
	}

	aParam.Code = "1"
	aParam.Message = "node is no exist"
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiDeleteNode(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	fmt.Println(string(bodyBuf[:n]))

	nodeInfo := &struct {
		InterfaceName string    `json:"CollInterfaceName"`
		DAddr         string 	`json:"Addr"`
		DType         string 	`json:"Type"`
		DName         string    `json:"Name"`
	}{}

	err := json.Unmarshal(bodyBuf[:n], nodeInfo)
	if err != nil {
		fmt.Println("nodeInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	for _, v := range device.CollectInterfaceMap {
		if v.CollInterfaceName == nodeInfo.InterfaceName {
			for _, n := range v.DeviceNodeMap {
				if n.Addr == nodeInfo.DAddr {
					v.DeleteDeviceNode(nodeInfo.DName,nodeInfo.DAddr, nodeInfo.DType)
					device.WriteCollectInterfaceManageToJson()

					aParam.Code = "0"
					aParam.Message = ""
					sJson, _ := json.Marshal(aParam)
					context.String(http.StatusOK, string(sJson))
					return
				}
			}
		}
	}

	aParam.Code = "1"
	aParam.Data = ""
	aParam.Message = "addr is not exist"

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

/**
从缓存中获取设备变量
*/
func apiGetNodeVariableFromCache(context *gin.Context) {

	sName := context.Query("CollInterfaceName")
	sAddr := context.Query("Addr")

	aParam := &struct {
		Code    string
		Message string
		Data    []api.VariableTemplate
	}{}

	for _, v := range device.CollectInterfaceMap {
		if v.CollInterfaceName == sName {
			for _, v := range v.DeviceNodeMap {
				if v.Addr == sAddr {

					aParam.Code = "0"
					aParam.Message = ""
					aParam.Data = v.VariableMap
					sJson, _ := json.Marshal(aParam)
					context.String(http.StatusOK, string(sJson))
					return
				}
			}
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
func apiGetNodeVariableFromDevice(context *gin.Context) {

	//sName := context.Query("interfaceName")
	//sAddr := context.Query("addr")
	//
	//aParam := &struct {
	//	Code    string
	//	Message string
	//	Data    []api.VariableTemplate
	//}{}
	//
	//
	//for _,v := range device.CollectInterfaceMap {
	//	if v.CollInterfaceName == nodeInfo.InterfaceName {
	//
	//	}
	//}
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

func apiAddTemplate(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	fmt.Println(string(bodyBuf[:n]))

	typeInfo := &struct {
		TemplateName    string `json:"TemplateName"`    //模板名称
		TemplateType    string `json:"TemplateType"`    //模板型号
		TemplateMessage string `json:"TemplateMessage"` //备注信息
	}{}

	err := json.Unmarshal(bodyBuf[:n], typeInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	index := len(device.DeviceNodeTypeMap.DeviceNodeType)
	template := device.DeviceNodeTypeTemplate{
		TemplateName:    typeInfo.TemplateName,
		TemplateType:    typeInfo.TemplateType,
		TemplateID:      index,
		TemplateMessage: typeInfo.TemplateMessage,
	}

	device.DeviceNodeTypeMap.DeviceNodeType = append(device.DeviceNodeTypeMap.DeviceNodeType, template)

	device.WriteDeviceNodeTypeMapToJson()

	aParam.Code = "0"
	aParam.Data = ""

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiGetTemplate(context *gin.Context) {

	aParam := &struct {
		Code    string
		Message string
		Data    []device.DeviceNodeTypeTemplate
	}{}

	if len(device.DeviceNodeTypeMap.DeviceNodeType) > 0 {
		aParam.Code = "0"
		aParam.Message = ""
		aParam.Data = device.DeviceNodeTypeMap.DeviceNodeType
	} else {
		aParam.Code = "1"
		aParam.Message = "nodeTypeCnt is 0"
		aParam.Data = device.DeviceNodeTypeMap.DeviceNodeType
	}

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

func apiAddCommInterface(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	//fmt.Println(string(bodyBuf[:n]))

	interfaceInfo := struct {
		Name  string           `json:"Name"` //接口名称
		Type  string           `json:"Type"` //接口类型,比如serial,tcp,udp,http
		Param *json.RawMessage `json:"Param"`
	}{}

	err := json.Unmarshal(bodyBuf[:n], &interfaceInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
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
	switch interfaceInfo.Type {
	case "serial":
		serial := &device.SerialInterfaceParam{}
		err := json.Unmarshal(msg, serial)
		if err != nil {
			log.Println("CommunicationSerialInterface json unMarshall err,", err)
			break
		}
		log.Printf("type %+v\n", serial)
		//device.CommInterfaceList.AddCommInterface(serial.Name,serial.Type,serial.Param)
	case "tcp":
	}

	//device.WriteCommInterfaceListToJson()

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = ""

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiGetCommInterface(context *gin.Context) {

	type CommunicationInterfaceTemplate struct {
		Name  string      `json:"Name"`  //接口名称
		Type  string      `json:"Type"`  //接口类型,比如serial,tcp,udp,http
		Param interface{} `json:"Param"` //接口参数
	}

	type CommunicationInterfaceManageTemplate struct {
		InterfaceCnt int
		InterfaceMap []CommunicationInterfaceTemplate
	}

	aParam := &struct {
		Code    string
		Message string
		Data    CommunicationInterfaceManageTemplate
	}{}

	CommunicationInterfaceManage := CommunicationInterfaceManageTemplate{
		InterfaceCnt: 0,
		InterfaceMap: make([]CommunicationInterfaceTemplate, 0),
	}

	aParam.Code = "0"
	aParam.Message = ""
	for _, v := range device.CommunicationSerialMap {

		CommunicationInterface := CommunicationInterfaceTemplate{
			Name:  v.Name,
			Type:  v.Type,
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

func apiAddCommSerialInterface(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	//fmt.Println(string(bodyBuf[:n]))

	interfaceInfo := struct {
		Name  string                      `json:"Name"` //接口名称
		Type  string                      `json:"Type"` //接口类型,比如serial,tcp,udp,http
		Param device.SerialInterfaceParam `json:"Param"`
	}{}

	err := json.Unmarshal(bodyBuf[:n], &interfaceInfo)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	for _, v := range device.CommunicationSerialMap {
		//判断通信接口名称是否一致
		if (v.Name==interfaceInfo.Name) || (v.Param.Name == interfaceInfo.Param.Name){
			aParam.Code = "1"
			aParam.Message = "name is exist"
			aParam.Data = ""

			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	SerialInterface := device.CommunicationSerialTemplate{
		Param: interfaceInfo.Param,
		CommunicationTemplate: device.CommunicationTemplate{
			Name: interfaceInfo.Name,
			Type: interfaceInfo.Type,
		},
	}

	device.CommunicationSerialMap = append(device.CommunicationSerialMap,SerialInterface)
	device.WriteCommSerialInterfaceListToJson()

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = ""
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiModifyCommSerialInterface(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)

	interfaceInfo := struct {
		Name  string                      `json:"Name"` //接口名称
		Type  string                      `json:"Type"` //接口类型,比如serial,tcp,udp,http
		Param device.SerialInterfaceParam `json:"Param"`
	}{}

	err := json.Unmarshal(bodyBuf[:n], &interfaceInfo)
	if err != nil {
		fmt.Println("CommSerialInterface json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	for k, v := range device.CommunicationSerialMap {
		//判断名称是否一致
		if v.Name == interfaceInfo.Name {
			device.CommunicationSerialMap[k].Type = interfaceInfo.Type
			device.CommunicationSerialMap[k].Param = interfaceInfo.Param
			device.WriteCommSerialInterfaceListToJson()

			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = ""
			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	aParam.Code = "1"
	aParam.Message = "addr is not exist"
	aParam.Data = ""

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiDeleteCommSerialInterface(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)

	interfaceInfo := struct {
		Name string `json:"Name"` //接口名称
		Type string `json:"Type"` //接口类型,比如serial,tcp,udp,http
	}{}

	err := json.Unmarshal(bodyBuf[:n], &interfaceInfo)
	if err != nil {
		fmt.Println("CommSerialInterface json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	for k, v := range device.CommunicationSerialMap {
		//判断名称是否一致
		if v.Name == interfaceInfo.Name {

			device.CommunicationSerialMap = append(device.CommunicationSerialMap[:k],
				device.CommunicationSerialMap[k+1:]...)
			device.WriteCommSerialInterfaceListToJson()

			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = ""
			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	aParam.Code = "1"
	aParam.Message = "addr is not exist"
	aParam.Data = ""

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}
