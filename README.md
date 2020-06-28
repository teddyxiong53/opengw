# goAdapter

#### 介绍
基于golang，适用于物联网项目中协议转换器

#### 软件架构
软件架构说明



#### 安装教程

1.  编译
ARM版本：GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-s -w"


#### 使用说明

1、将生成的可执行文件拷贝到板子内，同时把selfpara webroot这2个文件夹内全部文件也拷贝到板子内
2、修改selfpara文件夹中的networkpara.json,修改IP地址
3、在浏览器中输入修改后的IP地址，同时加上端口即可，比如192.168.1.1:8090


