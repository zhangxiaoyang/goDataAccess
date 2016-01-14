#!/usr/bin/env bash
BASEDIR="`pwd`"
DATA="$BASEDIR/data"
BIN="$BASEDIR/bin"
OUTPUT="$BASEDIR/output"
LOG="$BASEDIR/log"

function init() {
    ulimit -n 100000

    mkdir -p "$OUTPUT"
    mkdir -p "$LOG"
    mkdir -p "$DATA"
    touch "$LOG/spider.log"
    touch "$DATA/output.txt"
    touch "$DATA/status.txt"

    if [ ! -f "$DATA/top-1m.txt" ]
    then
        cat > "$DATA/top-1m.txt" <<EOM
        qq.com
        163.com
EOM
    fi
}

function logger()
{
    echo -e "`date '+%Y-%m-%d %H:%M:%S'` $1" >> "$2"
}

init
logger "Started" "$LOG/run.log"
go run "$BIN/spider.go" "$BIN/spider.json" "$DATA/top-1m.txt" "$OUTPUT/output.txt" "$OUTPUT/status.txt" "$LOG/spider.log"
logger "Finished" "$LOG/run.log"
