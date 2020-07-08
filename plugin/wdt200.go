package main

import (
	"encoding/binary"
	"goAdapter/api"
	"log"
	"strconv"
	"sync"
	"time"
)

type DeviceNodeTemplate struct{
	Variables    []api.VariableTemplate
	TemplateName string					//模板名称
	TemplateType string					//模板型号
	TemplateID   string					//模板ID
	TemplateMessage string              //备注信息
}

type crc struct {
	once  sync.Once
	table []uint16
}

var crcTb crc


// initTable 初始化表
func (c *crc) initTable() {
	crcPoly16 := uint16(0xa001)
	c.table = make([]uint16, 256)

	for i := uint16(0); i < 256; i++ {
		crc := uint16(0)
		b := i

		for j := uint16(0); j < 8; j++ {
			if ((crc ^ b) & 0x0001) > 0 {
				crc = (crc >> 1) ^ crcPoly16
			} else {
				crc = crc >> 1
			}
			b = b >> 1
		}
		c.table[i] = crc
	}
}

func crc16(bs []byte) uint16 {
	crcTb.once.Do(crcTb.initTable)

	val := uint16(0xFFFF)
	for _, v := range bs {
		val = (val >> 8) ^ crcTb.table[(val^uint16(v))&0x00FF]
	}
	return val
}

func AddVariable(vindex int,vname string,vlable string,vtype string) api.VariableTemplate{

	variable := api.VariableTemplate{}
	variable.Index = vindex
	variable.Name = vname
	variable.Lable = vlable
	variable.Type = vtype

	return variable
}

func NewVariables() []api.VariableTemplate{

	VariableMap := make([]api.VariableTemplate,0)

	VariableMap = append(VariableMap,AddVariable(0,"Addr","通信地址","string"))
	VariableMap = append(VariableMap,AddVariable(1,"DeviceType","设备类型","string"))
	VariableMap = append(VariableMap,AddVariable(2,"SoftVer","软件版本","string"))
	VariableMap = append(VariableMap,AddVariable(3,"SerialNumber","设备编码","string"))
	VariableMap = append(VariableMap,AddVariable(4,"RTC","设备时钟","string"))
	VariableMap = append(VariableMap,AddVariable(5,"RoomTemp","室内温度","string"))
	VariableMap = append(VariableMap,AddVariable(6,"RoomHumi","室内湿度","string"))
	VariableMap = append(VariableMap,AddVariable(7,"AirStatus","风机状态","string"))
	VariableMap = append(VariableMap,AddVariable(8,"RelayStatus","阀门状态","string"))
	VariableMap = append(VariableMap,AddVariable(9,"CurModeSet","运行模式","string"))
	VariableMap = append(VariableMap,AddVariable(10,"SeaconStatus","季节","string"))
	VariableMap = append(VariableMap,AddVariable(11,"CurTempStep","设定温度","byte"))
	VariableMap = append(VariableMap,AddVariable(12,"LockStatus","锁定状态","string"))
	VariableMap = append(VariableMap,AddVariable(13,"EnergySavingStatus","节能状态","string"))
	VariableMap = append(VariableMap,AddVariable(14,"RSSI","信号强度","byte"))
	VariableMap = append(VariableMap,AddVariable(15,"BLKStatus","背光状态","string"))
	VariableMap = append(VariableMap,AddVariable(16,"ErrorCode","故障代码","uint16"))
	VariableMap = append(VariableMap,AddVariable(17,"TotalTime","累计时间","uint32"))
	VariableMap = append(VariableMap,AddVariable(18,"LowTotalTime","低速累计时间","uint32"))
	VariableMap = append(VariableMap,AddVariable(19,"MiddleTotalTime","中速累计时间","uint32"))
	VariableMap = append(VariableMap,AddVariable(20,"HighTotalTime","高速累计时间","uint32"))

	return VariableMap
}

func GenerateGetRealVariables(sAddr string) []byte{

	addr,_ := strconv.Atoi(sAddr)
	requestAdu := make([]byte,0)

	requestAdu = append(requestAdu,byte(addr))
	requestAdu = append(requestAdu,0x03)
	requestAdu = append(requestAdu,0x00,0x01)
	requestAdu = append(requestAdu,0x00,0x02)

	checksum := crc16(requestAdu)
	requestAdu = append(requestAdu,byte(checksum),byte(checksum>>8))

	return requestAdu
}

func AnalysisRx(sAddr string,variables []api.VariableTemplate,rxBuf []byte,rxBufCnt int) chan bool{

	status := make(chan bool,1)

	addr,_ := strconv.Atoi(sAddr)

	index := 0
	for {
		if index < rxBufCnt{
			if rxBufCnt < 4{
				return status
			}
			crc := crc16(rxBuf[0:len(rxBuf)-2])
			expect := binary.LittleEndian.Uint16(rxBuf[len(rxBuf)-2:])
			if crc != expect{
				return status
			}
			if rxBuf[0] != byte(addr){
				return status
			}
			if rxBuf[1] != 0x03{
				return status
			}
			if rxBuf[2] != 4{
				return status
			}
			log.Println("processRx ok")

			timeNowStr := time.Now().Format("2006-01-02 15:04:05")

			variables[0].TimeStamp = timeNowStr
			status<-true
			return status
		}else{
			break
		}
		index++
	}
	return status
}