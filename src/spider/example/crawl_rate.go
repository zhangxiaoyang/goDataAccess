package main

import (
	"fmt"
	"spider/common"
	"spider/core/engine"
	//"spider/core/pipeline"
	"time"
)

type MyProcesser struct {
	baseUrl string
}

func NewMyProcesser() *MyProcesser {
	return &MyProcesser{}
}

var startingTime = time.Now()
var crawledCount = 0

func (this *MyProcesser) Process(resp *common.Response, y *common.Yield) {
	crawledCount++
	t := float64(time.Now().Sub(startingTime).Minutes())
	if t > 0 {
		fmt.Printf("%1.0f pages/min\n", float64(crawledCount)/t)
	}
}

func genUrls() []string {
	var urls = []string{}
	for i := 0; i < 1000; i++ {
		u := fmt.Sprintf("http://baike.baidu.com/view/%d.htm", i)
		urls = append(urls, u)
	}
	return urls
}

func main() {
	config := common.NewConfig().
		SetConcurrency(1000).
		SetWaitTime(10 * time.Millisecond).
		SetPollingTime(10 * time.Millisecond)

	engine.NewEngine("crawl_rate").
		SetStartUrls(genUrls()).
		SetProcesser(NewMyProcesser()).
		SetConfig(config).
		Start()
}
