/*
@Description: This is auto comment by koroFileHeader.
@Author: Linn
@Date: 2021-09-13 11:28:56
@LastEditors: WalkMiao
@LastEditTime: 2021-09-14 14:56:21
@FilePath: /goAdapter-Raw/device/helper.go
*/
package device

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"goAdapter/pkg/luautils"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	lua "github.com/yuin/gopher-lua"
)

type LUAFUNC string

const (
	GETREAL      LUAFUNC = "GetRealVariables"         //获取实时值
	GENERATEREAL LUAFUNC = "GenerateGetRealVariables" //生成发送字节
	//GETDEVICEREAL LUAFUNC = "GetDeviceRealVariables"
	DEVICECUSTOMCMD LUAFUNC = "DeviceCustomCmd"
	RXBUF           LUAFUNC = "rxBuf"
	TXBUF           LUAFUNC = "txBuf"
	ANALYSISRX      LUAFUNC = "AnalysisRx"
	NewVariables    LUAFUNC = "NewVariables"
)

const (
	PLUGINPATH        = "./plugin"
	COMMJSON          = "commInterface.json"
	NETWORKJSON       = "commTcpClientInterface.json"
	IOINJSON          = "commIoInInterface.json"
	IOOUTJSON         = "commIoOutInterface.json"
	COLLINTERFACEJSON = "collInterface.json"
)

const (
	SERIALTYPE    = "LocalSerial"
	TCPCLIENTTYPE = "TcpClient"
	IOINTYPE      = "IoIn"
	IOOUTTYPE     = "IoOut"
)

const (
	ONLINE  = "onLine"
	OFFLINE = "offLine"
)

func disPatchCommonFunction(state *lua.LState) {
	for name, fn := range luautils.CommonLuaFuncMap {
		state.SetGlobal(name, state.NewFunction(fn))
	}
}

func parseJson(jsonFile string, index int) (err error) {
	fp, err := os.Open(jsonFile)
	if err != nil {
		return fmt.Errorf("open json config file %s Err:%v", jsonFile, err)
	}
	defer fp.Close()
	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return
	}
	var devTemp = DeviceTemplate{
		TemplateID: index,
	}
	err = json.Unmarshal(data, &devTemp)

	if err != nil {
		return
	}

	DeviceTemplateMap[devTemp.Type] = &devTemp
	return nil
}

func parseLua(luaFile string, masterFile string) (err error) {
	tName := strings.Split(path.Base(masterFile), ".")
	if len(tName) != 2 {
		return fmt.Errorf("%s is not valid pattern", masterFile)
	}
	t, ok := DeviceTemplateMap[tName[0]]
	if !ok {
		return fmt.Errorf("template %s haven't reload from json", tName[0])
	}
	if luaFile == masterFile {
		state, err := luautils.LuaOpenFile(luaFile)
		if err != nil {
			return fmt.Errorf("open lua file %s error:%v", luaFile, err)
		}
		disPatchCommonFunction(state)
		t.LuaState = state
		log.Println(color.CyanString("模板主lua【%s】初始化成功", luaFile))

	} else {
		state := t.LuaState
		if state != nil {
			err = state.DoFile(luaFile)
			if err != nil {
				log.Println(color.RedString("模板lua加载辅助lua【%s】失败", luaFile, err))
				return err
			}
			log.Println(color.CyanString("模板lua加载辅助lua【%s】成功", luaFile))
		} else {
			log.Println(color.RedString("模板主lua【%s】还未加载！"))
		}
	}

	return
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func WriteJsonErrorHandler(ctx *gin.Context, cfg string, errCode, succCode int, succInfo string) {
	if err := WriteToJson(cfg); err != nil {
		ctx.JSON(errCode, struct {
			Code    string
			Message string
		}{
			Code:    "1",
			Message: err.Error(),
		})
	}
	ctx.JSON(succCode, struct {
		Code    string
		Message string
	}{
		Code:    "0",
		Message: succInfo,
	})
}

//所有不同类型的通讯接口都放在一个json里
func WriteToJson(cfg string) (err error) {
	fileDir, err := filepath.Abs(path.Join("./selfpara", cfg))
	if err != nil {
		err = fmt.Errorf("load file 【%s】 error:%v", cfg, err)
		return
	}

	fp, err := os.OpenFile(fileDir, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return
	}
	defer fp.Close()
	var data []byte
	switch cfg {
	case COMMJSON:
		data, err = json.Marshal(CommunicationInterfaceMap)
	case COLLINTERFACEJSON:
		data, err = json.Marshal(CollectInterfaceMap)
		//TODO 其他case
	}
	if err != nil {
		return fmt.Errorf("marshal %s error:%v", cfg, err)
	}
	_, err = fp.Write(data)
	if err != nil {
		return err
	}
	return
}

func LoadJsonFile(cfg string) error {
	fileDir, err := filepath.Abs(path.Join("./selfpara", cfg))
	if err != nil {
		return fmt.Errorf("load file 【%s】 error:%v", cfg, err)
	}

	if FileExist(fileDir) {
		fp, err := os.OpenFile(fileDir, os.O_RDONLY, 0777)
		if err != nil {
			return fmt.Errorf("open file 【%s】  err:%v", fileDir, err)
		}
		defer fp.Close()

		data, err := ioutil.ReadAll(fp)
		if err != nil {
			return fmt.Errorf("ioutil readall 【%s】 error:%v", fileDir, err)

		}
		switch cfg {
		case COMMJSON:
			return CommInterfaceInit(data)
		case COLLINTERFACEJSON:
			if err := json.Unmarshal(data, &CollectInterfaceMap); err != nil {
				return err
			}
			for k, v := range CollectInterfaceMap {
				if v.DeviceNodes == nil {
					v.DeviceNodes = make([]*DeviceNodeTemplate, 0, 10)
				} else {
					for _, node := range v.DeviceNodes {
						v.InitDeviceNode(node)
					}
				}
				var param = &CollectInterfaceParamTemplate{
					CollInterfaceName: v.CollInterfaceName,
					CommInterfaceName: v.CommInterfaceName,
					PollPeriod:        v.PollPeriod,
					OfflinePeriod:     v.OfflinePeriod,
					DeviceNodeCnt:     v.DeviceNodeCnt,
					DeviceNodes:       v.DeviceNodes,
				}

				collect, err := NewCollectInterface(param)
				if err != nil {
					return err
				}
				CollectInterfaceMap[k] = collect
				CommunicationManage.Collectors <- &CollectInterfaceStatus{
					Tmp: collect,
					ACT: ADD,
				}
			}
		default:
			return errors.New("unsupported cfg file")
		}
	}
	return nil
}

func ReadPlugins(plugPath string) error {
	base, err := filepath.Abs(plugPath)
	if err != nil {
		return fmt.Errorf("get absolute dir path error:%v", err)
	}

	if _, err := os.Stat(base); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(base, 0666); err != nil {
				return err
			}

		}
		return err
	}

	//遍历json和so文件
	entries, err := ioutil.ReadDir(base)
	if err != nil {
		return err
	}
	//遍历plugin 目录
	for index, entry := range entries {
		if entry.IsDir() {
			pluginDir := path.Join(base, entry.Name())
			fs, err := ioutil.ReadDir(pluginDir)
			if err != nil {
				return err
			}

			var priorityMap = make(map[int]string)
			for index, f := range fs {
				if !f.IsDir() {
					switch ext := path.Ext(f.Name()); ext {
					case ".json":
						priorityMap[-1] = path.Join(pluginDir, f.Name())
					case ".lua":
						if strings.HasPrefix(f.Name(), entry.Name()) {
							priorityMap[-2] = path.Join(pluginDir, f.Name())
						} else {
							priorityMap[index] = path.Join(pluginDir, f.Name())
						}
					}

				}
			}
			if f, ok := priorityMap[-1]; ok {
				if err = parseJson(f, index); err != nil {
					return err
				}
			}
			if f, ok := priorityMap[-2]; ok {
				if err = parseLua(f, priorityMap[-2]); err != nil {
					return err
				}
			}
			for priority, f := range priorityMap {
				if priority >= 0 {
					if err := parseLua(f, priorityMap[-2]); err != nil {
						return err
					}
				}
			}
			log.Println(color.CyanString("解析模板【%s】成功！", entry.Name()))

		}
	}
	return nil
}
