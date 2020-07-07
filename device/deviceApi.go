package device

type DeviceNodeTypeTemplate struct{
	TemplateID   int											`json:"templateID"`			//模板ID
	TemplateName string											`json:"templateName"`		//模板名称
	TemplateType string											`json:"templateType"`		//模板型号
	TemplateMessage string              						`json:"templateMessage"`	//备注信息
}

//变量标签模版
type VariableTemplate struct{
	Index   	int      										`json:"index"`			//变量偏移量
	Name 		string											`json:"name"`			//变量名
	Lable 		string											`json:"lable"`			//变量标签
	Value 		interface{}										`json:"value"`			//变量值
	TimeStamp   string											`json:"timestamp"`		//变量时间戳
	Type    	string                  						`json:"type"`			//变量类型
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

type DeviceNodeInterface interface {

	//解析接收到的数据
	ProcessRx(rxChan chan bool,rxBuf []byte,rxBufCnt int) chan bool
	//生成读变量的数据帧
	GetDeviceRealVariables() []byte
	//生成写变量的数据帧
	SetDeviceRealVariables() int
	//创建设备变量表
	NewVariables()
	//获取设备变量值
	GetDeviceVariablesValue() []VariableTemplate
}

type Build interface{
	New(index int,dAddr string,dType string)DeviceNodeInterface
}

var DeviceTemplateMap = map[string]Build{
	"modbus":&DeviceNodeModbusTemplate{},
}