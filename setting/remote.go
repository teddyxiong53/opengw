package setting


type RemotePlatformTemplate struct{
	ID string					`json:"ID"`
	RemoteIP string				`json:"RemoteIP"`
	RemotePort string			`json:"RemotePort"`
	Protocal string         	`json:"Protocal"`
	ProtocalParam interface{}	`json:"ProtocalParam"`
}

type HttpRemoteTemplate struct{
	Method string				`json:"Method"`
	URL string					`json:"URL"`
	Timeout string				`json:"Timeout"`
}

type MQTTRemoteTemplate struct{
	UserName string
	Code string
	ClientID string
}

var RemotePlatform *RemotePlatformTemplate

func NewRemotePlatform(RemoteIP,RemotePort,Protocal string) *RemotePlatformTemplate{

	remote := &RemotePlatformTemplate{}

	if Protocal == "HTTP"{
		remote.RemoteIP = RemoteIP
		remote.RemotePort = RemotePort
		remote.Protocal = Protocal
		remote.ProtocalParam = HttpRemoteTemplate{}
	}

	return remote
}

func RemotePlatformInit(){

	RemotePlatform = NewRemotePlatform("192.168.1.1","60000","HTTP")
}

func (r *RemotePlatformTemplate)SetHTTPProtocalParam(method,url,timeout string){

	httpRemote := HttpRemoteTemplate{
		Method: method,
		URL: url,
		Timeout: timeout,
	}

	r.ProtocalParam = httpRemote
}