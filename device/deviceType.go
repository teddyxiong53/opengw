/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-09-13 15:53:35
@FilePath: /goAdapter-Raw/device/deviceType.go
*/
package device

import (
	lua "github.com/yuin/gopher-lua"
)

type DeviceTemplate struct {
	TemplateID int         `json:"TemplateID"`
	Name       string      `json:"TemplateName"`    //模板名称
	Type       string      `json:"TemplateType"`    //模板型号
	Message    string      `json:"TemplateMessage"` //备注信息
	LuaState   *lua.LState `json:"-"`
}

//type DeviceNodeTypeLuaState struct {
//	LuaState *lua.LState
//	TypeName string
//	CollName string
//}

var DeviceTemplateMap = make(map[string]*DeviceTemplate)
