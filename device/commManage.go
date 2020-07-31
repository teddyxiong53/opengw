package device

import (
	"log"
	"time"
)

type CommunicationCmdTemplate struct {
	CollInterfaceName string    //采集接口名称
	DeviceAddr    string 		//采集接口下设备地址
	FunName       string
	FunPara       interface{}
}

type CommunicationManageTemplate struct{
	EmergencyRequestChan chan CommunicationCmdTemplate
	CommonRequestChan    chan CommunicationCmdTemplate
	EmergencyAckChan     chan bool
	//RxChan               chan bool //接收正确
}

func NewCommunicationManageTemplate() *CommunicationManageTemplate{

	return &CommunicationManageTemplate{
		EmergencyRequestChan:make(chan CommunicationCmdTemplate, 100),
		CommonRequestChan: make(chan CommunicationCmdTemplate, 1),
		EmergencyAckChan: make(chan bool, 1),
	}
}

func CommunicationManageInit() {

	//go CommunicationManageDel()
}

func (c *CommunicationManageTemplate)CommunicationManageAddCommon(cmd CommunicationCmdTemplate) {

	c.CommonRequestChan <- cmd
}

func (c *CommunicationManageTemplate)CommunicationManageAddEmergency(cmd CommunicationCmdTemplate) bool {

	c.EmergencyRequestChan <- cmd

	return <-c.EmergencyAckChan
}

func (c *CommunicationManageTemplate)CommunicationManageDel() {

	for {
		select {
		case cmd := <-c.EmergencyRequestChan:
			{
				log.Println("emergency chan")
				log.Printf("funName %s\n", cmd.FunName)
				var status bool = false
				for _, c := range CollectInterfaceMap {
					if c.CollInterfaceName == cmd.CollInterfaceName {
						for k,v := range c.DeviceNodeMap{
							if v.Addr == cmd.DeviceAddr {
								log.Printf("index is %d\n", k)
								//--------------组包---------------------------
								txBuf := v.GenerateGetRealVariables(v.Addr)
								log.Printf("tx buf is %+v\n", txBuf)
								//---------------发送-------------------------
								for _,v := range CommunicationSerialMap{
									if v.Name == c.CommInterfaceName{
										v.WriteData(txBuf)
									}
								}
								//---------------等待接收----------------------
								//阻塞读
								rxBuf := make([]byte, 256)
								rxTotalBuf := make([]byte, 0)
								rxBufCnt := 0
								rxTotalBufCnt := 0
								//timeOut,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Timeout)
								timer := time.NewTimer(time.Duration(100) * time.Millisecond)
								for {
									select {
									//是否正确收到数据包
									case <-v.AnalysisRx(v.Addr, v.VariableMap, rxTotalBuf, rxTotalBufCnt):
										{
											log.Println("rx ok")
											//通信帧延时
											//interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Interval)
											//time.Sleep(time.Duration(interval)*time.Millisecond)
											status = true
											goto LoopEmerg
										}
									//是否接收超时
									case <-timer.C:
										{
											log.Println("rx timeout")
											//通信帧延时
											//interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Interval)
											//time.Sleep(time.Duration(interval)*time.Millisecond)
											status = false
											goto LoopEmerg
										}
									//继续接收数据
									default:
										{
											//rxBufCnt,_ = setting.SerialInterface.SerialPort[cmd.InterfaceID].Read(rxBuf)
											for _,v := range CommunicationSerialMap{
												if v.Name == c.CommInterfaceName{
													rxBufCnt = v.ReadData(rxBuf)
												}
											}
											if rxBufCnt > 0 {
												rxTotalBufCnt += rxBufCnt
												//追加接收的数据到接收缓冲区
												rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
												//清除本地接收数据
												rxBufCnt = 0
												log.Printf("rxbuf %+v\n", rxTotalBuf)
											}
										}
									}
								}
							LoopEmerg:
							}
						}
					}
				}
				c.EmergencyAckChan <- status
			}
		default:
			{
				select {
				case cmd := <-c.CommonRequestChan:
					{
						log.Println("common chan")
						log.Printf("funName %s\n", cmd.FunName)
						for _, v := range CollectInterfaceMap {
							if v.CollInterfaceName == cmd.CollInterfaceName {
								for k, v := range v.DeviceNodeMap {
									if v.Addr == cmd.DeviceAddr {
										log.Printf("index is %d\n", k)
										//--------------组包---------------------------
										txBuf := v.GenerateGetRealVariables(v.Addr)
										log.Printf("tx buf is %+v\n", txBuf)
										//---------------发送-------------------------
										//setting.SerialInterface.SerialPort[cmd.InterfaceID].Write(txBuf)
										//---------------等待接收----------------------
										//阻塞读
										rxBuf := make([]byte, 256)
										rxTotalBuf := make([]byte, 0)
										rxBufCnt := 0
										rxTotalBufCnt := 0
										//timeOut,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Timeout)
										timer := time.NewTimer(time.Duration(100) * time.Millisecond)
										for {
											select {
											//是否正确收到数据包
											case <-v.AnalysisRx(v.Addr, v.VariableMap, rxTotalBuf, rxTotalBufCnt):
												{
													log.Println("rx ok")
													//通信帧延时
													//interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Interval)
													//time.Sleep(time.Duration(interval)*time.Millisecond)
													goto Loop
												}
											//是否接收超时
											case <-timer.C:
												{
													log.Println("rx timeout")
													//通信帧延时
													//interval,_ := strconv.Atoi(setting.SerialInterface.SerialParam[cmd.InterfaceID].Interval)
													//time.Sleep(time.Duration(interval)*time.Millisecond)
													goto Loop
												}
											//继续接收数据
											default:
												{
													//rxBufCnt,_ = setting.SerialInterface.SerialPort[cmd.InterfaceID].Read(rxBuf)
													if rxBufCnt > 0 {
														rxTotalBufCnt += rxBufCnt
														//追加接收的数据到接收缓冲区
														rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
														//清除本地接收数据
														rxBufCnt = 0
														log.Printf("rxbuf %+v\n", rxTotalBuf)
													}
												}
											}
										}
									}
								}
							Loop:
							}
						}
					}
				default:
					time.Sleep(10 * time.Millisecond)
				}
			}
		}
	}
}

func (c *CommunicationManageTemplate)CommunicationManagePoll() {

	//cmd := CommunicationCmdTemplate{}
	//
	//for i:=0;i<CollectInterfaceMap[InterFaceID0].DeviceNodeCnt;i++{
	//
	//	cmd.InterfaceID = InterFaceID0
	//	cmd.DeviceAddr = CollectInterfaceMap[InterFaceID0].DeviceNodeMap[i].Addr
	//	cmd.DeviceType = CollectInterfaceMap[InterFaceID0].DeviceNodeMap[i].Type
	//	cmd.FunName = "GetDeviceRealVariables"
	//
	//	CommunicationManageAdd(cmd)
	//}
}
