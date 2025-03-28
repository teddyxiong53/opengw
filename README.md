#### 交流
qq群1028704210

#### 框架设计
<div align=center><img src="https://images.gitee.com/uploads/images/2021/0330/152140_b19a7690_1979498.png"/></div>

软件主要分成3层：
- 应用接口
> 用于与上层应用系统进行通信，可以设置定时上报硬件设备数据到物联网平台，或者接收物联网平台下发命令，转发给硬件设备；采用Json等格式数据与上层应用系统通信，对接更简单；
    
- 采集接口
> - 用于对硬件设备进行管理，支持对设备数量、设备类型、设备属性的增、删、查、改等操作，同时可以设置定时采集设备的属性并缓存，方便上层应用系统对硬件设备操作；
> - 支持采用Lua脚本实现对设备通信协议的编写，方便灵活；

    
- 通信接口
> 对物理通信接口的封装，比如串口、网络、GPIO等，封装接口后对上提供读取和写入2个接口，方便上层调用；


#### 功能特点
- 采用golang语言设计，运行效率高，跨平台方便；
- 内置WebServer，网页配置更方便、更快捷
- 采用Lua脚本，增加设备类型时不需要重新编码后台代码，更方便灵活；
- 支持MqttClient，ModbusTCPServer，OPCUaServer等通信，采用JSON格式通信，上层系统对接更快捷；
- 支持CSV文件导入功能，批量添加；
- 支持配置文件的备份和回复；


#### 编译运行
1、编译
大家可以参考网络上如何编译golang程序的帖子
[参考链接](http://my.oschina.net/u/4521128/blog/4521037)

比如交叉编译成linux系统下，armV7架构的
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o goAdapter -ldflags "-s -w"
    
2、拷贝文件
需要拷贝以下文件到运行环境中：
- goAdapter执行文件
- webroot整个文件夹
- config整个文件夹
注意：如果运行环境是Linux系统，记得对文件权限进行修改
    
3、运行
对执行程序运行即可，
linux系统：./goAdapter
然后在浏览器中输入127.0.0.1:8080，注意加上端口，即可正常访问页面

#### 功能介绍

1. 通信接口
2. 采集接口
3. 应用接口

