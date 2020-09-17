package device

import (
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"goAdapter/setting"
	"layeh.com/gopher-luar"
	"time"
)

var MaxDeviceNodeCnt int = 50

type ValueTemplate struct{
	Value       interface{}		//变量值，不可以是字符串
	Explain     interface{}     //变量值解释，必须是字符串
	TimeStamp   string
}

//变量标签模版
type VariableTemplate struct{
	Index   	int      										`json:"index"`			//变量偏移量
	Name 		string											`json:"name"`			//变量名
	Label 		string											`json:"lable"`			//变量标签
	Value 		[]ValueTemplate									`json:"value"`			//变量值
	Type    	string                  						`json:"type"`			//变量类型
}

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
	VariableMap    []VariableTemplate 	  `json:"-"`    		  //变量列表
}

func (d *DeviceNodeTemplate) NewVariables() []VariableTemplate {

	type LuaVariableTemplate struct{
		Index int
		Name string
		Label string
		Type string
	}

	type LuaVariableMapTemplate struct{
		Variable []*LuaVariableTemplate
	}

	for k,v := range DeviceNodeTypeMap.DeviceNodeType{
		if d.Type == v.TemplateType{

			//调用NewVariables
			err := DeviceTypePluginMap[k].CallByParam(lua.P{
				Fn:DeviceTypePluginMap[k].GetGlobal("NewVariables"),
				NRet:1,
				Protect: true,
			})
			if err != nil{
				setting.Logger.Warning("NewVariables err,",err)
			}

			//获取返回结果
			ret := DeviceTypePluginMap[k].Get(-1)
			DeviceTypePluginMap[k].Pop(1)

			LuaVariableMap := LuaVariableMapTemplate{}

			if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
				setting.Logger.Warning("NewVariables gluamapper.Map err,",err)
			}

			variables := make([]VariableTemplate,0)

			for _,v := range LuaVariableMap.Variable{
				variable := VariableTemplate{}
				variable.Index = v.Index
				variable.Name = v.Name
				variable.Label = v.Label
				variable.Type = v.Type

				variable.Value = make([]ValueTemplate,0)
				variables = append(variables,variable)
			}
			return variables
		}
	}

	//for k,v := range DeviceNodeTypeMap.DeviceNodeType{
	//	if d.Type == v.TemplateType{
	//		newVariablesFun, _ := DeviceTypePluginMap[k].Lookup("NewVariables")
	//		variables := newVariablesFun.(func() []VariableTemplate)()
	//		return variables
	//	}
	//}
	return nil
}

func (d *DeviceNodeTemplate) GenerateGetRealVariables(sAddr string,step int) ([]byte,bool) {

	type LuaVariableMapTemplate struct{
		Variable []*byte
	}

	for k,v := range DeviceNodeTypeMap.DeviceNodeType {
		if d.Type == v.TemplateType {

			//调用NewVariables
			err := DeviceTypePluginMap[k].CallByParam(lua.P{
				Fn:DeviceTypePluginMap[k].GetGlobal("GenerateGetRealVariables"),
				NRet:1,
				Protect: true,
			},lua.LString(sAddr),lua.LNumber(step))
			if err != nil{
				setting.Logger.Warning("GenerateGetRealVariables err,",err)
			}

			//获取返回结果
			ret := DeviceTypePluginMap[k].Get(-1)
			DeviceTypePluginMap[k].Pop(1)

			LuaVariableMap := LuaVariableMapTemplate{}
			if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
				setting.Logger.Warning("GenerateGetRealVariables gluamapper.Map err,",err)
			}

			ok := false
			nBytes := make([]byte,0)

			if len(LuaVariableMap.Variable) > 0{
				ok = true
				for _,v := range LuaVariableMap.Variable{
					nBytes = append(nBytes,*v)
				}
			}else{
				ok = false
			}

			return nBytes,ok
		}
	}

	//for k,v := range DeviceNodeTypeMap.DeviceNodeType {
	//	if d.Type == v.TemplateType {
	//		generateGetRealVariablesFun, _ := DeviceTypePluginMap[k].Lookup("GenerateGetRealVariables")
	//		nBytes,ok := generateGetRealVariablesFun.(func(string,int) ([]byte,bool))(sAddr,step)
	//		return nBytes,ok
	//	}
	//}
	return nil,false
}

func (d *DeviceNodeTemplate) AnalysisRx(sAddr string,variables []VariableTemplate,rxBuf []byte,rxBufCnt int) chan bool{

	status := make(chan bool,1)

	type LuaVariableTemplate struct{
		Index   int
		Name    string
		Label   string
		Type    string
		Value   interface{}
		Explain string
	}

	type LuaVariableMapTemplate struct{
		Variable []*LuaVariableTemplate
	}

	for k,v := range DeviceNodeTypeMap.DeviceNodeType{
		if d.Type == v.TemplateType{

			tbl := lua.LTable{}
			for _,v := range rxBuf{
				tbl.Append(lua.LNumber(v))
			}

			DeviceTypePluginMap[k].SetGlobal("rxBuf", luar.New(DeviceTypePluginMap[k], &tbl))

			//调用NewVariables
			err := DeviceTypePluginMap[k].CallByParam(lua.P{
				Fn:DeviceTypePluginMap[k].GetGlobal("AnalysisRx"),
				NRet:1,
				Protect: true,
			},lua.LString(sAddr),lua.LNumber(rxBufCnt))
			if err != nil{
				setting.Logger.Warning("AnalysisRx err,",err)
			}

			//获取返回结果
			ret := DeviceTypePluginMap[k].Get(-1)
			DeviceTypePluginMap[k].Pop(1)

			LuaVariableMap := LuaVariableMapTemplate{}

			if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
				setting.Logger.Warning("AnalysisRx gluamapper.Map err,",err)
			}

			timeNowStr := time.Now().Format("2006-01-02 15:04:05")
			value := ValueTemplate{}
			if len(LuaVariableMap.Variable) > 0{
				for _,lv := range LuaVariableMap.Variable{
					for k,v := range variables{
						if lv.Index == v.Index{

							variables[k].Index = lv.Index
							variables[k].Name = lv.Name
							variables[k].Label = lv.Label
							variables[k].Type = lv.Type

							value.Value = lv.Value
							value.Explain = lv.Explain
							value.TimeStamp = timeNowStr
							variables[k].Value = append(variables[k].Value,value)
						}
					}
				}
				status <-true
			}
		}
	}


	//for k,v := range DeviceNodeTypeMap.DeviceNodeType {
	//	if d.Type == v.TemplateType {
	//		analysisRxFun, _ := DeviceTypePluginMap[k].Lookup("AnalysisRx")
	//		status = analysisRxFun.(func(string,[]VariableTemplate,[]byte, int) chan bool)(sAddr, variables, rxBuf, rxBufCnt)
	//		return status
	//	}
	//}
	return status
}