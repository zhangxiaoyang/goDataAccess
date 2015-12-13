package main

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"github.com/zhangxiaoyang/goDataAccess/spider/plugin"
)

func main() {
	engine.
		NewQuickEngine("spider.json").
		GetEngine().
		AddPlugin(plugin.NewProxyPlugin()).
		Start()
}
