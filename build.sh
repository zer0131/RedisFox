#!/bin/sh

rm -rf output

mkdir -p output/bin
mkdir -p output/data
mkdir -p output/log

cp -r config tpl static tool/run.sh output/
cp tool/gosuv output/bin

curr_dir=$(cd `dirname $0`; pwd)

go build -o $curr_dir/output/bin/redisfox *go || exit -1

chmod -R 755 output/

echo "build complete!"

