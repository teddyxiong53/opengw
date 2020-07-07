package device

import (
	"goAdapter/setting"
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
	rxChan                  chan bool               //接收正确
)

func CommunicationManageInit(){

	commonChan 				= make(chan CommunicationCmd,100)
	emergencyRequestChan 	= make(chan CommunicationCmd,1)
	emergencyAckChan     	= make(chan bool,1)
	rxChan     				= make(chan bool,1)

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
				interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[0].Interval)
				time.Sleep(time.Duration(interval)*time.Millisecond)
			}
		default:
			{
				select {
				case cmd := <-commonChan:
					{
						log.Println("common chan")
						//log.Printf("funName %s\n", cmd.funName)


						for k,v := range DeviceInterfaceMap[cmd.interfaceID].DeviceNodeAddrMap{
							if v == cmd.deviceAddr{
								log.Printf("index is %d\n",k)
								//--------------组包---------------------------
								txBuf := DeviceInterfaceMap[cmd.interfaceID].DeviceNodeMap[k].GetDeviceRealVariables()
								log.Printf("tx buf is %+v\n",txBuf)
								//---------------发送-------------------------
								setting.SerialInterface.SerialPort[cmd.interfaceID].Write(txBuf)
								//---------------等待接收----------------------
								//阻塞读
								rxBuf  := make([]byte, 256)
								rxTotalBuf := make([]byte,0)
								rxBufCnt := 0
								rxTotalBufCnt := 0
								timeOut,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.interfaceID].Timeout)
								timer := time.NewTimer(time.Duration(timeOut)*time.Millisecond)
								for {
									select{
										//是否正确收到数据包
										case <-DeviceInterfaceMap[cmd.interfaceID].DeviceNodeMap[k].ProcessRx(rxChan,rxTotalBuf,rxTotalBufCnt):
										{
											log.Println("rx ok")
											//通信帧延时
											interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.interfaceID].Interval)
											time.Sleep(time.Duration(interval)*time.Millisecond)
											goto Loop
										}
										//是否接收超时
										case <-timer.C:
										{
											log.Println("rx timeout")
											//通信帧延时
											interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.interfaceID].Interval)
											time.Sleep(time.Duration(interval)*time.Millisecond)
											goto Loop
										}
										//继续接收数据
										default:
										{
											rxBufCnt,_ = setting.SerialInterface.SerialPort[cmd.interfaceID].Read(rxBuf)
											if rxBufCnt > 0{
												rxTotalBufCnt += rxBufCnt
												//追加接收的数据到接收缓冲区
												rxTotalBuf = append(rxTotalBuf,rxBuf[:rxBufCnt]...)
												//清除本地接收数据
												rxBufCnt = 0
												log.Printf("rxbuf %+v\n",rxTotalBuf)
											}
										}
									}
								}
								Loop:
							}
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

	//cmd := CommunicationCmd{}
	//
	//for i:=0;i<DeviceInterfaceMap[InterFaceID0].DeviceNodeCnt;i++{
	//
	//	cmd.interfaceID = InterFaceID0
	//	cmd.deviceAddr = DeviceInterfaceMap[InterFaceID0].DeviceNodeAddrMap[i]
	//	cmd.deviceType = DeviceInterfaceMap[InterFaceID0].DeviceNodeTypeMap[i]
	//	cmd.funName = "GetDeviceRealVariables"
	//
	//	CommunicationManageAdd(cmd)
	//}
}