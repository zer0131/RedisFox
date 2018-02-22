#!/bin/sh

if ! command -v glide >/dev/null 2>&1; then
  echo 'no exists glide'
  exit 1
fi

curr_dir=$(cd `dirname $0`; pwd)
old_gopath=$GOPATH
export GOPATH=$curr_dir

cd $curr_dir/src/redisfox

glide update

export GOPATH=$old_gopath

cd $curr_dir/