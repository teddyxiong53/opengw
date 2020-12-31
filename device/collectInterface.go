package device

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
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

type CommunicationMessageTemplate struct {
	CollName  string `json:"CollInterfaceName"` //接口名称
	TimeStamp string `json:"TimeStamp"`         //时间戳
	Direction string `json:"DataDirection"`     //数据方向
	Content   string `json:"DataContent"`       //数据内容
}

//采集接口模板
type CollectInterfaceTemplate struct {
	CollInterfaceName   string                         `json:"CollInterfaceName"` //采集接口
	CommInterfaceName   string                         `json:"CommInterfaceName"` //通信接口
	CommMessage         []CommunicationMessageTemplate `json:"-"`
	PollPeriod          int                            `json:"PollPeriod"`          //采集周期
	OfflinePeriod       int                            `json:"OfflinePeriod"`       //离线超时周期
	DeviceNodeCnt       int                            `json:"DeviceNodeCnt"`       //设备数量
	DeviceNodeOnlineCnt int                            `json:"DeviceNodeOnlineCnt"` //设备在线数量
	DeviceNodeMap       []*DeviceNodeTemplate          `json:"DeviceNodeMap"`       //节点表
	OnlineReportChan    chan string                    `json:"-"`
	OfflineReportChan   chan string                    `json:"-"`
	PropertyReportChan  chan string                    `json:"-"`
}

var CollectInterfaceMap = make([]*CollectInterfaceTemplate, 0)

func WriteCollectInterfaceManageToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/collInterface.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open collInterface.json err", err)
		return
	}
	defer fp.Close()

	//采集接口配置参数
	type CollectInterfaceParamTemplate struct {
		CollInterfaceName string   `json:"CollInterfaceName"` //采集接口
		CommInterfaceName string   `json:"CommInterfaceName"` //通信接口
		PollPeriod        int      `json:"PollPeriod"`        //采集周期
		OfflinePeriod     int      `json:"OfflinePeriod"`     //离线超时周期
		DeviceNodeCnt     int      `json:"DeviceNodeCnt"`     //设备数量
		DeviceNodeNameMap []string `json:"DeviceNodeNameMap"` //节点名称
		DeviceNodeAddrMap []string `json:"DeviceNodeAddrMap"` //节点地址
		DeviceNodeTypeMap []string `json:"DeviceNodeTypeMap"` //节点类型

	}

	//定义采集接口参数结构体
	CollectInterfaceParamMap := struct {
		CollectInterfaceParam []CollectInterfaceParamTemplate
	}{
		CollectInterfaceParam: make([]CollectInterfaceParamTemplate, 0),
	}

	for _, v := range CollectInterfaceMap {
		ParamTemplate := CollectInterfaceParamTemplate{
			CollInterfaceName: v.CollInterfaceName,
			CommInterfaceName: v.CommInterfaceName,
			PollPeriod:        v.PollPeriod,
			OfflinePeriod:     v.OfflinePeriod,
			DeviceNodeCnt:     v.DeviceNodeCnt,
		}

		ParamTemplate.DeviceNodeNameMap = make([]string, 0)
		ParamTemplate.DeviceNodeAddrMap = make([]string, 0)
		ParamTemplate.DeviceNodeTypeMap = make([]string, 0)

		for i := 0; i < v.DeviceNodeCnt; i++ {
			ParamTemplate.DeviceNodeNameMap = append(ParamTemplate.DeviceNodeNameMap, v.DeviceNodeMap[i].Name)
			ParamTemplate.DeviceNodeAddrMap = append(ParamTemplate.DeviceNodeAddrMap, v.DeviceNodeMap[i].Addr)
			ParamTemplate.DeviceNodeTypeMap = append(ParamTemplate.DeviceNodeTypeMap, v.DeviceNodeMap[i].Type)
		}
		CollectInterfaceParamMap.CollectInterfaceParam = append(CollectInterfaceParamMap.CollectInterfaceParam,
			ParamTemplate)
	}

	sJson, _ := json.Marshal(CollectInterfaceParamMap)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write collInterface.json err", err)
	}
	log.Println("write collInterface.json sucess")
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func ReadCollectInterfaceManageFromJson() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/collInterface.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open collInterface.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		//采集接口配置参数
		type CollectInterfaceParamTemplate struct {
			CollInterfaceName string   `json:"CollInterfaceName"` //采集接口
			CommInterfaceName string   `json:"CommInterfaceName"` //通信接口
			PollPeriod        int      `json:"PollPeriod"`        //采集周期
			OfflinePeriod     int      `json:"OfflinePeriod"`     //离线超时周期
			DeviceNodeCnt     int      `json:"DeviceNodeCnt"`     //设备数量
			DeviceNodeNameMap []string `json:"DeviceNodeNameMap"` //节点名称
			DeviceNodeAddrMap []string `json:"DeviceNodeAddrMap"` //节点地址
			DeviceNodeTypeMap []string `json:"DeviceNodeTypeMap"` //节点类型
		}

		//定义采集接口参数结构体
		CollectInterfaceParamMap := struct {
			CollectInterfaceParam []CollectInterfaceParamTemplate
		}{
			CollectInterfaceParam: make([]CollectInterfaceParamTemplate, 0),
		}

		err = json.Unmarshal(data[:dataCnt], &CollectInterfaceParamMap)
		if err != nil {
			log.Println("collInterface unmarshal err", err)

			return false
		}

		log.Printf("CollectInterfaceParamMap %+v\n", CollectInterfaceParamMap)
		for k, v := range CollectInterfaceParamMap.CollectInterfaceParam {

			//创建接口实例
			CollectInterfaceMap = append(CollectInterfaceMap, NewCollectInterface(v.CollInterfaceName,
				v.CommInterfaceName,
				v.PollPeriod,
				v.OfflinePeriod,
				v.DeviceNodeCnt))

			//创建设备实例
			for i := 0; i < v.DeviceNodeCnt; i++ {
				CollectInterfaceMap[k].NewDeviceNode(
					v.DeviceNodeNameMap[i],
					v.DeviceNodeTypeMap[i],
					v.DeviceNodeAddrMap[i])
			}
		}

		return true
	} else {
		log.Println("collInterface.json is not exist")

		return false
	}
}

func DeviceNodeManageInit() {

	//设备模版
	ReadDeviceNodeTypeMap()
	//通信接口
	CommInterfaceInit()
	//采集接口
	if ReadCollectInterfaceManageFromJson() == true {
		log.Println("read collectInterface json ok")
		//log.Printf("collectMInterfaceMap %+v\n",CollectInterfaceMap)
	} else {

	}
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
func NewCollectInterface(collInterfaceName, commInterfaceName string,
	pollPeriod, offlinePeriod int, deviceNodeCnt int) *CollectInterfaceTemplate {

	nodeManage := &CollectInterfaceTemplate{
		CollInterfaceName:  collInterfaceName,
		CommInterfaceName:  commInterfaceName,
		CommMessage:        make([]CommunicationMessageTemplate, 0),
		PollPeriod:         pollPeriod,
		OfflinePeriod:      offlinePeriod,
		DeviceNodeCnt:      deviceNodeCnt,
		DeviceNodeMap:      make([]*DeviceNodeTemplate, 0),
		OfflineReportChan:  make(chan string, 100),
		OnlineReportChan:   make(chan string, 100),
		PropertyReportChan: make(chan string, 100),
	}

	return nodeManage
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
func (d *CollectInterfaceTemplate) NewDeviceNode(dName string, dType string, dAddr string) {

	node := &DeviceNodeTemplate{}
	node.Index = len(d.DeviceNodeMap)
	node.Name = dName
	node.Addr = dAddr
	node.Type = dType
	node.LastCommRTC = "1970-01-01 00:00:00"
	node.CommTotalCnt = 0
	node.CommSuccessCnt = 0
	node.CurCommFailCnt = 0
	node.CommStatus = "offLine"
	node.VariableMap = make([]VariableTemplate, 0)
	variables := node.NewVariables()
	node.VariableMap = append(node.VariableMap, variables...)

	d.DeviceNodeMap = append(d.DeviceNodeMap, node)
}

func (d *CollectInterfaceTemplate) AddDeviceNode(dName string, dType string, dAddr string) (bool, string) {

	node := &DeviceNodeTemplate{}
	node.Index = len(d.DeviceNodeMap)
	node.Name = dName
	node.Addr = dAddr
	node.Type = dType
	node.LastCommRTC = "1970-01-01 00:00:00"
	node.CommTotalCnt = 0
	node.CommSuccessCnt = 0
	node.CurCommFailCnt = 0
	node.CommStatus = "offLine"
	node.VariableMap = make([]VariableTemplate, 0)
	variables := node.NewVariables()
	node.VariableMap = append(node.VariableMap, variables...)

	d.DeviceNodeMap = append(d.DeviceNodeMap, node)

	d.DeviceNodeCnt++

	return true, "add success"
}

func (d *CollectInterfaceTemplate) DeleteDeviceNode(dName string) {

	for k, v := range d.DeviceNodeMap {
		if v.Name == dName {
			//d.DeviceNodeMap = d.DeviceNodeMap[k : k+1]
			d.DeviceNodeMap = append(d.DeviceNodeMap[:k], d.DeviceNodeMap[k+1:]...)
			d.DeviceNodeCnt--
		}
	}
}

func (d *CollectInterfaceTemplate) GetDeviceNode(dAddr string) *DeviceNodeTemplate {

	for _, v := range d.DeviceNodeMap {
		if v.Addr == dAddr {
			return v
		}
	}

	return nil
}
