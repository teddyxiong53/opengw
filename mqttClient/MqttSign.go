package mqttClient

import (
    "crypto/hmac"
    "crypto/sha1"
    "fmt"
    "bytes"
    "log"
)

type AuthInfo struct {
    password        string
    username        string
    mqttClientId    string
}

func Calculate_sign(clientId, productKey, deviceName, deviceSecret, timeStamp string) AuthInfo {

    var raw_passwd bytes.Buffer

    raw_passwd.WriteString("clientId" + clientId)
    raw_passwd.WriteString("deviceName")
    raw_passwd.WriteString(deviceName)
    raw_passwd.WriteString("productKey")
    raw_passwd.WriteString(productKey)
    raw_passwd.WriteString("timestamp")
    raw_passwd.WriteString(timeStamp)
    log.Println(raw_passwd.String())

    // hmac, use sha1
    mac := hmac.New(sha1.New, []byte(deviceSecret))
    mac.Write([]byte(raw_passwd.String()))
    password:= fmt.Sprintf("%02x", mac.Sum(nil))
    log.Println(password)
    username := deviceName + "&" + productKey

    var MQTTClientId bytes.Buffer
    MQTTClientId.WriteString(clientId)
    // hmac, use sha1; securemode=2 means TLS connection 
    MQTTClientId.WriteString("|securemode=3,signmethod=hmacsha1,timestamp=")
    MQTTClientId.WriteString(timeStamp)
    MQTTClientId.WriteString("|")

    auth := AuthInfo{
        password:password,
        username:username,
        mqttClientId:MQTTClientId.String(),
    }
    return auth
}
