#!/usr/bin/env bash
BASEDIR="`pwd`"
LOG="$BASEDIR/log"

while true
do
    timeout 6h ./run_spider.sh > log.txt
    logger "Killed" "$LOG/run.log"
done
