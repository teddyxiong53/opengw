package setting

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var Logger = logrus.New()

func init() {

}

func LogerInit() {

	//log输出行号和ms
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

	if AppMode == "release" {
		// 设置日志格式为json格式　自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
		Logger.Formatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.999999999Z07:00",
		}
	} else if AppMode == "debug" {
		// 设置日志格式为json格式　自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
		Logger.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
			//TimestampFormat: "2006-01-02T15:04:05.999999999",
			TimestampFormat: "01-02T15:04:05.999999999",
			CallerPrettyfier: func(run *runtime.Frame) (function string, file string) {
				fileInfo := path.Base(run.File)
				//fileInfo := run.File
				lineInfo := strconv.Itoa(run.Line)
				return "", fileInfo + ":" + lineInfo
			},
		}
		// 设置输出文件名和行号
		Logger.ReportCaller = true
	}

	if LogSaveToFile == true {
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
			rotatelogs.WithRotationCount(LogFileMaxCnt),
			rotatelogs.WithRotationTime(time.Hour),
		)
		if err != nil {
			fmt.Println(err)
		}
		//Loger.SetOutput(writer)

		Logger.Out = writer
	} else {
		// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
		// 日志消息输出可以是任意的io.writer类型
		Logger.SetOutput(os.Stdout)
	}

	//日志的级别
	//- Fatal：挂了，或者极度不正常
	//- Error：跟遇到的用户说对不起，可能有bug
	//- Warn：记录一下，某事又发生了
	//- Info：提示一切正常
	//- debug：没问题，就看看堆栈

	switch LogLevel {
	case "TraceLevel": //用户级输出
		//Loger.SetLevel(logrus.DebugLevel)
		Logger.Level = logrus.TraceLevel
	case "DebugLevel": //用户级调试
		//Loger.SetLevel(logrus.DebugLevel)
		Logger.Level = logrus.DebugLevel
	case "InfoLevel": //用户级重要
		//Loger.SetLevel(logrus.InfoLevel)
		Logger.Level = logrus.InfoLevel
	case "WarnLevel": //用户级警告
		//Loger.SetLevel(logrus.WarnLevel)
		Logger.Level = logrus.WarnLevel
	case "ErrorLevel": //用户级错误
		//Loger.SetLevel(logrus.WarnLevel)
		Logger.Level = logrus.ErrorLevel
	}
}
