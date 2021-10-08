/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-10-07 11:41:29
@FilePath: /goAdapter-Raw/httpServer/controller/commessage.go
*/
package controller

import (
	"fmt"
	"goAdapter/device"
	"goAdapter/httpServer/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCommMessage(context *gin.Context) {
	interfaceName := struct {
		CollInterfaceNames []string `json:"CollInterfaceNames"` //接口名称
	}{}

	commMessageMap := make([]*device.CommunicationMessageTemplate, 0)

	err := context.ShouldBindJSON(&interfaceName)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,", err)

		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "json unMarshall err",
		})
		return
	}
	tmps := device.CollectInterfaceMap.GetAll()
	for _, name := range interfaceName.CollInterfaceNames {
		for _, v := range tmps {
			if name == v.CollInterfaceName {
				commMessageMap = append(commMessageMap, v.CommMessage...)
				// 清空map
				v.CommMessage = v.CommMessage[0:0]
			}
		}
	}
	context.JSON(http.StatusOK, model.Response{
		Code: "0",
		Data: commMessageMap,
	})
}
