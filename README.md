# goAdapter

#### 交流
qq群1028704210

#### 介绍
基于golang，适用于物联网项目中协议转换器

#### 软件架构
软件架构说明
![输入图片说明](https://images.gitee.com/uploads/images/2020/0904/151353_9a19564a_1979498.png "架构.png")

#### 安装教程

1.  编译
ARM版本：GOOS=linux GOARCH=arm GOARM=5 go build -o goAdapter-ldflags "-s -w"


#### 使用说明

1、将生成的可执行文件拷贝到板子内，同时把“config”、“webroot”这2个文件夹内全部文件也拷贝到板子内
2、在浏览器中输入127.0.0.1:8080，注意加上端口，即可正常访问页面


