package device

import "log"

type CommunicationInterface interface{
	Open()  				bool
	Close() 				bool
	WriteData(data []byte) 	int
	ReadData(data []byte)  	int
}

type CommunicationTemplate struct{
	Name 	string										`json:"Name"`			//接口名称
	Type    string          							`json:"Type"`			//接口类型,比如serial,tcp,udp,http
	Status  bool 										`json:"-"`			    //接口状态
}

type CommunicationInterfaceListTemplate struct{
	CommunicationInterfaceMap []CommunicationInterface
}

var CommunicationInterfaceList CommunicationInterfaceListTemplate

func CommInterfaceInit() {

	if ReadCommSerialInterfaceListFromJson() == false{

	}else{
		log.Println("read CommSerialInterfaceList.json ok")

		for _,v := range CommunicationSerialMap{

			CommunicationInterfaceList.CommunicationInterfaceMap = append(CommunicationInterfaceList.CommunicationInterfaceMap,&v)
		}
	}

	for _,v := range CommunicationInterfaceList.CommunicationInterfaceMap{

		v.Open()
	}
}