/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-30 11:43:43
@LastEditors: WalkMiao
@LastEditTime: 2021-10-07 22:07:09
@FilePath: /goAdapter-Raw/httpServer/model/deviceTSL.go
*/
package model

import (
	lua "github.com/yuin/gopher-lua"
)

type DeviceTSLPropertyParamTempate struct {
	Min             string `json:"Min"`                       //最小
	Max             string `json:"Max"`                       //最大
	MinMaxAlarm     bool   `json:"MinMaxAlarm"`               //范围报警
	Step            string `json:"Step"`                      //步长
	StepAlarm       bool   `json:"StepAlarm"`                 //阶跃报警
	Decimals        string `json:"Decimals,omitempty"`        //小数位数
	DataLength      string `json:"DataLength,omitempty"`      //字符串长度
	DataLengthAlarm bool   `json:"DataLengthAlarm,omitempty"` //字符长度报警
	Unit            string `json:"Unit"`                      //单位
}

type DeviceTSLPropertyValueTemplate struct {
	Index     int         `json:"Index"`
	Value     interface{} `json:"Value"`   //变量值，不可以是字符串
	Explain   interface{} `json:"Explain"` //变量值解释，必须是字符串
	TimeStamp string      `json:"TimeStamp"`
}

type DeviceTSLPropertyTemplate struct {
	Name       string                           `json:"Name"`       //属性名称，只可以是字母和数字的组合
	Explain    string                           `json:"Explain"`    //属性解释
	AccessMode int                              `json:"AccessMode"` //读写属性
	Type       int                              `json:"Type"`       //类型 uint32 int32 double string
	Params     DeviceTSLPropertyParamTempate    `json:"Params"`
	Value      []DeviceTSLPropertyValueTemplate `json:"-"`
}

type DeviceTSLServiceTempalte struct {
	Name     string                 `json:"Name"`     //服务名称
	Explain  string                 `json:"Explain"`  //服务名称说明
	CallType int                    `json:"CallType"` //服务调用方式
	Params   map[string]interface{} `json:"Params"`   //服务参数
}

//物模型 Thing Specification Language
type DeviceTSLTemplate struct {
	Name           string                      `json:"TSLName"`              //名称，只可以是字母和数字的组合
	Explain        string                      `json:"TSLExplain,omitempty"` //名称解释
	Plugin         string                      `json:"Plugin,omitempty"`
	Properties     []DeviceTSLPropertyTemplate `json:"Properties,omitempty"` //属性
	Services       []DeviceTSLServiceTempalte  `json:"Services,omitempty"`   //服务
	PluginTemplate *PluginTemplate             `json:"-"`
}

// plugin模板
type PluginTemplate struct {
	TemplateID int         `json:"TemplateID"`
	Name       string      `json:"TemplateName"`    //模板名称
	Type       string      `json:"TemplateType"`    //模板型号
	Message    string      `json:"TemplateMessage"` //备注信息
	LuaState   *lua.LState `json:"-"`
}
