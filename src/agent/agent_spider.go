package main

import (
	"spider/core/engine"
)

func main() {
	engine.NewQuickEngine("agent_spider.json").Start()
}
