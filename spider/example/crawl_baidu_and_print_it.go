package main

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/pipeline"
)

func main() {
	url := "http://m.baidu.com"
	engine.NewEngine("crawl_baidu_and_print_it").AddPipeline(pipeline.NewConsolePipeline()).SetStartUrl(url).Start()
}
