package serial

type SerialPortNameTemplate struct {
	Name []string `json:"Name"`
}

//var SerialPortNameTemplateMap = [...]SerialPortNameTemplate{
//	{Name:"/dev/ttyUSB0"},
//	{Name:"/dev/ttyUSB1"},
//	{Name:"/dev/ttyS0"},
//	{Name:"/dev/ttyS1"},
//	{Name:"/dev/tty.SLAB_USBtoUART"},
//}

var SerialPortNameTemplateMap = SerialPortNameTemplate{}
