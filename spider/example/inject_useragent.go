package main

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/pipeline"
	"github.com/zhangxiaoyang/goDataAccess/spider/plugin"
	"regexp"
)

type MyProcesser struct{}

func NewMyProcesser() *MyProcesser {
	return &MyProcesser{}
}

func (this *MyProcesser) Process(resp *common.Response, y *common.Yield) {
	m := regexp.MustCompile(`(?s)<div id="ua_string">.*?</span>(.*?)</div>`).FindAllStringSubmatch(resp.Body, -1)
	for _, v := range m {
		item := common.NewItem()
		item.Set("user-agent", v[1])
		y.AddItem(item)
	}
}

func main() {
	engine.NewEngine("inject_useragent").
		SetStartUrl("http://my-user-agent.com/").
		SetProcesser(NewMyProcesser()).
		AddPlugin(plugin.NewUserAgentPlugin()).
		AddPipeline(pipeline.NewConsolePipeline()).
		SetConfig(common.NewConfig().SetHeaders(map[string]string{"User-Agent": "golang spider"})).
		Start()
}
