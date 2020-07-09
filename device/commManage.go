package device

import (
	"goAdapter/setting"
	"log"
	"strconv"
	"time"
)

type CommunicationCmd struct{

	InterfaceID int    //接口ID
	DeviceAddr  string //接口下设备地址
	DeviceType  string
	FunName     string
	FunIndex    int
	FunPara     interface{}
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
				log.Printf("funName %s\n", cmd.FunName)
				var status bool = false
				for k,v := range DeviceInterfaceMap[cmd.InterfaceID].DeviceNodeMap{
					if v.Addr == cmd.DeviceAddr {
						log.Printf("index is %d\n",k)
						//--------------组包---------------------------
						txBuf := DeviceInterfaceMap[cmd.InterfaceID].DeviceNodeMap[k].GenerateGetRealVariables(v.Addr)
						log.Printf("tx buf is %+v\n",txBuf)
						//---------------发送-------------------------
						setting.SerialInterface.SerialPort[cmd.InterfaceID].Write(txBuf)
						//---------------等待接收----------------------
						//阻塞读
						rxBuf  := make([]byte, 256)
						rxTotalBuf := make([]byte,0)
						rxBufCnt := 0
						rxTotalBufCnt := 0
						timeOut,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Timeout)
						timer := time.NewTimer(time.Duration(timeOut)*time.Millisecond)
						for {
							select{
							//是否正确收到数据包
							case <-DeviceInterfaceMap[cmd.InterfaceID].DeviceNodeMap[k].AnalysisRx(v.Addr,v.VariableMap,rxTotalBuf,rxTotalBufCnt):
								{
									log.Println("rx ok")
									//通信帧延时
									interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Interval)
									time.Sleep(time.Duration(interval)*time.Millisecond)
									status = true
									goto LoopEmerg
								}
							//是否接收超时
							case <-timer.C:
								{
									log.Println("rx timeout")
									//通信帧延时
									interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Interval)
									time.Sleep(time.Duration(interval)*time.Millisecond)
									status = false
									goto LoopEmerg
								}
							//继续接收数据
							default:
								{
									rxBufCnt,_ = setting.SerialInterface.SerialPort[cmd.InterfaceID].Read(rxBuf)
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
					LoopEmerg:
					}
				}
				emergencyAckChan<- status
			}
		default:
			{
				select {
				case cmd := <-commonChan:
					{
						log.Println("common chan")
						//log.Printf("funName %s\n", cmd.funName)

						for k,v := range DeviceInterfaceMap[cmd.InterfaceID].DeviceNodeMap{
							if v.Addr == cmd.DeviceAddr {
								log.Printf("index is %d\n",k)
								//--------------组包---------------------------
								txBuf := DeviceInterfaceMap[cmd.InterfaceID].DeviceNodeMap[k].GenerateGetRealVariables(v.Addr)
								log.Printf("tx buf is %+v\n",txBuf)
								//---------------发送-------------------------
								setting.SerialInterface.SerialPort[cmd.InterfaceID].Write(txBuf)
								//---------------等待接收----------------------
								//阻塞读
								rxBuf  := make([]byte, 256)
								rxTotalBuf := make([]byte,0)
								rxBufCnt := 0
								rxTotalBufCnt := 0
								timeOut,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Timeout)
								timer := time.NewTimer(time.Duration(timeOut)*time.Millisecond)
								for {
									select{
										//是否正确收到数据包
										case <-DeviceInterfaceMap[cmd.InterfaceID].DeviceNodeMap[k].AnalysisRx(v.Addr,v.VariableMap,rxTotalBuf,rxTotalBufCnt):
										{
											log.Println("rx ok")
											//通信帧延时
											interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Interval)
											time.Sleep(time.Duration(interval)*time.Millisecond)
											goto Loop
										}
										//是否接收超时
										case <-timer.C:
										{
											log.Println("rx timeout")
											//通信帧延时
											interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Interval)
											time.Sleep(time.Duration(interval)*time.Millisecond)
											goto Loop
										}
										//继续接收数据
										default:
										{
											rxBufCnt,_ = setting.SerialInterface.SerialPort[cmd.InterfaceID].Read(rxBuf)
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

//func CommunicationManageAddEmergencyTest(){
//	cmd := CommunicationCmd{}
//
//	cmd.interfaceID = InterFaceID1
//	cmd.DeviceAddr = "2"
//	cmd.funName = "FCUGetRealData"
//	cmd.funPara = struct{
//		addr byte
//	}{0x02}
//
//	CommunicationManageAddEmergency(cmd)
//}

func CommunicationManagePoll(){

	cmd := CommunicationCmd{}

	for i:=0;i<DeviceInterfaceMap[InterFaceID0].DeviceNodeCnt;i++{

		cmd.InterfaceID = InterFaceID0
		cmd.DeviceAddr = DeviceInterfaceMap[InterFaceID0].DeviceNodeMap[i].Addr
		cmd.DeviceType = DeviceInterfaceMap[InterFaceID0].DeviceNodeMap[i].Type
		cmd.FunName = "GetDeviceRealVariables"

		CommunicationManageAdd(cmd)
	}
}