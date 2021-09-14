package controller

import (
	"encoding/binary"
	"fmt"
	"goAdapter/httpServer/model"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goburrow/serial"
	modbus "github.com/thinkgos/gomodbus/v2"
)

type MbParam struct {
	Addr    byte   `json:"Addr"`
	RegAddr uint16 `json:"RegAddr"`
	RegCnt  uint16 `json:"RegCnt"`
	Data    string `json:"Data"`
}

var client modbus.Client
var mbParam MbParam

func NewMbParam(addr byte, regAddr, regCnt uint16) *MbParam {
	return &MbParam{
		Addr:    addr,
		RegAddr: regAddr,
		RegCnt:  regCnt,
	}
}

func mbParamOpenPort(port string) bool {
	//调用RTUClientProvider的构造函数,返回结构体指针
	p := modbus.NewRTUClientProvider(modbus.WithSerialConfig(serial.Config{
		Address:  port,
		BaudRate: 115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   "N",
		Timeout:  100 * time.Millisecond,
	}))

	client = modbus.NewClient(p)
	client.LogMode(true)
	err := client.Connect()
	if err != nil {
		fmt.Println("start err,", err)
		return false
	}

	return true
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

func ReadHoldReg(context *gin.Context) {
	//获取读寄存器的参数
	rMbParam := &MbParam{}
	err := context.ShouldBindJSON(rMbParam)
	if err != nil {
		fmt.Println("rMbParam json unMarshall err,", err)

		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "json unMarshall err",
		})
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

	context.JSON(http.StatusOK, model.Response{
		Code:    "0",
		Message: "",
		Data:    sRegValue,
	})
}

func WriteMultiReg(context *gin.Context) {
	//获取写寄存器的参数
	rMbParam := &MbParam{}
	err := context.ShouldBindJSON(rMbParam)
	if err != nil {
		fmt.Println("rMbParam json unMarshall err,", err)

		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "json unMarshall err",
			Data:    "",
		})
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
		context.JSON(http.StatusOK, model.Response{
			Code:    "1",
			Message: "write reg timeout",
		})
		return
	}

	context.JSON(http.StatusOK, model.Response{
		Code:    "0",
		Message: "",
		Data:    "",
	})
}
