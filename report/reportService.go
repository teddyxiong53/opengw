package report

import (
	"goAdapter/report/mqttAliyun"
	mqttEmqx "goAdapter/report/mqttEMQX"
	"goAdapter/report/mqttHuawei"
)

type ReportServiceAPI interface {
	GWLogIn()
	GWLogOut()
	NodesLogIn()
	NodesLogOut()
	GWPropertyReport()
	NodesPropertyReport()
}

func init() {

}

func ReportServiceInit() {

	mqttAliyun.ReportServiceAliyunInit()

	mqttEmqx.ReportServiceEmqxInit()

	mqttHuawei.ReportServiceHuaweiInit()

}
