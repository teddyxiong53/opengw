package httpServer

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goburrow/serial"
	modbus "github.com/thinkgos/gomodbus/v2"
	"github.com/thinkgos/mb"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MbParam struct {
	Addr    byte   `json:"Addr"`
	RegAddr uint16 `json:"RegAddr"`
	RegCnt  uint16 `json:"RegCnt"`
	Data    string `json:"Data"`
}

var client *mb.Client
var mbParam MbParam

func NewMbParam(addr byte, regAddr, regCnt uint16) *MbParam {
	return &MbParam{
		Addr:    addr,
		RegAddr: regAddr,
		RegCnt:  regCnt,
	}
}

func mbParamOpenPort(port string) bool {
	status := false

	//调用RTUClientProvider的构造函数,返回结构体指针
	p := modbus.NewRTUClientProvider(modbus.WithSerialConfig(serial.Config{
		Address:  port,
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "N",
		Timeout:  100 * time.Millisecond,
	}))

	client = mb.New(p)
	client.LogMode(true)
	err := client.Start()
	if err != nil {
		fmt.Println("start err,", err)
		return status
	}

	status = true
	return status
}

func mbParamReadHoldReg(slaveAddr byte, regAddr uint16, regCnt uint16) []uint16 {
	value, err := client.ReadHoldingRegisters(slaveAddr, regAddr, regCnt)
	if err != nil {
		fmt.Println("readHoldErr,", err)
	} else {
		//fmt.Printf("%#v\n", value)
	}

	return value
}

func mbParamWriteMutilReg(slaveAddr byte, regAddr uint16, regCnt uint16, data []byte) error {
	err := client.WriteMultipleRegistersBytes(slaveAddr, regAddr, regCnt, data)
	if err != nil {
		fmt.Println("writeMulRegErr,", err)
	}

	return err
}

func apiReadHoldReg(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)

	fmt.Println(string(bodyBuf[:n]))

	//获取读寄存器的参数
	rMbParam := &MbParam{}
	err := json.Unmarshal(bodyBuf[:n], rMbParam)
	if err != nil {
		fmt.Println("rMbParam json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson, _ := json.Marshal(aParam)

		context.String(http.StatusOK, string(sJson))
		return
	}

	mRegValue := mbParamReadHoldReg(rMbParam.Addr, rMbParam.RegAddr, rMbParam.RegCnt)
	fmt.Printf("mRegValue %#v\n", mRegValue)
	var sRegValue string = ""

	for i := 0; i < len(mRegValue); i++ {
		sRegValue += strconv.FormatUint(uint64(mRegValue[i]), 10)

		if i != (len(mRegValue) - 1) {
			sRegValue = sRegValue + ","
		} else {
			sRegValue = sRegValue + ""
		}
	}
	fmt.Println(sRegValue)

	aParam.Code = "0"
	aParam.Message = ""

	aParam.Data = sRegValue
	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}

func apiWriteMultiReg(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)

	fmt.Println(string(bodyBuf[:n]))

	//获取写寄存器的参数
	rMbParam := &MbParam{}
	err := json.Unmarshal(bodyBuf[:n], rMbParam)
	if err != nil {
		fmt.Println("rMbParam json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)

		context.String(http.StatusOK, string(sJson))
		return
	}

	//将字符串中数值取出来
	sData := strings.Split(rMbParam.Data, ",")
	//fmt.Println(sData)
	bData := make([]byte, 0)
	bData2 := make([]byte, 2)
	for _, v := range sData {
		//tByte,err := strconv.ParseUint(v,10,16)
		tByte, _ := strconv.Atoi(v)
		binary.BigEndian.PutUint16(bData2, uint16(tByte))

		bData = append(bData, bData2...)
	}
	fmt.Printf("bData is %d\n", bData)

	err = mbParamWriteMutilReg(rMbParam.Addr, rMbParam.RegAddr,
		rMbParam.RegCnt, bData)
	if err != nil {
		aParam.Code = "1"
		aParam.Message = "write reg timeout"
		sJson, _ := json.Marshal(aParam)

		context.String(http.StatusOK, string(sJson))

	}
	aParam.Code = "0"
	aParam.Message = ""
	sJson, _ := json.Marshal(aParam)

	context.String(http.StatusOK, string(sJson))
}
