package main

import (
	"github.com/zhangxiaoyang/goDataAccess/agent/core"
)

func main() {
	dbDir := "../db/"
	ruleDir := "../rule/"
	a := core.NewAgent(ruleDir, dbDir)
	a.Update()
	//a.Validate("http://m.baidu.com", "baidu")
}
