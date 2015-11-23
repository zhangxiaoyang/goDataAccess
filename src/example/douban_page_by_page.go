package main

import (
	"regexp"
	"spider/engine"
	"spider/pipeline"
	"spider/util"
)

type MyProcesser struct {
	baseUrl string
}

func NewMyProcesser(baseUrl string) *MyProcesser {
	return &MyProcesser{baseUrl: baseUrl}
}

func (this *MyProcesser) processTitle(resp *util.Response, y *util.Yield) {
	m := regexp.MustCompile(`(?s)<div class="channel-item">.*?<h3><a href="(.*?)">(.*?)</a>`).FindAllStringSubmatch(resp.Body, -1)
	for _, v := range m {
		item := util.NewItem()
		item.Set("url", v[1])
		item.Set("title", v[2])
		y.AddItem(item)
	}
}

func (this *MyProcesser) processNext(resp *util.Response, y *util.Yield) {
	m := regexp.MustCompile(`(?s)<span class="next">.*?<a href="(.*?)"`).FindStringSubmatch(resp.Body)
	if len(m) > 0 {
		y.AddRequest(util.NewRequest(this.baseUrl + m[1]))
	}
}

func (this *MyProcesser) Process(resp *util.Response, y *util.Yield) {
	this.processTitle(resp, y)
	this.processNext(resp, y)
}

func main() {
	var baseUrl = "http://www.douban.com/group/explore/"

	engine.NewEngine("douban_page_by_page").
		SetStartUrl(baseUrl).
		SetPipeline(pipeline.NewConsolePipeline("\t")).
		SetProcesser(NewMyProcesser(baseUrl)).
		Start()
}
