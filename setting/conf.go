package setting

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
	"path/filepath"
)

var (
	AppMode string
	HttpPort string
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


/**************获取配置信息************************/
func GetConf() {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println(exeCurDir)

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

		cfg.SaveTo(path)
		return
	}

	LoadServer(iniFile)
	LoadSerial(iniFile)
	log.Printf("serial %+v\n",SerialPortNameTemplateMap)
}


