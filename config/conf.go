package config

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

/**************获取配置信息************************/
func GetConf() {
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println(exeCurDir)

	path := exeCurDir + "/config/config.ini"
	iniFile, err := ini.Load(path)
	if err != nil{
		log.Println("Load config.ini err,",err)
	}
	LoadServer(iniFile)
}


