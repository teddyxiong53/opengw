/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-10 09:28:15
@LastEditors: WalkMiao
@LastEditTime: 2021-10-08 08:39:40
@FilePath: /goAdapter-Raw/device/commInterface.go
*/
package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/leandro-lugaresi/hub"
)

type CommunicationInterface interface {
	Open() error
	io.ReadWriteCloser //和系统接口保持一致
	GetName() string
	GetType() string
	GetParam() interface{}
	GetTimeOut() string
	GetInterval() string
	Unique() string   //串口不能名字一致,网口不能ip和port一致
	Error() error     //打开或者关闭的时候是否产生错误，如果有错误就不执行read这些操作哦了
	Bind(name string) //和采集接口建立关系
	UnBind(name string)
	BindNames() []string //有哪些采集在用这个comm口
}

//通信接口Map
type CommunicationManager struct {
	m map[string]CommunicationInterface
	sync.RWMutex
	publisher *hub.Hub
	changed   bool
}

var CommunicationInterfaceMap = NewCommunicationManager()

func NewCommunicationManager() *CommunicationManager {
	mgr := &CommunicationManager{
		m:         make(map[string]CommunicationInterface),
		publisher: hub.New(),
	}
	return mgr
}

func (manager *CommunicationManager) Close() {
	manager.publisher.Close()
	//TODO 通信接口是不是还有其他需要关闭的资源？
}

func (manager *CommunicationManager) Changed() bool {
	var changed bool
	manager.RLock()
	changed = manager.changed
	manager.RUnlock()
	return changed
}

func (manager *CommunicationManager) Publish(name string, fields hub.Fields) {
	msg := hub.Message{
		Name:   name,
		Fields: fields,
	}
	manager.publisher.Publish(msg)
}

func (manager *CommunicationManager) SaveTo(f *os.File) error {
	manager.Lock()
	defer manager.Unlock()
	data, err := json.Marshal(manager.m)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	manager.changed = false
	return err
}

func (manager *CommunicationManager) Compare(comm CommunicationInterface) bool {
	manager.RLock()
	defer manager.RUnlock()
	for _, v := range manager.m {
		if v.Unique() == comm.Unique() {
			return false
		}
	}
	return true
}

func (manager *CommunicationManager) Add(comm CommunicationInterface) bool {
	manager.Lock()
	defer manager.Unlock()
	_, ok := manager.m[comm.GetName()]
	if ok {
		return false
	}
	manager.m[comm.GetName()] = comm
	manager.changed = true
	manager.Publish(CommAdd, hub.Fields{
		"Name": comm.GetName(),
	})
	return true
}

// 通信口的变动
func (manager *CommunicationManager) Update(comm CommunicationInterface) bool {
	manager.Lock()
	defer manager.Unlock()
	//之前的comm
	old, ok := manager.m[comm.GetName()]
	if !ok {
		return false
	}
	manager.m[comm.GetName()] = comm
	manager.Publish(CommUpdate, hub.Fields{
		"Name": comm.GetName(),
		"Old":  old,
		"New":  comm,
	})

	// for _, collectName := range old.BindNames() {
	// 	coll := CollectInterfaceMap.Get(collectName)
	// 	if coll != nil {
	// 		//和这个通信口有关的采集都要停掉
	// 		//delHandler(coll, true)
	// 		old.UnBind(collectName)
	// 		//替换新的comm口
	// 		if err := comm.Open(); err != nil {
	// 			mylog.ZAP.Error("新通讯口打开失败", zap.Error(err))
	// 		}
	// 		coll.CommInterface = comm
	// 		coll.CommInterfaceName = comm.GetName()
	// 		//新comm绑定这个采集
	// 		comm.Bind(collectName)

	// 		//如果之前打开就是错误的那么要重新add
	// 		if old.Error() != nil {
	// 			addHandler(coll, false)
	// 		}

	// 	}

	// }
	manager.changed = true
	return true
}

func (manager *CommunicationManager) Delete(commName string) bool {
	manager.Lock()
	defer manager.Unlock()
	if _, ok := manager.m[commName]; !ok {
		return false
	}
	delete(manager.m, commName)
	manager.Publish(CommDelete, hub.Fields{
		"Name": commName,
	})
	manager.changed = true
	return true
}

func (manager *CommunicationManager) Get(name string) CommunicationInterface {
	manager.RLock()
	v, ok := manager.m[name]
	if ok {
		manager.RUnlock()
		return v
	}
	manager.RUnlock()
	return nil
}

func (manager *CommunicationManager) GetAll() []CommunicationInterface {
	manager.RLock()
	comms := make([]CommunicationInterface, 0, len(manager.m))
	for _, v := range manager.m {
		comms = append(comms, v)
	}
	manager.RUnlock()
	return comms
}

func CommInterfaceInit(data []byte) error {
	var temp map[string]map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("comm json unmarshal error:%v", err)
	}
	CommunicationInterfaceMap.Lock()
	for k, i := range temp {

		typ, ok := i["Type"]
		if !ok {
			return errors.New("json file content is not include Type")
		}
		param, ok := i["Param"]
		if !ok {
			return errors.New("json file content is not include Param")
		}
		switch typ {
		case SERIALTYPE:
			sParam := CommunicationSerialTemplate{
				Name: i["Name"].(string),
			}
			sParam.Type = typ.(string)
			var serial SerialInterfaceParam
			data, _ := json.Marshal(param)
			if err := json.Unmarshal(data, &serial); err != nil {
				return err
			}
			sParam.Param = &serial
			bindData, ok := i["Bindings"]
			if !ok {
				return errors.New("json file content is not include Bindings")
			}
			var bindings = make([]string, 0)
			data, _ = json.Marshal(bindData)
			if err := json.Unmarshal(data, &bindings); err != nil {
				return err
			}
			sParam.Bindings = bindings
			CommunicationInterfaceMap.m[k] = &sParam
		//TODO 目前Binding只在串口类型上使用了 后续还可能扩展
		case TCPCLIENTTYPE:
			sParam := CommunicationTcpClientTemplate{
				Name: i["Name"].(string),
			}
			sParam.Type = typ.(string)
			var tcpClient TcpClientInterfaceParam
			data, _ := json.Marshal(param)
			if err := json.Unmarshal(data, &tcpClient); err != nil {
				return err
			}
			sParam.Param = &tcpClient
			CommunicationInterfaceMap.m[k] = &sParam
		case IOINTYPE:
			sParam := CommunicationIoInTemplate{
				Name: i["Name"].(string),
			}
			sParam.Type = typ.(string)
			var ioIn IoInInterfaceParam
			data, _ := json.Marshal(param)
			if err := json.Unmarshal(data, &ioIn); err != nil {
				return err
			}
			sParam.Param = &ioIn
			CommunicationInterfaceMap.m[k] = &sParam
		case IOOUTTYPE:
			sParam := CommunicationIoOutTemplate{
				Name: i["Name"].(string),
			}
			sParam.Type = typ.(string)
			var ioOut IoOutInterfaceParam
			data, _ := json.Marshal(param)
			if err := json.Unmarshal(data, &ioOut); err != nil {
				return err
			}
			sParam.Param = &ioOut
			CommunicationInterfaceMap.m[k] = &sParam
		default:
			return fmt.Errorf("unKnown type of json:%s", typ)
		}
	}
	CommunicationInterfaceMap.Unlock()
	return nil
}
