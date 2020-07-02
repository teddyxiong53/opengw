package main

import (
	"log"
	"strconv"
	"time"
)

type CommunicationCmd struct{

	interfaceID int					//接口ID
	deviceAddr  string				//接口下设备地址
	deviceType  string
	funName     string
	funIndex    int
	funPara     interface{}
}

var (
	emergencyRequestChan 	chan CommunicationCmd
	commonChan 				chan CommunicationCmd
	emergencyAckChan 		chan bool
)

func CommunicationManageInit(){

	commonChan 				= make(chan CommunicationCmd,100)
	emergencyRequestChan 	= make(chan CommunicationCmd,1)
	emergencyAckChan     	= make(chan bool,1)

	go CommunicationManageDel()
}

func CommunicationManageAdd(cmd CommunicationCmd){

	commonChan<- cmd
}

func CommunicationManageAddEmergency(cmd CommunicationCmd) bool{

	emergencyRequestChan<- cmd
	return <-emergencyAckChan
}

func CommunicationManageDel(){

	for {
		select {
		case cmd := <-emergencyRequestChan:
			{
				log.Println("emergency chan")
				log.Printf("funName %s\n", cmd.funName)
				var status bool = false
				//for _,v := range DeviceNodeManageMap[cmd.interfaceID].DeviceNodeMap{
				//	switch v.(type) {
				//	case DeviceNodeTemplate:
				//		{
				//			if v.(DeviceNodeTemplate).Addr == cmd.deviceAddr{
				//				//fcu := v.(DeviceNodeTemplate)
				//
				//			}
				//		}
				//	}
				//}
				emergencyAckChan<- status

				//通信帧延时
				interval,_ := strconv.Atoi(serialInterface.SerialParam[0].Interval)
				time.Sleep(time.Duration(interval)*time.Millisecond)
			}
		default:
			{
				select {
				case cmd := <-commonChan:
					{
						log.Println("common chan")
						log.Printf("funName %s\n", cmd.funName)

						for k,v := range DeviceInterfaceMap[cmd.interfaceID].DeviceNodeMap{
							log.Printf("index is %d\n",k)
							//组包
							txBuf := v.GetDeviceRealVariables(cmd.deviceAddr)
							log.Printf("buf is %+v\n",txBuf)
							//发送
							serialInterface.SerialPort[cmd.interfaceID].Write(txBuf)
							//等待接收
							rxBufTemp  := make([]byte, 256)
							rxBufTotal := make([]byte, 0)
							var (
								err           	error
								curRxBufCnt   	int 	= 0
								totalRxBufCnt 	int 	= 0
							)
							for {
								//阻塞读
								curRxBufCnt, err = serialInterface.SerialPort[cmd.interfaceID].Read(rxBufTemp)
								if err != nil{
									log.Println("threadModule read err,",err)
									break
								}
								if curRxBufCnt > 0 {
									totalRxBufCnt += curRxBufCnt
									//追加接收的数据到接收缓冲区
									rxBufTotal = append(rxBufTotal, rxBufTemp[:curRxBufCnt]...)
									//清除数据
									curRxBufCnt = 0
								}
							}


							//通信帧延时
							interval,_ := strconv.Atoi(serialInterface.SerialParam[0].Interval)
							time.Sleep(time.Duration(interval)*time.Millisecond)
						}
					}
				default:
					time.Sleep(10*time.Millisecond)
				}
			}
		}
	}
}

func CommunicationManageAddEmergencyTest(){
	cmd := CommunicationCmd{}

	cmd.interfaceID = InterFaceID1
	cmd.deviceAddr = "2"
	cmd.funName = "FCUGetRealData"
	cmd.funPara = struct{
		addr byte
	}{0x02}

	CommunicationManageAddEmergency(cmd)
}

func CommunicationManagePoll(){

	cmd := CommunicationCmd{}

	for i:=0;i<DeviceInterfaceMap[InterFaceID0].DeviceNodeCnt;i++{

		cmd.interfaceID = InterFaceID0
		cmd.deviceAddr = DeviceInterfaceMap[InterFaceID0].DeviceNodeAddrMap[i]
		cmd.deviceType = DeviceInterfaceMap[InterFaceID0].DeviceNodeTypeMap[i]
		cmd.funName = "GetDeviceRealVariables"

		CommunicationManageAdd(cmd)
	}
}