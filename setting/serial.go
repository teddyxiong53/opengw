package setting

type SerialPortNameTemplate struct{
	Name string			`json:"Name"`
}

var SerialPortNameTemplateMap = [...]SerialPortNameTemplate{
	{Name:"/dev/ttyUSB0"},
	{Name:"/dev/ttyUSB1"},
}
