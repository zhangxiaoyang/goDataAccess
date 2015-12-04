package main

import (
	"bufio"
	"os"
	"regexp"
	"spider/common"
	"spider/core/engine"
	"spider/core/pipeline"
	"strings"
)

type MyProcesser struct{}

func NewMyProcesser() *MyProcesser {
	return &MyProcesser{}
}

func (this *MyProcesser) Process(resp *common.Response, y *common.Yield) {
	y.AddItem(func() *common.Item {
		item := common.NewItem()
		item.Set("url", resp.Url)
		item.Set("title", func() string {
			m := regexp.MustCompile(`<title>(.*?)</title>`).FindStringSubmatch(resp.Body)
			if len(m) > 0 {
				return m[1]
			}
			return ""
		}())
		return item
	}())
}

func getUrlsFromFile(fileName string) []string {
	var urls = []string{}
	file, _ := os.Open(fileName)
	defer file.Close()

	r := bufio.NewReader(file)
	for i := 0; i < 10; i++ {
		line, _ := r.ReadString('\n')
		urls = append(urls, strings.TrimSpace(line))
	}
	return urls
}

func main() {
	file, _ := os.Create("crawl_baidubaike_and_store_it_output.txt")
	defer file.Close()

	engine.NewEngine("crawl_baidubaike_and_store_it").
		AddPipeline(pipeline.NewFilePipeline(file)).
		SetProcesser(NewMyProcesser()).
		SetStartUrls(getUrlsFromFile("test.url")).
		Start()
}
