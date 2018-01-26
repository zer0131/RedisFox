# RedisFox

## 简介

RedisFox是一款基于Golang开发的Redis可视化监控工具。

![redisfox](./redisfox.png)

## 安装及运行

假设你已经配置好Golang环境（作者用的是Go1.9.2环境）

1. 下载RedisFox

```
git clone https://github.com/zer0131/RedisFox.git
```

2. 安装golang.org/x/net/context包

去[https://gopm.io/](https://gopm.io/)下载压缩包，解压后，将文件夹中的内容复制到 **src/golang.org/x/net/** 目录下（没有目录自行创建）

3. 获取其他依赖包

```
sh getall.sh
```

4. 安装

```
sh install.sh
```

5. 运行

在 **conf/redis-fox.yaml** 配置redis服务器，并开启redis，然后执行start.sh脚本

```
sh start.sh
```

6. 访问

打开浏览器访问 **http://127.0.0.1:8080** 即可查看redis的监控状态

7. 停止

```
sh stop.sh
```


