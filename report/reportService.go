package report

import (
	"os"
)

//type ReportServiceParamListTemplate interface {
//	ReadParamFromJson() bool
//	WriteParamToJson()
//	AddReportService()
//}

func init() {

}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func ReportServiceInit() {

	for _, v := range ReportServiceParamListAliyun.ServiceList {

		go ReportServiceAliyunPoll(v)
	}
}