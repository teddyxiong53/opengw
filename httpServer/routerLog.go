package httpServer

import (
	"encoding/json"
	"fmt"
	"goAdapter/setting"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func apiGetLogParam(context *gin.Context) {

	aParam := &struct {
		Code    string
		Message string
		Data    setting.LoggerParamTemplate
	}{}

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = *setting.LoggerParam

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
	return
}

func apiSetLogParam(context *gin.Context) {

	aParam := &struct {
		Code    string
		Message string
		Data    string
	}{}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	fmt.Println(string(bodyBuf[:n]))

	logParam := &setting.LoggerParamTemplate{}

	err := json.Unmarshal(bodyBuf[:n], logParam)
	if err != nil {
		fmt.Println("logParam json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	setting.LoggerParam.WriteParamToJson(logParam)

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = ""

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
	return
}

func apiGetLogFilesInfo(context *gin.Context) {

	aParam := &struct {
		Code    string
		Message string
		Data    []string
	}{}

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	//遍历json和so文件
	path := exeCurDir + "/log"
	fileNameMap := make([]string, 0)
	fileInfoMap, err := ioutil.ReadDir(path)
	if err != nil {
		setting.Logger.Errorf("readLogFileDir err %v", err)
		aParam.Code = "1"
		aParam.Message = "readLogFileDir err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}
	for _, v := range fileInfoMap {
		fileNameMap = append(fileNameMap, v.Name())
	}

	aParam.Code = "0"
	aParam.Message = "readLogFilesInfo ok"
	aParam.Data = fileNameMap

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func apiDeleteLogFile(context *gin.Context) {
	aParam := &struct {
		Code    string
		Message string
		Data    string
	}{}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	fmt.Println(string(bodyBuf[:n]))

	filesName := &struct {
		Name []string
	}{}

	err := json.Unmarshal(bodyBuf[:n], filesName)
	if err != nil {
		fmt.Println("filesName json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "filesName json unMarshall err"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	//遍历log文件夹下的文件
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := exeCurDir + "/log/"
	fileInfoMap, err := ioutil.ReadDir(path)
	if err != nil {
		setting.Logger.Errorf("readLogFileDir err %v", err)
		aParam.Code = "1"
		aParam.Message = "readLogFileDir err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}
	for _, v := range filesName.Name {
		for _, f := range fileInfoMap {
			if v == f.Name() {
				err := os.Remove(path + f.Name())
				if err != nil {
					aParam.Code = "1"
					aParam.Message = "remove logFile err"
					sJson, _ := json.Marshal(aParam)
					context.String(http.StatusOK, string(sJson))
					return
				}
			}
		}
	}

	aParam.Code = "0"
	aParam.Message = "remove logFile ok"
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
	return
}

func apiGetLogFile(context *gin.Context) {

	aParam := &struct {
		Code    string
		Message string
		Data    string
	}{}

	fileName := context.Query("fileName")

	//遍历log文件夹下的文件
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	path := exeCurDir + "/log/"
	fileInfoMap, err := ioutil.ReadDir(path)
	if err != nil {
		setting.Logger.Errorf("readLogFileDir err %v", err)
		aParam.Code = "1"
		aParam.Message = "readLogFileDir err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	index := -1
	for k, f := range fileInfoMap {
		if fileName == f.Name() {
			index = k
			//返回文件流
			context.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
			context.File(path + fileName) //返回文件路径，自动调用http.ServeFile方法
		}
	}

	if index == -1 {
		aParam.Code = "1"
		aParam.Message = "file is not exist"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}
}
