package setting

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
)

var (
	AppMode 			string
	HttpPort 			string
	LogLevel 			string		//日志等级
	LogSaveToFile 		bool		//日志储存到文件中
	LogFileMaxCnt 		uint		//日志储存最多文件数量
)

func LoadServer(file *ini.File){
	 AppMode = file.Section("server").Key("AppMode").MustString("debug")
	 HttpPort = file.Section("server").Key("HttpPort").MustString(":8080")
}

func LoadSerial(file *ini.File){

	//type SerialPortTemplate struct{
	//	Name []string			`json:"Name"`
	//}
	//
	//SerialPortName := &SerialPortTemplate{}
	err := file.Section("serial").MapTo(&SerialPortNameTemplateMap)
	if err != nil{
		log.Println(err)
	}
}

func LoadLog(file *ini.File){
	LogLevel = file.Section("log").Key("Level").MustString("DebugLevel")
	LogSaveToFile = file.Section("log").Key("SaveToFile").MustBool(false)
	LogFileMaxCnt = uint(file.Section("log").Key("FileMaxCnt").MustInt(10))
}

func LoadNetwork(file *ini.File) {

	err := file.Section("network").MapTo(&NetworkNameList)
	if err != nil {
		log.Println(err)
	}
}

/**************获取配置信息************************/
func GetConf() {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//log.Println(exeCurDir)

	path := exeCurDir + "/config/config.ini"
	iniFile, err := ini.Load(path)
	if err != nil{
		log.Println("Load config.ini err,",err)

		cfg := ini.Empty()

		AppMode = "debug"
		HttpPort = ":8080"
		cfg.Section("server").Key("AppMode").SetValue("debug")
		cfg.Section("server").Key("HttpPort").SetValue(":8080")

		cfg.Section("serial").Key("serialPort").SetValue("/dev/ttyS0")

		cfg.Section("network").Key("name").SetValue("eth0")

		cfg.SaveTo(path)
		return
	}

	LoadServer(iniFile)
	LoadSerial(iniFile)
	LoadNetwork(iniFile)
	LoadLog(iniFile)
}


