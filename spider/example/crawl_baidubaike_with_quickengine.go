package main

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
)

func main() {
	engine.NewQuickEngine("crawl_baidubaike_with_quickengine.json").Start()
}
