package main

import (
	"encoding/binary"
	"fmt"
	"goAdapter/api"
	"log"
	"time"
)

const (
	INDEX_CmdHead = 0
	INDEX_Type = 2
	INDEX_Len = 3
	INDEX_Cmd = 4
	INDEX_ID = 4
	INDEX_Buff = 5
)

const (
	//数据模式
	TYPE_CMD  = 0  //命令
	TYPE_DATA = 1  //数据
	TYPE_ACK  = 0xFF
)

/**********************命令请求*****************************/
const (
	//允许注册
	CmdRequest_RegisterPermit       =0x04
	//Reset并重新建立网络
	CmdRequest_Reset                =0x08
	//设置RTC
	CmdRequest_SetRTC               =0x10
	//读取RTC
	CmdRequest_ReadRTC              =0x20
	//自组网系统状态
	CmdRequest_NetStatus            =0x11
	//读取版本信息
	CmdRequest_ReadVer              =0xC1
	//读取信道和PanID
	CmdRequest_ReadChanAndPanID     =0xE1
	//设置信道和PanID
	CmdRequest_WriteChanAndPanID    =0xE0
)


/**********************命令回应*****************************/
const (
	//允许注册
	CmdAck_RegisterPermit           =0x04
	//注册时间到
	CmdAck_RegisterTimeOut          =0x06
	//Reset并重新建立网络
	CmdAck_Reset                    =0x08
	//设置RTC
	CmdAck_WriteRTC                 =0x10
	//读取RTC
	CmdAck_ReadRTC                  =0x20
	//自组网系统状态
	CmdAck_ReadNetStatus            =0x11
	//读取版本信息
	CmdAck_ReadVer                  =0xC1
	//读取信道和PanID
	CmdAck_ReadChanAndPanID         =0xE1
	//设置信道和PanID
	CmdAck_WriteChanAndPanID        =0xE0
)

const (
	TotalStep = 2
)

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

	VariableMap = append(VariableMap, AddVariable(0,"Chan","信道","string"))
	VariableMap = append(VariableMap, AddVariable(1,"PanID","系统ID","string"))
	VariableMap = append(VariableMap, AddVariable(2,"SoftVer","软件版本","string"))
	VariableMap = append(VariableMap, AddVariable(3,"RTC","时钟","string"))

	return VariableMap
}

func TD200GetCRC(buf []byte,bufLen int) byte{

	var crc byte = 0
	var i   int

	for i=0;i<bufLen;i++{
		crc += buf[i]
	}

	return crc
}

/**
	生成读变量的数据包
*/
func GenerateGetRealVariables(sAddr string,step int) ([]byte,bool){

	requestAdu := make([]byte,0)

	if step == TotalStep{
		return requestAdu,false
	}

	if step == 0{
		requestAdu = append(requestAdu,0xFE)
		requestAdu = append(requestAdu,0xA5)
		requestAdu = append(requestAdu,TYPE_CMD)
		requestAdu = append(requestAdu,0x01)
		requestAdu = append(requestAdu,CmdRequest_ReadChanAndPanID)
		crc := TD200GetCRC(requestAdu[2:],3)
		requestAdu = append(requestAdu,crc)
	}else if step == 1{
		requestAdu = append(requestAdu,0xFE)
		requestAdu = append(requestAdu,0xA5)
		requestAdu = append(requestAdu,TYPE_CMD)
		requestAdu = append(requestAdu,0x01)
		requestAdu = append(requestAdu,CmdRequest_ReadRTC)
		crc := TD200GetCRC(requestAdu[2:],3)
		requestAdu = append(requestAdu,crc)
	}

	return requestAdu,true
}

func AnalysisRx(sAddr string,variables []api.VariableTemplate,rxBuf []byte,rxBufCnt int) chan bool{

	var (
		frameLen        byte    = 0
		frameHeader   	uint16 	= 0
		frameType    	byte 	= 0
		framePayloadLen byte  	= 0
		crc             byte
		framePayload            =  make([]byte,0)
	)

	status := make(chan bool,1)

	index := 0
	for {
		if index < rxBufCnt{
			if rxBufCnt < 5{
				return status
			}

			frameHeader = binary.BigEndian.Uint16(rxBuf[index:])
			frameType = rxBuf[index+2]
			framePayloadLen = rxBuf[index+3]
			if frameHeader != 0xFEA5{
				return status
			}
			//通信帧总长度
			frameLen = framePayloadLen + 5
			if int(frameLen)+index > rxBufCnt{
				return status
			}
			crc = TD200GetCRC(rxBuf[index+2:],int(framePayloadLen+2))
			if crc != rxBuf[index+int(framePayloadLen)+4]{
				return status
			}
			//获取到正确的数据
			framePayload = append(framePayload,rxBuf[index+4:index+4+int(framePayloadLen)]...)
			log.Printf("framePayLoad %X\n",framePayload)

			if frameType != TYPE_CMD{
				return status
			}
			if framePayload[0] == CmdAck_ReadChanAndPanID{
				timeNowStr := time.Now().Format("2006-01-02 15:04:05")

				variables[0].TimeStamp = timeNowStr
				variables[0].Value = framePayload[1]

				variables[1].TimeStamp = timeNowStr
				variables[1].Value = binary.BigEndian.Uint16(framePayload[2:])
			}else if framePayload[0] == CmdAck_ReadRTC {

				timeNowStr := time.Now().Format("2006-01-02 15:04:05")

				variables[3].TimeStamp = timeNowStr
				variables[3].Value = fmt.Sprintf("20%d-%02d-%02d %02d:%02d:%02d",framePayload[1],
					framePayload[2],
					framePayload[3],
					framePayload[4],
					framePayload[5],
					framePayload[6])
			}

			status<-true
			return status
		}else{
			break
		}
		index++
	}
	return status
}
