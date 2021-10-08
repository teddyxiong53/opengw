package controller

import (
	"encoding/json"
	"fmt"
	"goAdapter/device"
	"goAdapter/httpServer/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leandro-lugaresi/hub"
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
	if !device.CommunicationInterfaceMap.Compare(willAdd) {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("%s is already exists", willAdd.Unique()),
		})
		return
	}

	device.CommunicationInterfaceMap.Add(willAdd)
	device.CommunicationInterfaceMap.Publish(device.CommAdd, hub.Fields{
		"Name": willAdd.GetName(),
	})

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
		aParam.Message = "json unMarshal err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	var ok bool
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
		ok = device.CommunicationInterfaceMap.Update(SerialInterface)

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

		ok = device.CommunicationInterfaceMap.Update(TcpClientInterface)
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

		ok = device.CommunicationInterfaceMap.Update(IoOutInterface)
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

		ok = device.CommunicationInterfaceMap.Update(IoInInterface)
	}

	if err != nil {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintf("json unmarshal error:%v", err),
		})
		return
	}

	if !ok {
		context.JSON(200, model.Response{
			Code:    "1",
			Message: fmt.Sprintln("modify comm interface  error"),
		})
		return
	}

	// if err := device.WriteToJsonFile(device.COMMJSON); err != nil {
	// 	aParam.Code = "1"
	// 	aParam.Message = fmt.Sprintf("write to commjson error:%v", err)
	// 	sJson, _ := json.Marshal(aParam)
	// 	context.String(http.StatusOK, string(sJson))
	// 	return
	// }
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
	ok := device.CommunicationInterfaceMap.Delete(cName)
	if !ok {
		aParam.Code = "1"
		aParam.Message = fmt.Sprintf("comminterface %s is not exists", cName)
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}
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

	comms := device.CommunicationInterfaceMap.GetAll()
	for _, v := range comms {
		item := &CommunicationInterfaceTemplate{
			Name:  v.GetName(),
			Type:  v.GetType(),
			Param: v.GetParam(),
		}
		CommunicationInterfaceManage.InterfaceMap = append(CommunicationInterfaceManage.InterfaceMap, item)
	}
	CommunicationInterfaceManage.InterfaceCnt = len(comms)

	aParam.Data = CommunicationInterfaceManage
	aParam.Code = "0"
	context.JSON(http.StatusOK, aParam)
}
