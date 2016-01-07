package main

import (
	"github.com/zhangxiaoyang/goDataAccess/agent/core/server"
	"log"
	"os"
	"strconv"
)

func main() {
	var port int
	for i, _ := range os.Args {
		var err error
		port, err = strconv.Atoi(os.Args[i])
		if err != nil {
			port = 1234
		}
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	const dbPath = "db/"
	server.NewServer(dbPath, port).Start()
}
