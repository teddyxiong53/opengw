package system

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

func init() {

	GetTimeStart()
}

type SystemStateTemplate struct {
	MemTotal         string `json:"MemTotal"`
	MemUse           string `json:"MemUse"`
	DiskTotal        string `json:"DiskTotal"`
	DiskUse          string `json:"DiskUse"`
	Name             string `json:"Name"`
	SN               string `json:"SN"`
	HardVer          string `json:"HardVer"`
	SoftVer          string `json:"SoftVer"`
	SystemRTC        string `json:"SystemRTC"`
	RunTime          string `json:"RunTime"`          //累计时间
	DeviceOnline     string `json:"DeviceOnline"`     //设备在线率
	DevicePacketLoss string `json:"DevicePacketLoss"` //设备丢包率
}

type DataPointTemplate struct {
	Value string
	Time  string
}

type DataStreamTemplate struct {
	DataPoint    []DataPointTemplate `json:"DataPoint"`
	DataPointCnt int                 `json:"DataPointCnt"`
	Legend       string              `json:"Legend"` //别名
}

var SystemState = SystemStateTemplate{
	MemTotal:         "0",
	MemUse:           "0",
	DiskTotal:        "0",
	DiskUse:          "0",
	Name:             "openGW",
	SN:               "22005260001",
	HardVer:          "openGW-V.A",
	SoftVer:          "V1.0.0",
	SystemRTC:        "2020-05-26 12:00:00",
	RunTime:          "0",
	DeviceOnline:     "0",
	DevicePacketLoss: "0",
}

var timeStart time.Time
var (
	MemoryDataStream           = NewDataStreamTemplate("内存使用率")
	DiskDataStream             = NewDataStreamTemplate("硬盘使用率")
	DeviceOnlineDataStream     = NewDataStreamTemplate("设备在线率")
	DevicePacketLossDataStream = NewDataStreamTemplate("通信丢包率")
)

func SmoothReStart()error {
	if runtime.GOOS=="windows"{
		return fmt.Errorf("windows soft restart is not supported")
	}
	cmd := exec.Command("kill","-1",strconv.Itoa(syscall.Getpid()))
	_,err:=cmd.CombinedOutput()
	if err!=nil{
		return fmt.Errorf("excute command %s error:%v",cmd.String(),err)
	}
	return nil
}

func SystemSetRTC(rtc string) {
	cmd := exec.Command("date", "-s", rtc)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Start()
	if err != nil {
		return
	}

	//将时间写入硬件RTC中
	cmd = exec.Command("hwclock", "-w")
	cmd.Stdout = &out
	err = cmd.Start()
	if err != nil {
		return
	}
}

func GetMemState() {

	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	//log.Printf("Mem Total: %v, Free:%v, UsedPercent:%f%%\n",
	//					v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)

	SystemState.MemTotal = fmt.Sprintf("%d", v.Total/1024/1024)
	SystemState.MemUse = fmt.Sprintf("%3.1f", v.UsedPercent)
}

func GetDiskState() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	v, _ := disk.Usage(exeCurDir)

	// almost every return value is a struct
	//log.Printf("Disk Total: %v, Free:%v, UsedPercent:%f%%\n",
	//				v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)

	SystemState.DiskTotal = fmt.Sprintf("%d", v.Total/1024/1024)
	SystemState.DiskUse = fmt.Sprintf("%3.1f", v.UsedPercent)
}

func GetTimeStart() {

	timeStart = time.Now()
}

func GetRunTime() {

	elapsed := time.Since(timeStart)
	sec := int64(elapsed.Seconds())
	day := sec / 86400
	hour := sec % 86400 / 3600
	min := sec % 3600 / 60
	sec = sec % 60

	strRunTime := fmt.Sprintf("%d天%d时%d分%d秒", day, hour, min, sec)

	SystemState.SystemRTC = time.Now().Format("2006-01-02 15:04:05")
	SystemState.RunTime = strRunTime
}

func NewDataStreamTemplate(legend string) *DataStreamTemplate {

	return &DataStreamTemplate{
		DataPoint:    make([]DataPointTemplate, 0),
		DataPointCnt: 0,
		Legend:       legend,
	}
}

func (d *DataStreamTemplate) AddDataPoint(data DataPointTemplate) {

	if d.DataPointCnt < 2880 {
		d.DataPoint = append(d.DataPoint, data)
		d.DataPointCnt++
	} else {
		d.DataPoint = d.DataPoint[1:]
		d.DataPoint = append(d.DataPoint, data)
	}
}

func CollectSystemParam() {

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

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
