package httpServer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"goAdapter/device"
	"net/http"
)

func apiGetCommMessage(context *gin.Context) {

	aParam := &struct {
		Code    string
		Message string
		Data    []device.CommunicationMessageTemplate
	}{}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)

	interfaceName := struct {
		CollInterfaceNames  []string       `json:"CollInterfaceNames"` //接口名称
	}{}

	commMessageMap := make([]device.CommunicationMessageTemplate,0)

	err := json.Unmarshal(bodyBuf[:n], &interfaceName)
	if err != nil {
		fmt.Println("interfaceInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}else {
		for _,name := range interfaceName.CollInterfaceNames{
			for _,v := range device.CollectInterfaceMap{
				if name == v.CollInterfaceName{
					commMessageMap = append(commMessageMap,v.CommMessage...)
					// 清空map
					v.CommMessage = v.CommMessage[0:0]
				}
			}
		}
	}

	aParam.Code = "0"
	aParam.Message = ""
	aParam.Data = commMessageMap

	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}