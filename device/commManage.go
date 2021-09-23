package device

import (
	"errors"
	"fmt"
	"goAdapter/config"
	"goAdapter/pkg/mylog"
	"goAdapter/pkg/system"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/jasonlvhit/gocron"
	"go.uber.org/zap"
)

const (
	CommunicationManageMessageMaxCnt = 1024
)

type CommunicationCmdTemplate struct {
	CollInterfaceName string //采集接口名称
	DeviceName        string //采集接口下设备名称
	DeviceIndex       int
	FunName           LUAFUNC
	FunPara           string
}

type CommunicationRxTemplate struct {
	Err   error
	RxBuf []byte
}

type CommunicationManageTemplate struct {
	EmergencyRequestChan chan CommunicationCmdTemplate
	CommonRequestChan    chan CommunicationCmdTemplate
	EmergencyAckChan     chan CommunicationRxTemplate
	CollInterface        *CollectInterfaceTemplate
	PacketChan           chan []byte
	Signal               chan struct{}
	Ready                chan struct{}
}

type State uint8

const (
	Start State = iota
	Generate
	Send
	Wait
	WaitSuccess
	WaitFail
	Stop
)

var CommunicationManage = CommManger{
	ManagerTemp: make(map[string]*CommunicationManageTemplate),
	Collectors:  make(chan *CollectInterfaceStatus, 20),
}

type CommManger struct {
	ManagerTemp map[string]*CommunicationManageTemplate
	Collectors  chan *CollectInterfaceStatus
}

func NewCommunicationManageTemplate(coll *CollectInterfaceTemplate) *CommunicationManageTemplate {

	template := &CommunicationManageTemplate{
		EmergencyRequestChan: make(chan CommunicationCmdTemplate, 1),
		CommonRequestChan:    make(chan CommunicationCmdTemplate, 100),
		EmergencyAckChan:     make(chan CommunicationRxTemplate, 1),
		PacketChan:           make(chan []byte, 100), //最多连续接收100帧数据
		CollInterface:        coll,
		Ready:                make(chan struct{}, 1),
		Signal:               make(chan struct{}, 1),
	}
	if coll.CommInterface.Error() == nil {
		//启动接收协程
		mylog.ZAPS.Infof("采集接口【%s】打开成功，启动接收协程！", coll.CollInterfaceName)
		go template.ReadRx()
	}
	return template
}

func addHandler(scheduler *gocron.Scheduler, collect *CollectInterfaceStatus) {
	comm := collect.Tmp.CommInterface
	if comm == nil {
		mylog.ZAPS.Errorf("通讯口【%s】未绑定到接口【%s】",
			collect.Tmp.CommInterfaceName, collect.Tmp.CollInterfaceName)
		return
	}
	if err := comm.Open(); err != nil {
		mylog.ZAPS.Errorf("通讯口【%s】打开错误", comm.GetName())
		return
	}
	manager := NewCommunicationManageTemplate(collect.Tmp)
	CommunicationManage.ManagerTemp[collect.Tmp.CollInterfaceName] = manager

	go manager.CommunicationManagePoll(collect.Tmp.PollPeriod)
	go manager.CommunicationManageDel()
	mylog.ZAPS.Infof("添加采集【%s】到定时任务,Addr:%p", collect.Tmp.CollInterfaceName, manager)
}

func delHandler(scheduler *gocron.Scheduler, collect *CollectInterfaceStatus) {
	managerRemove := CommunicationManage.ManagerTemp[collect.Tmp.CollInterfaceName]
	if managerRemove != nil {
		//广播到所有监听管理模板的信号channel的goroutine
		close(managerRemove.Signal)
		managerRemove.CollInterface.CommInterface.Close()
		mylog.ZAPS.Infof("取消采集【%s】定时任务,Addr:%p", collect.Tmp.CollInterfaceName, managerRemove)
	}
}

func ScheduleJob(scheduler *gocron.Scheduler, quitChan chan struct{}) {
	go func() {
		for {
			select {
			case collect := <-CommunicationManage.Collectors:
				switch collect.ACT {
				case ADD:
					addHandler(scheduler, collect)
				case DELETE:
					delHandler(scheduler, collect)

				case UPDATE:
					//TODO 更新要稍微考虑一下，目前更新是先DELETE然后Add来实现

				}
			case <-quitChan:
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

}

func (c *CommunicationManageTemplate) CommunicationManageAddEmergency(cmd CommunicationCmdTemplate) CommunicationRxTemplate {

	//TODO 这里会不会阻塞住得看下超时的判断
	c.EmergencyRequestChan <- cmd

	return <-c.EmergencyAckChan
}

func (c *CommunicationManageTemplate) CommunicationManageMessageAdd(dir string, buf []byte) {
	CommunicationMessage := &CommunicationMessageTemplate{
		CollName:  c.CollInterface.CollInterfaceName,
		TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
		Direction: dir,
		Content:   fmt.Sprintf("%X", buf),
	}
	if len(c.CollInterface.CommMessage) < CommunicationManageMessageMaxCnt {
		c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, CommunicationMessage)
	} else {
		c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
		c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, CommunicationMessage)
	}
}

func (c *CommunicationManageTemplate) ReadRx() {

	//阻塞读
	rxBuf := make([]byte, 512)
	var rxBufCnt int
	var err error

	for range c.Ready {
		select {
		case <-c.Signal:
			return
		default:
			//等待数据已经完全写入缓冲区
			time.Sleep(time.Duration(config.Cfg.SerialCfg.BufferReadDelay) * time.Millisecond)

			//阻塞读
			rxBufCnt, err = c.CollInterface.CommInterface.Read(rxBuf)
			if err != nil && err != io.EOF {
				mylog.ZAP.Error("comm read error", zap.String("collinterface", c.CollInterface.CollInterfaceName), zap.String("comm", c.CollInterface.CommInterfaceName), zap.Error(err))
				continue
			}
			if rxBufCnt > 0 {

				c.PacketChan <- rxBuf[:rxBufCnt]
				//time.Sleep(time.Millisecond * 100)
			}

		}

	}
}

func (c *CommunicationManageTemplate) CommunicationStateMachine(cmd CommunicationCmdTemplate) (rxData CommunicationRxTemplate) {

	startT := time.Now() //计算当前时间

	if c.CollInterface == nil || c.CollInterface.DeviceNodes == nil {
		rxData.Err = errors.New("collect or devicenodes is nil")
		return
	}

	node := c.CollInterface.DeviceNodes[cmd.DeviceIndex]
	step := 0
	var txBuf []byte
	var hasFrame = true
	var err error
	var state State = Start
	var nodeIndex = -1
OUT:
	for {
		switch state {
		case Start:
			{
				for k, v := range c.CollInterface.DeviceNodes {
					if v.Name == cmd.DeviceName {
						nodeIndex = k
					}
				}
				if nodeIndex >= 0 {
					state = Generate
				} else {
					state = Stop
				}
			}
		case Generate:
			{
				if cmd.FunName == GETREAL {
					txBuf, hasFrame, err = node.GenerateGetRealVariables(node.Addr, step)
					if err != nil {
						rxData.Err = fmt.Errorf("lua %s generate error:%v", cmd.FunName, err)
						state = Stop
					} else {
						state = Send
						step++
					}

				} else {
					txBuf, hasFrame, err = node.DeviceCustomCmd(node.Addr, cmd.FunName, cmd.FunPara, step)
					if err != nil {
						rxData.Err = fmt.Errorf("device custom  cmd error:%v", err)
						state = Stop

					} else {
						state = Send
						step++
					}

				}
			}
		case Send:
			{ //---------------发送-------------------------
				//判断是否是串口采集
				_, err := c.CollInterface.CommInterface.Write(txBuf)
				//如果写入错误很有可能是串口关闭了
				if err != nil {
					mylog.ZAPS.Errorf("write data to comm %v error:%v", c.CollInterface.CommInterfaceName, err)
					return
				}
				c.Ready <- struct{}{}
				c.CommunicationManageMessageAdd("send", txBuf)
				node.CommTotalCnt++
				mylog.ZAPS.Debugf("【SEND】接口【%s】% X", cmd.CollInterfaceName, txBuf)
				state = Wait
			}

		case Wait:
			{
				var (
					rxBuf         []byte
					rxTotalBuf    []byte
					rxBufCnt      int
					rxTotalBufCnt int
					timeout, _    = strconv.Atoi(c.CollInterface.CommInterface.GetTimeOut())
					timer         = time.NewTimer(time.Duration(timeout) * time.Millisecond)
				)

				for {
					select {
					//是否接收超时
					case <-timer.C:
						{
							state = Stop
							CommunicationMessage := CommunicationMessageTemplate{
								CollName:  c.CollInterface.CollInterfaceName,
								TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
								Direction: "receive",
								Content:   fmt.Sprintf("接收数据超时了,超时阈值:%d ms", timeout),
							}
							if len(c.CollInterface.CommMessage) < 1024 {
								c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, &CommunicationMessage)
							} else {
								c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
								c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, &CommunicationMessage)
							}

							//通信帧延时
							//time.Sleep(time.Duration(interval) * time.Millisecond)

							node.CurCommFailCnt++
							rxTotalBufCnt = 0
							rxTotalBuf = []byte{}

							//如果失败次数大于offlinePeriod 就放弃这个设备了
							if node.CurCommFailCnt >= c.CollInterface.OfflinePeriod {
								node.CurCommFailCnt = 0
								//设备从上线变成离线
								if node.CommStatus == ONLINE {
									if len(c.CollInterface.OfflineReportChan) == 100 {
										<-c.CollInterface.OfflineReportChan
									}
									c.CollInterface.OfflineReportChan <- node.Name
									node.CommStatus = OFFLINE
									c.CollInterface.DeviceNodeOnlineCnt--
								}
								return
							}
							state = Start
							log.Println(color.MagentaString("采集器【%v】-设备【%s】接收数据超时，失败次数:%d 总次数:%d", c.CollInterface.CollInterfaceName, node.Name, node.CurCommFailCnt, c.CollInterface.OfflinePeriod))
							goto OUT
						}

					//继续接收数据
					case rxBuf = <-c.PacketChan:
						{
							mylog.ZAPS.Debugf("【RECV】 接口%s % X", c.CollInterface.CollInterfaceName, rxBuf)
							rxBufCnt = len(rxBuf)
							if rxBufCnt > 0 {
								rxTotalBufCnt += rxBufCnt
								//追加接收的数据到接收缓冲区
								rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)

							}
							err := node.AnalysisRx(node.Addr, node.VariableMap, rxTotalBuf, rxTotalBufCnt, txBuf)
							{
								if err != nil {
									rxData.Err = err
									return
								}
								state = WaitSuccess
								rxData.RxBuf = rxTotalBuf

								c.CommunicationManageMessageAdd("receive", rxTotalBuf)

								//设备从离线变成上线
								if node.CommStatus == OFFLINE {
									if len(c.CollInterface.OnlineReportChan) == 100 {
										<-c.CollInterface.OnlineReportChan
									}
									c.CollInterface.OnlineReportChan <- node.Name
									c.CollInterface.DeviceNodeOnlineCnt++
									node.CommStatus = ONLINE
								}

								//防止Chan阻塞
								if len(c.CollInterface.PropertyReportChan) >= 100 {
									<-c.CollInterface.PropertyReportChan
								}
								c.CollInterface.PropertyReportChan <- node.Addr
								//log.Printf("reportChan %v\n", len(c.CollInterface.PropertyReportChan))

								node.CommSuccessCnt++
								node.CurCommFailCnt = 0
								node.LastCommRTC = time.Now().Format("2006-01-02 15:04:05")

								rxTotalBufCnt = 0
								rxTotalBuf = rxTotalBuf[0:0]
								goto OUT
							}
						}

					}
				}

			}
		case WaitSuccess:
			//通信帧延时
			interval, _ := strconv.Atoi(c.CollInterface.CommInterface.GetInterval())
			time.Sleep(time.Duration(interval) * time.Millisecond)
			if hasFrame {
				log.Println(color.CyanString("采集接口【%s】还有后续帧,等待%d毫秒", c.CollInterface.CollInterfaceName, interval))
				state = Start
			} else {
				state = Stop
			}
		case Stop:
			{
				cost := time.Since(startT)
				log.Println(color.CyanString("接口【%s】 设备【%s】【第%d帧】 cost %v ", c.CollInterface.CollInterfaceName, node.Name, step, cost))
				c.CollInterface.DeviceNodeOnlineCnt = 0
				for _, v := range c.CollInterface.DeviceNodes {
					if v.CommStatus == ONLINE {
						c.CollInterface.DeviceNodeOnlineCnt++
					}
				}
				return
			}
		}
	}

}

func (c *CommunicationManageTemplate) CommunicationManageDel() {

	for {
		select {
		case cmd := <-c.EmergencyRequestChan:
			{
				mylog.Logger.Infof("emergency chan collName %v nodeName %v funName %v", c.CollInterface.CollInterfaceName, cmd.DeviceName, cmd.FunName)
				rxData := c.CommunicationStateMachine(cmd)

				GetDeviceOnline()
				GetDevicePacketLoss()
				c.EmergencyAckChan <- rxData
			}

		case cmd := <-c.CommonRequestChan:
			rxData := c.CommunicationStateMachine(cmd)
			GetDeviceOnline()
			GetDevicePacketLoss()
			if err := rxData.Err; err != nil {
				mylog.Logger.Debugf("get data from common request chan  error:%v", err)
			}
		case <-c.Signal:
			log.Println(color.RedString("停止接口【%s】采集协程", c.CollInterface.CollInterfaceName))
			return
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
		system.SystemState.DeviceOnline = "0"
	} else {
		system.SystemState.DeviceOnline = fmt.Sprintf("%2.1f", float32(deviceOnlineCnt*100.0/deviceTotalCnt))
	}
}

func GetDevicePacketLoss() {

	//更新设备丢包率
	deviceCommTotalCnt := 0
	deviceCommLossCnt := 0
	for _, v := range CollectInterfaceMap {
		for _, v := range v.DeviceNodes {
			deviceCommTotalCnt += v.CommTotalCnt
			deviceCommLossCnt += v.CommTotalCnt - v.CommSuccessCnt
		}
	}
	if deviceCommLossCnt == 0 {
		system.SystemState.DevicePacketLoss = "0"
	} else {
		system.SystemState.DevicePacketLoss = fmt.Sprintf("%2.1f", float32(deviceCommLossCnt*100.0/deviceCommTotalCnt))
	}
}

func (c *CommunicationManageTemplate) CommunicationManagePoll(polling int) {

	var first = make(chan struct{}, 1)
	for {
		select {
		case <-c.Signal:
			close(first)

			return
		case first <- struct{}{}:
			c.sendCmd()
		case <-time.After(time.Second * time.Duration(polling)):
			c.sendCmd()
		}
	}
}

func (c *CommunicationManageTemplate) sendCmd() {
	//对采集接口下设备进行遍历进行发送
	for index, node := range c.CollInterface.DeviceNodes {
		cmd := CommunicationCmdTemplate{
			CollInterfaceName: c.CollInterface.CollInterfaceName,
			DeviceName:        node.Name,
			FunName:           GETREAL,
			DeviceIndex:       index,
		}
		c.CommonRequestChan <- cmd
	}
}
