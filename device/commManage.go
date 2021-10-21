package device

import (
	"fmt"
	"goAdapter/config"
	"goAdapter/httpServer/model"
	"goAdapter/pkg/mylog"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/fatih/color"
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
	Done              chan struct{}
}

type CommunicationRxTemplate struct {
	Err   error
	RxBuf []byte
}

//透传model
type CommunicationDirectDataReqTemplate struct {
	CollInterfaceName string //采集接口名称
	Data              []byte
}
type CommunicationDirectDataAckTemplate struct {
	CollInterfaceName string //采集接口名称
	Error             error
	Data              []byte
}
type CommunicationManageTemplate struct {
	EmergencyRequestChan  chan CommunicationCmdTemplate
	CommonRequestChan     chan CommunicationCmdTemplate
	EmergencyAckChan      chan CommunicationRxTemplate
	DirectDataRequestChan chan CommunicationDirectDataReqTemplate
	DirectDataAckChan     chan CommunicationDirectDataAckTemplate
	CollInterface         *CollectInterfaceTemplate
	PacketChan            chan []byte
	Signal                chan struct{}
	Ready                 chan struct{}
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

func NewCommunicationManageTemplate(coll *CollectInterfaceTemplate) *CommunicationManageTemplate {

	template := &CommunicationManageTemplate{
		EmergencyRequestChan:  make(chan CommunicationCmdTemplate, 1),
		CommonRequestChan:     make(chan CommunicationCmdTemplate, 100),
		EmergencyAckChan:      make(chan CommunicationRxTemplate, 1),
		DirectDataRequestChan: make(chan CommunicationDirectDataReqTemplate, 10),
		DirectDataAckChan:     make(chan CommunicationDirectDataAckTemplate, 10),
		PacketChan:            make(chan []byte, 100), //最多连续接收100帧数据
		CollInterface:         coll,
		Ready:                 make(chan struct{}, 1),
		Signal:                make(chan struct{}, 1),
	}
	coll.CommunicationManager = template
	if coll.CommInterface.Error() == nil {
		//启动接收协程
		mylog.ZAPS.Infof("采集接口【%s】打开成功，启动接收协程！", coll.CollInterfaceName)
		go template.ReadRx()
	}
	return template
}

func addHandler(collect *CollectInterfaceTemplate, commChaned bool) {
	comm := collect.CommInterface
	if comm == nil {
		mylog.ZAPS.Errorf("通讯口【%s】未绑定到接口【%s】",
			collect.CommInterfaceName, collect.CollInterfaceName)
		return
	}
	if commChaned {
		if err := comm.Open(); err != nil {
			mylog.ZAPS.Errorf("通讯口【%s】打开【%s】错误", comm.GetName(), comm.Unique())
			return
		}
	}

	manager := NewCommunicationManageTemplate(collect)

	go manager.CommunicationManagePoll()
	go manager.CommunicationManageDel()
	mylog.ZAPS.Infof("添加采集【%s】到定时任务,Addr:%p", collect.CollInterfaceName, manager)
}

func delHandler(collect *CollectInterfaceTemplate, commChaned bool) {
	managerRemove := collect.CommunicationManager
	if managerRemove != nil {
		//广播到所有监听管理模板的信号channel的goroutine
		close(managerRemove.Signal)
		if commChaned {
			managerRemove.CollInterface.CommInterface.Close()
		}

		mylog.ZAPS.Infof("取消采集【%s】定时任务,Addr:%p", collect.CollInterfaceName, managerRemove)
	} else {
		mylog.ZAP.Error("采集接口未创建管理字段", zap.String("采集接口", collect.CollInterfaceName))
	}
}

func SubScribeCollect(topics string, quitChan chan struct{}) {
	mylog.ZAP.Debug("开始订阅采集器主题", zap.String("topics", topics))
	sub := CollectInterfaceMap.publisher.Subscribe(10, topics)
	go func() {
		for {
			select {
			case msg := <-sub.Receiver:
				collectName, ok := msg.Fields["Collect"].(*CollectInterfaceTemplate)
				if !ok {
					mylog.ZAP.Sugar().Errorf("this msg field Collect type error:%t", msg.Fields["Collect"])
					continue
				}
				commChaned, ok := msg.Fields["CommChange"].(bool)
				if !ok {
					mylog.ZAP.Sugar().Errorf("this msg field commchange error:%v", msg.Fields["CommChange"])
					continue
				}

				switch msg.Name {
				case CollectAdd:
					mylog.ZAP.Debug("添加采集接口", zap.String("name", collectName.CollInterfaceName), zap.String("comm", collectName.CommInterfaceName))
					addHandler(collectName, true)
				case CollectUpdate:
					mylog.ZAP.Debug("更新采集接口", zap.String("name", collectName.CollInterfaceName), zap.String("comm", collectName.CommInterfaceName))
					delHandler(collectName, commChaned)
					addHandler(collectName, commChaned)
				case CollectQuery:
					mylog.ZAP.Debug("查询采集接口", zap.String("name", collectName.CollInterfaceName), zap.String("comm", collectName.CommInterfaceName))
				case CollectDelete:
					mylog.ZAP.Debug("删除采集接口", zap.String("name", collectName.CollInterfaceName), zap.String("comm", collectName.CommInterfaceName))
					delHandler(collectName, commChaned)
				}
			case <-quitChan:
				CollectInterfaceMap.Close()
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func SubScribeComunication(topics string, quitChan chan struct{}) {
	mylog.ZAP.Debug("开始订阅通信接口主题", zap.String("topics", topics))
	sub := CommunicationInterfaceMap.publisher.Subscribe(10, topics)
	go func() {
		for {
			select {
			case msg := <-sub.Receiver:
				commName := msg.Fields["Name"].(string)
				// collectName, ok := msg.Fields["Collect"].(*CollectInterfaceTemplate)
				// if !ok {
				// 	mylog.ZAP.Sugar().Errorf("this msg field Collect type error:%t", msg.Fields["Collect"])
				// 	continue
				// }
				// commChaned, ok := msg.Fields["CommChange"].(bool)
				// if !ok {
				// 	mylog.ZAP.Sugar().Errorf("this msg field commchange error:%v", msg.Fields["CommChange"])
				// 	continue
				// }

				switch msg.Name {
				case CommAdd:
					mylog.ZAP.Debug("添加了通信接口", zap.String("name", commName))
				case CommUpdate:
					mylog.ZAP.Debug("更新了通信接口")
					oldCommField, ok := msg.Fields["Old"]
					if !ok {
						mylog.ZAPS.Errorf("msg have no field named Old")
						continue
					}
					oldComm, ok := oldCommField.(CommunicationInterface)
					if !ok {
						mylog.ZAPS.Errorf("msg field Old is not communicationinterface:%t", oldCommField)
						continue
					}
					newCommField, ok := msg.Fields["New"]
					if !ok {
						mylog.ZAPS.Errorf("msg have no field named New")
						continue
					}
					newComm, ok := newCommField.(CommunicationInterface)
					if !ok {
						mylog.ZAPS.Errorf("msg field New is not communicationinterface:%t", newCommField)
						continue
					}
					for _, collectName := range oldComm.BindNames() {
						coll := CollectInterfaceMap.Get(collectName)
						if coll != nil {
							//和这个通信口有关的采集都要停掉
							//delHandler(coll, true)
							oldComm.UnBind(collectName)
							//替换新的comm口
							if err := newComm.Open(); err != nil {
								mylog.ZAP.Error("新通讯口打开失败", zap.Error(err))
							}
							coll.CommInterface = newComm
							coll.CommInterfaceName = newComm.GetName()
							//新comm绑定这个采集
							newComm.Bind(collectName)

							//如果之前打开就是错误的那么要重新add
							if oldComm.Error() != nil {
								addHandler(coll, false)
							}

						}

					}
				case CommQuery:
				case CommDelete:
					mylog.ZAP.Debug("删除通信接口", zap.String("name", commName))
				}

			case <-quitChan:
				CollectInterfaceMap.Close()
				return
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}
func SubScribeTSL(topics string, quitChan chan struct{}) {
	mylog.ZAP.Debug("开始订阅设备模型主题", zap.String("topics", topics))
	sub := DeviceTSLMap.publisher.Subscribe(20, topics)
	go func() {
		for {
			select {
			case msg := <-sub.Receiver:

				switch msg.Name {
				case PropertyAdd:
					mylog.ZAPS.Debugf("添加单个属性%v", msg.Fields["name"])
					propertyValue, ok := msg.Fields["property"]
					if !ok {
						mylog.ZAP.Error("recv property msg not contain proerty field")
						continue
					}
					property, ok := propertyValue.(model.DeviceTSLPropertyTemplate)
					if !ok {
						mylog.ZAP.Error("recv property msg is not model.DeviceTSLPropertyTemplate")
						continue
					}
					colls := CollectInterfaceMap.GetAll()
					for _, v := range colls {
						for _, n := range v.DeviceNodes {
							if n.Type == msg.Fields["plugin"].(string) {
								n.Properties = append(n.Properties, property)
							}
						}
					}
				case PropertySync:
					mylog.ZAPS.Debugf("同步所有导入属性")
					propertiesValue, ok := msg.Fields["properties"]
					if !ok {
						mylog.ZAP.Error("recv property msg not contain proerty field")
						continue
					}
					properties, ok := propertiesValue.([]model.DeviceTSLPropertyTemplate)
					if !ok {
						mylog.ZAP.Error("recv property msg is not []*model.DeviceTSLPropertyTemplate")
						continue
					}
					colls := CollectInterfaceMap.GetAll()
					for _, v := range colls {
						for _, n := range v.DeviceNodes {
							if n.Type == msg.Fields["plugin"].(string) {
								n.Properties = make([]model.DeviceTSLPropertyTemplate, len(properties))
								copy(n.Properties, properties)
								ClearPropertyValue(n.Properties)
							}
						}
					}
				case PropertyUpdate:
					mylog.ZAPS.Debugf("更新属性%v", msg.Fields["name"])
					propertyValue, ok := msg.Fields["property"]
					if !ok {
						mylog.ZAP.Error("recv property msg not contain proerty field")
						continue
					}
					property, ok := propertyValue.(model.DeviceTSLPropertyTemplate)

					if !ok {
						mylog.ZAP.Error("recv property msg is not model.DeviceTSLPropertyTemplate")
						continue
					}
					colls := CollectInterfaceMap.GetAll()
					for _, v := range colls {
						for _, n := range v.DeviceNodes {
							if n.Type == msg.Fields["plugin"].(string) {
								for i := 0; i < len(n.Properties); i++ {
									if n.Properties[i].Name == property.Name {
										n.Properties[i].AccessMode = property.AccessMode
										n.Properties[i].Explain = property.Explain
										n.Properties[i].Params = property.Params
										n.Properties[i].Type = property.Type
									}
								}
							}
						}
					}

				case PropertyQuery:
				case PropertyDelete:
					mylog.ZAPS.Debugf("删除了属性%v", msg.Fields["name"])
					propertyValue, ok := msg.Fields["property"]
					if !ok {
						mylog.ZAP.Error("recv property msg not contain proerty field")
						continue
					}
					property, ok := propertyValue.(model.DeviceTSLPropertyTemplate)
					if !ok {
						mylog.ZAP.Error("recv property index  is not int")
						continue
					}
					colls := CollectInterfaceMap.GetAll()
					for _, v := range colls {
						for _, n := range v.DeviceNodes {
							if n.Type == msg.Fields["plugin"].(string) {
								var i = -1
								for index, v := range n.Properties {
									if v.Name == property.Name {
										i = index
										break
									}
								}
								if i != -1 {
									n.Properties = append(n.Properties[:i], n.Properties[i+1:]...)
								}
							}
						}
					}

				}

			case <-quitChan:
				CollectInterfaceMap.Close()
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
		TimeStamp: time.Now().Format("2006-01-02 15:04:05.000"),
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
			mylog.ZAPS.Debugf("%s关闭ReadRX goroutine", c.CollInterface.CollInterfaceName)
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

//设备采集状态机
func (c *CommunicationManageTemplate) CommunicationStateMachine(cmd CommunicationCmdTemplate) (rxData CommunicationRxTemplate) {

	//startT := time.Now() //计算当前时间
	if len(c.CollInterface.DeviceNodes) < cmd.DeviceIndex+1 {
		return
	}
	node := c.CollInterface.DeviceNodes[cmd.DeviceIndex]
	step := 0           //用来标识后续帧
	var txBuf []byte    //发送字节
	var hasFrame = true //是否有后续帧
	var err error
	var state = Start //初始状态
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
				if nodeIndex == -1 {
					return
				}
				state = Generate
			}
		case Generate:
			{
				if cmd.FunName == GETREAL {
					txBuf, hasFrame, err = node.GenerateGetRealVariables(node.Addr, step)
					if err != nil {
						rxData.Err = fmt.Errorf("device %s generate error:%v", node.Name, err)
						state = Stop
						goto OUT
					}
					if len(txBuf) <= 0 {
						rxData.Err = fmt.Errorf("device %s txbuf length <=0 ", node.Name)
						state = Stop
						goto OUT
					}
					state = Send
					step++

				} else {
					txBuf, hasFrame, err = node.DeviceCustomCmd(node.Addr, cmd.FunName, cmd.FunPara, step)

					if err != nil {
						rxData.Err = fmt.Errorf("device custom  cmd error:%v", err)
						state = Stop
						goto OUT

					}
					state = Send
					step++
				}
			}
		case Send:
			{ //---------------发送-------------------------
				//判断是否是串口采集
				_, err := c.CollInterface.CommInterface.Write(txBuf)
				//如果写入错误很有可能是串口关闭了
				if err != nil {
					rxData.Err = fmt.Errorf("device %s write data to comm %v error:%v", node.Name, c.CollInterface.CommInterfaceName, err)
					state = Stop
					goto OUT
				}
				c.Ready <- struct{}{}
				c.CommunicationManageMessageAdd("send", txBuf)
				node.CommTotalCnt++
				mylog.ZAPS.Debugf("【S-%s】% X", cmd.CollInterfaceName, txBuf)
				state = Wait

			}

		case Wait:
			{
				var (
					rxBuf         []byte
					rxTotalBuf    []byte
					rxBufCnt      int
					rxTotalBufCnt int
					//串口读超时
					commTimeout, _ = strconv.Atoi(c.CollInterface.CommInterface.GetTimeOut())
					timeout        = time.Duration(commTimeout) * time.Millisecond
					timer          = time.NewTimer(timeout)
				)

				select {
				//是否接收超时
				case <-timer.C:
					{
						CommunicationMessage := CommunicationMessageTemplate{
							CollName:  c.CollInterface.CollInterfaceName,
							TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
							Direction: "receive",
							Content:   fmt.Sprintf("timeout:%d ms", commTimeout),
						}
						if len(c.CollInterface.CommMessage) < 1024 {
							c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, &CommunicationMessage)
						} else {
							c.CollInterface.CommMessage = c.CollInterface.CommMessage[1:]
							c.CollInterface.CommMessage = append(c.CollInterface.CommMessage, &CommunicationMessage)
						}

						node.CurCommFailCnt++
						rxTotalBufCnt = 0
						// if count >= offlineperion,judge this node offline
						if node.CurCommFailCnt >= c.CollInterface.OfflinePeriod {
							log.Println(color.MagentaString("采集器【%s】-> 设备【%s】接收数据超时，失败次数已达到最大尝试次数:%d set offline", c.CollInterface.CollInterfaceName, node.Name, c.CollInterface.OfflinePeriod))
							node.CurCommFailCnt = 0
							//设备从上线变成离线
							if node.CommStatus == ONLINE {
								if len(c.CollInterface.OfflineReportChan) == 100 {
									<-c.CollInterface.OfflineReportChan
								}
								c.CollInterface.OfflineReportChan <- node.Name
								node.CommStatus = OFFLINE
								state = Stop
								goto OUT
							}
						}
						//没有达到最大失败次数如果有后续帧可以再尝试
						log.Println(color.MagentaString("采集器【%s】-> 设备【%s】接收数据超时，失败次数:%d 可尝试次数:%d", c.CollInterface.CollInterfaceName, node.Name, node.CurCommFailCnt, c.CollInterface.OfflinePeriod))
						state = WaitSuccess
					}

				//继续接收数据
				case rxBuf = <-c.PacketChan:
					{
						rxBufCnt = len(rxBuf)
						if rxBufCnt <= 0 {
							mylog.ZAPS.Errorf("采集器【%s】-> 设备【%s】接收数据rxBuf 长度小于等于0!", c.CollInterface.CollInterfaceName, node.Name)
							state = WaitSuccess
							goto OUT
						}

						mylog.ZAPS.Debugf("【R-%s】% X", c.CollInterface.CollInterfaceName, rxBuf)
						rxTotalBufCnt += rxBufCnt
						//追加接收的数据到接收缓冲区
						rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
						c.CommunicationManageMessageAdd("receive", rxTotalBuf)
						node.CommSuccessCnt++
						node.CurCommFailCnt = 0
						node.LastCommRTC = time.Now().Format("2006-01-02 15:04:05")
						err := node.AnalysisRx(node.Addr, node.Properties, rxTotalBuf, rxTotalBufCnt, txBuf)
						{
							if err != nil {
								mylog.ZAP.Error("解析帧数据错误", zap.String("device", node.Name), zap.Error(err), zap.Int("step", step))
								state = WaitSuccess
								goto OUT
							}
							state = WaitSuccess
							rxData.RxBuf = rxTotalBuf
							//设备从离线变成上线
							if node.CommStatus == OFFLINE {
								if len(c.CollInterface.OnlineReportChan) == 100 {
									<-c.CollInterface.OnlineReportChan
								}
								c.CollInterface.OnlineReportChan <- node.Name
								node.CommStatus = ONLINE
							}

							//防止Chan阻塞
							if len(c.CollInterface.PropertyReportChan) >= 100 {
								<-c.CollInterface.PropertyReportChan
							}
							c.CollInterface.PropertyReportChan <- node.Addr
							//log.Printf("reportChan %v\n", len(c.CollInterface.PropertyReportChan))
							rxTotalBufCnt = 0
						}
					}

				}

			}
		case WaitSuccess:
			//处理多帧的延时间隔
			interval, _ := strconv.Atoi(c.CollInterface.CommInterface.GetInterval())
			time.Sleep(time.Duration(interval) * time.Millisecond)
			if hasFrame {
				log.Println(color.CyanString("采集接口【%s】还有后续帧,等待%d毫秒", c.CollInterface.CollInterfaceName, interval))
				state = Generate
			} else {
				state = Stop
			}
		case Stop:
			//终止这个设备的采集,并计算设备在线数量,其他关于设备终止采集的代码也可以放这里
			{
				//cost := time.Since(startT)
				//log.Println(color.CyanString("接口【%s】 设备【%s】【第%d帧】 cost %v ", c.CollInterface.CollInterfaceName, node.Name, step, cost))
				c.CollInterface.DeviceNodeOnlineCnt = 0
				for _, v := range c.CollInterface.DeviceNodes {
					if v.CommStatus == ONLINE {
						c.CollInterface.DeviceNodeOnlineCnt++
					}
				}
				//完成标志
				cmd.Done <- struct{}{}
				return
			}
		}
	}

}

//设备透传报文
func (c *CommunicationManageTemplate) SendDirectData(req CommunicationDirectDataReqTemplate, commInterface CommunicationInterface) (ack CommunicationDirectDataAckTemplate) {

	if commInterface == nil {
		ack.Error = fmt.Errorf("comminterface is nil")
		return ack
	}
	var state = Send //初始状态

	for {
		switch state {
		case Send:
			{ //---------------发送-------------------------
				//判断是否是串口采集
				_, err := c.CollInterface.CommInterface.Write(req.Data)
				//如果写入错误很有可能是串口关闭了
				if err != nil {
					ack.Error = fmt.Errorf("write data to comm %v error:%v", c.CollInterface.CommInterfaceName, err)
					return ack
				}
				c.CommunicationManageMessageAdd("send", req.Data)
				c.Ready <- struct{}{}
				state = Wait

			}

		case Wait:
			{
				var (
					rxBuf         []byte
					rxTotalBuf    []byte
					rxBufCnt      int
					rxTotalBufCnt int
				)

				select {
				//是否接收超时
				case <-time.After(time.Second * 3):
					{
						ack.Error = fmt.Errorf("read from  comm %s timeout", commInterface.GetName())
						return

					}

				//继续接收数据
				case rxBuf = <-c.PacketChan:
					{
						rxBufCnt = len(rxBuf)
						rxTotalBufCnt += rxBufCnt
						//追加接收的数据到接收缓冲区
						rxTotalBuf = append(rxTotalBuf, rxBuf[:rxBufCnt]...)
						c.CommunicationManageMessageAdd("receive", rxTotalBuf)
						ack.Data = rxTotalBuf
						return
					}

				}

			}

		}
	}

}

func (c *CommunicationManageTemplate) CommunicationManageDel() {

	for {
		select {
		case <-c.Signal:
			log.Println(color.RedString("停止接口【%s】采集协程", c.CollInterface.CollInterfaceName))
			return
		case cmd := <-c.EmergencyRequestChan:
			{
				mylog.ZAPS.Infof("emergency chan collName %v nodeName %v funName %v", c.CollInterface.CollInterfaceName, cmd.DeviceName, cmd.FunName)
				rxData := c.CommunicationStateMachine(cmd)
				CollectInterfaceMap.Statics()
				c.EmergencyAckChan <- rxData
			}

		case cmd := <-c.CommonRequestChan:
			rxData := c.CommunicationStateMachine(cmd)
			CollectInterfaceMap.Statics()
			if err := rxData.Err; err != nil {
				mylog.ZAPS.Debugf("get data from common request chan  error:%v", err)
			}
		case req := <-c.DirectDataRequestChan:
			c.DirectDataAckChan <- c.SendDirectData(req, c.CollInterface.CommInterface)
		}
	}
}

func (c *CommunicationManageTemplate) CommunicationManagePoll() {

	//var first = make(chan struct{}, 1)
	//pollingTime:=time.Second*time.Duration(polling)
	//timer:=time.NewTimer(pollingTime)
	doneChan := make(chan struct{}, 1)
	for {
		select {
		case <-c.Signal:
			mylog.ZAPS.Debugf("采集接口【%s】 停止CommunicationManagePoll", c.CollInterface.CollInterfaceName)
			return
		default:
			c.sendCmd(doneChan)
		}
	}
}

func (c *CommunicationManageTemplate) sendCmd(doneChan chan struct{}) {
	if c.CollInterface == nil || c.CollInterface.DeviceNodes == nil {
		mylog.ZAP.Error("【sendcmd】采集接口或者设备节点未初始化")
		return
	}
	if c.CollInterface.CommInterface == nil {
		mylog.ZAP.Error("【sendcmd】通讯接口未初始化", zap.String("采集接口", c.CollInterface.CollInterfaceName))
		return
	}
	if err := c.CollInterface.CommInterface.Error(); err != nil {
		mylog.ZAP.Error("【sendcmd】通讯口打开错误", zap.Error(err))
		return
	}
	//对采集接口下设备进行遍历进行发送
	for index, node := range c.CollInterface.DeviceNodes {

		cmd := CommunicationCmdTemplate{
			CollInterfaceName: c.CollInterface.CollInterfaceName,
			DeviceName:        node.Name,
			FunName:           GETREAL,
			DeviceIndex:       index,
			Done:              doneChan,
		}
		c.CommonRequestChan <- cmd
		//有序完成
		<-cmd.Done
	}
}
