package device

import (
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
	"time"
	"goAdapter/setting"
)

var MaxDeviceNodeCnt int = 50

type ValueTemplate struct {
	Value     interface{} //变量值，不可以是字符串
	Explain   interface{} //变量值解释，必须是字符串
	TimeStamp string
}

//变量标签模版
type VariableTemplate struct {
	Index int             `json:"index"` //变量偏移量
	Name  string          `json:"name"`  //变量名
	Label string          `json:"lable"` //变量标签
	Value []ValueTemplate `json:"value"` //变量值
	Type  string          `json:"type"`  //变量类型
}

//设备模板
type DeviceNodeTemplate struct {
	Index          int                `json:"Index"`          //设备偏移量
	Name           string             `json:"Name"`           //设备名称
	Addr           string             `json:"Addr"`           //设备地址
	Type           string             `json:"Type"`           //设备类型
	LastCommRTC    string             `json:"LastCommRTC"`    //最后一次通信时间戳
	CommTotalCnt   int                `json:"CommTotalCnt"`   //通信总次数
	CommSuccessCnt int                `json:"CommSuccessCnt"` //通信成功次数
	CurCommFailCnt int                `json:"-"`              //当前通信失败次数
	CommStatus     string             `json:"CommStatus"`     //通信状态
	VariableMap    []VariableTemplate `json:"-"`              //变量列表
}

func (d *DeviceNodeTemplate) NewVariables() []VariableTemplate {

	type LuaVariableTemplate struct {
		Index int
		Name  string
		Label string
		Type  string
	}

	type LuaVariableMapTemplate struct {
		Variable []*LuaVariableTemplate
	}

	for _,c := range CollectInterfaceMap{
		for _,n := range c.DeviceNodeMap{
			if n.Name == d.Name{
				//调用NewVariables
				err := c.LuaState.CallByParam(lua.P{
					Fn:      c.LuaState.GetGlobal("NewVariables"),
					NRet:    1,
					Protect: true,
				})
				if err != nil {
					setting.Logger.Warning("NewVariables err,", err)
				}

				//获取返回结果
				ret := c.LuaState.Get(-1)
				c.LuaState.Pop(1)

				LuaVariableMap := LuaVariableMapTemplate{}

				if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
					setting.Logger.Warning("NewVariables gluamapper.Map err,", err)
				}

				variables := make([]VariableTemplate, 0)

				for _, v := range LuaVariableMap.Variable {
					variable := VariableTemplate{}
					variable.Index = v.Index
					variable.Name = v.Name
					variable.Label = v.Label
					variable.Type = v.Type

					variable.Value = make([]ValueTemplate, 0)
					variables = append(variables, variable)
				}
				return variables
			}
		}
	}
	return nil
}

func (d *DeviceNodeTemplate) GenerateGetRealVariables(sAddr string, step int) ([]byte, bool, bool) {

	type LuaVariableMapTemplate struct {
		Status   string `json:"Status"`
		Variable []*byte
	}

	for _,c := range CollectInterfaceMap{
		for _,n := range c.DeviceNodeMap {
			if n.Name == d.Name {
				//调用NewVariables
				err := c.LuaState.CallByParam(lua.P{
					Fn:      c.LuaState.GetGlobal("GenerateGetRealVariables"),
					NRet:    1,
					Protect: true,
				}, lua.LString(sAddr), lua.LNumber(step))
				if err != nil {
					setting.Logger.Warning("GenerateGetRealVariables err,", err)
				}

				//获取返回结果
				ret := c.LuaState.Get(-1)
				c.LuaState.Pop(1)

				LuaVariableMap := LuaVariableMapTemplate{}
				if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
					setting.Logger.Warning("GenerateGetRealVariables gluamapper.Map err,", err)
				}

				ok := false
				con := false //后续是否有报文
				if LuaVariableMap.Status == "0" {
					con = false
				} else {
					con = true
				}
				nBytes := make([]byte, 0)
				if len(LuaVariableMap.Variable) > 0 {
					ok = true
					for _, v := range LuaVariableMap.Variable {
						nBytes = append(nBytes, *v)
					}
				} else {
					ok = false
				}

				return nBytes, ok, con
			}
		}
	}
	return nil, false, false
}

func (d *DeviceNodeTemplate) DeviceCustomCmd(sAddr string, cmdName string, cmdParam string, step int) ([]byte, bool, bool) {

	type LuaVariableMapTemplate struct {
		Status   string  `json:"Status"`
		Variable []*byte `json:"Variable"`
	}

	//log.Printf("cmdParam %+v\n", cmdParam)
	for _,c := range CollectInterfaceMap{
		for _,n := range c.DeviceNodeMap {
			if n.Name == d.Name {
				var err error
				var ret lua.LValue

				//调用DeviceCustomCmd
				err = c.LuaState.CallByParam(lua.P{
					Fn:      c.LuaState.GetGlobal("DeviceCustomCmd"),
					NRet:    1,
					Protect: true,
				}, lua.LString(sAddr),
					lua.LString(cmdName),
					lua.LString(cmdParam),
					lua.LNumber(step))
				if err != nil {
					setting.Logger.Warning("DeviceCustomCmd err,", err)
					return nil, false, false
				}

				//获取返回结果
				ret = c.LuaState.Get(-1)
				c.LuaState.Pop(1)

				LuaVariableMap := LuaVariableMapTemplate{}
				if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
					setting.Logger.Warning("DeviceCustomCmd gluamapper.Map err,", err)
					return nil, false, false
				}

				ok := false
				con := false //后续是否有报文
				if LuaVariableMap.Status == "0" {
					con = false
				} else {
					con = true
				}
				nBytes := make([]byte, 0)
				if len(LuaVariableMap.Variable) > 0 {
					ok = true
					for _, v := range LuaVariableMap.Variable {
						nBytes = append(nBytes, *v)
					}
				} else {
					ok = false
				}

				return nBytes, ok, con
			}
		}
	}

	return nil, false, false
}

func (d *DeviceNodeTemplate) AnalysisRx(sAddr string, variables []VariableTemplate, rxBuf []byte, rxBufCnt int) chan bool {

	status := make(chan bool, 1)

	type LuaVariableTemplate struct {
		Index   int
		Name    string
		Label   string
		Type    string
		Value   interface{}
		Explain string
	}

	type LuaVariableMapTemplate struct {
		Status   string `json:"Status"`
		Variable []*LuaVariableTemplate
	}

	for _,c := range CollectInterfaceMap{
		for _,n := range c.DeviceNodeMap {
			if n.Name == d.Name {
				tbl := lua.LTable{}
				for _, v := range rxBuf {
					tbl.Append(lua.LNumber(v))
				}

				c.LuaState.SetGlobal("rxBuf", luar.New(c.LuaState, &tbl))

				//AnalysisRx
				err := c.LuaState.CallByParam(lua.P{
					Fn:      c.LuaState.GetGlobal("AnalysisRx"),
					NRet:    1,
					Protect: true,
				}, lua.LString(sAddr), lua.LNumber(rxBufCnt))
				if err != nil {
					setting.Logger.Warning("AnalysisRx err,", err)
				}

				//获取返回结果
				ret := c.LuaState.Get(-1)
				c.LuaState.Pop(1)

				LuaVariableMap := LuaVariableMapTemplate{}

				if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
					setting.Logger.Warning("AnalysisRx gluamapper.Map err,", err)
				}

				timeNowStr := time.Now().Format("2006-01-02 15:04:05")
				value := ValueTemplate{}
				//log.Printf("LuaVariableMap %+v\n", LuaVariableMap)
				if LuaVariableMap.Status == "0" {
					if len(LuaVariableMap.Variable) > 0 {
						for _, lv := range LuaVariableMap.Variable {
							for k, v := range variables {
								if lv.Index == v.Index {
									variables[k].Index = lv.Index
									variables[k].Name = lv.Name
									variables[k].Label = lv.Label
									variables[k].Type = lv.Type

									value.Value = lv.Value
									value.Explain = lv.Explain
									value.TimeStamp = timeNowStr

									if len(variables[k].Value) < 100 {
										variables[k].Value = append(variables[k].Value, value)
									} else {
										variables[k].Value = variables[k].Value[1:]
										variables[k].Value = append(variables[k].Value, value)
									}
								}
							}
						}
					}
					status <- true
				}
			}
		}
	}
	return status
}
