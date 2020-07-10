package device

type CommunicationTcpInterface struct{
	Name     	string     	`json:"Name"`
	IP 			string 		`json:"IP"`
	Port 		string		`json:"Port"`
	Timeout  	string     	`json:"Timeout"`		//通信超时
	Interval 	string		`json:"Interval"`		//通信间隔
}

func (c *CommunicationTcpInterface)Open(param interface{}) bool{

	return true
}

func (c *CommunicationTcpInterface)Close() bool{

	return true
}

func (c *CommunicationTcpInterface)WriteData(data []byte) int{

	return 0
}

func (c *CommunicationTcpInterface)ReadData(data []byte) int{

	return 0
}
