package device

type CommunicationTcpInterface struct{
	Name     	string     	`json:"Name"`
	IP 			string 		`json:"IP"`
	Port 		string		`json:"Port"`
	Timeout  	string     	`json:"Timeout"`		//通信超时
	Interval 	string		`json:"Interval"`		//通信间隔
}

func NewCommunicationTcpInterface(){


}

func (c *CommunicationTcpInterface)Open(param interface{}) bool{


}

func (c *CommunicationTcpInterface)Close() bool{


}

func (c *CommunicationTcpInterface)WriteWriteData() int{


}

func (c *CommunicationTcpInterface)ReadWriteData() int{


}
