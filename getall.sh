#!/bin/sh

curr_dir=$(cd `dirname $0`; pwd)
old_gopath=$GOPATH
export GOPATH=$curr_dir

src_dir="src"

if [ ! -d ${src_dir}"/github.com/gin-gonic/gin" ]; then
    go get github.com/gin-gonic/gin
    echo "install gin complete"
fi

if [ ! -d ${src_dir}"/github.com/garyburd/redigo" ]; then
    go get github.com/garyburd/redigo/redis
    echo "install redigo complete"
fi

if [ ! -d ${src_dir}"/github.com/go-yaml/yaml" ]; then
    go get github.com/go-yaml/yaml
    echo "install yaml complete"
fi

#先安装golang.org/x/net/context包
if [ ! -d ${src_dir}"/github.com/mattn/go-sqlite3" ]; then
    go get github.com/mattn/go-sqlite3
    echo "install go-sqlite3 complete"
fi

export GOPATH=$old_gopath
