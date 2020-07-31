package device

import "log"

type CommunicationInterface interface {
	Open() bool
	Close() bool
	WriteData(data []byte) int
	ReadData(data []byte) int
}

type CommunicationTemplate struct {
	Name   string `json:"Name"` //接口名称
	Type   string `json:"Type"` //接口类型,比如serial,tcp,udp,http
	Status bool   `json:"-"`    //接口状态
}

//通信接口Map
var CommunicationInterfaceMap = make([]CommunicationInterface, 0)

func CommInterfaceInit() {

	//获取串口通信接口参数
	if ReadCommSerialInterfaceListFromJson() == false {

	} else {
		log.Println("read CommSerialInterfaceList.json ok")

		//for _, v := range CommunicationSerialMap {
		//
		//	CommunicationInterfaceMap = append(CommunicationInterfaceMap, &v)
		//}
	}

	//打开串口通信
	for _, v := range CommunicationSerialMap {

		v.Open()
		//log.Printf("c %+v\n",v)
	}
}
