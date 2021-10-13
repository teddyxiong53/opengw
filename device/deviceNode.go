package device

import (
	"errors"
	"fmt"
	"goAdapter/httpServer/model"
	"goAdapter/pkg/mylog"
	"log"
	"reflect"

	"sync"
	"time"

	"github.com/5anthosh/chili/parser"
	"github.com/shopspring/decimal"
	"github.com/walkmiao/chili/environment"
	"github.com/walkmiao/chili/evaluator"

	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

var env *environment.Environment
var eval *evaluator.Evaluator

var MaxDeviceNodeCnt int = 50
var lock sync.Mutex

func initializeEval() {
	log.Println("初始化表达式....")
	env = environment.New()
	env.SetDefaultFunctions()
	env.SetDefaultVariables()
	eval = evaluator.New(env)
}

type ValueTemplate struct {
	Value     interface{} //变量值，不可以是字符串
	Explain   string      //变量值解释，必须是字符串
	TimeStamp string
}

//变量标签模版
type VariableTemplate struct {
	Index     int             `json:"index"`     //变量偏移量
	Name      string          `json:"name"`      //变量名
	Label     string          `json:"lable"`     //变量标签
	Values    []ValueTemplate `json:"value"`     //变量值
	Type      string          `json:"type"`      //变量类型
	ChannelNo uint32          `json:"channelNo"` //通道号
	Changed   bool            `json:"changed"`   //是否变化
}

//设备模板
type DeviceNodeTemplate struct {
	Index          int    `json:"Index"`          //设备偏移量
	Name           string `json:"Name"`           //设备名称
	Addr           string `json:"Addr"`           //设备地址
	Type           string `json:"Type"`           //设备类型
	LastCommRTC    string `json:"LastCommRTC"`    //最后一次通信时间戳
	CommTotalCnt   int    `json:"CommTotalCnt"`   //通信总次数
	CommSuccessCnt int    `json:"CommSuccessCnt"` //通信成功次数
	CurCommFailCnt int    `json:"-"`              //当前通信失败次数
	CommStatus     string `json:"CommStatus"`     //通信状态
	//VariableMap    []*VariableTemplate                `json:"-"`              //变量列表
	Properties []model.DeviceTSLPropertyTemplate `json:"-"` //属性列表
	Services   []model.DeviceTSLServiceTempalte  `json:"-"` //服务
	Parser     Parser                            `json:"-"` //表达式解析器
}

func ClearPropertyValue(properties []model.DeviceTSLPropertyTemplate) {
	for i := 0; i < len(properties); i++ {
		properties[i].Value = make([]model.DeviceTSLPropertyValueTemplate, 0)
	}
}

func (d *DeviceNodeTemplate) NewVariablesForTSL() error {
	tmps := DeviceTSLMap.GetAll()
	for _, v := range tmps {
		if v.Plugin == d.Type {
			d.Properties = make([]model.DeviceTSLPropertyTemplate, len(v.Properties))
			copy(d.Properties, v.Properties)
			ClearPropertyValue(d.Properties)
			d.Services = make([]model.DeviceTSLServiceTempalte, len(v.Services))
			copy(d.Properties, v.Properties)
			return nil
		}
	}
	return fmt.Errorf("tsl template not bind plugin %s", d.Type)
}

func (d *DeviceNodeTemplate) NewVariables() (variables []*VariableTemplate, err error) {

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
	defer lock.Unlock()
	template, ok := DeviceTemplateMap[d.Type]
	if !ok {
		err = fmt.Errorf("未发现设备模板: %s", d.Name)
		return
	}
	lState := template.LuaState
	if lState == nil {
		err = fmt.Errorf("device lua template is not initialized")
		return
	}

	err = lState.CallByParam(lua.P{
		Fn:      template.LuaState.GetGlobal(string(NewVariables)),
		NRet:    1,
		Protect: true,
	})

	if err != nil {
		err = fmt.Errorf("NewVariables Err:%v", err)
		return
	}

	//获取返回结果
	ret := lState.Get(-1)
	lState.Pop(1)
	//setting.ZAPS.Debugf("DeviceTemplateLuaMap Get,%v", ret)

	LuaVariableMap := LuaVariableMapTemplate{}

	table, ok := ret.(*lua.LTable)
	if !ok {
		return nil, errors.New("ret is not lua LTable")
	}

	if err = gluamapper.Map(table, &LuaVariableMap); err != nil {
		return
	}

	variables = make([]*VariableTemplate, 0, len(LuaVariableMap.Variable))

	for _, v := range LuaVariableMap.Variable {
		variable := &VariableTemplate{}
		variable.Index = v.Index
		variable.Name = v.Name
		variable.Label = v.Label
		variable.Type = v.Type

		variable.Values = make([]ValueTemplate, 0, 100)
		variables = append(variables, variable)
	}
	return
}

func (d *DeviceNodeTemplate) GenerateGetRealVariables(sAddr string, step int) (nBytes []byte, hasFrame bool, err error) {

	type LuaVariableMapTemplate struct {
		HasFrame bool `json:"HasFrame"` //设备是否还有后续帧
		Variable []*byte
	}

	lock.Lock()
	template, ok := DeviceTemplateMap[d.Type]
	if !ok {
		err = fmt.Errorf("设备模板【%s】未安装", d.Name)
		return
	}
	lState := template.LuaState
	if lState == nil {
		err = fmt.Errorf("device lua template is not initialized")
		return
	}

	err = lState.CallByParam(lua.P{
		Fn:      lState.GetGlobal(string(GENERATEREAL)),
		NRet:    1,
		Protect: true,
	}, lua.LString(sAddr), lua.LNumber(step))
	if err != nil {
		err = fmt.Errorf("GenerateGetRealVariables Err:%v", err)
		return
	}

	//获取返回结果
	ret := lState.Get(-1)
	lState.Pop(1)

	LuaVariableMap := LuaVariableMapTemplate{}
	table, ok := ret.(*lua.LTable)
	if !ok {
		err = fmt.Errorf("%v is not lua LTable", ret)
		return
	}
	if err = gluamapper.Map(table, &LuaVariableMap); err != nil {
		return
	}

	nBytes = make([]byte, 0, len(LuaVariableMap.Variable))
	if len(LuaVariableMap.Variable) > 0 {
		for _, v := range LuaVariableMap.Variable {
			nBytes = append(nBytes, *v)
		}
		if LuaVariableMap.HasFrame {
			hasFrame = true
		}
	}
	lock.Unlock()
	return
}

func (d *DeviceNodeTemplate) DeviceCustomCmd(
	sAddr string, cmdName LUAFUNC, cmdParam string, step int) (nBytes []byte, isContinue bool, err error) {

	type LuaVariableMapTemplate struct {
		Status   string  `json:"Status"`
		Variable []*byte `json:"Variable"`
	}

	lock.Lock()
	defer lock.Unlock()
	template, ok := DeviceTemplateMap[d.Type]
	if !ok {
		err = fmt.Errorf("no such device template %s", d.Type)
		return
	}
	state := template.LuaState
	if state == nil {
		err = fmt.Errorf("template %s is not initialized lua file ", d.Type)
		return
	}

	var ret lua.LValue

	//调用DeviceCustomCmd
	err = state.CallByParam(lua.P{
		Fn:      state.GetGlobal(string(DEVICECUSTOMCMD)),
		NRet:    1,
		Protect: true,
	}, lua.LString(sAddr),
		lua.LString(cmdName),
		lua.LString(cmdParam),
		lua.LNumber(step))
	if err != nil {
		err = fmt.Errorf("DeviceCustomCmd Err:%v", err)
		return
	}

	//获取返回结果
	ret = state.Get(-1)
	state.Pop(1)
	table, ok := ret.(*lua.LTable)
	if !ok {
		err = fmt.Errorf("ret is not type of lua LTable")
		return
	}
	LuaVariableMap := LuaVariableMapTemplate{}
	if err = gluamapper.Map(table, &LuaVariableMap); err != nil {
		return
	}

	if LuaVariableMap.Status != "0" {
		isContinue = true
	}

	nBytes = make([]byte, 0, len(LuaVariableMap.Variable))
	if len(LuaVariableMap.Variable) > 0 {
		for _, v := range LuaVariableMap.Variable {
			nBytes = append(nBytes, *v)
		}
	} else {
		err = fmt.Errorf("LuaVariableMap's Variable less than 0")
		return
	}
	return
}

func (d *DeviceNodeTemplate) AnalysisRx(sAddr string, properties []model.DeviceTSLPropertyTemplate, rxBuf []byte, rxBufCnt int, txBuf []byte) error {
	if len(properties) <= 0 {
		mylog.ZAPS.Warn("node %s properties type %s is not defined", d.Name, d.Type)
	}
	type LuaVariableTemplate struct {
		Index   int
		Name    string
		Label   string
		Type    string
		Value   interface{}
		Explain string
		Formula string
	}

	type LuaVariableMapTemplate struct {
		Status    string `json:"Status"`              //状态是否正常
		Formulaed bool   `json:"Formulaed,omitempty"` //是否有计算公式
		Variable  []*LuaVariableTemplate
	}

	lock.Lock()
	defer lock.Unlock()
	template := DeviceTSLMap.Get(d.Type)

	if template == nil {
		return fmt.Errorf("no such tsl template %s", d.Type)
	}
	if template.PluginTemplate == nil {
		return fmt.Errorf("tsl template %s haven't init plugin template", template.Name)
	}
	state := template.PluginTemplate.LuaState
	if state == nil {
		return errors.New("nil lua state")
	}
	tbl := lua.LTable{}
	for _, v := range rxBuf {
		tbl.Append(lua.LNumber(v))
	}
	state.SetGlobal(string(RXBUF), luar.New(state, &tbl))
	tbl_tx := lua.LTable{}
	for _, v := range txBuf {
		tbl_tx.Append(lua.LNumber(v))
	}
	state.SetGlobal(string(TXBUF), luar.New(state, &tbl_tx))
	//AnalysisRx
	err := state.CallByParam(lua.P{
		Fn:      state.GetGlobal(string(ANALYSISRX)),
		NRet:    1,
		Protect: true,
	}, lua.LString(sAddr), lua.LNumber(rxBufCnt))

	if err != nil {
		return err
	}

	//获取返回结果
	ret := state.Get(-1)
	state.Pop(1)

	LuaVariableMap := LuaVariableMapTemplate{}
	table, ok := ret.(*lua.LTable)
	if !ok {
		return errors.New("ret is not lua LTable")

	}
	if err = gluamapper.Map(table, &LuaVariableMap); err != nil {
		return fmt.Errorf("AnalysisRx gluamapper.Map error:%v", err)
	}

	//如果有公式的话
	if LuaVariableMap.Formulaed {
		if env == nil || eval == nil {
			initializeEval()
		}
	}

	timeNowStr := time.Now().Format("2006-01-02 15:04:05")
	value := model.DeviceTSLPropertyValueTemplate{}
	//正常
	if LuaVariableMap.Status != "0" {
		return fmt.Errorf("lua return  status  is not 0: %s", LuaVariableMap.Status)
	}

	if l := len(LuaVariableMap.Variable); l <= 0 {
		return fmt.Errorf("variable map is less than 0:%d", l)
	}
	var item float64
VLOOP:
	for _, lv := range LuaVariableMap.Variable {
		if lv == nil {
			mylog.ZAPS.Errorf("device %s variable is nil", d.Name)
			continue
		}
		if lv.Value == nil {
			mylog.ZAPS.Errorf("device %s variable %s value is nil", d.Name, lv.Label)
			continue
		}
		if len(properties) < lv.Index+1 {
			//mylog.ZAPS.Errorf("tsl template defined properties(%d) less than lua return(%d)", len(properties), len(LuaVariableMap.Variable))
			return nil
		}
		v := &properties[lv.Index]
		value.Index = lv.Index
		switch v.Type {
		case PropertyTypeInt32:
			item, ok = lv.Value.(float64)
			if ok {
				value.Value = (int32)(item)
			} else {
				mylog.ZAPS.Errorf("%v(%t) 不能转换为float64,将设置为NIL", lv.Value, reflect.TypeOf(lv.Value))
				value.Value = "NIL"
			}
		case PropertyTypeUInt32:
			item, ok = lv.Value.(float64)
			if ok {
				value.Value = (uint32)(item)
			} else {
				mylog.ZAPS.Errorf("%v(%t) 不能转换为float64,将设置为NIL", lv.Value, reflect.TypeOf(lv.Value))
				value.Value = "NIL"
			}
		case PropertyTypeDouble:
			item, ok = lv.Value.(float64)
			if ok {
				value.Value = item
			} else {
				mylog.ZAPS.Errorf("%v(%t) 不能转换为float64,将设置为NIL", lv.Value, reflect.TypeOf(lv.Value))
				value.Value = "NIL"
			}
		case PropertyTypeString:
			item, ok := lv.Value.(string)
			if ok {
				value.Value = item
			} else {
				mylog.ZAPS.Errorf("%v(%t) 不能转换为string,将设置为NIL", lv.Value, reflect.TypeOf(lv.Value))
				value.Value = "NIL"
			}

		default:
			mylog.ZAPS.Errorf("未识别的值类型:%s 将设置为NIL", v.Type)
			value.Value = "NIL"
		}

		//如果有表达式
		if lv.Formula != "" {
			//log.Printf("formula:%s\n", lv.Formula)
			if d.Parser == nil {
				d.Parser = &IndexParser{
					env: env,
				}
			}
			d.Parser.SetFormula(lv.Formula)
			if err := d.Parser.PreVarSet(properties); err != nil {
				mylog.ZAPS.Errorf("基础变量设置失败:%v", err)
				goto VLOOP
			}
			if err := d.Parser.VarSet(item); err != nil {
				mylog.ZAPS.Errorf("设置表达式val值错误:%v", err)
				goto VLOOP
			}

			parser := parser.New(lv.Formula)
			exp, err := parser.Parse()
			if err != nil {
				mylog.ZAPS.Errorf("解析表达式%s错误:%v", lv.Formula, err)
				goto VLOOP
			}
			ret, err := eval.Run(exp)
			if err != nil {
				mylog.ZAPS.Errorf("运行表达式%s错误:%v", lv.Formula, err)
				goto VLOOP
			}
			var endValue interface{}
			if d, ok := ret.(decimal.Decimal); ok {
				v, exact := d.Float64()
				if exact {
					endValue = v
				} else {
					endValue = d.String()
				}
			} else {
				mylog.ZAPS.Errorf("%s 不能转换为float64或者string", lv.Label)
			}
			value.Value = endValue

		}

		value.Explain = lv.Explain
		value.TimeStamp = timeNowStr

		if len(v.Value) < 100 {
			v.Value = append(v.Value, value)
		} else {
			v.Value = v.Value[1:]
			v.Value = append(v.Value, value)
		}
	}

	return nil
}
