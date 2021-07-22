package report

import (
	"goAdapter/report/mqttAliyun"
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

	mqttHuawei.ReportServiceHuaweiInit()
}
