package main

import (
	"regexp"
	"spider/common"
	"spider/core/engine"
	"spider/core/pipeline"
	"spider/plugin"
)

type MyProcesser struct{}

func NewMyProcesser() *MyProcesser {
	return &MyProcesser{}
}

func (this *MyProcesser) Process(resp *common.Response, y *common.Yield) {
	m := regexp.MustCompile(`(?s)<h2 class="info">(.*?)</h2>`).FindAllStringSubmatch(resp.Body, -1)
	for _, v := range m {
		item := common.NewItem()
		item.Set("user-agent", v[1])
		y.AddItem(item)
	}
}

func main() {
	engine.NewEngine("test_reseturl_plugin").
		SetStartUrl("http://whatsmyuseragent.com/").
		SetProcesser(NewMyProcesser()).
		AddPlugin(plugin.NewUserAgentPlugin()).
		AddPipeline(pipeline.NewConsolePipeline("\t")).
		Start()
}
