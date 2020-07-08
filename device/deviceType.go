package device

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"plugin"
)

//最大设备模板数
var MaxDeviceNodeTypeCnt int = 10

type DeviceNodeTypeTemplate struct {
	TemplateID      int    `json:"templateID"`      //模板ID
	TemplateName    string `json:"templateName"`    //模板名称
	TemplateType    string `json:"templateType"`    //模板型号
	TemplateMessage string `json:"templateMessage"` //备注信息
}

//配置参数
type DeviceNodeTypeMapStruct struct {
	DeviceNodeType    []DeviceNodeTypeTemplate
}

var DeviceNodeTypeMap   DeviceNodeTypeMapStruct
var DeviceTypePluginMap map[int]*plugin.Plugin

func WriteDeviceNodeTypeMapToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/deviceNodeType.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open DeviceNodeType.json err", err)
		return
	}
	defer fp.Close()

	sJson, _ := json.Marshal(DeviceNodeTypeMap)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write DeviceNodeType.json err", err)
	}
	log.Println("write DeviceNodeType.json sucess")
}

func ReadDeviceNodeTypeMapFromJson() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/deviceNodeType.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open DeviceNodeType.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		DeviceNodeTypeMap.DeviceNodeType = make([]DeviceNodeTypeTemplate, 0)

		err = json.Unmarshal(data[:dataCnt], &DeviceNodeTypeMap)
		if err != nil {
			log.Println("DeviceNodeType unmarshal err", err)

			return false
		}
		//创建设备模版
		DeviceTypePluginMap = make(map[int]*plugin.Plugin)
		//log.Printf("plugin %+v\n",DeviceNodeTypeMap)
		for k,v := range DeviceNodeTypeMap.DeviceNodeType{

			str := "plugin/" + v.TemplateType + ".so"
			log.Printf("pluginStr %s\n",str)
			log.Printf("plugin %+v\n",DeviceNodeTypeMap)
			template,pluginerr := plugin.Open(str)
			if pluginerr!=nil{
				log.Printf("open %s, %\n",pluginerr)
			}
			DeviceTypePluginMap[k] = template
		}
		//log.Printf("plugin %+v\n",DeviceTypePluginMap)

		return true
	} else {
		log.Println("DeviceNodeType.json is not exist")
		//创建设备模版
		DeviceTypePluginMap = make(map[int]*plugin.Plugin)

		return false
	}
}

