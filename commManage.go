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
							v.GetDeviceRealVariables(cmd.deviceAddr)
						}

						//通信帧延时
						interval,_ := strconv.Atoi(serialInterface.SerialParam[0].Interval)
						time.Sleep(time.Duration(interval)*time.Millisecond)
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