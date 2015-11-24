package main

import (
	"bufio"
	"os"
	"spider/common"
	"spider/core/engine"
	"spider/core/pipeline"
	"strings"
	"time"
)

type MyProcesser struct {
	baseUrl string
}

func NewMyProcesser() *MyProcesser {
	return &MyProcesser{}
}

func (this *MyProcesser) Process(resp *common.Response, y *common.Yield) {
	crawledCount++
	t := int(time.Now().Sub(startingTime).Seconds())
	if t > 0 {
		println(crawledCount / t)
	}
}

var startingTime = time.Now()
var crawledCount = 0

func getUrlsFromFile(fileName string) []string {
	var urls = []string{}
	file, _ := os.Open(fileName)
	r := bufio.NewReader(file)
	for i := 0; i < 1000; i++ {
		line, _ := r.ReadString('\n')
		urls = append(urls, strings.TrimSpace(line))
	}
	return urls
}

func main() {
	engine.NewEngine("test_speed").
		SetStartUrls(getUrlsFromFile("/home/zhangyang/baidubaike.url")).
		SetPipeline(pipeline.NewNullPipeline()).
		SetProcesser(NewMyProcesser()).
		SetConfig(common.NewConfig().SetConcurrency(10)).
		Start()
}
