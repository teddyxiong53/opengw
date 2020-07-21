package httpServer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"goAdapter/setting"
	"net/http"
)

func apiGetSerial(context *gin.Context) {

	aParam := struct {
		Code    string   `json:"Code"`
		Message string   `json:"Message"`
		Data    []string `json:"Data"`
	}{
		Code:    "0",
		Message: "",
		Data:    setting.SerialPortMap[:],
	}

	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}
