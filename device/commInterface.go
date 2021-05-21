package device

import "goAdapter/setting"

type CommunicationInterface interface {
	Open() bool
	Close() bool
	WriteData(data []byte) int
	ReadData(data []byte) int
	GetName() string
	GetTimeOut() string
	GetInterval() string
}

//通信接口Map
var CommunicationInterfaceMap = make([]CommunicationInterface, 0)

func CommInterfaceInit() {

	//获取串口通信接口参数
	if ReadCommSerialInterfaceListFromJson() == true {
		for _, v := range CommunicationSerialMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取TcpClient通信接口参数
	if ReadCommTcpClientInterfaceListFromJson() == true {
		for _, v := range CommunicationTcpClientMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取开关量输出通信接口参数
	if ReadCommIoOutInterfaceListFromJson() == true {
		for _, v := range CommunicationIoOutMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取开关量输入通信接口参数
	if ReadCommIoInInterfaceListFromJson() == true {
		for _, v := range CommunicationIoInMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	for _, v := range CommunicationInterfaceMap {
		setting.Logger.Debugf("commName %v,", v.GetName())
		v.Open()
	}
}
