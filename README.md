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

Suppose you have configured the Golang environment (the author uses the Go1.9.2 environment)

1. Download RedisFox

```
git clone https://github.com/zer0131/RedisFox.git
```

2. Dependency Package

**The project uses glide to manage dependencies. First, you need to install glide in your environment**

[https://glide.readthedocs.io/en/latest/](https://glide.readthedocs.io/en/latest/)

Glide.yaml is configured under the src/redisfox directory


```
sh pkg.sh
```

4. Compile and Install

```
sh build.sh
```

5. Run

Configure the redis server in **conf/redis-fox.yaml**, open redis, and then execute the start.sh script

```
cd output
sh start.sh
```

6. Visit

Open the browser to access **http://127.0.0.1:8080** and see the monitoring status of redis态

7. Stop

```
sh stop.sh
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
├─pkg.sh                 Get the Go dependency script
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

## Glide

```
package: RedisFox
import:
- package: github.com/garyburd/redigo
  version: ^1.4.0
  subpackages:
  - redis
- package: github.com/gin-gonic/gin
  version: ^1.2.0
- package: github.com/go-yaml/yaml
- package: github.com/mattn/go-sqlite3
  version: ^1.6.0
- package: golang.org/x/net
  repo: https://github.com/golang/net.git
- package: golang.org/x/sys
  repo: https://github.com/golang/sys.git
```


