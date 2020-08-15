package mqttClient

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	aiot "github.com/thinkgos/aliyun-iot"
	"github.com/thinkgos/aliyun-iot/dm"
	"github.com/thinkgos/aliyun-iot/infra"
	"github.com/thinkgos/aliyun-iot/sign"
	"log"
	"math/rand"
)

const (
	productKey    = "a1oSllgBCjt"
	productSecret = ""
	deviceName    = "1111"
	deviceSecret  = "2d7d200249a49568cfbdace0900e6dcd"
)

var dmClient *aiot.MQTTClient

type Properties struct{
	Value interface{}           `json:value`
}

type ServiceContent struct{
	ServiceID  string     `json:"service_id"`
	Properties Properties `json:"properties"`
}

type PublishContent struct{
	Service []ServiceContent `json:"services"`
}

type AliYunMqttClientTemplate struct{
	ClientId string				`json:"ClientID"`
	DeviceName string			`json:"DeviceName"`
	ProductKey string			`json:"ProductKey"`
	DeviceSecret string			`json:"DeviceSecret"`
	TimeStamp string			`json:"TimeStamp"`
}

var mClient mqtt.Client

var AliYunMqttClient = AliYunMqttClientTemplate{
	ClientId : "1111|securemode=3,signmethod=hmacsha1|",
	DeviceName : "1111",
	ProductKey : "a1oSllgBCjt",
	DeviceSecret : "2d7d200249a49568cfbdace0900e6dcd",
	TimeStamp : "1528018257135",
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Pub Client Topic : %s \n", msg.Topic())
	log.Printf("Pub Client msg : %s \n", msg.Payload())
}

//创建全局mqtt sub消息处理 handler
var messageSubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Sub Client Topic : %s \n", msg.Topic())
	log.Printf("Sub Client msg : %s \n", msg.Payload())
}

func MqttAppConnect() {

	//设备三元信息
	metaInfo := &infra.MetaInfo{
		ProductKey:    productKey,
		ProductSecret: productSecret,
		DeviceName:    deviceName,
		DeviceSecret:  deviceSecret,
	}
	//设备域名信息
	cloudRegionDomain := infra.CloudRegionDomain{
		Region: infra.CloudRegionShangHai,
	}

	mqttSign := sign.NewMQTTSign()
	mqttSign.SetSDKVersion(infra.IOTSDKVersion)
	mqttSignInfo,err := mqttSign.Generate(metaInfo,cloudRegionDomain)
	if err != nil {
		panic(err)
	}

	clientOpts := mqtt.NewClientOptions()
	clientOpts = clientOpts.AddBroker(fmt.Sprintf("%s:%d", mqttSignInfo.HostName, mqttSignInfo.Port))
	clientOpts.SetClientID(mqttSignInfo.ClientID)
	clientOpts.SetUsername(mqttSignInfo.UserName)
	clientOpts.SetPassword(mqttSignInfo.Password)
	clientOpts.SetCleanSession(true)
	clientOpts.SetAutoReconnect(true)
	clientOpts.SetOnConnectHandler(func(cli mqtt.Client) {
		log.Println("mqtt client connection success")
	})
	clientOpts.SetConnectionLostHandler(func(cli mqtt.Client, err error) {
		log.Println("mqtt client connection lost, ", err)
	})

	dmOpt := dm.NewConfig(productKey, deviceName, deviceSecret)
	dmOpt = dmOpt.Valid()

	dmClient = aiot.NewWithMQTT(dmOpt, mqtt.NewClient(clientOpts))
	dmClient.LogMode(true)

	dmClient.UnderlyingClient().Connect().Wait()
	if err = dmClient.AlinkConnect(); err != nil {
		panic(err)
	}

	/*
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)


	//opts := mqtt.NewClientOptions().AddBroker("iot-mqtts.cn-north-4.myhuaweicloud.com:1883")

	//brokerStr := AliYunMqttClient.ProductKey + ".iot-as-mqtt.cn-shanghai.aliyuncs.com:1883"
	//opts := mqtt.NewClientOptions().AddBroker(brokerStr)


	var raw_broker bytes.Buffer
	raw_broker.WriteString("tls://")
	raw_broker.WriteString(AliYunMqttClient.ProductKey)
	raw_broker.WriteString(".iot-as-mqtt.cn-shanghai.aliyuncs.com:1883")
	opts := mqtt.NewClientOptions().AddBroker(raw_broker.String())

	auth := Calculate_sign(AliYunMqttClient.ClientId,
							AliYunMqttClient.ProductKey,
								AliYunMqttClient.DeviceName,
									AliYunMqttClient.DeviceSecret,
										AliYunMqttClient.TimeStamp)
	log.Printf("auth %+v\n",auth)
	opts.SetClientID(auth.mqttClientId)
	opts.SetUsername(auth.username)
	opts.SetPassword(auth.password)
	opts.SetKeepAlive(60 * 2 * time.Second)

	opts.SetConnectTimeout(30*time.Second)
	//opts.SetKeepAlive(120*time.Second)
	opts.SetCleanSession(false)
	//opts.SetProtocolVersion(3)

	opts.SetDefaultPublishHandler(messagePubHandler)

	mClient = mqtt.NewClient(opts)

	if token := mClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//订阅消息
	subTopic := "/" + AliYunMqttClient.ProductKey + "/" + AliYunMqttClient.DeviceName + "/user/get";
	token := mClient.Subscribe(subTopic, 0, messageSubHandler)
	log.Printf("[Sub] end Subscribe msg to mqtt broker,token : %s \n", token)

	//
	//c.Disconnect(250)


	 */

}

/*
func mqttAppPublish(){

	var topic string = "$oc/devices/{" + "5ea67d5d58115909547f50e8_11111112" + "}/sys/properties/report"

	var publishContent = PublishContent{}
	serviceContent := make([]ServiceContent,0)

	properties := Properties{
		Value:23,
	}

	serviceContent = append(serviceContent, ServiceContent{
		ServiceID:  "Temp",
		Properties: properties,
	})

	publishContent.Service = serviceContent
	sJson,err := json.Marshal(publishContent)
	if err != nil{
		log.Println("publishContent json marshal err")
	}
	log.Println(string(sJson))

	token := mClient.Publish(topic, 0, false, string(sJson))
	token.Wait()
}

 */

func MqttAppPublish(){

	publishParam := map[string]interface{}{
		"Temp":         rand.Intn(200),
		"Humi":         rand.Intn(100),
	}

	err := dmClient.AlinkReport(dm.MsgTypeEventPropertyPost,
								dm.DevNodeLocal,
								publishParam)

	if err != nil {
		log.Printf("error: %#v", err)
	}
}

