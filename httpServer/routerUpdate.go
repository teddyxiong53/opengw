package httpServer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func apiUpdatePlugin(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
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
	// 获取文件名
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/plugin"
	fileName := fileDir + "/" + file.Filename
	log.Println(fileName)

	if fileExist(fileDir)==false{
		os.MkdirAll(fileDir, os.ModePerm)
	}

	//保存文件到服务器本地
	if err := context.SaveUploadedFile(file, fileName); err != nil {

		log.Println(err)
		aParam.Code = "1"
		aParam.Message = "save error"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	// 上传文件到指定的路径
	//context.SaveUploadedFile(file, "/opt/ibox/")

	aParam.Code = "0"
	aParam.Message = "save sucess"

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}