package main

import (
	"spider/core/engine"
	"spider/plugin"
)

func main() {
	var urls = []string{"http://m.qq.com", "http://m.baidu.com"}
	engine.NewEngine("test_reseturl_plugin").SetStartUrls(urls).AddPlugin(plugin.NewResetUrlPlugin()).Start()
}
