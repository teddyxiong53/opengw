package device

import (
	"context"
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
	Context              context.Context
	Cancel               context.CancelFunc
	Delay                time.Duration
	Waiting              bool
	Ready                chan struct{}
}

var CommunicationManage = CommManger{
	ManagerTemp: make(map[string]*CommunicationManageTemplate),
	Collectors:  make(chan *CollectInterfaceStatus, 20),
}

type CommManger struct {
	ManagerTemp map[string]*CommunicationManageTemplate
	Collectors  chan *CollectInterfaceStatus
}

func NewCommunicationManageTemplate(coll *CollectInterfaceTemplate) *CommunicationManageTemplate {

	ctx, cancel := context.WithCancel(context.Background())

	template := &CommunicationManageTemplate{
		EmergencyRequestChan: make(chan CommunicationCmdTemplate, 1),
		CommonRequestChan:    make(chan CommunicationCmdTemplate, 100),
		EmergencyAckChan:     make(chan CommunicationRxTemplate, 1),
		PacketChan:           make(chan []byte, 100), //最多连续接收100帧数据
		CollInterface:        coll,
		Context:              ctx,
		Cancel:               cancel,
		Delay:                time.Millisecond * 100,
		Ready:                make(chan struct{}, 1),
	}
	if coll.CommInterface.Error() == nil {
		//启动接收协程
		log.Println(color.CyanString("采集接口【%s】打开成功，启动接收协程！", coll.CollInterfaceName))
		go template.ReadRx(ctx)
	}
	return template
}
func addHandler(scheduler *gocron.Scheduler, collect *CollectInterfaceStatus) {
	comm := collect.Tmp.CommInterface
	if comm == nil {
		log.Println(color.YellowString("通讯口【%s】未绑定到接口【%s】",
			collect.Tmp.CommInterfaceName, collect.Tmp.CollInterfaceName))
		return
	}
	if err := comm.Open(); err != nil {
		log.Println(color.YellowString("通讯口【%s】打开错误", comm.GetName()))
		return
	}
	manager := NewCommunicationManageTemplate(collect.Tmp)
	CommunicationManage.ManagerTemp[collect.Tmp.CollInterfaceName] = manager

	scheduler.Every(uint64(collect.Tmp.PollPeriod)).Seconds().Do(manager.CommunicationManagePoll)
	go manager.CommunicationManageDel()
	log.Println(color.CyanString("添加采集【%s】到定时任务,Addr:%p", collect.Tmp.CollInterfaceName, manager))
}

func delHandler(scheduler *gocron.Scheduler, collect *CollectInterfaceStatus) {
	managerRemove := CommunicationManage.ManagerTemp[collect.Tmp.CollInterfaceName]
	if managerRemove != nil {
		scheduler.Remove(managerRemove.CommunicationManagePoll)
		managerRemove.Cancel()
		managerRemove.CollInterface.CommInterface.Close()
		log.Println(color.CyanString("取消采集【%s】定时任务,Addr:%p", collect.Tmp.CollInterfaceName, managerRemove))
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

func (c *CommunicationManageTemplate) ReadRx(ctx context.Context) {

	//阻塞读
	rxBuf := make([]byte, 512)
	var rxBufCnt int
	var err error

	for range c.Ready {
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
			rxBufCnt = 0
		}

	}
}

func (c *CommunicationManageTemplate) CommunicationStateMachine(cmd CommunicationCmdTemplate) (rxData CommunicationRxTemplate) {
	var timeout int
	var interval int
	timeout, _ = strconv.Atoi(c.CollInterface.CommInterface.GetTimeOut())
	interval, _ = strconv.Atoi(c.CollInterface.CommInterface.GetInterval())
	//通讯接口的超时
	ticker := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	defer ticker.Stop()
	defer time.Sleep(time.Duration(interval) * time.Millisecond)

	startT := time.Now() //计算当前时间
	if c.CollInterface != nil && c.CollInterface.DeviceNodes != nil {
		node := c.CollInterface.DeviceNodes[cmd.DeviceIndex]
		step := 0
		var txBuf []byte
		var hasFrame = true
		var err error
	OUT:
		//是否有后续帧
		for hasFrame {
			if cmd.FunName == GETREAL {
				txBuf, hasFrame, err = node.GenerateGetRealVariables(node.Addr, step)
				if err != nil {
					log.Printf("genterate get real error:%v\n", err)
					rxData.Err = err
					return
				}

			} else {
				txBuf, hasFrame, err = node.DeviceCustomCmd(node.Addr, cmd.FunName, cmd.FunPara, step)
				if err != nil {
					log.Printf("device custom  cmd error:%v\n", err)
					rxData.Err = err
					return
				}

			}
			step++
			mylog.ZAPS.Debugf("【SEND】接口【%s】% X", cmd.CollInterfaceName, txBuf)
			//setting.ZAPS.Debugf("interfaceName %s:txbuf %X", c.CollInterface.CollInterfaceName, txBuf)

			CommunicationMessage := CommunicationMessageTemplate{
				CollName:  c.CollInterface.CollInterfaceName,
				TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
				Direction: "send",
				Content:   fmt.Sprintf("%X", txBuf),
			}
			if len(c.CollInterface.CommMessage) < 1024 {
				c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, &CommunicationMessage)
			} else {
				c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
				c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, &CommunicationMessage)
			}

			//---------------发送-------------------------
			//判断是否是串口采集
			c.CollInterface.CommInterface.Write(txBuf)
			c.Ready <- struct{}{}
			node.CommTotalCnt++
			//---------------等待接收----------------------
			//阻塞读
			var (
				rxBuf         []byte
				rxTotalBuf    []byte
				rxBufCnt      int
				rxTotalBufCnt int
			)

			for {
				select {
				//是否接收超时
				case <-ticker.C:
					{
						CommunicationMessage := CommunicationMessageTemplate{
							CollName:  c.CollInterface.CollInterfaceName,
							TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
							Direction: "receive",
							Content:   fmt.Sprintf("接收数据超时了,超时阈值:%s ms", c.CollInterface.CommInterface.GetTimeOut()),
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
						log.Println(color.MagentaString("采集器【%v】-设备【%s】接收数据超时，失败次数:%d 总次数:%d", c.CollInterface.CollInterfaceName, node.Name, node.CurCommFailCnt, c.CollInterface.OfflinePeriod))
					}

				//继续接收数据
				case rxBuf = <-c.PacketChan:
					{
						rxBufCnt = len(rxBuf)
						if rxBufCnt > 0 {
							mylog.ZAPS.Debugf("【RECV】 接口%s % X", c.CollInterface.CollInterfaceName, rxBuf)
							rxTotalBufCnt += rxBufCnt
							//追加接收的数据到接收缓冲区
							rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
							err := node.AnalysisRx(node.Addr, node.VariableMap, rxTotalBuf, rxTotalBufCnt, txBuf)
							{
								ticker.Stop()
								if err != nil {
									log.Printf("analysisrx error:%v\n", err)
									rxData.Err = err
									return
								}

								rxData.RxBuf = rxTotalBuf

								CommunicationMessage := CommunicationMessageTemplate{
									CollName:  c.CollInterface.CollInterfaceName,
									TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
									Direction: "receive",
									Content:   fmt.Sprintf("%X", rxTotalBuf),
								}
								if len(c.CollInterface.CommMessage) < 1024 {
									c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, &CommunicationMessage)
								} else {
									c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
									c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, &CommunicationMessage)
								}

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
								//通信帧延时
								if hasFrame {
									log.Println(color.CyanString("采集接口【%s】还有后续帧,等待%d毫秒", c.CollInterface.CollInterfaceName, interval))
									time.Sleep(time.Duration(interval) * time.Millisecond)
									goto OUT
								} else {
									//清除本次接收数据
									rxBufCnt = 0
									rxBuf = rxBuf[0:0]
									return
								}

							}

							//log.Printf("rxTotalBuf %X\n", rxTotalBuf)
						}
					}
				}
			}

		}
	} else {
		rxData.Err = errors.New("interface or device nodes is nil")
	}
	tc := time.Since(startT) //计算耗时
	mylog.Logger.Debugf("%v time cost = %v", c.CollInterface.CollInterfaceName, tc)
	//更新设备在线数量
	c.CollInterface.DeviceNodeOnlineCnt = 0
	for _, v := range c.CollInterface.DeviceNodes {
		if v.CommStatus == "onLine" {
			c.CollInterface.DeviceNodeOnlineCnt++
		}
	}

	return
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
		case <-c.Context.Done():
			log.Println(color.RedString("停止接口【%s】采集协程", c.CollInterface.CollInterfaceName))
			return

		case cmd := <-c.CommonRequestChan:
			rxData := c.CommunicationStateMachine(cmd)
			GetDeviceOnline()
			GetDevicePacketLoss()
			if err := rxData.Err; err != nil {
				mylog.Logger.Debugf("get data from common request chan  error:%v", err)
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

func (c *CommunicationManageTemplate) CommunicationManagePoll() {

	if c.Waiting {
		time.Sleep(c.Delay)
	}
	if c.CollInterface.DeviceNodeCnt <= 0 {
		log.Println(color.YellowString("采集接口【%s】目前还未下挂设备", c.CollInterface.CollInterfaceName))
		c.Waiting = true
		c.Delay *= 2
		return
	} else {
		c.Waiting = false
		c.Delay = time.Millisecond * 100
	}

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
