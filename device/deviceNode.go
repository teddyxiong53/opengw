package device

import (
	"bytes"
	"goAdapter/setting"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

var MaxDeviceNodeCnt int = 50
var lock sync.Mutex

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

	lock.Lock()
	//setting.Logger.Debugf("DeviceTypePluginMap %v", DeviceTypePluginMap)
	for k, v := range DeviceNodeTypeMap.DeviceNodeType {
		if d.Type == v.TemplateType {
			//调用NewVariables
			//setting.Logger.Debugf("TemplateType %v", v.TemplateType)
			err := DeviceTypePluginMap[k].CallByParam(lua.P{
				Fn:      DeviceTypePluginMap[k].GetGlobal("NewVariables"),
				NRet:    1,
				Protect: true,
			})
			if err != nil {
				setting.Logger.Warning("NewVariables err,", err)
			}

			//获取返回结果
			ret := DeviceTypePluginMap[k].Get(-1)
			DeviceTypePluginMap[k].Pop(1)
			//setting.Logger.Debugf("DeviceTypePluginMap Get,%v", ret)

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
			lock.Unlock()
			//log.Printf("variables %v\n",variables)
			return variables
		}
	}
	lock.Unlock()
	return nil
}

func (d *DeviceNodeTemplate) GenerateGetRealVariables(sAddr string, step int) ([]byte, bool, bool) {

	type LuaVariableMapTemplate struct {
		Status   string `json:"Status"`
		Variable []*byte
	}

	lock.Lock()
	for k, v := range DeviceNodeTypeMap.DeviceNodeType {
		if d.Type == v.TemplateType {
			//调用NewVariables
			err := DeviceTypePluginMap[k].CallByParam(lua.P{
				Fn:      DeviceTypePluginMap[k].GetGlobal("GenerateGetRealVariables"),
				NRet:    1,
				Protect: true,
			}, lua.LString(sAddr), lua.LNumber(step))
			if err != nil {
				setting.Logger.Warning("GenerateGetRealVariables err,", err)
			}

			//获取返回结果
			ret := DeviceTypePluginMap[k].Get(-1)
			DeviceTypePluginMap[k].Pop(1)

			LuaVariableMap := LuaVariableMapTemplate{}
			if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
				setting.Logger.Warning("GenerateGetRealVariables gluamapper.Map err,", err)
			}

			ok := false
			con := false //后续是否有报文
			nBytes := make([]byte, 0)
			if len(LuaVariableMap.Variable) > 0 {
				ok = true
				for _, v := range LuaVariableMap.Variable {
					nBytes = append(nBytes, *v)
				}
				if LuaVariableMap.Status == "0" {
					con = false
				} else {
					con = true
				}
			} else {
				ok = true
			}
			lock.Unlock()
			return nBytes, ok, con
		}
	}
	lock.Unlock()
	return nil, false, false
}

func (d *DeviceNodeTemplate) DeviceCustomCmd(sAddr string, cmdName string, cmdParam string, step int) ([]byte, bool, bool) {

	type LuaVariableMapTemplate struct {
		Status   string  `json:"Status"`
		Variable []*byte `json:"Variable"`
	}

	lock.Lock()
	//log.Printf("cmdParam %+v\n", cmdParam)
	for k, v := range DeviceNodeTypeMap.DeviceNodeType {
		if d.Type == v.TemplateType {
			var err error
			var ret lua.LValue

			//调用DeviceCustomCmd
			err = DeviceTypePluginMap[k].CallByParam(lua.P{
				Fn:      DeviceTypePluginMap[k].GetGlobal("DeviceCustomCmd"),
				NRet:    1,
				Protect: true,
			}, lua.LString(sAddr),
				lua.LString(cmdName),
				lua.LString(cmdParam),
				lua.LNumber(step))
			if err != nil {
				setting.Logger.Warning("DeviceCustomCmd err,", err)
				lock.Unlock()
				return nil, false, false
			}

			//获取返回结果
			ret = DeviceTypePluginMap[k].Get(-1)
			DeviceTypePluginMap[k].Pop(1)

			LuaVariableMap := LuaVariableMapTemplate{}
			if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
				setting.Logger.Warning("DeviceCustomCmd gluamapper.Map err,", err)
				lock.Unlock()
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
			lock.Unlock()
			return nBytes, ok, con
		}
	}
	lock.Unlock()
	return nil, false, false
}

func getGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
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

	lock.Lock()
	for k, v := range DeviceNodeTypeMap.DeviceNodeType {
		if d.Type == v.TemplateType {
			tbl := lua.LTable{}
			for _, v := range rxBuf {
				tbl.Append(lua.LNumber(v))
			}
			DeviceTypePluginMap[k].SetGlobal("rxBuf", luar.New(DeviceTypePluginMap[k], &tbl))

			//AnalysisRx
			err := DeviceTypePluginMap[k].CallByParam(lua.P{
				Fn:      DeviceTypePluginMap[k].GetGlobal("AnalysisRx"),
				NRet:    1,
				Protect: true,
			}, lua.LString(sAddr), lua.LNumber(rxBufCnt))
			if err != nil {
				setting.Logger.Warning("AnalysisRx err,", err)
			}

			//获取返回结果
			ret := DeviceTypePluginMap[k].Get(-1)
			DeviceTypePluginMap[k].Pop(1)

			LuaVariableMap := LuaVariableMapTemplate{}

			if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
				setting.Logger.Warning("AnalysisRx gluamapper.Map err,", err)
			}

			timeNowStr := time.Now().Format("2006-01-02 15:04:05")
			value := ValueTemplate{}
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
								//log.Printf("LuaVariables %+v\n", variables[k])
							}
						}
					}
				}
				status <- true
			}
		}
	}
	lock.Unlock()
	return status
}
