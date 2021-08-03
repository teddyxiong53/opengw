package mqttEmqx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
)

type AuthInfo struct {
	password     string
	username     string
	mqttClientId string
}

func MqttClient_CalculateSign(productKey, deviceName, deviceSecret, timeStamp string) AuthInfo {

	clientId := productKey + "&" + deviceName

	var raw_passwd bytes.Buffer
	raw_passwd.WriteString("clientId")
	raw_passwd.WriteString(clientId)
	raw_passwd.WriteString("deviceName")
	raw_passwd.WriteString(deviceName)
	raw_passwd.WriteString("productKey")
	raw_passwd.WriteString(productKey)
	raw_passwd.WriteString("timestamp")
	raw_passwd.WriteString(timeStamp)
	//log.Println(raw_passwd.String())

	// hmac, use sha1
	mac := hmac.New(sha1.New, []byte(deviceSecret))
	mac.Write([]byte(raw_passwd.String()))
	password := fmt.Sprintf("%02x", mac.Sum(nil))
	//log.Println(password)
	username := deviceName + "&" + productKey

	var MQTTClientId bytes.Buffer
	MQTTClientId.WriteString(clientId)
	MQTTClientId.WriteString("|securemode=3,signmethod=hmacsha1,timestamp=")
	MQTTClientId.WriteString(timeStamp)
	MQTTClientId.WriteString("|")

	auth := AuthInfo{
		password:     password,
		username:     username,
		mqttClientId: MQTTClientId.String(),
	}
	return auth
}
