package main

import (
	"github.com/zhangxiaoyang/goDataAccess/agent/core/agent"
	"io"
	"log"
	"os"
)

func main() {
	for i, _ := range os.Args {
		if os.Args[i] == "--log" && i+1 < len(os.Args) && os.Args[i+1] != "" {
			fileName := os.Args[i+1]
			f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			defer f.Close()
			log.SetOutput(io.MultiWriter(f, os.Stdout))
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			break
		}
	}

	const (
		dbPath   = "db/"
		rulePath = "rule/"
	)
	agent.NewAgent(dbPath, rulePath).Start()
}
