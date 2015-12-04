package main

import (
	"agent/core/agent"
)

func main() {
	a := agent.NewAgent("../db/", "../rule/", 1000)
	//a.Update()
	a.Validate("http://m.baidu.com", "baidu")
}
