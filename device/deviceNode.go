package device

import (
	"goAdapter/api"
)

var MaxDeviceNodeCnt int = 50

//设备模板
type DeviceNodeTemplate struct {
	Index          int                    `json:"Index"`          //设备偏移量
	Name           string  				  `json:"Name"`           //设备名称
	Addr           string                 `json:"Addr"`           //设备地址
	Type           string                 `json:"Type"`           //设备类型
	LastCommRTC    string                 `json:"LastCommRTC"`    //最后一次通信时间戳
	CommTotalCnt   int                    `json:"CommTotalCnt"`   //通信总次数
	CommSuccessCnt int                    `json:"CommSuccessCnt"` //通信成功次数
	CurCommFailCnt int 				      `json:"-"` 			  //当前通信失败次数
	CommStatus     string                 `json:"CommStatus"`     //通信状态
	VariableMap    []api.VariableTemplate `json:"-"`    //变量列表
}

func (d *DeviceNodeTemplate) NewVariables() []api.VariableTemplate {

	for k,v := range DeviceNodeTypeMap.DeviceNodeType{
		if d.Type == v.TemplateType{
			newVariablesFun, _ := DeviceTypePluginMap[k].Lookup("NewVariables")
			variables := newVariablesFun.(func() []api.VariableTemplate)()
			return variables
		}
	}
	return nil
}

func (d *DeviceNodeTemplate) GenerateGetRealVariables(sAddr string,step int) ([]byte,bool) {

	for k,v := range DeviceNodeTypeMap.DeviceNodeType {
		if d.Type == v.TemplateType {
			generateGetRealVariablesFun, _ := DeviceTypePluginMap[k].Lookup("GenerateGetRealVariables")
			nBytes,ok := generateGetRealVariablesFun.(func(string,int) ([]byte,bool))(sAddr,step)
			return nBytes,ok
		}
	}
	return nil,false
}

func (d *DeviceNodeTemplate) AnalysisRx(sAddr string,variables []api.VariableTemplate,rxBuf []byte,rxBufCnt int) chan bool{

	status := make(chan bool,1)

	for k,v := range DeviceNodeTypeMap.DeviceNodeType {
		if d.Type == v.TemplateType {
			analysisRxFun, _ := DeviceTypePluginMap[k].Lookup("AnalysisRx")
			status = analysisRxFun.(func(string,[]api.VariableTemplate,[]byte, int) chan bool)(sAddr, variables, rxBuf, rxBufCnt)
			return status
		}
	}
	return status
}