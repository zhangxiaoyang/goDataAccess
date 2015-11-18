package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"spider/engine"
	"spider/pipeline"
	"spider/util"
	"strings"
	"time"
)

type MyProcesser struct{}

func NewMyProcesser() *MyProcesser {
	return &MyProcesser{}
}

func (this *MyProcesser) Process(resp *util.Response) *util.Items {
	items := util.NewItems()
	items.Set("url", resp.Url)
	items.Set("title", func(s string) string {
		m := regexp.MustCompile(`<title>(.*?)</title>`).FindStringSubmatch(s)
		if len(m) > 0 {
			return m[0]
		}
		return ""
	}(resp.Body))
	return items
}

func getUrlsFromFile(fileName string) []string {
	var urls = []string{}
	file, _ := os.Open("test.url")
	r := bufio.NewReader(file)
	for i := 0; i < 10; i++ {
		line, _ := r.ReadString('\n')
		urls = append(urls, strings.TrimSpace(line))
	}
	return urls
}
func main() {
	fmt.Println(time.Now())

	engine.NewEngine("crawl_and_store_file_output.txt").
		SetPipeline(pipeline.NewFilePipeline("output.txt")).
		SetProcesser(NewMyProcesser()).
		SetStartUrls(getUrlsFromFile("test.url")).
		Start()

	fmt.Println(time.Now())
}
