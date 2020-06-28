package main

import (
	"encoding/json"
	"log"
	"os"
)

type VariableTemplate struct{
	Index   	int      				`json:"index"`			//变量偏移量
	Name 		string					`json:"name"`			//变量名
	Lable 		string					`json:"lable"`			//变量标签
	Value 		interface{}				`json:"value"`			//变量值
	TimeStamp   string					`json:"timestamp"`		//变量时间戳
	Type    	string                  `json:"type"`			//变量类型
}

//FCU设备模板
type FCUDeviceNodeTemplate struct{
	Index               int                     				`json:"Index"`					//设备偏移量
	Addr 				string									`json:"Addr"`					//设备地址
	//Type 				string       							`json:"Type"`					//设备类型
	LastCommRTC 		string      							`json:"LastCommRTC"`          	//最后一次通信时间戳
	CommTotalCnt 		int										`json:"CommTotalCnt"`			//通信总次数
	CommSuccessCnt 		int										`json:"CommSuccessCnt"`			//通信成功次数
	CommStatus 			string         							`json:"CommStatus"`				//通信状态
	VariableMap    		[]VariableTemplate						`json:"-"`						//变量列表
}

//通信接口的设备链表
type DeviceNodeManage struct{
	InterfaceID 		int										`json:"InterfaceID"`			//通信接口
	PollPeriod 			int										`json:"PollPeriod"`				//采集周期
	OfflinePeriod 		int      								`json:"OfflinePeriod"`			//离线超时周期
	DeviceNodeCnt       int                 					`json:"DeviceNodeCnt"`			//设备数量
	DeviceNodeMap       [MaxDeviceNodeCnt]interface{} 			`json:"DeviceNodeMap"`			//节点表
	DeviceNodeUseMap    [MaxDeviceNodeCnt]bool                  `json:"DeviceNodeUseMap"`		//节点占用表
	DeviceNodeTypeMap   [MaxDeviceNodeCnt]string				`json:"DeviceNodeTypeMap"`		//节点类型表
	DeviceNodeAddrMap   [MaxDeviceNodeCnt]string				`json:"DeviceNodeAddrMap"`		//节点地址表
}

type DeviceInterfaceManage struct{
	InterfaceManage 	[MaxDeviceNodeManage]*DeviceNodeManage
	DeviceNodeTypeMap   [MaxDeviceNodeCnt]string                `json:"DeviceNodeTypeMap"`		//节点类型链表
	DeviceNodeAddrMap   [MaxDeviceNodeCnt]string                `json:"DeviceNodeAddrMap"`		//节点地址链表
}

const (
	MaxDeviceNodeManage int = 8

	InterFaceID0 int = 0
	InterFaceID1 int = 1
	InterFaceID2 int = 2
	InterFaceID3 int = 3
	InterFaceID4 int = 4
	InterFaceID5 int = 5
	InterFaceID6 int = 6
	InterFaceID7 int = 7

	MaxDeviceNodeCnt int = 50
)

const (

	DeviceType_FCU200 int = 0
	DeviceType_TD200  int = 1

)


var DeviceNodeManageMap	[MaxDeviceNodeManage]*DeviceNodeManage
var deviceInterfaceManage DeviceInterfaceManage

func DeviceNodeManageInit(){

	if ReadDeviceInterfaceManageFromJson() == true{

		for i:=0;i<MaxDeviceNodeManage;i++{
			//创建设备表实例
			DeviceNodeManageMap[i] = NewDeviceNodeManage(i,
				deviceInterfaceManage.InterfaceManage[i].PollPeriod,
				deviceInterfaceManage.InterfaceManage[i].OfflinePeriod,
				deviceInterfaceManage.InterfaceManage[i].DeviceNodeCnt)

			//创建节点实例
			for y:=0;y<deviceInterfaceManage.InterfaceManage[i].DeviceNodeCnt;y++ {
				DeviceNodeManageMap[i].NewDeviceNode(
					deviceInterfaceManage.DeviceNodeAddrMap[y],
					deviceInterfaceManage.DeviceNodeTypeMap[y])
			}
		}
	}else{
		for i:=0;i<MaxDeviceNodeManage;i++{

			//添加设备表
			DeviceNodeManageMap[i] = NewDeviceNodeManage(i,
				60,
				3,
				0)

			for y:=0;y<MaxDeviceNodeCnt;y++{
				DeviceNodeManageMap[i].DeviceNodeUseMap[y] = false
			}
		}
	}



	//for i:=0;i<MaxDeviceNodeManage;i++{
	//	//添加设备表
	//	DeviceNodeManageMap[i] = NewDeviceNodeManage(i,
	//		60,
	//		3,
	//		0)
	//
	//	for y:=0;y<MaxDeviceNodeCnt;y++{
	//		DeviceNodeManageMap[i].DeviceNodeUseMap[y] = false
	//	}
	//}
}

/********************************************************
功能描述：	增加节点链表
参数说明：
返回说明：
调用方式：
全局变量：
读写时间：
注意事项：
日期    ：
********************************************************/
func NewDeviceNodeManage(interfaceID,pollPeriod,offlinePeriod int,deviceNodeCnt int) *DeviceNodeManage{

	nodeManage := &DeviceNodeManage{
		InterfaceID		: interfaceID,
		PollPeriod		: pollPeriod,
		OfflinePeriod	: offlinePeriod,
		DeviceNodeCnt	: deviceNodeCnt,
	}

	return nodeManage
}

/********************************************************
功能描述：	修改节点链表
参数说明：
返回说明：
调用方式：
全局变量：
读写时间：
注意事项：
日期    ：
********************************************************/
func (d *DeviceNodeManage)ModifyDeviceNodeManage(pollPeriod,offlinePeriod int){

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

	for i:=0;i<MaxDeviceNodeManage;i++{

		deviceInterfaceManage.InterfaceManage[i] = DeviceNodeManageMap[i]
		for y:=0;y<deviceInterfaceManage.InterfaceManage[i].DeviceNodeCnt;y++{

			deviceInterfaceManage.DeviceNodeAddrMap[y] = DeviceNodeManageMap[i].DeviceNodeAddrMap[y]
			deviceInterfaceManage.DeviceNodeTypeMap[y] = DeviceNodeManageMap[i].DeviceNodeTypeMap[y]
		}
	}

	sJson,_ := json.Marshal(deviceInterfaceManage)

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

		err = json.Unmarshal(data[:dataCnt],&deviceInterfaceManage)
		if err != nil {
			log.Println("deviceNodeManage unmarshal err", err)

			return false
		}

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
func (d *DeviceNodeManage)NewDeviceNode(dAddr string,dType string){

	for k,v := range d.DeviceNodeUseMap{
		if v == false{
			node := FCUDeviceNodeTemplate{}
			node.Addr = dAddr
			//node.Type = dType
			node.Index = k
			//node.NewFCUVariableTemplate()

			d.DeviceNodeMap[k] = node
			d.DeviceNodeAddrMap[k] = dAddr
			d.DeviceNodeTypeMap[k] = dType
			//置为占用
			d.DeviceNodeUseMap[k] = true
			break
		}
	}
}

func (d *DeviceNodeManage)AddDeviceNode(dAddr string,dType string) (bool,string){

	for k,v := range d.DeviceNodeUseMap{
		if d.DeviceNodeAddrMap[k] == dAddr{
			return false,"addr is exist"
		}else{
			if v == false {
				node := FCUDeviceNodeTemplate{}
				node.Addr = dAddr
				//node.Type = dType
				node.Index = k
				//node.NewFCUVariableTemplate()

				d.DeviceNodeMap[k] = node
				d.DeviceNodeAddrMap[k] = dAddr
				d.DeviceNodeTypeMap[k] = dType
				d.DeviceNodeCnt++
				//置为占用
				d.DeviceNodeUseMap[k] = true

				return true, "add sucess"
			}
		}
	}

	return false,"type id no exist"
}

func (d *DeviceNodeManage)DeleteDeviceNode(dAddr string,dType string){

	log.Printf("addr %s\n",dAddr)
	log.Printf("type %s\n",dType)

	for k,v := range d.DeviceNodeUseMap{
		if v == true{
			if d.DeviceNodeAddrMap[k] == dAddr{
				log.Printf("addr find ok\n")

				//置为空闲
				d.DeviceNodeUseMap[k] = false
				d.DeviceNodeAddrMap[k] = ""
				d.DeviceNodeTypeMap[k] = ""
				d.DeviceNodeCnt--
				d.DeviceNodeMap[k] = FCUDeviceNodeTemplate{}
			}
			break
		}
	}
}

func (d *DeviceNodeManage)GetDeviceNode(dAddr string) interface{} {

	for k,v := range d.DeviceNodeAddrMap{

		if v == dAddr{
			return d.DeviceNodeMap[k]
		}
	}

	return nil
}

func (d *DeviceNodeManage)ModifyDeviceNode(index int,dType string){

}

