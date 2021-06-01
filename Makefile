#.PHONY用来定义伪目标。不创建目标文件，而是去执行这个目标下面的命令
.PHONY: linux-armv5 linux-armv7

BINARY="openGW"

linux-armv5:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o openGW -ldflags "-s -w"
linux-armv7:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o openGW -ldflags "-s -w"
linux-386:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o openGW -ldflags "-s -w"	
linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o openGW -ldflags "-s -w"
windows-386:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o openGW -ldflags "-s -w"  
windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o openGW -ldflags "-s -w" 