#!/bin/sh

curr_dir=$(cd `dirname $0`; pwd)
old_gopath=$GOPATH
export GOPATH=$curr_dir

go install redisfox

export GOPATH=$old_gopath
