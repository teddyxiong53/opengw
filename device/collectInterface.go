package device

import "fmt"

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

const (
	ADD CollectAction = iota
	DELETE
	UPDATE
	GET
)

type CollectInterfaceStatus struct {
	Tmp *CollectInterfaceTemplate
	ACT CollectAction
}

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
	CollInterfaceName   string                          `json:"CollInterfaceName"` //采集接口
	CommInterfaceName   string                          `json:"CommInterfaceName"` //通信接口
	CommInterface       CommunicationInterface          `json:"-"`
	CommMessage         []*CommunicationMessageTemplate `json:"-"`
	PollPeriod          int                             `json:"PollPeriod"`          //采集周期
	OfflinePeriod       int                             `json:"OfflinePeriod"`       //离线超时周期
	DeviceNodeCnt       int                             `json:"DeviceNodeCnt"`       //设备数量
	DeviceNodeOnlineCnt int                             `json:"DeviceNodeOnlineCnt"` //设备在线数量
	DeviceNodes         []*DeviceNodeTemplate           `json:"DeviceNodeMap"`
	OnlineReportChan    chan string                     `json:"-"`
	OfflineReportChan   chan string                     `json:"-"`
	PropertyReportChan  chan string                     `json:"-"`
}

var CollectInterfaceMap = make(map[string]*CollectInterfaceTemplate)

func NodeManageInit() (err error) {

	//设备模版
	if err = ReadPlugins(PLUGINPATH); err != nil {
		return err
	}
	//通信接口
	if err = LoadJsonFile(COMMJSON); err != nil {
		return err
	}
	//采集接口
	return LoadJsonFile(COLLINTERFACEJSON)
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

	comm, ok := CommunicationInterfaceMap[coll.CommInterfaceName]
	if !ok {
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
