package device

import (
	"encoding/json"
	"goAdapter/setting"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

//最大设备模板数
var MaxDeviceNodeTypeCnt int = 10

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

var DeviceNodeTypeMap = DeviceNodeTypeMapStruct{
	DeviceNodeType: make([]DeviceNodeTypeTemplate, 0),
}

//var DeviceTypePluginMap = make(map[int]*plugin.Plugin)
var DeviceTypePluginMap = make(map[int]*lua.LState)

//var DeviceTypePluginMap = make([]DeviceNodeTypeLuaState,0)

func init() {

}

func updataDeviceType(path string, fileName []string) ([]string, error) {

	rd, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("readDir err,", err)
		return fileName, err
	}

	for _, fi := range rd {
		setting.Logger.Debugf("fi %v", fi.Name())
		if fi.IsDir() {
			fullDir := path + "/" + fi.Name()
			fileName, _ = updataDeviceType(fullDir, fileName)
		} else {
			fullName := path + "/" + fi.Name()
			if strings.Contains(fi.Name(), ".json") {
				//log.Println("fullName ",fullName)
				fileName = append(fileName, fullName)
			} else if strings.Contains(fi.Name(), ".lua") {
				//log.Println("fullName ",fullName)
				fileName = append(fileName, fullName)
			}
		}
	}

	return fileName, nil
}

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
			setting.Logger.Debugf("fileDirInfo %v", v.Name())
			fileDirName := pluginPath + "/" + v.Name()
			fileMap, err := ioutil.ReadDir(fileDirName)
			if err != nil {
				log.Println("readDir err,", err)
				return false
			}
			for _, f := range fileMap {
				setting.Logger.Debugf("fileName %v", f.Name())
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
					setting.Logger.Debugf("DeviceNodeType  %v", DeviceNodeTypeMap.DeviceNodeType)
				} else if strings.Contains(f.Name(), ".lua") {
					if strings.Contains(f.Name(), v.Name()) == true { //lua文件和设备模版名字一样
						template, err := setting.LuaOpenFile(fileFullName)
						if err != nil {
							setting.Logger.Errorf("openPlug %s err,%s\n", v.Name(), err)
						} else {
							setting.Logger.Debugf("openPlug  %s ok\n", f.Name())
						}
						for k, d := range DeviceNodeTypeMap.DeviceNodeType {
							if d.TemplateType == v.Name() {
								DeviceTypePluginMap[k] = template
								DeviceTypePluginMap[k].SetGlobal("GetCRCModbus", DeviceTypePluginMap[k].NewFunction(setting.GetCRCModbus))
								DeviceTypePluginMap[k].SetGlobal("CheckCRCModbus", DeviceTypePluginMap[k].NewFunction(setting.CheckCRCModbus))
							}
						}
					} else {
						for k, d := range DeviceNodeTypeMap.DeviceNodeType {
							if d.TemplateType == v.Name() {
								err = DeviceTypePluginMap[k].DoFile(fileFullName)
								if err != nil {
									setting.Logger.Errorf("%s Lua DoFile err, %v", fileFullName, err)
									return false
								} else {
									setting.Logger.Debugf("%s Lua DoFile ok", f.Name())
								}
							}
						}
					}
				}
			}
		}
	}

	return true
}

//func ReadDeviceNodeTypeMap() bool {
//
//	deviceTypeTemplate := struct {
//		TemplateName    string `json:"TemplateName"`    //模板名称
//		TemplateType    string `json:"TemplateType"`    //模板型号
//		TemplateMessage string `json:"TemplateMessage"` //备注信息
//	}{}
//
//	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
//
//	//遍历json和so文件
//	pluginPath := exeCurDir + "/plugin"
//	fileNameMap := make([]string, 0)
//	fileNameMap, _ = updataDeviceType(pluginPath, fileNameMap)
//	for _, v := range fileNameMap {
//		if strings.Contains(v, ".json") {
//			fp, err := os.OpenFile(v, os.O_RDONLY, 0777)
//			if err != nil {
//				setting.Logger.Errorf("open %s err", v)
//				return false
//			}
//			defer fp.Close()
//
//			data := make([]byte, 2048)
//			dataCnt, err := fp.Read(data)
//
//			err = json.Unmarshal(data[:dataCnt], &deviceTypeTemplate)
//			if err != nil {
//				log.Println("deviceTypeTemplate unmarshal err", err)
//				return false
//			}
//
//			nodeType := DeviceNodeTypeTemplate{}
//			nodeType.TemplateID = len(DeviceNodeTypeMap.DeviceNodeType)
//			nodeType.TemplateType = deviceTypeTemplate.TemplateType
//			nodeType.TemplateName = deviceTypeTemplate.TemplateName
//			nodeType.TemplateMessage = deviceTypeTemplate.TemplateMessage
//
//			DeviceNodeTypeMap.DeviceNodeType = append(DeviceNodeTypeMap.DeviceNodeType, nodeType)
//		}
//	}
//
//	//打开lua文件
//	for k, v := range DeviceNodeTypeMap.DeviceNodeType {
//		for _, fileName := range fileNameMap {
//			if strings.Contains(fileName, ".lua") {
//				if strings.Contains(fileName, v.TemplateType) {
//					template, err := setting.LuaOpenFile(fileName)
//					if err != nil {
//						setting.Logger.Errorf("openPlug %s err,%s\n", fileName, err)
//					} else {
//						setting.Logger.Debugf("openPlug  %s ok\n", fileName)
//					}
//					DeviceTypePluginMap[k] = template
//					DeviceTypePluginMap[k].SetGlobal("GetCRCModbus", DeviceTypePluginMap[k].NewFunction(setting.GetCRCModbus))
//					DeviceTypePluginMap[k].SetGlobal("CheckCRCModbus", DeviceTypePluginMap[k].NewFunction(setting.CheckCRCModbus))
//				}
//			}
//		}
//	}
//
//	return true
//}

//func UpdateDeviceNodeType(collName string) {
//
//	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
//
//	//遍历json和so文件
//	pluginPath := exeCurDir + "/plugin"
//	fileNameMap := make([]string,0)
//	fileNameMap,_ = updataDeviceType(pluginPath,fileNameMap)
//
//	var filenameWithSuffix string
//	var fileSuffix string
//	for _,fileName := range fileNameMap{
//		if strings.Contains(fileName,".lua") {
//			typeLuaState := DeviceNodeTypeLuaState{}
//			typeLuaState.CollName = collName
//			filenameWithSuffix = path.Base(fileName)
//			fileSuffix = path.Ext(filenameWithSuffix)
//			typeLuaState.TypeName = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
//			typeLuaState.LuaState = lua.NewState()
//
//			typeLuaState.LuaState.SetGlobal("GetCRCModbus", typeLuaState.LuaState.NewFunction(setting.GetCRCModbus))
//			typeLuaState.LuaState.SetGlobal("CheckCRCModbus", typeLuaState.LuaState.NewFunction(setting.CheckCRCModbus))
//			err := typeLuaState.LuaState.DoFile(fileName)
//			if err != nil {
//				log.Printf("openLua %s err,%s\n",fileName,err)
//			}else{
//				log.Printf("openLua  %s ok\n", fileName)
//				DeviceTypePluginMap = append(DeviceTypePluginMap,typeLuaState)
//			}
//		}
//	}
//}
