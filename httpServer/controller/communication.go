package controller

import (
	"encoding/json"
	"fmt"
	"goAdapter/device"
	"goAdapter/httpServer/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddCommInterface(context *gin.Context) {
	data, err := context.GetRawData()
	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}
	var Param json.RawMessage
	interfaceInfo := struct {
		Name  string           `json:"Name"` // 接口名称
		Type  string           `json:"Type"` // 接口类型,比如serial,TcpClient,udp,http
		Param *json.RawMessage `json:"Param"`
	}{
		Param: &Param,
	}

	err = json.Unmarshal(data, &interfaceInfo)
	if err != nil {
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}
	var willAdd device.CommunicationInterface
	switch typ := interfaceInfo.Type; typ {
	case device.SERIALTYPE:
		serial := &device.SerialInterfaceParam{}
		err = json.Unmarshal(Param, &serial)
		if err != nil {
			break
		}
		willAdd = &device.CommunicationSerialTemplate{
			Param: serial,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

	case device.TCPCLIENTTYPE:
		TcpClient := device.TcpClientInterfaceParam{}
		err = json.Unmarshal(Param, &TcpClient)
		if err != nil {
			break
		}
		willAdd = &device.CommunicationTcpClientTemplate{
			Param: &TcpClient,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

	case device.IOOUTTYPE:
		IoOut := device.IoOutInterfaceParam{}
		err = json.Unmarshal(Param, &IoOut)
		if err != nil {
			break
		}
		willAdd = &device.CommunicationIoOutTemplate{
			Param: &IoOut,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

	case device.IOINTYPE:
		IoIn := device.IoInInterfaceParam{}
		err = json.Unmarshal(Param, &IoIn)
		if err != nil {
			break
		}
		willAdd = &device.CommunicationIoInTemplate{
			Param: &IoIn,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}
	}

	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("json unmarshal error:%v", err),
		})
		return
	}
	for _, v := range device.CommunicationInterfaceMap {
		if v.Unique() == willAdd.Unique() {
			context.JSON(200, model.Response{
				Code:    "1",
				Message: fmt.Sprintf("%s is already exists", v.Unique()),
			})
			return
		}
	}

	device.CommunicationInterfaceMap[interfaceInfo.Name] = willAdd
	if err = device.WriteToJson(device.COMMJSON); err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, model.Response{
		Code: "0",
	})
}

func ModifyCommInterface(context *gin.Context) {
	var aParam = model.Response{
		Code: "0",
	}

	data, err := context.GetRawData()
	if err != nil {
		aParam.Code = "1"
		aParam.Message = err.Error()
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}
	var Param json.RawMessage
	interfaceInfo := struct {
		Name  string           `json:"Name"` // 接口名称
		Type  string           `json:"Type"` // 接口类型,比如serial,TcpClient,udp,http
		Param *json.RawMessage `json:"Param"`
	}{
		Param: &Param,
	}

	err = json.Unmarshal(data, &interfaceInfo)
	if err != nil {
		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	switch interfaceInfo.Type {
	case "LocalSerial":
		serial := &device.SerialInterfaceParam{}
		err = json.Unmarshal(Param, &serial)
		if err != nil {
			break
		}
		SerialInterface := &device.CommunicationSerialTemplate{
			Param: serial,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}
		device.CommunicationInterfaceMap[interfaceInfo.Name] = SerialInterface

	case "TcpClient":
		TcpClient := device.TcpClientInterfaceParam{}
		err = json.Unmarshal(Param, &TcpClient)
		if err != nil {
			break
		}

		TcpClientInterface := &device.CommunicationTcpClientTemplate{
			Param: &TcpClient,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		device.CommunicationInterfaceMap[interfaceInfo.Name] = TcpClientInterface
	case "IoOut":
		IoOut := device.IoOutInterfaceParam{}
		err = json.Unmarshal(Param, &IoOut)
		if err != nil {
			break
		}
		IoOutInterface := &device.CommunicationIoOutTemplate{
			Param: &IoOut,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		device.CommunicationInterfaceMap[interfaceInfo.Name] = IoOutInterface
	case "IoIn":
		IoIn := device.IoInInterfaceParam{}
		err = json.Unmarshal(Param, &IoIn)
		if err != nil {
			break
		}
		IoInInterface := &device.CommunicationIoInTemplate{
			Param: &IoIn,
			Name:  interfaceInfo.Name,
			Type:  interfaceInfo.Type,
		}

		device.CommunicationInterfaceMap[interfaceInfo.Name] = IoInInterface
	}

	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("json unmarshal error:%v", err),
		})
		return
	}
	if err := device.WriteToJson(device.COMMJSON); err != nil {
		aParam.Code = "1"
		aParam.Message = fmt.Sprintf("write to commjson error:%v", err)
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func DeleteCommInterface(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code: "0",
	}

	cName := context.Query("commInterface")
	_, ok := device.CommunicationInterfaceMap[cName]
	if !ok {
		aParam.Code = "1"
		aParam.Message = fmt.Sprintf("comminterface %s is not exists", cName)
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}
	delete(device.CommunicationInterfaceMap, cName)
	device.WriteToJson(device.COMMJSON)
	aParam.Message = fmt.Sprintf("delete comminterface %s success", cName)
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func GetCommInterface(context *gin.Context) {

	type CommunicationInterfaceTemplate struct {
		Name  string      `json:"Name"`  // 接口名称
		Type  string      `json:"Type"`  // 接口类型,比如serial,TcpClient,udp,http
		Param interface{} `json:"Param"` // 接口参数
	}

	type CommunicationInterfaceManageTemplate struct {
		InterfaceCnt int
		InterfaceMap []*CommunicationInterfaceTemplate
	}

	aParam := &struct {
		Code    string
		Message string
		Data    CommunicationInterfaceManageTemplate
	}{}

	CommunicationInterfaceManage := CommunicationInterfaceManageTemplate{
		InterfaceCnt: 0,
		InterfaceMap: make([]*CommunicationInterfaceTemplate, 0),
	}

	for _, v := range device.CommunicationInterfaceMap {
		item := &CommunicationInterfaceTemplate{
			Name:  v.GetName(),
			Type:  v.GetType(),
			Param: v.GetParam(),
		}
		CommunicationInterfaceManage.InterfaceMap = append(CommunicationInterfaceManage.InterfaceMap, item)
	}
	CommunicationInterfaceManage.InterfaceCnt = len(device.CommunicationInterfaceMap)

	aParam.Data = CommunicationInterfaceManage
	aParam.Code = "0"
	context.JSON(http.StatusOK, aParam)
}