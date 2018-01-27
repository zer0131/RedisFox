#!/bin/sh

curr_dir=$(cd `dirname $0`; pwd)
old_gopath=$GOPATH
export GOPATH=$curr_dir

go install -v redisfox

echo "install complete"

export GOPATH=$old_gopath
