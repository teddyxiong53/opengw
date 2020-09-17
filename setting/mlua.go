package setting

import (
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"os"
	"path/filepath"
)

func LuaCallNewVariables(L *lua.LState){
	//调用NewVariables
	err := L.CallByParam(lua.P{
		Fn:L.GetGlobal("NewVariables"),
		NRet:1,
		Protect: true,
	})
	if err != nil{
		panic(err)
	}
	//获取返回结果
	ret := L.Get(-1)
	L.Pop(1)
	switch ret.(type) {
	case lua.LString:
		Logger.Info("string")
	case *lua.LTable:
		Logger.Info("table")
	}

	type VariableTemplate struct{
		Index int
		Name string
		Label string
		Type string
	}

	type VariableMapTemplate struct{
		Variable []*VariableTemplate
	}

	VariableMap :=VariableMapTemplate{}

	if err := gluamapper.Map(ret.(*lua.LTable), &VariableMap); err != nil {
		panic(err)
	}

	for _,v := range VariableMap.Variable{
		Logger.Infof("%+v",v.Label)
	}
}

func LuaInit() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	path := exeCurDir + "/plugin/"

	L := lua.NewState()
	defer L.Close()
	//加载Lua
	if err := L.DoFile(path+"td200.lua"); err != nil {
		Logger.Warning("open td200.lua fail",err)
	}
	Logger.Info("open TD200.lua OK")

	LuaCallNewVariables(L)
}

func LuaOpenFile(filePath string) (*lua.LState,error){

	lState := lua.NewState()
	//defer L.Close()

	//加载Lua
	err := lState.DoFile(filePath)

	return lState,err
}

