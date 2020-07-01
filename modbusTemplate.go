package main

import (
	"log"
)

type DeviceNodeModbusTemplate struct{
	DeviceNodeTemplate
	TemplateName string					//模板名称
	TemplateType string					//模板型号
	TemplateID   string					//模板ID
	TemplateMessage string              //备注信息
}

func (d DeviceNodeModbusTemplate)GetDeviceRealVariables(deviceAddr string) int{

	log.Printf("ModbusTemplate %s\n",d.Type)

	return 0
}

type DeviceNodeModbus2Template struct{
	DeviceNodeTemplate
	TemplateName string					//模板名称
	TemplateType string					//模板型号
	TemplateID   string					//模板ID
	TemplateMessage string              //备注信息
}

func (d DeviceNodeModbus2Template)SetDeviceRealVariables(deviceAddr string) int{



	return 0
}



func (d DeviceNodeModbusTemplate)SetDeviceRealVariables(deviceAddr string) int{



	return 0
}

func (d DeviceNodeModbus2Template)GetDeviceRealVariables(deviceAddr string) int{

	log.Printf("Modbus2Template %s\n",d.Type)

	return 0
}

