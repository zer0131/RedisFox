#!/bin/sh

curr_dir=$(cd `dirname $0`; pwd)
old_gopath=$GOPATH
export GOPATH=$curr_dir

cd $curr_dir/src/redisfox

glide update

export GOPATH=$old_gopath
