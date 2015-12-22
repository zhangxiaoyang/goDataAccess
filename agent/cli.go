package main

import (
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/agent/core"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	"time"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  Step 1: fetching free proxies")
	fmt.Println("    go run cli.go [update/u]")
	fmt.Println("  Step 2: picking available proxies which can be used to visit `validateUrl`(response bodies contains `succ`)")
	fmt.Println("    go run cli.go [validate/v] [validateUrl] [succ]")
	fmt.Println("  Step 3: RPC service for spiders")
	fmt.Println("    go run cli.go [serve/s]")
	fmt.Println()
}

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
			break
		}
	}

	if len(os.Args) > 1 {
		dbDir := "db/"
		ruleDir := "rule/"
		agent := core.NewAgent(ruleDir, dbDir)
		op := os.Args[1]
		startTime := time.Now()
		defer func() {
			log.Printf("[cli.go] took %s to complete\n", time.Since(startTime))
		}()

		switch strings.ToLower(op) {
		case "u":
			fallthrough
		case "update":
			if len(os.Args) >= 2 {
				agent.Update()
			}
			return
		case "v":
			fallthrough
		case "validate":
			if len(os.Args) >= 4 {
				validateUrl, succ := os.Args[2], os.Args[3]
				agent.Validate(validateUrl, succ)
				return
			}
		case "s":
			fallthrough
		case "serve":
			rpc.Register(core.NewAgentServer(dbDir))
			rpc.HandleHTTP()
			listen, err := net.Listen("tcp", ":1234")
			if err != nil {
				log.Printf("listen error %s\n", err)
				return
			}
			http.Serve(listen, nil)
			return
		}
	}
	usage()
}
