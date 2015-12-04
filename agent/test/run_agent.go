package main

import (
	"github.com/zhangxiaoyang/goDataAccess/agent/core/agent"
)

func main() {
	a := agent.NewAgent("../db/", "../rule/", 1000)
	a.Update()
	//a.Validate("http://m.baidu.com", "baidu")
}
