package luautils

import (
<<<<<<< HEAD:setting/mlua.go
<<<<<<< Updated upstream:setting/mlua.go
	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
=======
	"goAdapter/pkg/mylog"
>>>>>>> Stashed changes:pkg/luautils/mlua.go
=======
	"log"
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包):pkg/luautils/mlua.go
	"os"
	"path/filepath"
	"sync"

	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

type LuaFuncMap map[string]lua.LGFunction
type TemplateID string

var (
	CommonLuaFuncMap = LuaFuncMap{
		"GetCRCModbus":   GetCRCModbus,
		"CheckCRCModbus": CheckCRCModbus,
		"ParseByte":      ParseByte,
	}
)

type crc struct {
	once  sync.Once
	table []uint16
}

var crc16tab = [256]uint16{
	0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50a5, 0x60c6, 0x70e7,
	0x8108, 0x9129, 0xa14a, 0xb16b, 0xc18c, 0xd1ad, 0xe1ce, 0xf1ef,
	0x1231, 0x0210, 0x3273, 0x2252, 0x52b5, 0x4294, 0x72f7, 0x62d6,
	0x9339, 0x8318, 0xb37b, 0xa35a, 0xd3bd, 0xc39c, 0xf3ff, 0xe3de,
	0x2462, 0x3443, 0x0420, 0x1401, 0x64e6, 0x74c7, 0x44a4, 0x5485,
	0xa56a, 0xb54b, 0x8528, 0x9509, 0xe5ee, 0xf5cf, 0xc5ac, 0xd58d,
	0x3653, 0x2672, 0x1611, 0x0630, 0x76d7, 0x66f6, 0x5695, 0x46b4,
	0xb75b, 0xa77a, 0x9719, 0x8738, 0xf7df, 0xe7fe, 0xd79d, 0xc7bc,
	0x48c4, 0x58e5, 0x6886, 0x78a7, 0x0840, 0x1861, 0x2802, 0x3823,
	0xc9cc, 0xd9ed, 0xe98e, 0xf9af, 0x8948, 0x9969, 0xa90a, 0xb92b,
	0x5af5, 0x4ad4, 0x7ab7, 0x6a96, 0x1a71, 0x0a50, 0x3a33, 0x2a12,
	0xdbfd, 0xcbdc, 0xfbbf, 0xeb9e, 0x9b79, 0x8b58, 0xbb3b, 0xab1a,
	0x6ca6, 0x7c87, 0x4ce4, 0x5cc5, 0x2c22, 0x3c03, 0x0c60, 0x1c41,
	0xedae, 0xfd8f, 0xcdec, 0xddcd, 0xad2a, 0xbd0b, 0x8d68, 0x9d49,
	0x7e97, 0x6eb6, 0x5ed5, 0x4ef4, 0x3e13, 0x2e32, 0x1e51, 0x0e70,
	0xff9f, 0xefbe, 0xdfdd, 0xcffc, 0xbf1b, 0xaf3a, 0x9f59, 0x8f78,
	0x9188, 0x81a9, 0xb1ca, 0xa1eb, 0xd10c, 0xc12d, 0xf14e, 0xe16f,
	0x1080, 0x00a1, 0x30c2, 0x20e3, 0x5004, 0x4025, 0x7046, 0x6067,
	0x83b9, 0x9398, 0xa3fb, 0xb3da, 0xc33d, 0xd31c, 0xe37f, 0xf35e,
	0x02b1, 0x1290, 0x22f3, 0x32d2, 0x4235, 0x5214, 0x6277, 0x7256,
	0xb5ea, 0xa5cb, 0x95a8, 0x8589, 0xf56e, 0xe54f, 0xd52c, 0xc50d,
	0x34e2, 0x24c3, 0x14a0, 0x0481, 0x7466, 0x6447, 0x5424, 0x4405,
	0xa7db, 0xb7fa, 0x8799, 0x97b8, 0xe75f, 0xf77e, 0xc71d, 0xd73c,
	0x26d3, 0x36f2, 0x0691, 0x16b0, 0x6657, 0x7676, 0x4615, 0x5634,
	0xd94c, 0xc96d, 0xf90e, 0xe92f, 0x99c8, 0x89e9, 0xb98a, 0xa9ab,
	0x5844, 0x4865, 0x7806, 0x6827, 0x18c0, 0x08e1, 0x3882, 0x28a3,
	0xcb7d, 0xdb5c, 0xeb3f, 0xfb1e, 0x8bf9, 0x9bd8, 0xabbb, 0xbb9a,
	0x4a75, 0x5a54, 0x6a37, 0x7a16, 0x0af1, 0x1ad0, 0x2ab3, 0x3a92,
	0xfd2e, 0xed0f, 0xdd6c, 0xcd4d, 0xbdaa, 0xad8b, 0x9de8, 0x8dc9,
	0x7c26, 0x6c07, 0x5c64, 0x4c45, 0x3ca2, 0x2c83, 0x1ce0, 0x0cc1,
	0xef1f, 0xff3e, 0xcf5d, 0xdf7c, 0xaf9b, 0xbfba, 0x8fd9, 0x9ff8,
	0x6e17, 0x7e36, 0x4e55, 0x5e74, 0x2e93, 0x3eb2, 0x0ed1, 0x1ef0,
}

var crcTb crc

// initTable 初始化表
func (c *crc) initTable() {
	crcPoly16 := uint16(0xa001)
	c.table = make([]uint16, 256)

	for i := uint16(0); i < 256; i++ {
		crc := uint16(0)
		b := i

		for j := uint16(0); j < 8; j++ {
			if ((crc ^ b) & 0x0001) > 0 {
				crc = (crc >> 1) ^ crcPoly16
			} else {
				crc = crc >> 1
			}
			b = b >> 1
		}
		c.table[i] = crc
	}
}

func crc16(bs []byte) uint16 {
	crcTb.once.Do(crcTb.initTable)

	val := uint16(0xFFFF)
	for _, v := range bs {
		val = (val >> 8) ^ crcTb.table[(val^uint16(v))&0x00FF]
	}
	return val
}

func LuaCallNewVariables(L *lua.LState) {
	//调用NewVariables
	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("NewVariables"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		panic(err)
	}
	//获取返回结果
	ret := L.Get(-1)
	L.Pop(1)
	switch ret.(type) {
	case lua.LString:
<<<<<<< HEAD:setting/mlua.go
<<<<<<< Updated upstream:setting/mlua.go
		Logger.Info("string")
	case *lua.LTable:
		Logger.Info("table")
=======
		mylog.Logger.Info("string")
	case *lua.LTable:
		mylog.Logger.Info("table")
>>>>>>> Stashed changes:pkg/luautils/mlua.go
=======
		log.Logger.Info("string")
	case *lua.LTable:
		log.Logger.Info("table")
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包):pkg/luautils/mlua.go
	}

	type VariableTemplate struct {
		Index int
		Name  string
		Label string
		Type  string
	}

	type VariableMapTemplate struct {
		Variable []*VariableTemplate
	}

	VariableMap := VariableMapTemplate{}

	if err := gluamapper.Map(ret.(*lua.LTable), &VariableMap); err != nil {
		panic(err)
	}

	for _, v := range VariableMap.Variable {
<<<<<<< HEAD:setting/mlua.go
<<<<<<< Updated upstream:setting/mlua.go
		Logger.Infof("%+v", v.Label)
=======
		mylog.Logger.Infof("%+v", v.Label)
>>>>>>> Stashed changes:pkg/luautils/mlua.go
=======
		log.Logger.Infof("%+v", v.Label)
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包):pkg/luautils/mlua.go
	}
}

func LuaInit() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	path := exeCurDir + "/plugin/"

	L := lua.NewState()
	defer L.Close()
	//加载Lua
	if err := L.DoFile(path + "td200.lua"); err != nil {
<<<<<<< HEAD:setting/mlua.go
<<<<<<< Updated upstream:setting/mlua.go
		Logger.Warning("open td200.lua fail", err)
	}
	Logger.Info("open TD200.lua OK")
=======
		mylog.Logger.Warning("open td200.lua fail", err)
	}
	mylog.Logger.Info("open TD200.lua OK")
>>>>>>> Stashed changes:pkg/luautils/mlua.go
=======
		log.Logger.Warning("open td200.lua fail", err)
	}
	log.Logger.Info("open TD200.lua OK")
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包):pkg/luautils/mlua.go

	LuaCallNewVariables(L)
}

func GetCRCModbus(L *lua.LState) int {

	type LuaVariableMapTemplate struct {
		Variable []*byte
	}

	lv := L.ToTable(1)

	LuaVariableMap := LuaVariableMapTemplate{}
	if err := gluamapper.Map(lv, &LuaVariableMap); err != nil {
<<<<<<< HEAD:setting/mlua.go
<<<<<<< Updated upstream:setting/mlua.go
		Logger.Warning("GetCRC16 gluamapper.Map err,", err)
=======
		mylog.Logger.Warning("GetCRC16 gluamapper.Map err,", err)
>>>>>>> Stashed changes:pkg/luautils/mlua.go
=======
		log.Logger.Warning("GetCRC16 gluamapper.Map err,", err)
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包):pkg/luautils/mlua.go
	}

	nBytes := make([]byte, 0)
	if len(LuaVariableMap.Variable) > 0 {
		for _, v := range LuaVariableMap.Variable {
			nBytes = append(nBytes, *v)
		}
	}

	//log.Printf("crcBytes %x\n", nBytes)
	//lenCRC := len(nBytes)
	crc := crc16(nBytes)
	//log.Printf("crcValue %v\n", crc)
	L.Push(lua.LNumber(crc)) /* push result */

	return 1 /* number of results */
}

func CheckCRCModbus(L *lua.LState) int {

	type LuaVariableMapTemplate struct {
		Variable []*byte
	}

	lv := L.ToTable(1)

	LuaVariableMap := LuaVariableMapTemplate{}
	if err := gluamapper.Map(lv, &LuaVariableMap); err != nil {
<<<<<<< HEAD:setting/mlua.go
<<<<<<< Updated upstream:setting/mlua.go
		Logger.Warning("GetCRC16 gluamapper.Map err,", err)
=======
		mylog.Logger.Warning("GetCRC16 gluamapper.Map err,", err)
>>>>>>> Stashed changes:pkg/luautils/mlua.go
=======
		log.Logger.Warning("GetCRC16 gluamapper.Map err,", err)
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包):pkg/luautils/mlua.go
	}

	nBytes := make([]byte, 0)
	if len(LuaVariableMap.Variable) > 0 {
		for _, v := range LuaVariableMap.Variable {
			nBytes = append(nBytes, *v)
		}
	}

	//log.Printf("crcBytes %x\n", nBytes)
	crc := crc16(nBytes)
	//log.Printf("crcValue %v\n", crc)
	L.Push(lua.LNumber(crc)) /* push result */

	return 1 /* number of results */
}

func LuaOpenFile(filePath string) (*lua.LState, error) {

	lState := lua.NewState()
	//defer L.Close()

	//加载Lua
	err := lState.DoFile(filePath)

	return lState, err
}

func ParseByte(L *lua.LState) int {
	byteTemp := struct {
		Value *byte
	}{}
	lv := L.ToTable(1)
	if err := gluamapper.Map(lv, &byteTemp); err != nil {
<<<<<<< HEAD:setting/mlua.go
<<<<<<< Updated upstream:setting/mlua.go
		Logger.Error("parseByte gluamapper Map error:%v", err)
=======
		mylog.Logger.Error("parseByte gluamapper Map error:%v", err)
>>>>>>> Stashed changes:pkg/luautils/mlua.go
=======
		log.Logger.Error("parseByte gluamapper Map error:%v", err)
>>>>>>> 137f07b (新增 1.添加viper配置框架以及zap日志框架 修改 1.将原setting包改为pkg包,并将不同功能的模块函数各自独立为一个单独包):pkg/luautils/mlua.go
	}

	var t = L.NewTable()
	v := byteTemp.Value
	for i := 0; i < 8; i++ {
		n := lua.LNumber((*v) >> i & 0x1)
		t.Append(n)
	}
	L.Push(t)
	return 1
}
