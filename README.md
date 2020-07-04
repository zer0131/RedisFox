# RedisFox

## Introduction

[[Simplified Chinese]](./tool/zh.md)

RedisFox is a visual monitoring tool for Redis based on Golang development

![redisfox](./tool/redisfox.png)

## Latest Version Download

File Name|Kind|OS|Size
------|------|------|------
[redisfox1.0.2.darwin-amd64.tar.gz](http://resource.zhangenrui.cn/redisfox/redisfox1.0.2.darwin-amd64.tar.gz)|Archive|MacOS|6M
[redisfox1.0.2.linux-amd64.tar.gz](http://resource.zhangenrui.cn/redisfox/redisfox1.0.2.linux-amd64.tar.gz)|Archive|Linux|6M

## Instructions

### Special description 
* Go version requires **>1.12**, and use go modlue as package dependency management(The author used Go1.12.9)
* Run with gosuv driver 

1. Download RedisFox

```
git clone https://github.com/zer0131/RedisFox.git
```

2. Compile and Install

```
sh build.sh
```

3. Run

Adjustment programs.yml **directory**
Configure the redis server in **conf/redis-fox.yaml**, open redis, and then execute the run.sh script

```
cd output
sh run.sh start
```

4. Visit

Open the browser to access **http://127.0.0.1:8080** and see the monitoring status of redis态

5. Stop

```
sh run.sh stop
```

## Directory Structure

```
├─config                 Config directory
│  ├─redis-fox.yaml      Config file
├─log                    Log directory
├─data                   Data directory
├─static                 Resource directory
├─tpl                    Template directory
├─tool                   Tool directory
├─conf                   Source Code conf
├─dataprovider           Source Code dataprovider
├─process                Source Code process
├─server                 Source Code server
├─util                   Source Code util
├─main.go                Source Code main file
└─build.sh               Compile and install the script
```

## Configuration Description

Configuration based on yaml syntax, do not understand the stamp here[http://www.ruanyifeng.com/blog/2016/07/yaml.html](http://www.ruanyifeng.com/blog/2016/07/yaml.html)

```
#!The same level of field indentation is the same, and can only be indented with space.
#!The relative "-" or ":" must have spaces before each field value value

#[Redis service configuration]

#redis server
servers:
 - server: 127.0.0.1
   port: 6379
   conntype: tcp
   password: 123456 #passport

#Get redis information interval time (second)
sleeptime: 30

#The maximum number of connections in the redis connection pool
maxidle: 3

#The maximum number of active redis connection pools, 0 unrestricted
maxactive: 3

#The redis connection pool connects the timeout time, and 0 indicates no timeout
idletimeout: 0

#Storage data type
datatype: sqlite
#Data storage path
datapath: ./data/redisfox.db

#Log
logpath: ./log/
logname: redisfox.log
loglevel: 4

#Web
serverip: 127.0.0.1
serverport: 8080
debugmode: 0

#Resource
staticdir: ./static/
tpldir: ./tpl/
```

## Nginx

```
server {
    server_name wwww.xxxx.com;
    listen 80; # or 443
    # ssl on; Whether to enable encrypted connections
    # If you use HTTPS, you also need to fill in ssl_certificate and ssl_certificate_key

    location / { # If you want to access the subpath, this is changed to a subpath, pay attention to / begin and end / end
        proxy_pass http://127.0.0.1:8080/;
    }
    access_log  /your-path/nginx/logs/redisfox.log;
}
```

