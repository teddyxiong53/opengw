package main

import (
	"encoding/binary"
	"log"
	"strconv"
	"sync"
)

type DeviceNodeModbusTemplate struct{
	DeviceNodeTemplate
	TemplateName string					//模板名称
	TemplateType string					//模板型号
	TemplateID   string					//模板ID
	TemplateMessage string              //备注信息
}

type crc struct {
	once  sync.Once
	table []uint16
}

var (
	crcTb 				crc
)

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

func (d *DeviceNodeModbusTemplate)ProcessRx(){


	var (
		curRxBufCnt   	int 	= 0
		totalRxBufCnt 	int 	= 0
		err           	error
		index         	int 	= 0
		frameLen        byte    = 0
		frameHeader   	uint16 	= 0
		frameType    	byte 	= 0
		framePayloadLen byte  	= 0
		frameParam      ThreadModuleFrameParam
		crc             byte
	)

	rxBufTemp := make([]byte, 256)
	rxBufTotal := make([]byte, 0)

	framePayload := make([]byte,0)

	for {
		//阻塞读
		curRxBufCnt, err = serialInterface.SerialPort[0].Read(rxBufTemp)
		if err != nil{
			log.Println("threadModule read err,",err)
			return
		}
		if curRxBufCnt > 0{
			totalRxBufCnt += curRxBufCnt

			//追加接收的数据到接收缓冲区
			rxBufTotal = append(rxBufTotal,rxBufTemp[:curRxBufCnt]...)
			//清除数据
			curRxBufCnt = 0

			//log.Println("totalRxBufCnt:",totalRxBufCnt)
			//threadModuleLogHex(1,rxBufTotal)

			index = 0
			for {
				//log.Printf("index:%d\n",index)
				if index < totalRxBufCnt{
					if (totalRxBufCnt>=5) && (index+5<=totalRxBufCnt) {
						frameHeader = binary.BigEndian.Uint16(rxBufTotal[index:])
						frameType = rxBufTotal[2+index]
						framePayloadLen = rxBufTotal[3+index]
						if frameHeader == 0xFEA5{
							//log.Println("header sucess")
							//估算完整数据帧的长度
							frameLen = framePayloadLen + 5
							if int(frameLen)+index <= totalRxBufCnt{
								crc = threadModuleGetCRC(rxBufTotal[2+index:],int(framePayloadLen+2))
								if crc == rxBufTotal[index+int(framePayloadLen)+4]{
									//log.Println("crc sucess")

									framePayload = append(framePayload,rxBufTotal[index+4:index+4+int(framePayloadLen)]...)

									frameParam.FrameType = frameType
									frameParam.FramePayloadLen = framePayloadLen
									frameParam.FramePayload = framePayload
									//log.Printf("FrameType:%v,FramePayloadLen:%v,FramePayload:%v \n",
									//	frameParam.FrameType,frameParam.FramePayloadLen,frameParam.FramePayload)

									FrameParamChan <- frameParam

									framePayload = framePayload[0:0]
									rxBufTotal = rxBufTotal[0:0]
									framePayloadLen = 0
									//totalRxBufCnt -= int(frameLen)
									totalRxBufCnt = 0

									break
								}
							}
						}
					}
					index++
				}else{
					break
				}
			}
		}
	}
}

func (d *DeviceNodeModbusTemplate)GetDeviceRealVariables(deviceAddr string) []byte{

	addr,_ := strconv.Atoi(d.Addr)
	requestAdu := make([]byte,0)

	requestAdu = append(requestAdu,byte(addr))
	requestAdu = append(requestAdu,0x03)
	requestAdu = append(requestAdu,0x00,0x01)
	requestAdu = append(requestAdu,0x00,0x02)

	checksum := crc16(requestAdu)
	requestAdu = append(requestAdu,byte(checksum),byte(checksum>>8))

	return requestAdu
}

func (d *DeviceNodeModbusTemplate)SetDeviceRealVariables(deviceAddr string) int{


	return 0
}

func (d *DeviceNodeModbusTemplate)AddVariable(vindex int,vname string,vlable string,vtype string){

	variable := VariableTemplate{}
	variable.Index = vindex
	variable.Name = vname
	variable.Lable = vlable
	variable.Type = vtype

	d.VariableMap = append(d.VariableMap,variable)
}

func (d *DeviceNodeModbusTemplate)NewVariables() {

	d.VariableMap = make([]VariableTemplate,0)

	d.AddVariable(0,"Addr","通信地址","string")
	d.AddVariable(1,"DeviceType","设备类型","string")
	d.AddVariable(2,"SoftVer","软件版本","string")
	d.AddVariable(3,"SerialNumber","设备编码","string")
	d.AddVariable(4,"RTC","设备时钟","string")
	d.AddVariable(5,"RoomTemp","室内温度","string")
	d.AddVariable(6,"RoomHumi","室内湿度","string")
	d.AddVariable(7,"AirStatus","风机状态","string")
	d.AddVariable(8,"RelayStatus","阀门状态","string")
	d.AddVariable(9,"CurModeSet","运行模式","string")
	d.AddVariable(10,"SeaconStatus","季节","string")
	d.AddVariable(11,"CurTempStep","设定温度","byte")
	d.AddVariable(12,"LockStatus","锁定状态","string")
	d.AddVariable(13,"EnergySavingStatus","节能状态","string")
	d.AddVariable(14,"RSSI","信号强度","byte")
	d.AddVariable(15,"BLKStatus","背光状态","string")
	d.AddVariable(16,"ErrorCode","故障代码","uint16")
	d.AddVariable(17,"TotalTime","累计时间","uint32")
	d.AddVariable(18,"LowTotalTime","低速累计时间","uint32")
	d.AddVariable(19,"MiddleTotalTime","中速累计时间","uint32")
	d.AddVariable(20,"HighTotalTime","高速累计时间","uint32")
}



