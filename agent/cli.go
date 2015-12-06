package main

import (
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/agent/core"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		dbDir := "db/"
		ruleDir := "rule/"
		agent := core.NewAgent(ruleDir, dbDir)
		op := os.Args[1]

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
		}

		fmt.Println("Usage")
		fmt.Println("go run cli.go [update/u]")
		fmt.Println("go run cli.go [validate/v] [validateUrl] [succ]")
		fmt.Println()
	}
}
