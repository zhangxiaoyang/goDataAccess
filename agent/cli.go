package main

import (
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/agent/core"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	"time"
)

func usage() {
	fmt.Println("Usage")
	fmt.Println("go run cli.go [update/u]")
	fmt.Println("go run cli.go [validate/v] [validateUrl] [succ]")
	fmt.Println("go run cli.go [serve/s]")
	fmt.Println()
}

func main() {
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
			if len(os.Args) == 2 {
				agent.Update()
			}
			return
		case "v":
			fallthrough
		case "validate":
			if len(os.Args) == 4 {
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
