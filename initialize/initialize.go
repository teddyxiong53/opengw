<<<<<<< HEAD
/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-14 13:50:01
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 14:50:44
@FilePath: /goAdapter-Raw/initialize/initialize.go
*/
=======
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包)
package initialize

import (
	"goAdapter/config"
<<<<<<< HEAD
	"goAdapter/pkg/mylog"
=======
	"goAdapter/pkg/log"
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包)
)

func Init() {
	config.InitConfig("")
<<<<<<< HEAD
	mylog.InitLogger()
=======
	log.InitLogger()
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包)
}
