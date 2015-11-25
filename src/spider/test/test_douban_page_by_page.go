package main

import (
	"regexp"
	"spider/common"
	"spider/core/engine"
	"spider/core/pipeline"
)

type MyProcesser struct {
	baseUrl string
}

func NewMyProcesser(baseUrl string) *MyProcesser {
	return &MyProcesser{baseUrl: baseUrl}
}

func (this *MyProcesser) processTitle(resp *common.Response, y *common.Yield) {
	m := regexp.MustCompile(`(?s)<div class="channel-item">.*?<h3><a href="(.*?)">(.*?)</a>`).FindAllStringSubmatch(resp.Body, -1)
	for _, v := range m {
		item := common.NewItem()
		item.Set("url", v[1])
		item.Set("title", v[2])
		y.AddItem(item)
	}
}

func (this *MyProcesser) processNext(resp *common.Response, y *common.Yield) {
	m := regexp.MustCompile(`(?s)<span class="next">.*?<a href="(.*?)"`).FindStringSubmatch(resp.Body)
	if len(m) > 0 {
		y.AddRequest(common.NewRequest(this.baseUrl + m[1]))
	}
}

func (this *MyProcesser) Process(resp *common.Response, y *common.Yield) {
	this.processTitle(resp, y)
	this.processNext(resp, y)
}

func main() {
	var baseUrl = "http://www.douban.com/group/explore/"

	engine.NewEngine("test_douban_page_by_page").
		SetStartUrl(baseUrl).
		SetPipeline(pipeline.NewConsolePipeline("\t")).
		SetProcesser(NewMyProcesser(baseUrl)).
		Start()
}
