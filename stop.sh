#!/bin/sh

PID=`tail -1 run_redisfox.pid|awk '{print $NF}'`
kill  -9  $PID
rm -f run_redisfox.pid