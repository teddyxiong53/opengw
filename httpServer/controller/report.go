package controller

import (
	"encoding/json"
	"fmt"
	"goAdapter/pkg/csv"
	"goAdapter/pkg/mylog"
	"goAdapter/report/mqttAliyun"
	mqttEmqx "goAdapter/report/mqttEMQX"
	"goAdapter/report/mqttHuawei"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func SetReportGWParam(context *gin.Context) {

	type ReportServiceTemplate struct {
		ServiceName string
		IP          string
		Port        string
		ReportTime  int
		Protocol    string
		Param       interface{}
	}

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

	var Param json.RawMessage
	param := ReportServiceTemplate{
		Param: &Param,
	}

	err := json.Unmarshal(bodyBuf[:n], &param)
	if err != nil {
		fmt.Println("param json unMarshall err,", err)

		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)

		context.String(http.StatusOK, string(sJson))
		return
	}

	switch param.Protocol {
	case "Aliyun.MQTT":
		ReportServiceGWParamAliyun := mqttAliyun.ReportServiceGWParamAliyunTemplate{}
		if err := json.Unmarshal(bodyBuf[:n], &ReportServiceGWParamAliyun); err != nil {
			fmt.Println("ReportServiceGWParamAliyun json unMarshall err,", err)
		}
		mqttAliyun.ReportServiceParamListAliyun.AddReportService(ReportServiceGWParamAliyun)
	case "EMQX.MQTT":
		ReportServiceGWParamEmqx := mqttEmqx.ReportServiceGWParamEmqxTemplate{}
		if err := json.Unmarshal(bodyBuf[:n], &ReportServiceGWParamEmqx); err != nil {
			fmt.Println("ReportServiceGWParamEmqx json unMarshall err,", err)
		}
		mqttEmqx.ReportServiceParamListEmqx.AddReportService(ReportServiceGWParamEmqx)
	case "Huawei.MQTT":
		ReportServiceGWParamHuawei := mqttHuawei.ReportServiceGWParamHuaweiTemplate{}
		if err := json.Unmarshal(bodyBuf[:n], &ReportServiceGWParamHuawei); err != nil {
			fmt.Println("ReportServiceGWParamAliyun json unMarshall err,", err)
		}
		mqttHuawei.ReportServiceParamListHuawei.AddReportService(ReportServiceGWParamHuawei)
	default:
		mylog.Logger.Errorf("unknown param.Protocol")
		aParam.Code = "1"
		aParam.Message = "unknown param.Protocol"
		aParam.Data = ""
		sJson, _ := json.Marshal(aParam)

		context.String(http.StatusOK, string(sJson))
		return
	}

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = ""
	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

func GetReportGWParam(context *gin.Context) {

	type ReportServiceTemplate struct {
		ServiceName string
		IP          string
		Port        string
		ReportTime  int
		CommStatus  string
		Protocol    string
		Param       interface{}
	}

	aParam := struct {
		Code    string                  `json:"Code"`
		Message string                  `json:"Message"`
		Data    []ReportServiceTemplate `json:"Data"`
	}{
		Data: make([]ReportServiceTemplate, 0),
	}

	aParam.Code = "0"
	aParam.Message = ""

	for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {

		ReportService := ReportServiceTemplate{}
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = v.GWParam.Param
		ReportService.CommStatus = v.GWParam.ReportStatus

		aParam.Data = append(aParam.Data, ReportService)
	}

	for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {

		ReportService := ReportServiceTemplate{}
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = v.GWParam.Param
		ReportService.CommStatus = v.GWParam.ReportStatus

		aParam.Data = append(aParam.Data, ReportService)
	}

	for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {

		ReportService := ReportServiceTemplate{}
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = v.GWParam.Param
		ReportService.CommStatus = v.GWParam.ReportStatus

		aParam.Data = append(aParam.Data, ReportService)
	}

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

func DeleteReportGWParam(context *gin.Context) {

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

	param := struct {
		ServiceName string
	}{}

	err := json.Unmarshal(bodyBuf[:n], &param)
	if err != nil {
		fmt.Println("param json unMarshall err,", err)

		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)

		context.String(http.StatusOK, string(sJson))
		return
	}

	//查看Aliyun
	for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
		if v.GWParam.ServiceName == param.ServiceName {
			mqttAliyun.ReportServiceParamListAliyun.DeleteReportService(param.ServiceName)

			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = ""
			sJson, _ := json.Marshal(aParam)

			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	//查看Emqx
	for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
		if v.GWParam.ServiceName == param.ServiceName {
			mqttEmqx.ReportServiceParamListEmqx.DeleteReportService(param.ServiceName)

			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = ""
			sJson, _ := json.Marshal(aParam)

			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {
		if v.GWParam.ServiceName == param.ServiceName {
			mqttHuawei.ReportServiceParamListHuawei.DeleteReportService(param.ServiceName)

			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = ""
			sJson, _ := json.Marshal(aParam)

			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	aParam.Code = "1"
	aParam.Message = "serviceName is not exist"
	aParam.Data = ""
	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

func SetReportNodeWParam(context *gin.Context) {

	type ReportServiceNodeTemplate struct {
		ServiceName       string
		CollInterfaceName string
		Addr              string
		Name              string
		CommStatus        string
		ReportStatus      string
		Protocol          string
		Param             interface{}
	}

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

	var Param json.RawMessage
	param := ReportServiceNodeTemplate{
		Param: &Param,
	}

	err := json.Unmarshal(bodyBuf[:n], &param)
	if err != nil {
		fmt.Println("param json unMarshall err,", err)

		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)

		context.String(http.StatusOK, string(sJson))
		return
	}

	switch param.Protocol {
	case "Aliyun.MQTT":
		ReportServiceNodeParamAliyun := mqttAliyun.ReportServiceNodeParamAliyunTemplate{}
		if err := json.Unmarshal(bodyBuf[:n], &ReportServiceNodeParamAliyun); err != nil {
			mylog.Logger.Errorf("ReportServiceNodeParamAliyun json unMarshall err,", err)
		}
		for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				v.AddReportNode(ReportServiceNodeParamAliyun)
			}
		}
		mylog.Logger.Debugf("ParamListAliyun %v", mqttAliyun.ReportServiceParamListAliyun.ServiceList)
	case "EMQX.MQTT":
		ReportServiceNodeParamEmqx := mqttEmqx.ReportServiceNodeParamEmqxTemplate{}
		if err := json.Unmarshal(bodyBuf[:n], &ReportServiceNodeParamEmqx); err != nil {
			mylog.Logger.Errorf("ReportServiceNodeParamEmqx json unMarshall err,", err)
		}
		for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				v.AddReportNode(ReportServiceNodeParamEmqx)
			}
		}
		mylog.Logger.Debugf("ParamListAliyun %v", mqttAliyun.ReportServiceParamListAliyun.ServiceList)
	case "Huawei.MQTT":
		ReportServiceNodeParamHuawei := mqttHuawei.ReportServiceNodeParamHuaweiTemplate{}
		if err := json.Unmarshal(bodyBuf[:n], &ReportServiceNodeParamHuawei); err != nil {
			mylog.Logger.Errorf("ReportServiceNodeParamHuawei json unMarshall err,", err)
		}
		for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				v.AddReportNode(ReportServiceNodeParamHuawei)
			}
		}
	default:
		mylog.Logger.Errorf("unknown param.Protocol")
		aParam.Code = "1"
		aParam.Message = "unknown param.Protocol"
		aParam.Data = ""
		sJson, _ := json.Marshal(aParam)

		context.String(http.StatusOK, string(sJson))
		return
	}

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = ""
	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

func BatchAddReportNodeParam(context *gin.Context) {

	aParam := struct {
		Code    string
		Message string
		Data    string
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	// 获取文件头
	file, err := context.FormFile("file")
	if err != nil {
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileName := exeCurDir + "/config/" + file.Filename

	//保存文件到服务器本地
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		aParam.Code = "1"
		aParam.Message = "save File Error"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	result := csv.LoadCsvCfg(fileName, 1)
	if result == nil {
		return
	}

	for _, record := range result.Records {
		protocol := record.GetString("Protocol")
		mylog.Logger.Debugf("protocal %v\n", protocol)
		switch protocol {
		case "Aliyun.MQTT":
			{
				ReportServiceNodeParamAliyun := mqttAliyun.ReportServiceNodeParamAliyunTemplate{}
				ReportServiceNodeParamAliyun.ServiceName = record.GetString("ServiceName")
				ReportServiceNodeParamAliyun.CollInterfaceName = record.GetString("CollInterfaceName")
				ReportServiceNodeParamAliyun.Name = record.GetString("Name")
				ReportServiceNodeParamAliyun.Addr = record.GetString("Addr")
				ReportServiceNodeParamAliyun.Protocol = record.GetString("Protocol")
				ReportServiceNodeParamAliyun.Param.ProductKey = record.GetString("ProductKey")
				ReportServiceNodeParamAliyun.Param.DeviceName = record.GetString("DeviceName")
				ReportServiceNodeParamAliyun.Param.DeviceSecret = record.GetString("DeviceSecret")

				for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
					if v.GWParam.ServiceName == ReportServiceNodeParamAliyun.ServiceName {
						v.AddReportNode(ReportServiceNodeParamAliyun)
					}
				}
			}
		}
	}

	aParam.Code = "0"
	aParam.Message = ""

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func GetReportNodeWParam(context *gin.Context) {

	type ReportServiceNodeTemplate struct {
		ServiceName       string
		CollInterfaceName string
		Name              string
		Addr              string
		CommStatus        string
		ReportStatus      string
		Protocol          string
		Param             interface{}
	}

	aParam := struct {
		Code    string                      `json:"Code"`
		Message string                      `json:"Message"`
		Data    []ReportServiceNodeTemplate `json:"Data"`
	}{
		Data: make([]ReportServiceNodeTemplate, 0),
	}

	ServiceName := context.Query("ServiceName")

	for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			ReportServiceNode := ReportServiceNodeTemplate{}
			for _, d := range v.NodeList {
				ReportServiceNode.ServiceName = d.ServiceName
				ReportServiceNode.CollInterfaceName = d.CollInterfaceName
				ReportServiceNode.Name = d.Name
				ReportServiceNode.Addr = d.Addr
				ReportServiceNode.Protocol = d.Protocol
				ReportServiceNode.CommStatus = d.CommStatus
				ReportServiceNode.ReportStatus = d.ReportStatus
				ReportServiceNode.Param = d.Param
				aParam.Data = append(aParam.Data, ReportServiceNode)
			}
			aParam.Code = "0"
			aParam.Message = ""
			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			ReportServiceNode := ReportServiceNodeTemplate{}
			for _, d := range v.NodeList {
				ReportServiceNode.ServiceName = d.ServiceName
				ReportServiceNode.CollInterfaceName = d.CollInterfaceName
				ReportServiceNode.Name = d.Name
				ReportServiceNode.Addr = d.Addr
				ReportServiceNode.Protocol = d.Protocol
				ReportServiceNode.CommStatus = d.CommStatus
				ReportServiceNode.ReportStatus = d.ReportStatus
				ReportServiceNode.Param = d.Param
				aParam.Data = append(aParam.Data, ReportServiceNode)
			}
			aParam.Code = "0"
			aParam.Message = ""
			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {
		if v.GWParam.ServiceName == ServiceName {
			ReportServiceNode := ReportServiceNodeTemplate{}
			for _, d := range v.NodeList {
				ReportServiceNode.ServiceName = d.ServiceName
				ReportServiceNode.CollInterfaceName = d.CollInterfaceName
				ReportServiceNode.Name = d.Name
				ReportServiceNode.Addr = d.Addr
				ReportServiceNode.Protocol = d.Protocol
				ReportServiceNode.CommStatus = d.CommStatus
				ReportServiceNode.ReportStatus = d.ReportStatus
				ReportServiceNode.Param = d.Param
				aParam.Data = append(aParam.Data, ReportServiceNode)
			}
			aParam.Code = "0"
			aParam.Message = ""
			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}

	aParam.Code = "1"
	aParam.Message = "ServiceName is not correct"
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func DeleteReportNodeWParam(context *gin.Context) {

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

	param := struct {
		ServiceName       string
		CollInterfaceName string
		Addr              string
	}{}

	err := json.Unmarshal(bodyBuf[:n], &param)
	if err != nil {
		fmt.Println("param json unMarshall err,", err)

		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)

		context.String(http.StatusOK, string(sJson))
		return
	}

	//查看Aliyun
	for _, v := range mqttAliyun.ReportServiceParamListAliyun.ServiceList {
		for _, n := range v.NodeList {
			if (n.ServiceName == param.ServiceName) &&
				(n.CollInterfaceName == param.CollInterfaceName) &&
				(n.Addr == param.Addr) {

				v.DeleteReportNode(param.Addr)

				aParam.Code = "0"
				aParam.Message = ""
				aParam.Data = ""
				sJson, _ := json.Marshal(aParam)

				context.String(http.StatusOK, string(sJson))
				return
			}
		}
	}

	//查看Emqx
	for _, v := range mqttEmqx.ReportServiceParamListEmqx.ServiceList {
		for _, n := range v.NodeList {
			if (n.ServiceName == param.ServiceName) &&
				(n.CollInterfaceName == param.CollInterfaceName) &&
				(n.Addr == param.Addr) {

				v.DeleteReportNode(param.Addr)

				aParam.Code = "0"
				aParam.Message = ""
				aParam.Data = ""
				sJson, _ := json.Marshal(aParam)

				context.String(http.StatusOK, string(sJson))
				return
			}
		}
	}

	for _, v := range mqttHuawei.ReportServiceParamListHuawei.ServiceList {
		for _, n := range v.NodeList {
			if (n.ServiceName == param.ServiceName) &&
				(n.CollInterfaceName == param.CollInterfaceName) &&
				(n.Addr == param.Addr) {

				v.DeleteReportNode(param.Addr)

				aParam.Code = "0"
				aParam.Message = ""
				aParam.Data = ""
				sJson, _ := json.Marshal(aParam)

				context.String(http.StatusOK, string(sJson))
				return
			}
		}
	}

	aParam.Code = "1"
	aParam.Message = "node is not exist"
	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}
