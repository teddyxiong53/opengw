package main

import (
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
)

type Properties struct{
	Value interface{}           `json:value`
}

type ServiceContent struct{
	ServiceID 	string			`json:"service_id"`
	Properties  Properties		`json:"properties"`
}

type PublishContent struct{
	Service []ServiceContent   	`json:"services"`
}

var mClient mqtt.Client
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Pub Client Topic : %s \n", msg.Topic())
	fmt.Printf("Pub Client msg : %s \n", msg.Payload())
}

//创建全局mqtt sub消息处理 handler
var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Sub Client Topic : %s \n", msg.Topic())
	fmt.Printf("Sub Client msg : %s \n", msg.Payload())
}

func mqttAppConnect(){

	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	opts := mqtt.NewClientOptions().AddBroker("iot-mqtts.cn-north-4.myhuaweicloud.com:1883")

	opts.SetConnectTimeout(30*time.Second)
	opts.SetKeepAlive(120*time.Second)
	opts.SetCleanSession(false)
	opts.SetProtocolVersion(4)

	opts.SetClientID("5ea67d5d58115909547f50e8_11111112_0_0_2020042810")
	opts.SetUsername("5ea67d5d58115909547f50e8_11111112")
	opts.SetPassword("8689ac2c6207f1580d08182c2ed9d4129b76790a59ea4e0e811d361381335372")

	opts.SetDefaultPublishHandler(messagePubHandler)

	mClient = mqtt.NewClient(opts)

	//订阅消息
	token := mClient.Subscribe("go-test-topic", 0, messageSubHandler)
	fmt.Printf("[Sub] end Subscribe msg to mqtt broker,token : %s \n", token)

	if token := mClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}


	//
	//c.Disconnect(250)
}

func mqttAppPublish(){

	var topic string = "$oc/devices/{" + "5ea67d5d58115909547f50e8_11111112" + "}/sys/properties/report"

	var publishContent = PublishContent{}
	serviceContent := make([]ServiceContent,0)

	properties := Properties{
		Value:23,
	}

	serviceContent = append(serviceContent,ServiceContent{
		ServiceID:  "Temp",
		Properties: properties,
	})

	publishContent.Service = serviceContent
	sJson,err := json.Marshal(publishContent)
	if err != nil{
		fmt.Println("publishContent json marshal err")
	}
	log.Println(string(sJson))

	token := mClient.Publish(topic, 0, false, string(sJson))
	token.Wait()
}

