#!/bin/sh

rm -rf output
rm -rf src

mkdir -p src/RedisFox
mkdir -p output/bin
mkdir -p output/data
mkdir -p output/log

cp -r config output/
cp -r tpl output/
cp -r static output/
cp tool/start.sh output/
cp tool/stop.sh output/

cp -r conf dataprovider process server util glide.yaml main.go src/RedisFox

cd src/RedisFox
glide update -v || exit -1

curr_dir=$(cd `dirname $0`; pwd)
old_gopath=$GOPATH
export GOPATH=$curr_dir


go build -o $curr_dir/output/bin/redisfox *go || exit -1

export GOPATH=$old_gopath
cd $curr_dir
#rm -rf src

echo "build complete!"

