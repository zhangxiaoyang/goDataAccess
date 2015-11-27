package main

import (
	"spider/core/engine"
)

func main() {
	engine.NewQuickEngine("crawl_baidubaike_with_quickengine.json").Start()
}
