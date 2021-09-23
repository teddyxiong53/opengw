/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-14 13:54:05
@LastEditors: WalkMiao
@LastEditTime: 2021-09-17 08:25:53
@FilePath: /goAdapter-Raw/config/config.go
*/
package config

import (
	"log"

	"github.com/spf13/viper"
)

var Cfg *Config

type Config struct {
	ServerCfg HttpServerConfig `mapstructure:"server"`
	UdpCfg    UdpServerConfig  `mapstructure:"udpserver"`
	LogCfg    LoggerConfig     `mapstructure:"log"`
	SerialCfg SerialConfig     `mapstructure:"serial"`
}

type HttpServerConfig struct {
	Port    string `mapstructure:"port"`
	Timeout int    `mapstructure:"timeout"`
	GinMode string `mapstructure:"mode"`
}

type UdpServerConfig struct {
	Addr string `mapstructure:"addr"`
}

type SerialConfig struct {
	ReadTimeOut     int `mapstructure:"readtimeout"`
	BufferReadDelay int `mapstructure:"bufferReadDelay"`
}
type LoggerConfig struct {
	Level string `mapstructure:"level"`
	Dir   string `mapstructure:"dir"` //日志存放目录
	File  string `mapstructure:"file"`
}

func InitConfig(path string) {
	log.Println("开始初始化配置文件...")
	vp := viper.New()
	if path != "" {
		vp.SetConfigFile(path)
	} else {
		vp.AddConfigPath(".")
		vp.AddConfigPath("..")
		vp.SetConfigName("config")
		vp.SetConfigType("yaml")
		vp.SetConfigFile("config.yaml")
	}

	if err := vp.ReadInConfig(); err != nil {
		panic(err)
	}
	var settings Config
	if err := vp.Unmarshal(&settings); err != nil {
		panic(err)
	}
	Cfg = &settings
	log.Printf("初始化配置文件成功...")
	log.Printf("串口读超时设置为: %d ms delay:%d ms\n", Cfg.SerialCfg.ReadTimeOut, Cfg.SerialCfg.BufferReadDelay)
}
