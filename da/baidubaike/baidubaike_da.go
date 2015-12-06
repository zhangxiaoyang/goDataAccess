package main

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/core/downloader"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
)

func main() {
	engine.
		NewQuickEngine("baidubaike_da.json").
		GetEngine().
		SetDownloader(downloader.NewAgentDownloader()).
		Start()
}
