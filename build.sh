#!/bin/sh

curr_dir=$(cd `dirname $0`; pwd)
old_gopath=$GOPATH
export GOPATH=$curr_dir

go install -v redisfox

rm -rf output

mkdir -p output/bin
mkdir -p output/data
mkdir -p output/log

cp bin/redisfox output/bin/
cp -r conf output/
cp -r tpl output/
cp -r static output/
cp tool/start.sh output/
cp tool/stop.sh output/

echo "build complete!"

export GOPATH=$old_gopath
