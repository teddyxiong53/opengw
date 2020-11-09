package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"goAdapter/report"
)

func apiSetReportGWParam(context *gin.Context) {

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
		ReportServiceGWParamAliyun := report.ReportServiceGWParamAliyunTemplate{}
		if err := json.Unmarshal(bodyBuf[:n], &ReportServiceGWParamAliyun); err != nil {
			fmt.Println("ReportServiceGWParamAliyun json unMarshall err,", err)
		}
		report.ReportServiceParamListAliyun.AddReportService(ReportServiceGWParamAliyun)
	case "Emqx":

	case "Huawei":
	default:
		log.Printf("unknown param.Protocol")
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

func apiGetReportGWParam(context *gin.Context) {

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

	for _, v := range report.ReportServiceParamListAliyun.ServiceList {

		ReportService := ReportServiceTemplate{}
		ReportService.ServiceName = v.GWParam.ServiceName
		ReportService.IP = v.GWParam.IP
		ReportService.Port = v.GWParam.Port
		ReportService.ReportTime = v.GWParam.ReportTime
		ReportService.Protocol = v.GWParam.Protocol
		ReportService.Param = v.GWParam.Param
		ReportService.CommStatus = v.CommStatus

		aParam.Data = append(aParam.Data, ReportService)
	}

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

func apiDeleteReportGWParam(context *gin.Context) {

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
	for _, v := range report.ReportServiceParamListAliyun.ServiceList {
		if v.GWParam.ServiceName == param.ServiceName {
			report.ReportServiceParamListAliyun.DeleteReportService(param.ServiceName)

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

func apiSetReportNodeWParam(context *gin.Context) {

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
		ReportServiceNodeParamAliyun := report.ReportServiceNodeParamAliyunTemplate{}
		if err := json.Unmarshal(bodyBuf[:n], &ReportServiceNodeParamAliyun); err != nil {
			fmt.Println("ReportServiceNodeParamAliyun json unMarshall err,", err)
		}
		for _, v := range report.ReportServiceParamListAliyun.ServiceList {
			if v.GWParam.ServiceName == param.ServiceName {
				v.AddReportNode(ReportServiceNodeParamAliyun)
			}
		}
	case "Emqx.MQTT":

	case "Huawei":
	default:
		log.Printf("unknown param.Protocol")
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

func apiGetReportNodeWParam(context *gin.Context) {

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
		Data:make([]ReportServiceNodeTemplate,0),
	}

	ServiceName := context.Query("ServiceName")

	for _, v := range report.ReportServiceParamListAliyun.ServiceList {
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

				aParam.Code = "0"
				aParam.Message = ""
				sJson, _ := json.Marshal(aParam)
				context.String(http.StatusOK, string(sJson))
				return
			}
		}
	}

	aParam.Code = "1"
	aParam.Message = ""
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiDeleteReportNodeWParam(context *gin.Context) {

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
	for _, v := range report.ReportServiceParamListAliyun.ServiceList {
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
	aParam.Data = ""
	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}
