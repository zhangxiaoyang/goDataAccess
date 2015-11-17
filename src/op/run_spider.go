package main

import (
	"spider/engine"
)

func main() {
	var urls = []string{"http://m.qq.com", "http://m.baidu.com"}
	engine.NewEngine("test").SetStartUrls(urls).Start()
}
