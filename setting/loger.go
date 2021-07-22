package setting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

type LoggerParamTemplate struct {
	Level      string `json:"Level"`
	Output     int    `json:"Output"`     //输出方向：0表示输出到终端，1表示同时输出到终端和文件，2表示同时输出到终端和以太网
	FileMaxCnt uint   `json:"FileMaxCnt"` //日志文件最大存储数量
	FileDir    string `json:"FileDir"`    //日志文件存储路径
}

type DefaultFieldsHook struct {
}

var Logger = logrus.New()
var LoggerParam *LoggerParamTemplate

func init() {
	//log输出行号和ms
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
}

func (df *DefaultFieldsHook) Fire(entry *logrus.Entry) error {

	switch entry.Level {
	case logrus.TraceLevel:
		fallthrough
	case logrus.DebugLevel:
		fallthrough
	case logrus.InfoLevel:
		fallthrough
	case logrus.WarnLevel:
		fallthrough
	case logrus.ErrorLevel:
		{
			var line int
			if entry.Caller != nil {
				line = entry.Caller.Line
			}
			entry.Data["timeStamp"] = time.Now().Local().Format("01-02T15:04:05.000")
			entry.Data["goRoutineID"] = getGoRoutineID()
			entry.Data["fileLine"] = fmt.Sprintf("%d", line)
		}
	}

	return nil
}

func (df *DefaultFieldsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func getGoRoutineID() uint64 {
	b := make([]byte, 64)

	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)

	return n
}

func (lg *LoggerParamTemplate) ReadParamFromJson() bool {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/config/loggerParam.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open loggerParam.json err,", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt], lg)
		if err != nil {
			log.Println("loggerParam unmarshal err", err)
			return false
		}

		log.Print("read loggerParam.json ok")
		return true
	} else {
		//Logger.Print("loggerParam.json is not exist")
		return false
	}
}

func (lg *LoggerParamTemplate) WriteParamToJson(param *LoggerParamTemplate) {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/config/loggerParam.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open loggerParam.json err", err)
		return
	}
	defer fp.Close()

	lg.Level = param.Level
	lg.Output = param.Output
	lg.FileMaxCnt = param.FileMaxCnt
	lg.FileDir = param.FileDir

	sJson, err := json.Marshal(*lg)
	if err != nil {
		log.Printf("loggerParam Marshalerr %v", err)
	}

	_, err = fp.Write(sJson)
	if err != nil {
		log.Printf("write loggerParam.json err %v", err)
	}
	log.Printf("write loggerParam.json success")
}

func LogerInit() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	LoggerParam = &LoggerParamTemplate{
		Level:      "DebugLevel",
		Output:     1,
		FileMaxCnt: 5,
		FileDir:    exeCurDir + "/log/",
	}
	LoggerParam.ReadParamFromJson()

	formatter := &logrus.JSONFormatter{
		DisableTimestamp: true,
		CallerPrettyfier: func(run *runtime.Frame) (function string, file string) {
			fileInfo := path.Base(run.File)
			return "", fileInfo
		},
		PrettyPrint: false,
	}
	Logger.SetFormatter(formatter)
	// 添加自己实现的Hook
	Logger.AddHook(&DefaultFieldsHook{})

	switch LoggerParam.Output {
	case 0: //0表示输出到终端
		{
			loggerWriters := []io.Writer{
				os.Stdout,
			}
			fileAndStdoutWriter := io.MultiWriter(loggerWriters...)
			Logger.SetOutput(fileAndStdoutWriter)
		}
	case 1: //1表示同时输出到终端和文件
		{
			/* 日志轮转相关函数
			`WithLinkName` 为最新的日志建立软连接
			`WithRotationTime` 设置日志分割的时间，隔多久分割一次
			`WithMaxAge 和 WithRotationCount二者只能设置一个
			`WithMaxAge` 设置文件清理前的最长保存时间
			`WithRotationCount` 设置文件清理前最多保存的个数
			*/
			// 下面配置日志每隔 24小时轮转一个新文件，保留最近3天的日志文件，多余的自动清理掉。
			creatTime := time.Now().Local().Format("01-02T15")
			file, _ := rotatelogs.New(
				LoggerParam.FileDir+creatTime+".csv",
				//rotatelogs.WithLinkName(path),
				rotatelogs.WithRotationTime(24*time.Hour),
				rotatelogs.WithRotationCount(LoggerParam.FileMaxCnt),
			)
			loggerWriters := []io.Writer{
				file,
				os.Stdout,
			}
			fileAndStdoutWriter := io.MultiWriter(loggerWriters...)
			Logger.SetOutput(fileAndStdoutWriter)
		}
	case 2: //2表示同时输出到终端和以太网
		{

		}
	}

	//日志的级别
	//- Fatal：挂了，或者极度不正常
	//- Error：跟遇到的用户说对不起，可能有bug
	//- Warn：记录一下，某事又发生了
	//- Info：提示一切正常
	//- debug：没问题，就看看堆栈
	switch LoggerParam.Level {
	case "TraceLevel": //用户级输出
		//设置输出文件名和行号
		Logger.SetReportCaller(true)
		Logger.Level = logrus.TraceLevel
	case "DebugLevel": //用户级调试
		//设置输出文件名和行号
		Logger.SetReportCaller(true)
		Logger.Level = logrus.DebugLevel
	case "InfoLevel": //用户级重要
		//设置输出文件名和行号
		Logger.SetReportCaller(false)
		Logger.Level = logrus.InfoLevel
	case "WarnLevel": //用户级警告
		//设置输出文件名和行号
		Logger.SetReportCaller(false)
		Logger.Level = logrus.WarnLevel
	case "ErrorLevel": //用户级错误
		//设置输出文件名和行号
		Logger.SetReportCaller(false)
		Logger.Level = logrus.ErrorLevel
	}
}
