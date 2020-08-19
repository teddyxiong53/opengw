package setting

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

var Loger = logrus.New()

func LogerInit(level string,save bool,cnt uint){

	//log输出行号和ms
	//log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	// 设置日志格式为json格式　自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	Loger.Formatter = &logrus.JSONFormatter{}
	//fmt.Printf("level %v\n",level)
	//fmt.Printf("save %v\n",save)
	if save == true{

		exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

		path := exeCurDir + "/log/"
		/* 日志轮转相关函数
		`WithLinkName` 为最新的日志建立软连接
		`WithRotationTime` 设置日志分割的时间，隔多久分割一次
		`WithMaxAge 和 WithRotationCount二者只能设置一个
		`WithMaxAge` 设置文件清理前的最长保存时间
		`WithRotationCount` 设置文件清理前最多保存的个数
		*/
		// 下面配置日志每隔 60 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
		writer, err := rotatelogs.New(
			path+"%Y%m%d%H%M.txt",
			//rotatelogs.WithLinkName(path),
			rotatelogs.WithRotationCount(cnt),
			rotatelogs.WithRotationTime(time.Hour),
		)
		if err != nil{
			fmt.Println(err)
		}
		//Loger.SetOutput(writer)

		Loger.Out = writer
	}else{
		// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
		// 日志消息输出可以是任意的io.writer类型
		Loger.SetOutput(os.Stdout)
	}

	// 设置日志级别为warn以上
	switch level {
	case "DebugLevel":
		//Loger.SetLevel(logrus.DebugLevel)
		Loger.Level = logrus.DebugLevel
	case "InfoLevel":
		//Loger.SetLevel(logrus.InfoLevel)
		Loger.Level = logrus.InfoLevel
	case "WarnLevel":
		//Loger.SetLevel(logrus.WarnLevel)
		Loger.Level = logrus.WarnLevel
	}
}
