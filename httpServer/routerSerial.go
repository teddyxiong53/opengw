package httpServer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"goAdapter/setting"
	"net/http"
)

func apiGetSerial(context *gin.Context) {

	type SerialPortNameTemplate struct{
		Name string			`json:"Name"`
	}

	aParam := struct {
		Code    string   							`json:"Code"`
		Message string   							`json:"Message"`
		Data    []SerialPortNameTemplate 			`json:"Data"`
	}{
		Code:    "0",
		Message: "",
		Data:    make([]SerialPortNameTemplate,0),
	}

	SerialPortName := SerialPortNameTemplate{}
	for _,v := range setting.SerialPortNameTemplateMap.Name{
		SerialPortName.Name = v
		aParam.Data = append(aParam.Data,SerialPortName)
	}

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}
