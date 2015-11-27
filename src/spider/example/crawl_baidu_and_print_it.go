package main

import (
	"spider/core/engine"
	"spider/core/pipeline"
)

func main() {
	url := "http://m.baidu.com"
	engine.NewEngine("crawl_baidu_and_print_it").AddPipeline(pipeline.NewConsolePipeline()).SetStartUrl(url).Start()
}
