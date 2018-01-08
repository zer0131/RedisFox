#!/bin/sh

if [ -z "$1" ];then
    echo "input pkg name"
    exit
fi

curr_dir=$(cd `dirname $0`; pwd)
old_gopath=$GOPATH
export GOPATH=$curr_dir

go install $1

export GOPATH=$old_gopath
