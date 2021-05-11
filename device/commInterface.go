package device

type CommunicationInterface interface {
	Open() bool
	Close() bool
	WriteData(data []byte) int
	ReadData(data []byte) int
	GetName() string
	GetTimeOut() string
	GetInterval() string
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
	if ReadCommSerialInterfaceListFromJson() == true {
		for _, v := range CommunicationSerialMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	//获取TCP通信接口参数
	if ReadCommTcpInterfaceListFromJson() == true {
		for _, v := range CommunicationTcpMap {
			CommunicationInterfaceMap = append(CommunicationInterfaceMap, v)
		}
	}

	for _, v := range CommunicationInterfaceMap {
		v.Open()
	}
}
