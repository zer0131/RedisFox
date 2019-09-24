#!/bin/sh

PID=`tail -1 run_redisfox.pid|awk '{print $NF}'`
kill $PID
rm -f run_redisfox.pid
echo "RedisFox stop!"