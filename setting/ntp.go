package setting

import (
	"encoding/json"
	"github.com/beevik/ntp"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
)


type NTPHostAddrTemplate struct {
	Status   bool		`json:"Status"`
	HostAddr []string	`json:"HostAddr"`
}

var NTPHostAddr = NTPHostAddrTemplate{
	Status:false,
	HostAddr: make([]string,0),
}

func NTPInit(){
	ReadNTPHostAddrFromJson()
}

func NTPGetTime() bool{

	if NTPHostAddr.Status == true{
		//多个服务器只要有一个能获取到时间就退出
		for _,v := range NTPHostAddr.HostAddr{
			ntpTime, err := ntp.Time(v)
			Logger.WithFields(logrus.Fields{
				"host"   : v,
				"err"    : err,
				"ntpTime":ntpTime,
			}).Warning("getNTPTime err")
			if err != nil{
				return false
			}else{
				SystemSetRTC(ntpTime.String())
				return true
			}
		}
	}

	return false
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
		log.Println("ntpHostAddr.json is not exist")

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