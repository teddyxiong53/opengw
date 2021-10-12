/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-29 07:46:59
@LastEditors: WalkMiao
@LastEditTime: 2021-10-04 16:20:39
@FilePath: /goAdapter-Raw/device/deviceTSL.go
*/
package device

import (
	"encoding/json"
	"fmt"
	"goAdapter/httpServer/model"
	"goAdapter/pkg/backup"
	"os"
	"path"
	"sync"

	"github.com/leandro-lugaresi/hub"
)

const (
	TSLAccessModeRead int = iota
	TSLAccessModeWrite
	TSLAccessModeReadWrite
)

const (
	PropertyTypeUInt32 int = iota
	PropertyTypeInt32
	PropertyTypeDouble
	PropertyTypeString
)

type DeviceTSLManager struct {
	m map[string]*model.DeviceTSLTemplate
	sync.RWMutex
	publisher *hub.Hub
	changed   bool
}

var DeviceTSLMap = NewDeviceTSLManager()

func NewDeviceTSLManager() *DeviceTSLManager {
	return &DeviceTSLManager{
		m:         make(map[string]*model.DeviceTSLTemplate),
		publisher: hub.New(),
	}
}

func (tslManager *DeviceTSLManager) Init() error {
	tslManager.Lock()
	defer tslManager.Unlock()
	for _, v := range tslManager.m {
		if v.Plugin != "" {
			tmp, ok := DeviceTemplateMap[v.Plugin]
			if !ok {
				return fmt.Errorf("tsl template %s binding plugin %s is not exists", v.Name, v.Plugin)
			}
			v.PluginTemplate = tmp
		}
	}
	return nil
}

func (tslManager *DeviceTSLManager) Publish(name string, fields hub.Fields) {
	msg := hub.Message{
		Name:   name,
		Fields: fields,
	}
	tslManager.publisher.Publish(msg)
}

func (tslManager *DeviceTSLManager) SetChanged(r bool) {

	tslManager.changed = r

}

func (tslManager *DeviceTSLManager) Changed() bool {
	var r bool
	tslManager.RLock()
	r = tslManager.changed
	tslManager.RUnlock()
	return r
}

func (tslManager *DeviceTSLManager) AddTSL(tsl *model.DeviceTSLTemplate) bool {
	if tslManager.m == nil {
		tslManager.m = make(map[string]*model.DeviceTSLTemplate)
	}
	tslManager.Lock()
	defer tslManager.Unlock()
	if _, ok := tslManager.m[tsl.Name]; ok {
		return false
	}

	tslManager.m[tsl.Name] = tsl
	tslManager.changed = true
	return true
}

func (tslManager *DeviceTSLManager) DeleteTSL(tsl *model.DeviceTSLTemplate) bool {
	tslManager.Lock()
	defer tslManager.Unlock()
	if _, ok := tslManager.m[tsl.Name]; !ok {
		return false
	}
	delete(tslManager.m, tsl.Name)
	tslManager.changed = true
	return true
}

func (tslManager *DeviceTSLManager) Get(tslName string) *model.DeviceTSLTemplate {
	tslManager.RLock()
	tmp := tslManager.m[tslName]
	tslManager.RUnlock()
	return tmp
}

func (tslManager *DeviceTSLManager) GetAll() []*model.DeviceTSLTemplate {
	tslManager.RLock()
	var temps = make([]*model.DeviceTSLTemplate, 0, len(tslManager.m))
	for _, v := range tslManager.m {
		temps = append(temps, v)
	}
	tslManager.RUnlock()
	return temps
}

func (tslManager *DeviceTSLManager) SaveTo(f *os.File) error {
	tslManager.RLock()
	defer tslManager.RUnlock()
	data, err := json.Marshal(tslManager.m)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

// plugin
func (tslManager *DeviceTSLManager) ModifyPlugin(name string, plugin string) error {

	tmp := tslManager.Get(name)
	if tmp == nil {
		return fmt.Errorf("no such tsl template %s", name)
	}
	tmp.Plugin = plugin
	tslManager.changed = true
	return nil
}

func DeviceTSLExportPlugin(pluginName string) (zipPath string, err error) {

	//遍历文件
	pluginPath := path.Join(PLUGINPATH, pluginName)

	//保留原来文件的结构
	zipPath = path.Join(PLUGINPATH, pluginName+".zip")
	err = backup.Zip(pluginPath, zipPath)
	if err != nil {
		return zipPath, fmt.Errorf("zipFile err %v", err)
	}

	return
}
