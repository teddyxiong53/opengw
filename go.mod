module goAdapter

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/gin-gonic/gin v1.6.3
	github.com/robfig/cron v1.2.0
	github.com/safchain/ethtool v0.0.0-20200609180057-ab94f15152e7
	github.com/shirou/gopsutil v2.20.5+incompatible
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
	github.com/thinkgos/gomodbus v1.5.2
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	gopkg.in/ini.v1 v1.57.0
	deviceAPI v0.0.0
)

replace (
	deviceAPI => ../../deviceAPI
)