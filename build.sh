#!/bin/sh

CURDIR=`pwd`
OLD_GOPATH=$GOPATH
export GOPATH=$CURDIR

go get "github.com/astaxie/beego"

export GOPATH=$OLD_GOPATH
