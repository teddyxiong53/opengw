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

func (d *DeviceNodeModbusTemplate)GetDeviceRealVariables(deviceAddr string) int{

	log.Printf("ModbusTemplate %s\n",d.Type)
	//log.Printf("varibales %+v\n",d.VariableMap)

	return 0
}

func (d *DeviceNodeModbusTemplate)SetDeviceRealVariables(deviceAddr string) int{


	return 0
}

func (d *DeviceNodeModbusTemplate)AddVariable(vindex int,vname string,vlable string,vtype string){

	variable := VariableTemplate{}
	variable.Index = vindex
	variable.Name = vname
	variable.Lable = vlable
	variable.Type = vtype

	d.VariableMap = append(d.VariableMap,variable)
}

func (d *DeviceNodeModbusTemplate)NewVariables() {

	d.VariableMap = make([]VariableTemplate,0)

	d.AddVariable(0,"Addr","通信地址","string")
	d.AddVariable(1,"DeviceType","设备类型","string")
	d.AddVariable(2,"SoftVer","软件版本","string")
	d.AddVariable(3,"SerialNumber","设备编码","string")
	d.AddVariable(4,"RTC","设备时钟","string")
	d.AddVariable(5,"RoomTemp","室内温度","string")
	d.AddVariable(6,"RoomHumi","室内湿度","string")
	d.AddVariable(7,"AirStatus","风机状态","string")
	d.AddVariable(8,"RelayStatus","阀门状态","string")
	d.AddVariable(9,"CurModeSet","运行模式","string")
	d.AddVariable(10,"SeaconStatus","季节","string")
	d.AddVariable(11,"CurTempStep","设定温度","byte")
	d.AddVariable(12,"LockStatus","锁定状态","string")
	d.AddVariable(13,"EnergySavingStatus","节能状态","string")
	d.AddVariable(14,"RSSI","信号强度","byte")
	d.AddVariable(15,"BLKStatus","背光状态","string")
	d.AddVariable(16,"ErrorCode","故障代码","uint16")
	d.AddVariable(17,"TotalTime","累计时间","uint32")
	d.AddVariable(18,"LowTotalTime","低速累计时间","uint32")
	d.AddVariable(19,"MiddleTotalTime","中速累计时间","uint32")
	d.AddVariable(20,"HighTotalTime","高速累计时间","uint32")
}

