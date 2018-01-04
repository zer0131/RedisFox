#!/bin/sh

if [ -z "$1" ];then
    echo "input pkg name"
    exit
fi

CURDIR=`pwd`
OLD_GOPATH=$GOPATH
export GOPATH=$CURDIR

go install $1

export GOPATH=$OLD_GOPATH
