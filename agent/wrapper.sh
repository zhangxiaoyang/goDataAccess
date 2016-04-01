#!/usr/bin/env bash

function prepare()
{
    rm -rf db/

    ps aux | grep run_agent | awk '{print $2}' | xargs kill -9
    ps aux | grep run_server | awk '{print $2}' | xargs kill -9

    go run run_agent.go >/dev/null 2>&1 &

    sleep 1m
    go run run_server.go >/dev/null 2>&1 &
}

function wait_server()
{
    until [ "`curl --silent --show-error --connect-timeout 1 -I http://127.0.0.1:1234 | grep 'OK'`" != "" ];
    do
        echo "Waiting server"
        sleep 20s
    done
}

echo "Initializing..."
prepare 
wait_server
echo "Served at :1234
