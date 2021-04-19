package device

import (
	"fmt"
	"goAdapter/setting"
	"log"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type CommunicationCmdTemplate struct {
	CollInterfaceName string //采集接口名称
	DeviceName        string //采集接口下设备名称
	FunName           string
	FunPara           string
}

type CommunicationManageTemplate struct {
	EmergencyRequestChan chan CommunicationCmdTemplate
	CommonRequestChan    chan CommunicationCmdTemplate
	EmergencyAckChan     chan bool
	CollInterface        *CollectInterfaceTemplate
	PacketChan           chan []byte
}

var CommunicationManage = make([]*CommunicationManageTemplate, 0)

func NewCommunicationManageTemplate(coll *CollectInterfaceTemplate) *CommunicationManageTemplate {

	template := &CommunicationManageTemplate{
		EmergencyRequestChan: make(chan CommunicationCmdTemplate, 1),
		CommonRequestChan:    make(chan CommunicationCmdTemplate, 100),
		EmergencyAckChan:     make(chan bool, 1),
		PacketChan:           make(chan []byte, 100), //最多连续接收100帧数据
		CollInterface:        coll,
	}
	//启动接收协程
	go template.AnalysisRx()

	return template
}

func (c *CommunicationManageTemplate) CommunicationManageAddCommon(cmd CommunicationCmdTemplate) {

	c.CommonRequestChan <- cmd
}

func (c *CommunicationManageTemplate) CommunicationManageAddEmergency(cmd CommunicationCmdTemplate) bool {

	c.EmergencyRequestChan <- cmd

	return <-c.EmergencyAckChan
}

func (c *CommunicationManageTemplate) AnalysisRx() {

	//阻塞读
	rxBuf := make([]byte, 1024)
	rxBufCnt := 0

	serialPort := &CommunicationSerialTemplate{}

	for k, v := range CommunicationSerialMap {
		if v.Name == c.CollInterface.CommInterfaceName {
			serialPort = CommunicationSerialMap[k]
		}
	}

	for {
		//阻塞读
		rxBufCnt = serialPort.ReadData(rxBuf)
		if rxBufCnt > 0 {
			//log.Printf("curRxBufCnt %v,", rxBufCnt)
			//log.Printf("CurRxBuf %X\n", rxBuf[:rxBufCnt])

			//rxTotalBufCnt += rxBufCnt
			//追加接收的数据到接收缓冲区
			//rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
			c.PacketChan <- rxBuf[:rxBufCnt]
			//log.Printf("chanLen %v\n", len(c.PacketChan))

			//清除本次接收数据
			rxBufCnt = 0
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *CommunicationManageTemplate) CommunicationStateMachine(cmd CommunicationCmdTemplate) bool {

	status := false

	startT := time.Now() //计算当前时间
	for _, v := range c.CollInterface.DeviceNodeMap {
		if v.Name == cmd.DeviceName {
			setting.Logger.Debugf("%v:name %v\n", c.CollInterface.CollInterfaceName, v.Name)
			step := 0
			for {
				//--------------组包---------------------------
				txBuf := make([]byte, 0)
				ok := false
				con := false
				if cmd.FunName == "GetDeviceRealVariables" {
					txBuf, ok, con = v.GenerateGetRealVariables(v.Addr, step)
				} else {
					txBuf, ok, con = v.DeviceCustomCmd(cmd.DeviceName, cmd.FunName, cmd.FunPara, step)
					if ok == false {
						setting.Logger.Errorln("DeviceCustomCmd false")
						goto LoopCommon
					}
				}
				if ok == false {
					setting.Logger.Debugf("%v:cmd finish\n", c.CollInterface.CollInterfaceName)
					goto LoopCommon
				}

				step++
				setting.Logger.Debugf("%v:txbuf %X\n", c.CollInterface.CollInterfaceName, txBuf)

				CommunicationMessage := CommunicationMessageTemplate{
					CollName:  c.CollInterface.CollInterfaceName,
					TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
					Direction: "send",
					Content:   fmt.Sprintf("%X", txBuf),
				}
				if len(c.CollInterface.CommMessage) < 1024 {
					c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, CommunicationMessage)
				} else {
					c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
					c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, CommunicationMessage)
				}

				//---------------发送-------------------------
				var timeout int
				var interval int
				//判断是否是串口采集
				for _, v := range CommunicationSerialMap {
					if v.Name == c.CollInterface.CommInterfaceName {
						v.WriteData(txBuf)
						timeout, _ = strconv.Atoi(v.Param.Timeout)
						interval, _ = strconv.Atoi(v.Param.Interval)
					}
				}
				timerOut := time.NewTimer(time.Duration(timeout) * time.Millisecond)
				v.CommTotalCnt++
				//---------------等待接收----------------------
				//阻塞读
				rxBuf := make([]byte, 256)
				rxTotalBuf := make([]byte, 0)
				rxBufCnt := 0
				rxTotalBufCnt := 0
				for {
					select {
					//是否接收超时
					case <-timerOut.C:
						{
							timerOut.Stop()
							CommunicationMessage := CommunicationMessageTemplate{
								CollName:  c.CollInterface.CollInterfaceName,
								TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
								Direction: "receive",
								Content:   fmt.Sprintf("%X", rxTotalBuf),
							}
							if len(c.CollInterface.CommMessage) < 1024 {
								c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, CommunicationMessage)
							} else {
								c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
								c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, CommunicationMessage)
							}

							setting.Logger.Debugf("%v,rx timeout\n", c.CollInterface.CollInterfaceName)
							setting.Logger.Debugf("%v:rxbuf %X\n", c.CollInterface.CollInterfaceName, rxTotalBuf)
							//通信帧延时
							//time.Sleep(time.Duration(interval) * time.Millisecond)

							v.CurCommFailCnt++
							if v.CurCommFailCnt >= c.CollInterface.OfflinePeriod {
								v.CurCommFailCnt = 0
								//设备从上线变成离线
								if v.CommStatus == "onLine" {
									c.CollInterface.OfflineReportChan <- v.Addr
									v.CommStatus = "offLine"
								}
							}
							rxTotalBufCnt = 0
							rxTotalBuf = rxTotalBuf[0:0]
							status = false
							goto LoopCommonStep
						}
					//是否正确收到数据包
					case <-v.AnalysisRx(v.Addr, v.VariableMap, rxTotalBuf, rxTotalBufCnt):
						{
							timerOut.Stop()
							setting.Logger.Debugf("%v:rx ok\n", c.CollInterface.CollInterfaceName)
							setting.Logger.Debugf("%v:rxbuf %X\n", c.CollInterface.CollInterfaceName, rxTotalBuf)

							CommunicationMessage := CommunicationMessageTemplate{
								CollName:  c.CollInterface.CollInterfaceName,
								TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
								Direction: "receive",
								Content:   fmt.Sprintf("%X", rxTotalBuf),
							}
							if len(c.CollInterface.CommMessage) < 1024 {
								c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, CommunicationMessage)
							} else {
								c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
								c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, CommunicationMessage)
							}
							if len(c.CommonRequestChan) > 0 {
								//通信帧延时
								time.Sleep(time.Duration(interval) * time.Millisecond)
							} else {
								//是否后续有通信帧
								if con == true {
									//通信帧延时
									time.Sleep(time.Duration(interval) * time.Millisecond)
								}
							}

							//设备从离线变成上线
							if v.CommStatus == "offLine" {
								c.CollInterface.OnlineReportChan <- v.Addr
							}

							//防止Chan阻塞
							if len(c.CollInterface.PropertyReportChan) >= 100 {
								<-c.CollInterface.PropertyReportChan
							}
							c.CollInterface.PropertyReportChan <- v.Addr
							//log.Printf("reportChan %v\n", len(c.CollInterface.PropertyReportChan))

							v.CommSuccessCnt++
							v.CurCommFailCnt = 0
							v.CommStatus = "onLine"
							v.LastCommRTC = time.Now().Format("2006-01-02 15:04:05")

							rxTotalBufCnt = 0
							rxTotalBuf = rxTotalBuf[0:0]
							status = true
							goto LoopCommonStep
						}
					//继续接收数据
					case rxBuf = <-c.PacketChan:
						{
							rxBufCnt = len(rxBuf)
							if rxBufCnt > 0 {
								rxTotalBufCnt += rxBufCnt
								//追加接收的数据到接收缓冲区
								rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
								//清除本次接收数据
								rxBufCnt = 0
								rxBuf = rxBuf[0:0]
								//log.Printf("rxTotalBuf %X\n", rxTotalBuf)
							}
						}
					}
				}
				//}
			LoopCommonStep:
			}
		LoopCommon:
		}
	}
	tc := time.Since(startT) //计算耗时
	setting.Logger.Debugf("%v: ,time cost = %v\n", c.CollInterface.CollInterfaceName, tc)

	//更新设备在线数量
	c.CollInterface.DeviceNodeOnlineCnt = 0
	for _, v := range c.CollInterface.DeviceNodeMap {
		if v.CommStatus == "onLine" {
			c.CollInterface.DeviceNodeOnlineCnt++
		}
	}

	return status
}

func (c *CommunicationManageTemplate) CommunicationManageDel() {

	for {
		select {
		case cmd := <-c.EmergencyRequestChan:
			{
				log.Println("emergency chan")

				setting.Logger.WithFields(logrus.Fields{
					"collName":   c.CollInterface.CollInterfaceName,
					"deviceName": cmd.DeviceName,
					"funName":    cmd.FunName,
				}).Info("emergency chan")
				status := false

				status = c.CommunicationStateMachine(cmd)

				GetDeviceOnline()
				GetDevicePacketLoss()

				c.EmergencyAckChan <- status
			}
		default:
			{
				select {
				case cmd := <-c.CommonRequestChan:
					{
						setting.Logger.Debugf("%v:,commChanLen %v\n", c.CollInterface.CollInterfaceName, len(c.CommonRequestChan))
						c.CommunicationStateMachine(cmd)

						GetDeviceOnline()
						GetDevicePacketLoss()
					}
				default:
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}
}

func GetDeviceOnline() {

	//更新设备在线率
	deviceTotalCnt := 0
	deviceOnlineCnt := 0
	for _, v := range CollectInterfaceMap {
		deviceTotalCnt += v.DeviceNodeCnt
		deviceOnlineCnt += v.DeviceNodeOnlineCnt
	}
	if deviceOnlineCnt == 0 {
		setting.SystemState.DeviceOnline = "0"
	} else {
		setting.SystemState.DeviceOnline = fmt.Sprintf("%2.1f", float32(deviceOnlineCnt*100.0/deviceTotalCnt))
	}
}

func GetDevicePacketLoss() {

	//更新设备丢包率
	deviceCommTotalCnt := 0
	deviceCommLossCnt := 0
	for _, v := range CollectInterfaceMap {
		for _, v := range v.DeviceNodeMap {
			deviceCommTotalCnt += v.CommTotalCnt
			deviceCommLossCnt += v.CommTotalCnt - v.CommSuccessCnt
		}
	}
	if deviceCommLossCnt == 0 {
		setting.SystemState.DevicePacketLoss = "0"
	} else {
		setting.SystemState.DevicePacketLoss = fmt.Sprintf("%2.1f", float32(deviceCommLossCnt*100.0/deviceCommTotalCnt))
	}
}

func (c *CommunicationManageTemplate) CommunicationManagePoll() {

	cmd := CommunicationCmdTemplate{}
	//对采集接口进行遍历
	for _, coll := range CollectInterfaceMap {
		if coll.CollInterfaceName == c.CollInterface.CollInterfaceName {
			//对采集接口下设备进行遍历
			for _, v := range coll.DeviceNodeMap {
				cmd.CollInterfaceName = coll.CollInterfaceName
				cmd.DeviceName = v.Name
				cmd.FunName = "GetDeviceRealVariables"
				c.CommunicationManageAddCommon(cmd)
			}
			setting.Logger.Debugf("commChanTotalLen %v\n", len(c.CommonRequestChan))
		}
	}
}
