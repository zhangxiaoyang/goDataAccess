package main

import (
	"bufio"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"github.com/zhangxiaoyang/goDataAccess/spider/plugin"
	"os"
	"strings"
)

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
	statusFile, _ := os.Create("status.url")
	defer statusFile.Close()

	engine.NewEngine("crawling_status").
		SetStartUrls(getUrlsFromFile("test.url")).
		AddPlugin(plugin.NewStatusPlugin(statusFile)).
		Start()
}
