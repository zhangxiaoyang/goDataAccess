package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"spider/common"
	"spider/core/engine"
	"spider/core/pipeline"
	"strings"
	"time"
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
	r := bufio.NewReader(file)
	for i := 0; i < 10; i++ {
		line, _ := r.ReadString('\n')
		urls = append(urls, strings.TrimSpace(line))
	}
	return urls
}

func main() {
	fmt.Println(time.Now())

	file, _ := os.Create("output.txt")
	defer file.Close()

	engine.NewEngine("test_store_in_file").
		SetPipeline(pipeline.NewFilePipeline(file, "\t")).
		SetProcesser(NewMyProcesser()).
		SetStartUrls(getUrlsFromFile("test.url")).
		Start()

	fmt.Println(time.Now())
}
