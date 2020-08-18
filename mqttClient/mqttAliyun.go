package mqttClient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	productKey 		string = "a1oSllgBCjt"
	deviceName 		string = "1111"
	deviceSecret 	string = "2d7d200249a49568cfbdace0900e6dcd"
	clientId 		string = "1111"
	timeStamp 		string = "1528018257135"
)

var (
	subTopic 		string = "/" + productKey + "/" + deviceName + "/user/get"
	pubTopic 		string = "/" + productKey + "/" + deviceName + "/user/update"
	registerTopic 	string = "/" + "ext/session/" + "a1Hhs4E2xUG" + "/" + "2001" + "/combine/login"
)

// define a function for the default message handler
var publishHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func MQTTClient_Init() {

	// set the login broker url
	var raw_broker bytes.Buffer
	raw_broker.WriteString("tls://")
	raw_broker.WriteString(productKey)
	raw_broker.WriteString(".iot-as-mqtt.cn-shanghai.aliyuncs.com:1883")
	opts := MQTT.NewClientOptions().AddBroker(raw_broker.String())

	// calculate the login auth info, and set it into the connection options
	auth := Calculate_sign(clientId, productKey, deviceName, deviceSecret, timeStamp)
	opts.SetClientID(auth.mqttClientId)
	opts.SetUsername(auth.username)
	opts.SetPassword(auth.password)
	opts.SetKeepAlive(60 * 2 * time.Second)
	opts.SetDefaultPublishHandler(publishHandler)

	// set the tls configuration
	//tlsconfig := NewTLSConfig()
	//opts.SetTLSConfig(tlsconfig)

	// create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Printf("Connect aliyun IoT Cloud Sucess\n")

	// subscribe to subTopic("/a1Zd7n5yTt8/deng/user/get") and request messages to be delivered
	if token := c.Subscribe(subTopic, 0, nil); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}
	log.Printf("Subscribe topic " + subTopic + " success\n")

	//// publish 5 messages to pubTopic("/a1Zd7n5yTt8/deng/user/update")
	//for i := 0; i < 5; i++ {
	//	log.Println("publish msg:", i)
	//	text := fmt.Sprintf("ABC #%d", i)
	//	token := c.Publish(pubTopic, 0, false, text)
	//	log.Println("publish msg: ", text)
	//	token.Wait()
	//	time.Sleep(2 * time.Second)
	//}

	authDevice := Calculate_sign("2001", "a1Hhs4E2xUG", "2001", "470c8262d6cee4bf440cf66758c2f3b8", timeStamp)
	log.Printf("authDevice %v\n",authDevice.password)
	registerText := `{
	  "id": "123",
	  "params": {
		"productKey": "a1Hhs4E2xUG",
		"deviceName": "2001",
		"clientId": "2001",
		"timestamp": "1528018257135",
		"signMethod": "hmacsha1",
		"sign": "edabfd45469dec12cf3432dd6083c58dadffb4ca",
		"cleanSession": "true"
	  }
	}`

	token := c.Publish(registerTopic, 0, false, registerText)
	log.Println("publish msg: ", registerText)
	token.Wait()

	/*
	// unsubscribe from subTopic("/a1Zd7n5yTt8/deng/user/get")
	if token := c.Unsubscribe(subTopic);token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
	 */
}



func NewTLSConfig() *tls.Config {
	// Import trusted certificates from CAfile.pem.
	// Alternatively, manually add CA certificates to default openssl CA bundle.
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile("./x509/root.pem")
	if err != nil {
		fmt.Println("0. read file error, game over!!")

	}

	certpool.AppendCertsFromPEM(pemCerts)

	// Create tls.Config with desired tls properties
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: false,
		// Certificates = list of certs client sends to server.
		// Certificates: []tls.Certificate{cert},
	}
}





