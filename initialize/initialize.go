/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-15 10:15:30
@LastEditors: WalkMiao
@LastEditTime: 2021-09-15 11:32:46
@FilePath: /goAdapter-Raw/initialize/initialize.go
*/
package initialize

import (
	"goAdapter/config"
	"goAdapter/pkg/mylog"
)

func Init() {
	config.InitConfig("")
	mylog.InitLogger()

}
