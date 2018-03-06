#!/bin/sh

curr=$(cd `dirname $0`; pwd)
nohup ${curr}/bin/redisfox -config ${curr}/conf/redis-fox.yaml >/dev/null 2>&1 &
echo "redisfox start..."