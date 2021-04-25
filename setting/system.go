package setting

import (
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type SystemStateTemplate struct{
	MemTotal  string		`json:"MemTotal"`
	MemUse    string		`json:"MemUse"`
	DiskTotal string        `json:"DiskTotal"`
	DiskUse   string		`json:"DiskUse"`
	Name      string		`json:"Name"`
	SN        string		`json:"SN"`
	HardVer   string		`json:"HardVer"`
	SoftVer   string		`json:"SoftVer"`
	SystemRTC string		`json:"SystemRTC"`
	RunTime   string		`json:"RunTime"`			//累计时间
	DeviceOnline string     `json:"DeviceOnline"`		//设备在线率
	DevicePacketLoss string `json:"DevicePacketLoss"`	//设备丢包率
}

type DataPointTemplate struct{
	Value string
	Time string
}

type DataStreamTemplate struct{
	DataPoint 		[]DataPointTemplate		`json:"DataPoint"`
	DataPointCnt 	int						`json:"DataPointCnt"`
	Legend 			string 					`json:"Legend"`			//别名
}

var SystemState = SystemStateTemplate{
	MemTotal		:"0",
	MemUse			:"0",
	DiskTotal		:"0",
	DiskUse			:"0",
	Name			:"goteway",
	SN				:"22005260001",
	HardVer			:"goteway-V.A",
	SoftVer			:"V0.0.1",
	SystemRTC		:"2020-05-26 12:00:00",
	RunTime			:"0",
	DeviceOnline    :"0",
	DevicePacketLoss : "0",
}

var timeStart time.Time
var (

	MemoryDataStream 			*DataStreamTemplate
	DiskDataStream 				*DataStreamTemplate
	DeviceOnlineDataStream		*DataStreamTemplate
	DevicePacketLossDataStream	*DataStreamTemplate
)


func SystemReboot() {
	cmd := exec.Command("reboot")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Start()

	str := out.String()
	fmt.Println(str)
}

func SystemSetRTC(rtc string){

	cmd := exec.Command("date","-s",rtc)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Start()

	//将时间写入硬件RTC中
	cmd = exec.Command("hwclock", "-w")
	cmd.Stdout = &out
	cmd.Start()
}

func GetMemState(){

	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	//log.Printf("Mem Total: %v, Free:%v, UsedPercent:%f%%\n",
	//					v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)

	SystemState.MemTotal = fmt.Sprintf("%d",v.Total/1024/1024)
	SystemState.MemUse = fmt.Sprintf("%3.1f",v.UsedPercent)
}

func GetDiskState(){

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	v, _ := disk.Usage(exeCurDir)

	// almost every return value is a struct
	//log.Printf("Disk Total: %v, Free:%v, UsedPercent:%f%%\n",
	//				v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)

	SystemState.DiskTotal = fmt.Sprintf("%d",v.Total/1024/1024)
	SystemState.DiskUse = fmt.Sprintf("%3.1f",v.UsedPercent)
}

func GetTimeStart(){

	timeStart = time.Now()
}

func GetRunTime(){

	elapsed := time.Since(timeStart)
	sec := int64(elapsed.Seconds())
	day := sec/86400
	hour := sec%86400/3600
	min := sec%3600/60
	sec = sec % 60

	strRunTime := fmt.Sprintf("%d天%d时%d分%d秒",day,hour,min,sec)

	SystemState.SystemRTC = time.Now().Format("2006-01-02 15:04:05")
	SystemState.RunTime = strRunTime
}

func NewDataStreamTemplate(legend string) *DataStreamTemplate{

	return &DataStreamTemplate{
		DataPoint: make([]DataPointTemplate,0),
		DataPointCnt: 0,
		Legend: legend,
	}
}

func (d *DataStreamTemplate)AddDataPoint(data DataPointTemplate){

	if d.DataPointCnt < 2880{
		d.DataPoint = append(d.DataPoint,data)
		d.DataPointCnt++
	}else{
		d.DataPoint = d.DataPoint[1:]
		d.DataPoint = append(d.DataPoint,data)
	}
}

func CollectSystemParam(){

	GetMemState()
	GetRunTime()

	point := DataPointTemplate{}

	point.Value = SystemState.MemUse
	point.Time = SystemState.SystemRTC
	MemoryDataStream.AddDataPoint(point)

	point.Value = SystemState.DiskUse
	point.Time = SystemState.SystemRTC
	DiskDataStream.AddDataPoint(point)

	point.Value = SystemState.DeviceOnline
	point.Time = SystemState.SystemRTC
	DeviceOnlineDataStream.AddDataPoint(point)

	point.Value = SystemState.DevicePacketLoss
	point.Time = SystemState.SystemRTC
	DevicePacketLossDataStream.AddDataPoint(point)
}

