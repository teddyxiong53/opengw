package device

type DeviceNodeTypeTemplate struct{
	TemplateID   int											`json:"templateID"`			//模板ID
	TemplateName string											`json:"templateName"`		//模板名称
	TemplateType string											`json:"templateType"`		//模板型号
	TemplateMessage string              						`json:"templateMessage"`	//备注信息
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