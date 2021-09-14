/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 19:24:41
@FilePath: /goAdapter-Raw/httpServer/controller/serial.go
*/
package controller

import (
	"net/http"

	"goAdapter/httpServer/model"
	"goAdapter/pkg/serial"

	"github.com/gin-gonic/gin"
)

func GetSerial(context *gin.Context) {
	type SerialPortNameTemplate struct {
		Name string `json:"Name"`
	}

	data := make([]SerialPortNameTemplate, 0)

	SerialPortName := SerialPortNameTemplate{}
	for _, v := range serial.SerialPortNameTemplateMap.Name {
		SerialPortName.Name = v
		data = append(data, SerialPortName)
	}

	context.JSON(http.StatusOK, model.Response{
		Code:    "0",
		Message: "",
		Data:    data,
	})
}
