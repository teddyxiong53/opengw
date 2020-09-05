package device

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"goAdapter/setting"
	"log"
	"strconv"
	"time"
)

type CommunicationCmdTemplate struct {
	CollInterfaceName    	string    //采集接口名称
	DeviceAddr    			string 		//采集接口下设备地址
	FunName       			string
	FunPara       			interface{}
}

type CommunicationManageTemplate struct{
	EmergencyRequestChan chan CommunicationCmdTemplate
	CommonRequestChan    chan CommunicationCmdTemplate
	EmergencyAckChan     chan bool
	CollInterface        *CollectInterfaceTemplate
}

func NewCommunicationManageTemplate(coll *CollectInterfaceTemplate) *CommunicationManageTemplate {

	template := &CommunicationManageTemplate{
		EmergencyRequestChan: make(chan CommunicationCmdTemplate, 1),
		CommonRequestChan:    make(chan CommunicationCmdTemplate, 100),
		EmergencyAckChan:     make(chan bool, 1),
		CollInterface:        coll,
	}

	return template
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
				status := false
				for _, c := range CollectInterfaceMap {
					if c.CollInterfaceName == cmd.CollInterfaceName {
						for k,v := range c.DeviceNodeMap{
							if v.Addr == cmd.DeviceAddr {
								log.Printf("index is %d\n", k)
								step := 0
								for{
									//--------------组包---------------------------
									txBuf,ok := v.GenerateGetRealVariables(v.Addr,step)
									if ok == false{
										//log.Printf("getVariables false\n")
										goto LoopEmergency
									}
									step++
									log.Printf("tx buf is %X\n", txBuf)
									//---------------发送-------------------------
									var timeout int
									var interval int
									for _,v := range CommunicationSerialMap{
										if v.Name == c.CommInterfaceName{
											v.WriteData(txBuf)
											timeout,_ = strconv.Atoi(v.Param.Timeout)
											interval,_ = strconv.Atoi(v.Param.Interval)
										}
									}
									v.CommTotalCnt++
									//---------------等待接收----------------------
									//阻塞读
									rxBuf := make([]byte, 256)
									rxTotalBuf := make([]byte, 0)
									rxBufCnt := 0
									rxTotalBufCnt := 0
									timerOut := time.NewTimer(time.Duration(timeout) * time.Millisecond)
									for {
										select {
										//是否正确收到数据包
										case <-v.AnalysisRx(v.Addr, v.VariableMap, rxTotalBuf, rxTotalBufCnt):
											{
												log.Println("rx ok")
												log.Printf("rxbuf %X\n", rxTotalBuf)
												//通信帧延时
												time.Sleep(time.Duration(interval)*time.Millisecond)
												v.CommSuccessCnt++
												v.CurCommFailCnt = 0
												v.CommStatus = "onLine"
												v.LastCommRTC = time.Now().Format("2006-01-02 15:04:05")
												rxTotalBufCnt = 0
												rxTotalBuf = rxTotalBuf[0:0]
												goto LoopEmergencyStep
											}
										//是否接收超时
										case <-timerOut.C:
											{
												log.Println("rx timeout")
												//通信帧延时
												time.Sleep(time.Duration(interval)*time.Millisecond)
												v.CurCommFailCnt++
												if v.CurCommFailCnt >= c.OfflinePeriod{
													v.CurCommFailCnt = 0
													v.CommStatus = "offLine"
												}
												rxTotalBufCnt = 0
												rxTotalBuf = rxTotalBuf[0:0]
												goto LoopEmergencyStep
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
												}
											}
										}
									}
								LoopEmergencyStep:
								}
							LoopEmergency:
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
						//setting.Logrus.Printf("%v:common chan\n",c.CollInterfaceName)
						setting.Loger.WithFields(logrus.Fields{
							"collName": c.CollInterface.CollInterfaceName,
						}).Info("common chan")

						for _, coll := range CollectInterfaceMap {
							if coll.CollInterfaceName == cmd.CollInterfaceName {
								for _,v := range coll.DeviceNodeMap{
									if v.Addr == cmd.DeviceAddr {
										log.Printf("%v:addr %v\n", coll.CollInterfaceName,v.Addr)
										step := 0
										for{
											//--------------组包---------------------------
											txBuf,ok := v.GenerateGetRealVariables(v.Addr,step)
											if ok == false{
												log.Printf("%v:getVariables finish\n",coll.CollInterfaceName)
												goto LoopCommon
											}
											step++
											log.Printf("%v:txbuf %X\n", coll.CollInterfaceName,txBuf)

											CommunicationMessage := CommunicationMessageTemplate{
												CollName: coll.CollInterfaceName,
												TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
												Direction: "send",
												Content: fmt.Sprintf("%X",txBuf),
											}
											if len(c.CollInterface.CommMessage) < 1024{
												c.CollInterface.CommMessage = append(c.CollInterface.CommMessage,CommunicationMessage)
											}else{
												c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
												c.CollInterface.CommMessage = append(c.CollInterface.CommMessage,CommunicationMessage)
											}

											//---------------发送-------------------------
											var timeout int
											var interval int
											//判断是否是串口采集
											for _,v := range CommunicationSerialMap{
												if v.Name == coll.CommInterfaceName{
													v.WriteData(txBuf)
													timeout,_ = strconv.Atoi(v.Param.Timeout)
													interval,_ = strconv.Atoi(v.Param.Interval)
												}
											}
											v.CommTotalCnt++
											//---------------等待接收----------------------
											//阻塞读
											rxBuf := make([]byte, 256)
											rxTotalBuf := make([]byte, 0)
											rxBufCnt := 0
											rxTotalBufCnt := 0
											timerOut := time.NewTimer(time.Duration(timeout) * time.Millisecond)
											for {
												select {
												//是否正确收到数据包
												case <-v.AnalysisRx(v.Addr, v.VariableMap, rxTotalBuf, rxTotalBufCnt):
													{
														log.Printf("%v:rx ok\n",coll.CollInterfaceName)
														log.Printf("%v:rxbuf %X\n", coll.CollInterfaceName,rxTotalBuf)

														CommunicationMessage := CommunicationMessageTemplate{
															CollName: coll.CollInterfaceName,
															TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
															Direction: "receive",
															Content: fmt.Sprintf("%X",rxTotalBuf),
														}
														if len(c.CollInterface.CommMessage) < 1024{
															c.CollInterface.CommMessage = append(c.CollInterface.CommMessage,CommunicationMessage)
														}else{
															c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
															c.CollInterface.CommMessage = append(c.CollInterface.CommMessage,CommunicationMessage)
														}

														//通信帧延时
														time.Sleep(time.Duration(interval)*time.Millisecond)
														v.CommSuccessCnt++
														v.CurCommFailCnt = 0
														v.CommStatus = "onLine"
														v.LastCommRTC = time.Now().Format("2006-01-02 15:04:05")
														rxTotalBufCnt = 0
														rxTotalBuf = rxTotalBuf[0:0]
														goto LoopCommonStep
													}
												//是否接收超时
												case <-timerOut.C:
													{
														CommunicationMessage := CommunicationMessageTemplate{
															CollName: coll.CollInterfaceName,
															TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
															Direction: "receive",
															Content: fmt.Sprintf("%X",rxTotalBuf),
														}
														if len(c.CollInterface.CommMessage) < 1024{
															c.CollInterface.CommMessage = append(c.CollInterface.CommMessage,CommunicationMessage)
														}else{
															c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
															c.CollInterface.CommMessage = append(c.CollInterface.CommMessage,CommunicationMessage)
														}

													    log.Printf("%v,rx timeout\n",coll.CollInterfaceName)
														//通信帧延时
														time.Sleep(time.Duration(interval)*time.Millisecond)
														v.CurCommFailCnt++
														if v.CurCommFailCnt >= coll.OfflinePeriod{
															v.CurCommFailCnt = 0
															v.CommStatus = "offLine"
														}
														rxTotalBufCnt = 0
														rxTotalBuf = rxTotalBuf[0:0]
														goto LoopCommonStep
													}
												//继续接收数据
												default:
													{
														for _,v := range CommunicationSerialMap{
															if v.Name == coll.CommInterfaceName{
																rxBufCnt = v.ReadData(rxBuf)
															}
														}
														if rxBufCnt > 0 {
															rxTotalBufCnt += rxBufCnt
															//追加接收的数据到接收缓冲区
															rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
															//清除本次接收数据
															rxBufCnt = 0

															//log.Printf("rxbuf %X\n",rxTotalBuf)
														}
													}
												}
											}
											LoopCommonStep:
										}
										LoopCommon:
									}
								}
								coll.DeviceNodeOnlineCnt = 0
								for _,v := range coll.DeviceNodeMap{
									if v.CommStatus == "onLine"{
										coll.DeviceNodeOnlineCnt++
									}
								}
							}
						}

						//更新设备在线率
						deviceTotalCnt := 0
						deviceOnlineCnt := 0
						for _,v := range CollectInterfaceMap{
							deviceTotalCnt += v.DeviceNodeCnt
							deviceOnlineCnt += v.DeviceNodeOnlineCnt
						}
						if deviceOnlineCnt == 0{
							setting.SystemState.DeviceOnline = "0"
						}else{
							setting.SystemState.DeviceOnline = fmt.Sprintf("%2.1f",float32(deviceOnlineCnt*100.0/deviceTotalCnt))
						}

						//更新设备丢包率
						deviceCommTotalCnt := 0
						deviceCommLossCnt := 0
						for _,v := range CollectInterfaceMap{
							for _,v := range v.DeviceNodeMap{
								deviceCommTotalCnt += v.CommTotalCnt
								deviceCommLossCnt += v.CommTotalCnt-v.CommSuccessCnt
							}
						}
						if deviceCommLossCnt == 0{
							setting.SystemState.DevicePacketLoss = "0"
						}else{
							setting.SystemState.DevicePacketLoss = fmt.Sprintf("%2.1f",float32(deviceCommLossCnt*100.0/deviceCommTotalCnt))
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

	cmd := CommunicationCmdTemplate{}
	//对采集接口进行遍历
	for _, coll := range CollectInterfaceMap {
		if coll.CollInterfaceName == c.CollInterface.CollInterfaceName {
			//对采集接口下设备进行遍历
			for _,v := range coll.DeviceNodeMap{
				cmd.CollInterfaceName = coll.CollInterfaceName
				cmd.DeviceAddr = v.Addr
				cmd.FunName = "GetDeviceRealVariables"
				c.CommunicationManageAddCommon(cmd)
			}
		}
	}
}
