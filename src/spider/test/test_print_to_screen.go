package main

import (
	"spider/core/engine"
	"spider/core/pipeline"
)

func main() {
	var urls = []string{"http://m.qq.com", "http://m.baidu.com"}
	engine.NewEngine("test_print_to_screen").AddPipeline(pipeline.NewConsolePipeline("\t")).SetStartUrls(urls).Start()
}
