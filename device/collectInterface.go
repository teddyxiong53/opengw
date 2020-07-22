package device

import (
	"encoding/json"
	"goAdapter/api"
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

//采集接口模板
type CollectInterfaceTemplate struct {
	InterfaceID   int                   `json:"InterfaceID"`   //接口ID
	PollPeriod    int                   `json:"PollPeriod"`    //采集周期
	OfflinePeriod int                   `json:"OfflinePeriod"` //离线超时周期
	DeviceNodeCnt int                   `json:"DeviceNodeCnt"` //设备数量
	DeviceNodeMap []*DeviceNodeTemplate `json:"DeviceNodeMap"` //节点表
}

//采集接口配置参数
type CollectInterfaceParamTemplate struct {
	InterfaceID       int      `json:"InterfaceID"`   //接口ID
	PollPeriod        int      `json:"PollPeriod"`    //采集周期
	OfflinePeriod     int      `json:"OfflinePeriod"` //离线超时周期
	DeviceNodeCnt     int      `json:"DeviceNodeCnt"` //设备数量
	DeviceNodeAddrMap []string `json:"DeviceAddrMap"` //节点地址
	DeviceNodeTypeMap []string `json:"DeviceTypeMap"` //节点类型
}

//配置参数
//type CollectInterfaceParamMapTemplate struct {
//	CollectInterfaceParam [MaxCollectInterfaceManage]CollectInterfaceParamTemplate
//}

var CollectInterfaceMap [MaxCollectInterfaceManage]*CollectInterfaceTemplate
//var CollectInterfaceParamMap CollectInterfaceParamMapTemplate

func WriteCollectInterfaceManageToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/deviceNodeManage.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open deviceNodeManage.json err", err)
		return
	}
	defer fp.Close()

	//定义采集接口参数结构体
	CollectInterfaceParamMap := struct {
		CollectInterfaceParam []CollectInterfaceParamTemplate
	}{
		CollectInterfaceParam: make([]CollectInterfaceParamTemplate, 0),
	}

	for k, v := range CollectInterfaceMap {
		CollectInterfaceParamMap.CollectInterfaceParam[k].InterfaceID = v.InterfaceID
		CollectInterfaceParamMap.CollectInterfaceParam[k].PollPeriod = v.PollPeriod
		CollectInterfaceParamMap.CollectInterfaceParam[k].OfflinePeriod = v.OfflinePeriod
		CollectInterfaceParamMap.CollectInterfaceParam[k].DeviceNodeCnt = v.DeviceNodeCnt

		CollectInterfaceParamMap.CollectInterfaceParam[k].DeviceNodeAddrMap = CollectInterfaceParamMap.CollectInterfaceParam[k].DeviceNodeAddrMap[0:0]
		CollectInterfaceParamMap.CollectInterfaceParam[k].DeviceNodeTypeMap = CollectInterfaceParamMap.CollectInterfaceParam[k].DeviceNodeTypeMap[0:0]
		for i := 0; i < v.DeviceNodeCnt; i++ {
			CollectInterfaceParamMap.CollectInterfaceParam[k].DeviceNodeAddrMap = append(CollectInterfaceParamMap.CollectInterfaceParam[k].DeviceNodeAddrMap, v.DeviceNodeMap[i].Addr)
			CollectInterfaceParamMap.CollectInterfaceParam[k].DeviceNodeTypeMap = append(CollectInterfaceParamMap.CollectInterfaceParam[k].DeviceNodeTypeMap, v.DeviceNodeMap[i].Type)
		}
	}

	sJson, _ := json.Marshal(CollectInterfaceParamMap)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write deviceNodeManage.json err", err)
	}
	log.Println("write deviceNodeManage.json sucess")
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func ReadCollectInterfaceManageFromJson() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/deviceNodeManage.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open deviceNodeManage.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		//定义采集接口参数结构体
		CollectInterfaceParamMap := struct {
			CollectInterfaceParam []CollectInterfaceParamTemplate
		}{
			CollectInterfaceParam: make([]CollectInterfaceParamTemplate, 0),
		}

		for _, v := range CollectInterfaceParamMap.CollectInterfaceParam {
			v.DeviceNodeAddrMap = make([]string, 0)
			v.DeviceNodeTypeMap = make([]string, 0)
		}

		err = json.Unmarshal(data[:dataCnt], &CollectInterfaceParamMap)
		if err != nil {
			log.Println("deviceNodeManage unmarshal err", err)

			return false
		}

		log.Printf("CollectInterfaceParamMap %+v\n",CollectInterfaceParamMap)
		for k, v := range CollectInterfaceParamMap.CollectInterfaceParam {

			//创建接口实例
			CollectInterfaceMap[k] = NewCollectInterface(k,
				v.PollPeriod,
				v.OfflinePeriod,
				v.DeviceNodeCnt)

			//创建设备实例
			for i := 0; i < v.DeviceNodeCnt; i++ {
				CollectInterfaceMap[k].NewDeviceNode(
					v.DeviceNodeTypeMap[i],
					v.DeviceNodeAddrMap[i])
			}
		}

		return true
	} else {
		log.Println("deviceNodeManage.json is not exist")

		return false
	}
}

func DeviceNodeManageInit() {
	//设备模版
	ReadDeviceNodeTypeMapFromJson()
	//通信接口
	CommInterfaceInit()
	//采集接口
	if ReadCollectInterfaceManageFromJson() == true {
		log.Println("read collectInterface json ok")
		//log.Printf("collectMInterfaceMap %+v\n",CollectInterfaceMap)
	} else {

		for i := 0; i < MaxCollectInterfaceManage; i++ {
			//创建接口实例
			CollectInterfaceMap[i] = NewCollectInterface(i,
				60,
				180,
				0)
		}
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
func NewCollectInterface(interfaceID, pollPeriod, offlinePeriod int, deviceNodeCnt int) *CollectInterfaceTemplate {

	nodeManage := &CollectInterfaceTemplate{
		InterfaceID:   interfaceID,
		PollPeriod:    pollPeriod,
		OfflinePeriod: offlinePeriod,
		DeviceNodeCnt: deviceNodeCnt,
		DeviceNodeMap: make([]*DeviceNodeTemplate, 0),
	}

	//打开串口
	//setting.SerialOpen(nodeManage.InterfaceID)

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
func (d *CollectInterfaceTemplate) NewDeviceNode(dType string, dAddr string) {

	//builder,ok := DeviceTemplateMap[dType]
	//if !ok{
	//	panic("deviceNodeType is not exist")
	//}
	//
	//index := len(d.DeviceNodeMap)
	//node := builder.New(index,dType,dAddr)
	//d.DeviceNodeMap = append(d.DeviceNodeMap,node)

	node := &DeviceNodeTemplate{}
	node.Type = dType
	node.Addr = dAddr
	node.Index = len(d.DeviceNodeMap)
	node.VariableMap = make([]api.VariableTemplate, 0)
	variables := node.NewVariables()
	node.VariableMap = append(node.VariableMap, variables...)

	d.DeviceNodeMap = append(d.DeviceNodeMap, node)
}

func (d *CollectInterfaceTemplate) AddDeviceNode(dType string, dAddr string) (bool, string) {

	node := &DeviceNodeTemplate{}
	node.Type = dType
	node.Addr = dAddr
	node.Index = len(d.DeviceNodeMap)
	node.VariableMap = make([]api.VariableTemplate, 0)
	variables := node.NewVariables()
	node.VariableMap = append(node.VariableMap, variables...)

	d.DeviceNodeMap = append(d.DeviceNodeMap, node)
	d.DeviceNodeCnt++

	return true, "add success"
}

func (d *CollectInterfaceTemplate) DeleteDeviceNode(dAddr string, dType string) {

	log.Printf("addr %s\n", dAddr)
	log.Printf("type %s\n", dType)

	for k, v := range d.DeviceNodeMap {
		if v.Addr == dAddr {
			d.DeviceNodeMap = d.DeviceNodeMap[k : k+1]
			d.DeviceNodeMap = append(d.DeviceNodeMap[:k], d.DeviceNodeMap[k+1:]...)
			d.DeviceNodeCnt--
		}
	}
}

func (d *CollectInterfaceTemplate) GetDeviceNode(dAddr string) interface{} {

	for _, v := range d.DeviceNodeMap {
		if v.Addr == dAddr {
			return v
		}
	}

	return nil
}
