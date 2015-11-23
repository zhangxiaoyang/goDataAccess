package main

import (
	"spider/engine"
)

func main() {
	var urls = []string{"http://m.qq.com", "http://m.baidu.com"}
	engine.NewEngine("crawl_and_print_to_screen").SetStartUrls(urls).Start()
}
