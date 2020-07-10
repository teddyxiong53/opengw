package device

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type CommunicationInterface interface{
	Open(interface{})  		bool
	Close() 				bool
	WriteData(data []byte) 	int
	ReadData(data []byte)  	int
}

type CommInterfaceTemplate struct{
	Name 	string										`json:"Name"`			//接口名称
	Type    string          							`json:"Type"`			//接口类型,比如serial,tcp,udp,http
	Param   interface{}     							`json:"Param"`			//接口参数
	Status  bool 										`json:"Status"`			//接口状态
	Port    CommunicationInterface						`json:"-"`
}

type CommInterfaceListTemplate struct{
	InterfaceCnt 	 int								`json:"InterfaceCnt"`
	InterfaceMap 	 []CommInterfaceTemplate
}

var CommInterfaceList  *CommInterfaceListTemplate


func NewCommInterfaceList() *CommInterfaceListTemplate{

	return &CommInterfaceListTemplate{
		InterfaceCnt: 0,
		InterfaceMap: make([]CommInterfaceTemplate,0),
	}
}

func NewCommInterface(Name string,Type string,Param interface{}) CommInterfaceTemplate{

	comm := CommInterfaceTemplate{
		Name:Name,
		Type:Type,
		Param:Param,
	}

	return comm
}

func (c *CommInterfaceTemplate)NewCommInterfacePort(){

	switch c.Type{
		case "serial":
			c.Port = &CommunicationSerialInterface{}
		case "tcp":
			c.Port = &CommunicationTcpInterface{}
	}
}

//func (c *CommInterfaceListTemplate)NewCommInterfacePort(){
//
//	for _,v := range c.InterfaceMap{
//		log.Printf("type is %s\n",v.Type)
//		switch v.Type{
//		case "serial":
//			c.InterfacePortMap = append(c.InterfacePortMap,&CommunicationSerialInterface{})
//		case "tcp":
//			c.InterfacePortMap = append(c.InterfacePortMap,&CommunicationTcpInterface{})
//		}
//	}
//}

func (c *CommInterfaceListTemplate)AddCommInterface(Name string,Type string,Param interface{}){

	comm := NewCommInterface(Name,Type,Param)

	c.InterfaceMap = append(c.InterfaceMap,comm)
	c.InterfaceCnt++
}

func (c *CommInterfaceListTemplate)ModifyCommInterface(Name string,Type string,Param interface{}) bool{

	for k,v := range c.InterfaceMap{
		if v.Name == Name{
			c.InterfaceMap[k].Name = Name
			c.InterfaceMap[k].Type = Type
			c.InterfaceMap[k].Param = Param
			return true
		}
	}
	return false
}

func (c *CommInterfaceListTemplate)GetCommInterface(Name string) CommInterfaceTemplate{

	for k,v := range c.InterfaceMap{
		if v.Name == Name{
			return c.InterfaceMap[k]
		}
	}
	return CommInterfaceTemplate{}
}

func (c *CommInterfaceListTemplate)DeleteCommInterface(Name string) bool{

	for k,v := range c.InterfaceMap{
		if v.Name == Name{
			c.InterfaceMap = append(c.InterfaceMap[:k],c.InterfaceMap[k+1:]...)
			return true
		}
	}
	return false
}

func WriteCommInterfaceListToJson() {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	fileDir := exeCurDir + "/selfpara/commInterface.json"

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("open commInterface.json err", err)
		return
	}
	defer fp.Close()

	sJson, _ := json.Marshal(CommInterfaceList)

	_, err = fp.Write(sJson)
	if err != nil {
		log.Println("write commInterface.json err", err)
	}
	log.Println("write commInterface.json sucess")
}

func ReadCommInterfaceListFromJson() bool {

	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileDir := exeCurDir + "/selfpara/commInterface.json"

	if fileExist(fileDir) == true {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			log.Println("open commInterface.json err", err)
			return false
		}
		defer fp.Close()

		data := make([]byte, 20480)
		dataCnt, err := fp.Read(data)

		err = json.Unmarshal(data[:dataCnt], CommInterfaceList)
		if err != nil {
			log.Println("commInterface unmarshal err", err)
			return false
		}

		return true
	} else {
		log.Println("commInterface.json is not exist")

		return false
	}
}

func CommInterfaceInit() {

	CommInterfaceList = NewCommInterfaceList()
	if ReadCommInterfaceListFromJson() == true{
		log.Println("read commInterface.json ok")
		log.Printf("%+v\n",CommInterfaceList)
		for _,v := range CommInterfaceList.InterfaceMap{
			v.NewCommInterfacePort()
			v.Status = v.Port.Open(v.Param)
			if v.Status != true{
				log.Printf("%s open err\n",v.Param)
			}
		}
	}
}