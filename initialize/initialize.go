package initialize

import (
	"goAdapter/config"
	"goAdapter/pkg/mylog"
)

func Init() {
	config.InitConfig("")

	mylog.InitLogger()

}
