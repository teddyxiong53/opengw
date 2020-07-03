package main

import (
	"encoding/json"
	"log"
	"os"
)

type DeviceNodeTypeTemplate struct{
	TemplateID   int					`json:"templateID"`			//模板ID
	TemplateName string					`json:"templateName"`		//模板名称
	TemplateType string					`json:"templateType"`		//模板型号
	TemplateMessage string              `json:"templateMessage"`	//备注信息
}

type VariableTemplate struct{
	Index   	int      										`json:"index"`			//变量偏移量
	Name 		string											`json:"name"`			//变量名
	Lable 		string											`json:"lable"`			//变量标签
	Value 		interface{}										`json:"value"`			//变量值
	TimeStamp   string											`json:"timestamp"`		//变量时间戳
	Type    	string                  						`json:"type"`			//变量类型
}

type DeviceNodeInterface interface {
	//获取设备变量
	GetDeviceRealVariables(deviceAddr string) []byte
	//设置设备变量
	SetDeviceRealVariables(deviceAddr string) int
	//创建设备变量
	NewVariables()
	//接收解析
	ProcessRx(rxChan chan bool,rxBuf []byte,rxBufCnt int) chan bool
	//获取设备变量值
	GetDeviceVariablesValue() []VariableTemplate
}

//设备模板
type DeviceNodeTemplate struct{
	Index               int                     				`json:"Index"`					//设备偏移量
	Addr 				string									`json:"Addr"`					//设备地址
	Type 				string       							`json:"Type"`					//设备类型
	LastCommRTC 		string      							`json:"LastCommRTC"`          	//最后一次通信时间戳
	CommTotalCnt 		int										`json:"CommTotalCnt"`			//通信总次数
	CommSuccessCnt 		int										`json:"CommSuccessCnt"`			//通信成功次数
	CommStatus 			string         							`json:"CommStatus"`				//通信状态
	VariableMap    		[]VariableTemplate						`json:"-"`						//变量列表
}

//接口模板
type DeviceInterfaceTemplate struct{
	InterfaceID 		int										`json:"InterfaceID"`			//通信接口
	PollPeriod 			int										`json:"PollPeriod"`				//采集周期
	OfflinePeriod 		int      								`json:"OfflinePeriod"`			//离线超时周期
	DeviceNodeCnt       int                 					`json:"DeviceNodeCnt"`			//设备数量
	DeviceNodeMap       []DeviceNodeInterface 					`json:"DeviceNodeMap"`			//节点表
	DeviceNodeAddrMap   []string 								`json:"DeviceAddrMap"`			//节点地址
	DeviceNodeTypeMap   []string 								`json:"DeviceTypeMap"`			//节点类型
}

//配置参数
type DeviceInterfaceParamTemplate struct{
	InterfaceID 		int										`json:"InterfaceID"`			//通信接口
	PollPeriod 			int										`json:"PollPeriod"`				//采集周期
	OfflinePeriod 		int      								`json:"OfflinePeriod"`			//离线超时周期
	DeviceNodeCnt       int										`json:"DeviceNodeCnt"`			//设备数量
	DeviceNodeAddrMap   []string 								`json:"DeviceAddrMap"`			//节点地址
	DeviceNodeTypeMap   []string 								`json:"DeviceTypeMap"`			//节点类型
}

//配置参数
type DeviceInterfaceParamMapTemplate struct {
	DeviceInterfaceParam [MaxDeviceInterfaceManage]DeviceInterfaceParamTemplate
	DeviceNodeTypeMap 	 []DeviceNodeTypeTemplate
}


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

	MaxDeviceNodeCnt int = 50

	//最大设备模板数
	MaxDeviceNodeTypeCnt int = 10
)

//var DeviceNodeTypeMap [MaxDeviceNodeTypeCnt]DeviceNodeTypeTemplate
var DeviceInterfaceMap	[MaxDeviceInterfaceManage]*DeviceInterfaceTemplate
var DeviceInterfaceParamMap DeviceInterfaceParamMapTemplate

func DeviceNodeManageInit(){

	if ReadDeviceInterfaceManageFromJson() == true{
		log.Println("read interface json ok")

		for k,v := range DeviceInterfaceParamMap.DeviceInterfaceParam{

			//创建接口实例
			DeviceInterfaceMap[k] = NewDeviceInterface(k,
				v.PollPeriod,
				v.OfflinePeriod,
				v.DeviceNodeCnt)

			DeviceInterfaceMap[k].DeviceNodeAddrMap = v.DeviceNodeAddrMap
			DeviceInterfaceMap[k].DeviceNodeTypeMap = v.DeviceNodeTypeMap

			//创建设备实例
			for i:=0;i<v.DeviceNodeCnt;i++{
				DeviceInterfaceMap[k].NewDeviceNode(
					v.DeviceNodeAddrMap[i],
					v.DeviceNodeTypeMap[i])
			}
		}
	}else{
		//初始化设备模板
		DeviceInterfaceParamMap.DeviceNodeTypeMap = make([]DeviceNodeTypeTemplate,0)

		for i:=0;i<MaxDeviceInterfaceManage;i++{
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
func NewDeviceInterface(interfaceID,pollPeriod,offlinePeriod int,deviceNodeCnt int) *DeviceInterfaceTemplate{

	nodeManage := &DeviceInterfaceTemplate{
		InterfaceID		: interfaceID,
		PollPeriod		: pollPeriod,
		OfflinePeriod	: offlinePeriod,
		DeviceNodeCnt	: deviceNodeCnt,
		DeviceNodeMap   : make([]DeviceNodeInterface,0),
		DeviceNodeAddrMap   : make([]string,0),
		DeviceNodeTypeMap   : make([]string,0),
	}

	//打开串口
	serialOpen(nodeManage.InterfaceID)

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
func (d *DeviceInterfaceTemplate)ModifyDeviceInterface(pollPeriod,offlinePeriod int){

	d.PollPeriod = pollPeriod
	d.OfflinePeriod = offlinePeriod
}

func WriteDeviceInterfaceManageToJson(){

	fileDir := exeCurDir + "/selfpara/deviceNodeManage.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open deviceNodeManage.json err",err)
		return
	}
	defer fp.Close()


	for k,v := range DeviceInterfaceMap{
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].InterfaceID = v.InterfaceID
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].PollPeriod = v.PollPeriod
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].OfflinePeriod = v.OfflinePeriod
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeCnt = v.DeviceNodeCnt
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeAddrMap = v.DeviceNodeAddrMap
		DeviceInterfaceParamMap.DeviceInterfaceParam[k].DeviceNodeTypeMap = v.DeviceNodeTypeMap


	}

	sJson,_ := json.Marshal(DeviceInterfaceParamMap)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write deviceNodeManage.json err",err)
	}
	log.Println("write deviceNodeManage.json sucess")
}

func ReadDeviceInterfaceManageFromJson() bool{

	fileDir := exeCurDir + "/selfpara/deviceNodeManage.json"

	if FileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open deviceNodeManage.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		DeviceInterfaceParamMap.DeviceNodeTypeMap = make([]DeviceNodeTypeTemplate,0)

		err = json.Unmarshal(data[:dataCnt],&DeviceInterfaceParamMap)
		if err != nil {
			log.Println("deviceNodeManage unmarshal err", err)

			return false
		}

		log.Printf("%+v\n",DeviceInterfaceParamMap)

		return true
	}else{
		log.Println("deviceNodeManage.json is not exist")

		return false
	}
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
func (d *DeviceInterfaceTemplate)NewDeviceNode(dAddr string,dType string){


	if dType == "modbus"{
		index := len(d.DeviceNodeMap)
		node := &DeviceNodeModbusTemplate{}
		node.Addr = dAddr
		node.Type = dType
		node.Index = index
		d.DeviceNodeMap = append(d.DeviceNodeMap,node)
		d.DeviceNodeMap[index].NewVariables()
	}

}

func (d *DeviceInterfaceTemplate)AddDeviceNode(dAddr string,dType string) (bool,string){

	for _,v := range d.DeviceNodeAddrMap{
		if v == dAddr{
			return false,"addr is exist"
		}
	}
	if dType == "modbus"{
		index := len(d.DeviceNodeMap)
		node := &DeviceNodeModbusTemplate{}
		node.Addr = dAddr
		node.Type = dType
		node.Index = index
		d.DeviceNodeMap = append(d.DeviceNodeMap,node)
		d.DeviceNodeAddrMap = append(d.DeviceNodeAddrMap,dAddr)
		d.DeviceNodeTypeMap = append(d.DeviceNodeTypeMap,dType)
		d.DeviceNodeMap[index].NewVariables()
		d.DeviceNodeCnt++
	}else if dType == "modbus2"{

	}


	return true,"add success"
}

func (d *DeviceInterfaceTemplate)DeleteDeviceNode(dAddr string,dType string){

	log.Printf("addr %s\n",dAddr)
	log.Printf("type %s\n",dType)

	for k,v := range d.DeviceNodeAddrMap{
		if v == dAddr{
			d.DeviceNodeMap = d.DeviceNodeMap[k:k+1]
			d.DeviceNodeMap = append(d.DeviceNodeMap[:k],d.DeviceNodeMap[k+1:]...)
			d.DeviceNodeAddrMap = append(d.DeviceNodeAddrMap[:k],d.DeviceNodeAddrMap[k+1:]...)
			d.DeviceNodeTypeMap = append(d.DeviceNodeTypeMap[:k],d.DeviceNodeTypeMap[k+1:]...)
			d.DeviceNodeCnt--
		}
	}
}

func (d *DeviceInterfaceTemplate)GetDeviceNode(dAddr string) interface{} {

	for _,v := range d.DeviceNodeAddrMap{
		if v == dAddr{
			return v
		}
	}

	return nil
}


