package device

import (
	"encoding/json"
	"goAdapter/setting"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/gluamapper"

	lua "github.com/yuin/gopher-lua"
)

type DeviceNodeTypeTemplate struct {
	TemplateID      int    `json:"TemplateID"`      //模板ID
	TemplateName    string `json:"TemplateName"`    //模板名称
	TemplateType    string `json:"TemplateType"`    //模板型号
	TemplateMessage string `json:"TemplateMessage"` //备注信息
}

//配置参数
type DeviceNodeTypeMapStruct struct {
	DeviceNodeType []DeviceNodeTypeTemplate
}

type DeviceNodeTypeLuaState struct {
	LuaState *lua.LState
	TypeName string
	CollName string
}

type DeviceNodeTypeVariableTemplate struct {
	Index int
	Name  string
	Label string
	Type  string
}

type DeviceNodeTypeVariableMapTemplate struct {
	TemplateType string
	Variable     []DeviceNodeTypeVariableTemplate
}

var DeviceNodeTypeMap = DeviceNodeTypeMapStruct{
	DeviceNodeType: make([]DeviceNodeTypeTemplate, 0),
}
var DeviceTypePluginMap = make(map[int]*lua.LState)
var DeviceNodeTypeVariableMap = make([]DeviceNodeTypeVariableMapTemplate, 0)

func ReadDeviceNodeTypeMap() bool {

	deviceTypeTemplate := struct {
		TemplateName    string `json:"TemplateName"`    //模板名称
		TemplateType    string `json:"TemplateType"`    //模板型号
		TemplateMessage string `json:"TemplateMessage"` //备注信息
	}{}

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	//遍历json和so文件
	pluginPath := exeCurDir + "/plugin"
	fileInfoMap, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		log.Println("readDir err,", err)
		return false
	}
	for _, v := range fileInfoMap {
		//文件夹
		if v.IsDir() == true {
			//setting.Logger.Debugf("fileDirInfo %v", v.Name())
			fileDirName := pluginPath + "/" + v.Name()
			fileMap, err := ioutil.ReadDir(fileDirName)
			if err != nil {
				log.Println("readDir err,", err)
				return false
			}
			//遍历json文件，查找设备模版
			for _, f := range fileMap {
				fileFullName := fileDirName + "/" + f.Name()
				if strings.Contains(f.Name(), ".json") {
					fp, err := os.OpenFile(fileFullName, os.O_RDONLY, 0777)
					if err != nil {
						setting.Logger.Errorf("open %s err", f)
						return false
					}
					defer fp.Close()

					data := make([]byte, 2048)
					dataCnt, err := fp.Read(data)

					err = json.Unmarshal(data[:dataCnt], &deviceTypeTemplate)
					if err != nil {
						log.Println("deviceTypeTemplate unmarshal err", err)
						return false
					}

					nodeType := DeviceNodeTypeTemplate{}
					nodeType.TemplateID = len(DeviceNodeTypeMap.DeviceNodeType)
					nodeType.TemplateType = deviceTypeTemplate.TemplateType
					nodeType.TemplateName = deviceTypeTemplate.TemplateName
					nodeType.TemplateMessage = deviceTypeTemplate.TemplateMessage
					DeviceNodeTypeMap.DeviceNodeType = append(DeviceNodeTypeMap.DeviceNodeType, nodeType)

				}
			}
			setting.Logger.Debugf("DeviceNodeType %v", DeviceNodeTypeMap.DeviceNodeType)

			index := -1
			for _, f := range fileMap {
				//setting.Logger.Debugf("fileName %v", f.Name())
				fileFullName := fileDirName + "/" + f.Name()
				if strings.Contains(f.Name(), ".lua") {
					if strings.Contains(f.Name(), v.Name()) == true { //lua文件和设备模版名字一样
						template, err := setting.LuaOpenFile(fileFullName)
						if err != nil {
							setting.Logger.Errorf("openPlug %s err %v", v.Name(), err)
						} else {
							setting.Logger.Debugf("openPlug  %s ok", f.Name())
						}
						for k, d := range DeviceNodeTypeMap.DeviceNodeType {
							if d.TemplateType == v.Name() {
								index = k
								DeviceTypePluginMap[k] = template
								DeviceTypePluginMap[k].SetGlobal("GetCRCModbus", DeviceTypePluginMap[k].NewFunction(setting.GetCRCModbus))
								DeviceTypePluginMap[k].SetGlobal("CheckCRCModbus", DeviceTypePluginMap[k].NewFunction(setting.CheckCRCModbus))
							}
						}
						break
					}
				}
			}
			if index == -1 {
				continue
			}

			for _, f := range fileMap {
				fileFullName := fileDirName + "/" + f.Name()
				if strings.Contains(f.Name(), ".lua") {
					if strings.Contains(f.Name(), v.Name()) == false { //lua文件和设备模版名字不一样
						err = DeviceTypePluginMap[index].DoFile(fileFullName)
						if err != nil {
							setting.Logger.Errorf("%s Lua DoFile err, %v", fileFullName, err)
							return false
						} else {
							setting.Logger.Debugf("%s Lua DoFile ok", f.Name())
						}
					}
				}
			}

			for _, f := range fileMap {
				if strings.Contains(f.Name(), ".lua") {
					if strings.Contains(f.Name(), v.Name()) == false { //lua文件和设备模版名字不一样
						//获取设备模板中的变量
						ReadDeviceNodeTypeVariableMap(v.Name(), DeviceTypePluginMap[index])
					}
				}
			}
		}
	}

	return true
}

func ReadDeviceNodeTypeVariableMap(templateType string, l *lua.LState) {

	type LuaVariableTemplate struct {
		Index int
		Name  string
		Label string
		Type  string
	}

	type LuaVariableMapTemplate struct {
		Variable []*LuaVariableTemplate
	}

	//调用NewVariables
	err := l.CallByParam(lua.P{
		Fn:      l.GetGlobal("NewVariables"),
		NRet:    1,
		Protect: true,
	})
	if err != nil {
		setting.Logger.Warning("NewVariables err,", err)
	}

	//获取返回结果
	ret := l.Get(-1)
	l.Pop(1)

	LuaVariableMap := DeviceNodeTypeVariableMapTemplate{}

	if err := gluamapper.Map(ret.(*lua.LTable), &LuaVariableMap); err != nil {
		setting.Logger.Warning("NewVariables gluamapper.Map err,", err)
	}

	TypeVariableMap := DeviceNodeTypeVariableMapTemplate{}
	TypeVariableMap.TemplateType = templateType
	for _, v := range LuaVariableMap.Variable {
		variable := DeviceNodeTypeVariableTemplate{}
		variable.Index = v.Index
		variable.Name = v.Name
		variable.Label = v.Label
		variable.Type = v.Type
		TypeVariableMap.Variable = append(TypeVariableMap.Variable, variable)
	}
	DeviceNodeTypeVariableMap = append(DeviceNodeTypeVariableMap, TypeVariableMap)
}
