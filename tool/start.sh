#!/bin/sh

curr=$(cd `dirname $0`; pwd)
nohup ${curr}/bin/redisfox -config ${curr}/config/redis-fox.yaml >/dev/null 2>&1 &
echo "RedisFox start..."