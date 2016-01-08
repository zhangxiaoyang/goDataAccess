#!/usr/bin/env bash
BASEDIR="`pwd`"
DATA="$BASEDIR/data"
BIN="$BASEDIR/bin"
OUTPUT="$BASEDIR/output"
LOG="$BASEDIR/log"
LIMIT=10000

function init() {
    ulimit -n 4000

    rm -rf "$SPLIT"
    rm -rf "$OUTPUT"
    rm -rf "$LOG"
    mkdir -p "$OUTPUT"
    mkdir -p "$LOG"
    mkdir -p "$DATA"

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
