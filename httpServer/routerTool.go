package httpServer

import (
	"fmt"
	"goAdapter/device"
	"net/http"

	"github.com/gin-gonic/gin"
)

func apiGetCommMessage(context *gin.Context) {
	interfaceName := struct {
		CollInterfaceNames []string `json:"CollInterfaceNames"` //接口名称
	}{}

	commMessageMap := make([]device.CommunicationMessageTemplate, 0)

	err := context.ShouldBindJSON(&interfaceName)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,", err)

		context.JSON(http.StatusOK, Response{
			Code:    "1",
			Message: "json unMarshall err",
		})
		return
	}
	for _, name := range interfaceName.CollInterfaceNames {
		for _, v := range device.CollectInterfaceMap {
			if name == v.CollInterfaceName {
				commMessageMap = append(commMessageMap, v.CommMessage...)
				// 清空map
				v.CommMessage = v.CommMessage[0:0]
			}
		}
	}
	context.JSON(http.StatusOK, ResponseData{
		Code: "1",
		Data: commMessageMap,
	})
}
