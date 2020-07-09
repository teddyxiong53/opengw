package device

import (
	"encoding/json"
	"goAdapter/api"
	"goAdapter/setting"
	"log"
	"os"
	"path/filepath"
)

const (
	MaxDeviceInterfaceManage int = 2

	InterFaceID0 int = 0
	InterFaceID1 int = 1
	InterFaceID2 int = 2
	InterFaceID3 int = 3
	InterFaceID4 int = 4
	InterFaceID5 int = 5
	InterFaceID6 int = 6
	InterFaceID7 int = 7

)

//通信接口模板
type DeviceInterfaceTemplate struct {
	InterfaceID       int                   `json:"InterfaceID"`   //通信接口
	PollPeriod        int                   `json:"PollPeriod"`    //采集周期
	OfflinePeriod     int                   `json:"OfflinePeriod"` //离线超时周期
	DeviceNodeCnt     int                   `json:"DeviceNodeCnt"` //设备数量
	DeviceNodeMap     []*DeviceNodeTemplate `json:"DeviceNodeMap"` //节点表
}

//配置参数
type DeviceInterfaceParamTemplate struct {
	InterfaceID       int      `json:"InterfaceID"`   //通信接口
	PollPeriod        int      `json:"PollPeriod"`    //采集周期
	OfflinePeriod     int      `json:"OfflinePeriod"` //离线超时周期
	DeviceNodeCnt     int      `json:"DeviceNodeCnt"` //设备数量
	DeviceNodeAddrMap []string `json:"DeviceAddrMap"` //节点地址
	DeviceNodeTypeMap []string `json:"DeviceTypeMap"` //节点类型
}

//配置参数
type DeviceInterfaceParamMapTemplate struct {
	DeviceInterfaceParam [MaxDeviceInterfaceManage]DeviceInterfaceParamTemplate
}

//var DeviceNodeTypeMap [MaxDeviceNodeTypeCnt]DeviceNodeTypeTemplate
var DeviceInterfaceMap [MaxDeviceInterfaceManage]*DeviceInterfaceTemplate
var DeviceInterfaceParamMap DeviceInterfaceParamMapTemplate


func WriteDeviceInterfaceManageToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/deviceNodeManage.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open deviceNodeManage.json err", err)
		return
	}
	defer fp.Close()

	for k, v := range DeviceInterfaceMap {
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].InterfaceID = v.InterfaceID
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].PollPeriod = v.PollPeriod
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].OfflinePeriod = v.OfflinePeriod
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeCnt = v.DeviceNodeCnt

		DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeAddrMap = DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeAddrMap[0:0]
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeTypeMap = DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeTypeMap[0:0]
		for i:=0;i<v.DeviceNodeCnt;i++{
			DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeAddrMap = append(DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeAddrMap,v.DeviceNodeMap[i].Addr)
			DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeTypeMap = append(DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeTypeMap,v.DeviceNodeMap[i].Type)
		}
	}

	sJson, _ := json.Marshal(DeviceInterfaceParamMap)

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

func ReadDeviceInterfaceManageFromJson() bool {

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

		for _,v := range DeviceInterfaceParamMap.DeviceInterfaceParam{

			v.DeviceNodeAddrMap = make([]string,0)
			v.DeviceNodeTypeMap = make([]string,0)
		}

		err = json.Unmarshal(data[:dataCnt], &DeviceInterfaceParamMap)
		if err != nil {
			log.Println("deviceNodeManage unmarshal err", err)

			return false
		}

		return true
	} else {
		log.Println("deviceNodeManage.json is not exist")

		return false
	}
}

func DeviceNodeManageInit() {

	ReadDeviceNodeTypeMapFromJson()

	CommInterfaceInit()

	if ReadDeviceInterfaceManageFromJson() == true {
		log.Println("read interface json ok")

		for k, v := range DeviceInterfaceParamMap.DeviceInterfaceParam {

			//创建接口实例
			DeviceInterfaceMap[k] = NewDeviceInterface(k,
				v.PollPeriod,
				v.OfflinePeriod,
				v.DeviceNodeCnt)

			//创建设备实例
			for i := 0; i < v.DeviceNodeCnt; i++ {
				DeviceInterfaceMap[k].NewDeviceNode(
					v.DeviceNodeTypeMap[i],
					v.DeviceNodeAddrMap[i])
			}
		}
	} else {

		for i := 0; i < MaxDeviceInterfaceManage; i++ {
			//创建接口实例
			DeviceInterfaceMap[i] = NewDeviceInterface(i,
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
func NewDeviceInterface(interfaceID, pollPeriod, offlinePeriod int, deviceNodeCnt int) *DeviceInterfaceTemplate {

	nodeManage := &DeviceInterfaceTemplate{
		InterfaceID:       interfaceID,
		PollPeriod:        pollPeriod,
		OfflinePeriod:     offlinePeriod,
		DeviceNodeCnt:     deviceNodeCnt,
		DeviceNodeMap:     make([]*DeviceNodeTemplate, 0),
	}

	//打开串口
	setting.SerialOpen(nodeManage.InterfaceID)

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
func (d *DeviceInterfaceTemplate) ModifyDeviceInterface(pollPeriod, offlinePeriod int) {

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
func (d *DeviceInterfaceTemplate) NewDeviceNode(dType string, dAddr string) {

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
	node.VariableMap = append(node.VariableMap,variables...)

	d.DeviceNodeMap = append(d.DeviceNodeMap, node)
}

func (d *DeviceInterfaceTemplate) AddDeviceNode(dType string, dAddr string) (bool, string) {

	node := &DeviceNodeTemplate{}
	node.Type = dType
	node.Addr = dAddr
	node.Index = len(d.DeviceNodeMap)
	node.VariableMap = make([]api.VariableTemplate, 0)
	variables := node.NewVariables()
	node.VariableMap = append(node.VariableMap,variables...)

	d.DeviceNodeMap = append(d.DeviceNodeMap, node)
	d.DeviceNodeCnt++

	return true, "add success"
}

func (d *DeviceInterfaceTemplate) DeleteDeviceNode(dAddr string, dType string) {

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

func (d *DeviceInterfaceTemplate) GetDeviceNode(dAddr string) interface{} {

	for _, v := range d.DeviceNodeMap {
		if v.Addr == dAddr {
			return v
		}
	}

	return nil
}