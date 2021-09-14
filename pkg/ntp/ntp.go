<<<<<<< HEAD:setting/ntp.go
<<<<<<< Updated upstream:setting/ntp.go
package setting
=======
=======
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包):pkg/ntp/ntp.go
/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
<<<<<<< HEAD:setting/ntp.go
@LastEditTime: 2021-09-14 14:47:01
@FilePath: /goAdapter-Raw/pkg/ntp/ntp.go
*/
package ntp
>>>>>>> Stashed changes:pkg/ntp/ntp.go
=======
@LastEditTime: 2021-09-14 14:15:46
@FilePath: /goAdapter-Raw/pkg/ntp/ntp.go
*/
package ntp
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包):pkg/ntp/ntp.go

import (
	"encoding/json"
	"goAdapter/pkg/mylog"
	"goAdapter/pkg/system"
	"log"
	"os"
	"path/filepath"

	"github.com/beevik/ntp"
	"github.com/sirupsen/logrus"
)

type NTPHostAddrTemplate struct {
	Status   bool     `json:"Status"`
	HostAddr []string `json:"HostAddr"`
}

var NTPHostAddr = NTPHostAddrTemplate{
	Status:   false,
	HostAddr: make([]string, 0),
}

func init() {
	ReadNTPHostAddrFromJson()
}

func NTPGetTime() bool {

	if NTPHostAddr.Status == true {
		//多个服务器只要有一个能获取到时间就退出
		for _, v := range NTPHostAddr.HostAddr {
			ntpTime, err := ntp.Time(v)
			mylog.Logger.WithFields(logrus.Fields{
				"host":    v,
				"err":     err,
				"ntpTime": ntpTime,
			}).Warning("getNTPTime err")
			if err != nil {
				return false
			} else {
				system.SystemSetRTC(ntpTime.String())
				return true
			}
		}
	}

	return false
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func ReadNTPHostAddrFromJson() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/ntpHostAddr.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open ntpHostAddr.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt], &NTPHostAddr)
		if err != nil {
			log.Println("ntpHostAddr unmarshal err", err)
			return false
		}
		return true
	} else {
		mylog.Logger.Infof("ntpHostAddr.json is not exist")

		return false
	}
}

func WriteNTPHostAddrToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/ntpHostAddr.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open ntpHostAddr.json err", err)
		return
	}
	defer fp.Close()

	sJson, _ := json.Marshal(NTPHostAddr)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write ntpHostAddr.json err", err)
	}
	log.Println("write ntpHostAddr.json sucess")
}
