# RedisFox

## 简介

RedisFox是一款基于Golang开发的Redis可视化监控工具。

![redisfox](./redisfox.png)

## 最新版本下载

File Name|Kind|OS|Size
------|------|------|------
[redisfox1.0.2.darwin-amd64.tar.gz](http://7xkyq4.com1.z0.glb.clouddn.com/redisfox/redisfox1.0.2.darwin-amd64.tar.gz)|Archive|MacOS|6M
[redisfox1.0.2.linux-amd64.tar.gz](http://7xkyq4.com1.z0.glb.clouddn.com/redisfox/redisfox1.0.2.linux-amd64.tar.gz)|Archive|Linux|6M

## 编译安装及运行

### 特别说明：
* Go版本需要 **>v1.12**，并且使用Go Modlue做为包依赖管理(作者用的是Go1.12.9)
* 使用gosuv驱动程序运行

1. 下载RedisFox

```
git clone https://github.com/zer0131/RedisFox.git
```

2. 编译安装

```
sh build.sh
```

3. 运行

调整programs.yml中的 **directory**
在 **conf/redis-fox.yaml** 配置redis服务器，并开启redis，然后执行run.sh

```
cd output
sh start.sh start
```

4. 访问

打开浏览器访问 **http://127.0.0.1:8080** 即可查看redis的监控状态

5. 停止

```
sh run.sh stop
```

## 目录介绍

```
├─config                 配置文件目录
│  ├─redis-fox.yaml      配置文件
├─log                    日志目录
├─data                   数据目录
├─static                 静态资源目录
├─tpl                    模板目录
├─tool                   工具目录
├─conf                   源码 conf
├─dataprovider           源码 dataprovider
├─process                源码 process
├─server                 源码 server
├─util                   源码 util
├─main.go                源码 main 文件
└─build.sh               程序编译安装脚本
```

## 配置说明

配置基于yaml语法，不懂戳这里[http://www.ruanyifeng.com/blog/2016/07/yaml.html](http://www.ruanyifeng.com/blog/2016/07/yaml.html)

```
#!同级的字段缩进相同，且只能用空格缩进
#!每个字段value值前相对“-”或者“：”必须有空格

#[redis服务相关配置]

#redis服务器
servers:
 - server: 127.0.0.1
   port: 6379
   conntype: tcp
   password: 123456 #redis密码，根据自己的需求更改，无密码不用配置

#获取redis信息间隔时间(秒)
sleeptime: 30

#redis连接池最大连接数
maxidle: 3

#redis连接池最大活跃数,0表示无限制
maxactive: 3

#redis连接池连接超时时间,0表示不超时
idletimeout: 0

#存储数据类型
datatype: sqlite
#数据存储路径
datapath: ./data/redisfox.db

#日志相关
logpath: ./log/
logname: redisfox.log
loglevel: 4

#web访问ip
serverip: 127.0.0.1
#web访问端口
serverport: 8080
#web调试模式是否开启
debugmode: 0

#静态文件目录
staticdir: ./static/
#模板文件目录
tpldir: ./tpl/
```

## Nginx反向代理

```
server {
    server_name 域名或IP;
    listen 80; # 或者 443，如果你使用 HTTPS 的话
    # ssl on; 是否启用加密连接
    # 如果你使用 HTTPS，还需要填写 ssl_certificate 和 ssl_certificate_key

    location / { # 如果你希望通过子路径访问，此处修改为子路径，注意以 / 开头并以 / 结束
        proxy_pass http://127.0.0.1:8080/;
    }
    access_log  /your-path/nginx/logs/redisfox.log;
}
```


