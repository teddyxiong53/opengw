package device

import (
	"encoding/json"
	"fmt"
	"goAdapter/pkg/system"
	"os"
	"sync"

	"github.com/leandro-lugaresi/hub"
)

const (
	MaxCollectInterfaceManage int = 2

	InterFaceID0 int = 0
	InterFaceID1 int = 1
	InterFaceID2 int = 2
	InterFaceID3 int = 3
	InterFaceID4 int = 4
	InterFaceID5 int = 5
	InterFaceID6 int = 6
	InterFaceID7 int = 7
)

//映射采集接口的CURD
type CollectAction uint8

type CommunicationMessageTemplate struct {
	CollName  string `json:"CollInterfaceName"` //接口名称
	TimeStamp string `json:"TimeStamp"`         //时间戳
	Direction string `json:"DataDirection"`     //数据方向
	Content   string `json:"DataContent"`       //数据内容
}

//采集接口配置参数
type CollectInterfaceParamTemplate struct {
	CollInterfaceName string                `json:"CollInterfaceName"`       //采集接口
	CommInterfaceName string                `json:"CommInterfaceName"`       //通信接口
	PollPeriod        int                   `json:"PollPeriod"`              //采集周期
	OfflinePeriod     int                   `json:"OfflinePeriod"`           //离线超时周期
	DeviceNodeCnt     int                   `json:"DeviceNodeCnt,omitempty"` //设备数量
	DeviceNodes       []*DeviceNodeTemplate `json:"DeviceNodeMap,omitempty"` //节点名称
}

//采集接口模板
type CollectInterfaceTemplate struct {
	CollInterfaceName    string                          `json:"CollInterfaceName"` //采集接口
	CommInterfaceName    string                          `json:"CommInterfaceName"` //通信接口
	CommInterface        CommunicationInterface          `json:"-"`
	CommMessage          []*CommunicationMessageTemplate `json:"-"`
	PollPeriod           int                             `json:"PollPeriod"`          //采集周期
	OfflinePeriod        int                             `json:"OfflinePeriod"`       //离线超时周期
	DeviceNodeCnt        int                             `json:"DeviceNodeCnt"`       //设备数量
	DeviceNodeOnlineCnt  int                             `json:"DeviceNodeOnlineCnt"` //设备在线数量
	DeviceNodes          []*DeviceNodeTemplate           `json:"DeviceNodeMap"`
	OnlineReportChan     chan string                     `json:"-"`
	OfflineReportChan    chan string                     `json:"-"`
	PropertyReportChan   chan string                     `json:"-"`
	CommunicationManager *CommunicationManageTemplate    `json:"-"`
}

type CollectInterfaceManager struct {
	m map[string]*CollectInterfaceTemplate
	sync.RWMutex
	publisher *hub.Hub
	Changed   bool
}

func NewCollectInterfaceManager() *CollectInterfaceManager {
	mgr := &CollectInterfaceManager{
		m:         make(map[string]*CollectInterfaceTemplate),
		publisher: hub.New(),
	}
	return mgr
}

var CollectInterfaceMap = NewCollectInterfaceManager()

func (cim *CollectInterfaceManager) Add(collect *CollectInterfaceTemplate) {
	cim.Lock()
	cim.m[collect.CollInterfaceName] = collect
	cim.Changed = true
	cim.Publish(CollectAdd, hub.Fields{"Name": collect.CollInterfaceName, "Collect": collect, "CommChange": true})
	cim.Unlock()
}

func (cim *CollectInterfaceManager) Update(newCollect CollectInterfaceParamTemplate) bool {
	old := cim.Get(newCollect.CollInterfaceName)
	var commChanged bool
	if old == nil {
		return false
	}
	if old.CommInterfaceName != newCollect.CommInterfaceName {
		commChanged = true
		old.CommInterface.UnBind(old.CollInterfaceName)
		old.CommInterfaceName = newCollect.CommInterfaceName
		old.CommInterface = CommunicationInterfaceMap.Get(newCollect.CommInterfaceName)
		old.CommInterface.Bind(old.CollInterfaceName)
	}
	old.PollPeriod = newCollect.PollPeriod
	old.OfflinePeriod = newCollect.OfflinePeriod
	cim.Changed = true
	cim.Publish(CollectUpdate, hub.Fields{"Name": newCollect.CollInterfaceName, "Collect": old, "CommChange": commChanged})
	return true
}

func (cim *CollectInterfaceManager) Delete(name string) bool {
	old := cim.Get(name)
	if old == nil {
		return false
	}
	cim.Lock()
	delete(cim.m, name)
	old.CommInterface.UnBind(old.CollInterfaceName)
	cim.Changed = true
	cim.Publish(CollectDelete, hub.Fields{"Name": old.CollInterfaceName, "Collect": old, "CommChange": true})
	cim.Unlock()
	return true
}

func (cim *CollectInterfaceManager) Get(name string) *CollectInterfaceTemplate {
	cim.RLock()
	if v, ok := cim.m[name]; ok {
		cim.RUnlock()
		return v
	}
	cim.RUnlock()
	return nil
}

func (cim *CollectInterfaceManager) GetAll() []*CollectInterfaceTemplate {
	cim.RLock()
	r := make([]*CollectInterfaceTemplate, 0, len(cim.m))
	for _, v := range cim.m {
		r = append(r, v)
	}
	cim.RUnlock()
	return r

}

func (cim *CollectInterfaceManager) SaveTo(f *os.File) error {
	cim.RLock()
	defer cim.RUnlock()
	if !cim.Changed {
		return nil
	}
	data, err := json.Marshal(cim.m)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func (cim *CollectInterfaceManager) CommCheck(comName string) (used bool, collectName string) {
	cim.RLock()
	for k, v := range cim.m {
		if v.CommInterfaceName == comName {
			cim.RUnlock()
			return true, k
		}
	}
	cim.RUnlock()
	return
}

//在线和丢包率统计
func (cim *CollectInterfaceManager) Statics() {
	var online, total, commLoss, commTotal int
	cim.RLock()
	for _, v := range CollectInterfaceMap.m {
		total += v.DeviceNodeCnt
		online += v.DeviceNodeOnlineCnt
		for _, node := range v.DeviceNodes {
			commTotal += node.CommTotalCnt
			commLoss += node.CommTotalCnt - node.CommSuccessCnt
		}

	}
	cim.RUnlock()
	if online == 0 {
		system.SystemState.DeviceOnline = "0"
	} else {
		system.SystemState.DeviceOnline = fmt.Sprintf("%2.1f", float32(online*100.0/total))
	}

	if commLoss == 0 {
		system.SystemState.DevicePacketLoss = "0"
	} else {
		system.SystemState.DevicePacketLoss = fmt.Sprintf("%2.1f", float32(commLoss*100.0/commTotal))
	}
}

func (cim *CollectInterfaceManager) Publish(name string, fields hub.Fields) {
	if cim.publisher == nil {
		cim.publisher = hub.New()
	}
	cim.publisher.Publish(hub.Message{
		Name:   name,
		Fields: fields,
	})
}

func (cim *CollectInterfaceManager) Close() {
	cim.publisher.Close()
	//TODO 其他资源的回收？
}

func (cim *CollectInterfaceManager) Init() error {
	cim.Lock()
	defer cim.Unlock()
	for k, v := range CollectInterfaceMap.m {
		if v.DeviceNodes == nil {
			v.DeviceNodes = make([]*DeviceNodeTemplate, 0, 10)
		}

		for _, node := range v.DeviceNodes {
			v.InitDeviceNode(node)
		}

		var param = &CollectInterfaceParamTemplate{
			CollInterfaceName: v.CollInterfaceName,
			CommInterfaceName: v.CommInterfaceName,
			PollPeriod:        v.PollPeriod,
			OfflinePeriod:     v.OfflinePeriod,
			DeviceNodeCnt:     v.DeviceNodeCnt,
			DeviceNodes:       v.DeviceNodes,
		}

		collect, err := NewCollectInterface(param)
		if err != nil {
			return err
		}
		CollectInterfaceMap.m[k] = collect
		cim.Publish(CollectAdd, hub.Fields{
			"Name":       k,
			"Collect":    collect,
			"CommChange": true,
		})
	}
	return nil
}

func NodeManageInit() (err error) {

	//设备模版
	if err = ReadPlugins(PLUGINPATH); err != nil {
		return err
	}

	//通信接口,采集接口，物模型
	if err = LoadAllCfg(); err != nil {
		return err
	}
	return
}

/********************************************************
功能描述：	增加接口
参数说明：
返回说明：
调用方式：
全局变量：
读写时间：
注意事项：
日期    ：
********************************************************/
func NewCollectInterface(coll *CollectInterfaceParamTemplate) (cit *CollectInterfaceTemplate, err error) {

	comm := CommunicationInterfaceMap.Get(coll.CommInterfaceName)

	if comm == nil {
		err = fmt.Errorf("no such comm interface %s", coll.CommInterfaceName)
		return
	}
	cit = &CollectInterfaceTemplate{
		CollInterfaceName:  coll.CollInterfaceName,
		CommInterfaceName:  coll.CommInterfaceName,
		CommInterface:      comm,
		CommMessage:        make([]*CommunicationMessageTemplate, 0),
		PollPeriod:         coll.PollPeriod,
		OfflinePeriod:      coll.OfflinePeriod,
		DeviceNodeCnt:      coll.DeviceNodeCnt,
		DeviceNodes:        coll.DeviceNodes,
		OfflineReportChan:  make(chan string, 100),
		OnlineReportChan:   make(chan string, 100),
		PropertyReportChan: make(chan string, 100),
	}
	comm.Bind(coll.CollInterfaceName)
	return
}

/********************************************************
功能描述：	修改接口
参数说明：
返回说明：
调用方式：
全局变量：
读写时间：
注意事项：
日期    ：
********************************************************/
func (d *CollectInterfaceTemplate) ModifyCollectInterface(pollPeriod, offlinePeriod int) {

	d.PollPeriod = pollPeriod
	d.OfflinePeriod = offlinePeriod
}

/********************************************************
功能描述：	增加单个节点
参数说明：
返回说明：
调用方式：
全局变量：
读写时间：
注意事项：
日期    ：
********************************************************/

func (d *CollectInterfaceTemplate) AddDeviceNode(dName string, dType string, dAddr string) error {

	node := &DeviceNodeTemplate{}
	node.Index = d.DeviceNodeCnt + 1
	node.Name = dName
	node.Addr = dAddr
	node.Type = dType
	node.LastCommRTC = "1970-01-01 00:00:00"
	node.CommStatus = OFFLINE
	variables, err := node.NewVariables()
	if err != nil {
		return err
	}
	node.VariableMap = make([]*VariableTemplate, 0, len(variables))
	node.VariableMap = append(node.VariableMap, variables...)
	d.DeviceNodes = append(d.DeviceNodes, node)
	d.DeviceNodeCnt++

	return nil
}

func (d *CollectInterfaceTemplate) InitDeviceNode(node *DeviceNodeTemplate) error {

	variables, err := node.NewVariables()
	if err != nil {
		return err
	}
	node.VariableMap = make([]*VariableTemplate, 0, len(variables))
	node.VariableMap = append(node.VariableMap, variables...)
	node.CommStatus = OFFLINE

	return nil
}

func (d *CollectInterfaceTemplate) DeleteDeviceNode(dName string) {

	var index int
	for k, node := range d.DeviceNodes {
		if node.Name == dName {
			index = k
			break
		}
	}
	d.DeviceNodes = append(d.DeviceNodes[:index], d.DeviceNodes[index+1:]...)
	d.DeviceNodeCnt--
}

func (d *CollectInterfaceTemplate) GetDeviceNode(dAddr string) *DeviceNodeTemplate {

	for _, v := range d.DeviceNodes {
		if v.Addr == dAddr {
			return v
		}
	}

	return nil
}
