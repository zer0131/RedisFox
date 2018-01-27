#!/bin/sh

curr_dir=$(cd `dirname $0`; pwd)
old_gopath=$GOPATH
export GOPATH=$curr_dir

#用于更新
if [ "$1" != "" ]; then
    go get -v -u $1
    export GOPATH=$old_gopath
    echo $1" update complete!"
    exit
fi

src_dir="src"

if [ ! -d ${src_dir}"/github.com/go-yaml/yaml" ]; then
    go get -v github.com/go-yaml/yaml
    echo "install yaml complete!"
fi

if [ ! -d ${src_dir}"/github.com/garyburd/redigo" ]; then
    go get -v github.com/garyburd/redigo/redis
    echo "install redigo complete!"
fi

#先安装golang.org/x/net/context包
if [ ! -d ${src_dir}"/github.com/mattn/go-sqlite3" ]; then
    go get -v github.com/mattn/go-sqlite3
    echo "install go-sqlite3 complete!"
fi

if [ ! -d ${src_dir}"/github.com/gin-gonic/gin" ]; then
    go get -v github.com/gin-gonic/gin
    echo "install gin complete!"
fi

export GOPATH=$old_gopath
