package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"spider/common"
	"spider/core/engine"
	"spider/core/extractor"
	"spider/core/pipeline"
	"strings"
)

type Rules struct {
	ScopeRule string            `json:"scope_rule"`
	Rules     map[string]string `json:"rules"`
	TrimFunc  string            `json:"trim_func"`
}

type Spider struct {
	TaskName   string   `json:"task_name"`
	BaseUrl    string   `json:"base_url"`
	MaxDepth   int      `json:"max_depth"`
	StartUrls  []string `json:"start_urls"`
	Items      Rules    `json:"items"`
	Urls       Rules    `json:"urls"`
	Merge      bool     `json:"merge"`
	OutputFile string   `json:"output_file"`
}

type MyProcesser struct {
	spider *Spider
	depth  int
}

func NewMyProcesser(spider *Spider) *MyProcesser {
	return &MyProcesser{spider: spider, depth: 0}
}

func (this *MyProcesser) processItems(resp *common.Response, y *common.Yield) {
	var trimFunc extractor.TrimFunc
	switch this.spider.Items.TrimFunc {
	case "trim_html_tags":
		trimFunc = extractor.TrimHtmlTags
	case "trim_blank":
		trimFunc = extractor.TrimBlank
	}

	items := extractor.NewExtractor().
		SetScopeRule(this.spider.Items.ScopeRule).
		SetRules(this.spider.Items.Rules).
		SetTrimFunc(trimFunc).
		Extract(resp.Body)
	for _, item := range items {
		y.AddItem(item)
	}
}

func (this *MyProcesser) processRequests(resp *common.Response, y *common.Yield) {
	var trimFunc extractor.TrimFunc
	switch this.spider.Urls.TrimFunc {
	case "trim_html_tags":
		trimFunc = extractor.TrimHtmlTags
	case "trim_blank":
		trimFunc = extractor.TrimBlank
	}

	items := extractor.NewExtractor().
		SetScopeRule(this.spider.Urls.ScopeRule).
		SetRules(this.spider.Urls.Rules).
		SetTrimFunc(trimFunc).
		Extract(resp.Body)
	for _, item := range items {
		for _, url := range item.GetAll() {
			if strings.HasPrefix(url, "http://") {
				y.AddRequest(common.NewRequest(url))
			} else {
				y.AddRequest(common.NewRequest(this.spider.BaseUrl + url))
			}
		}
	}
}

func (this *MyProcesser) Process(resp *common.Response, y *common.Yield) {
	if this.depth > this.spider.MaxDepth {
		return
	}

	this.depth++
	this.processItems(resp, y)
	this.processRequests(resp, y)
	y.SetMerge(this.spider.Merge)
}

func NewSpider(fileName string) *Spider {
	s := Spider{}
	file, _ := ioutil.ReadFile("spider.json")
	json.Unmarshal(file, &s)
	return &s
}

func main() {
	s := NewSpider("spider.json")

	file, _ := os.Create(s.OutputFile)
	defer file.Close()

	engine.NewEngine(s.TaskName).
		AddPipeline(pipeline.NewFilePipeline(file)).
		SetProcesser(NewMyProcesser(s)).
		SetStartUrls(s.StartUrls).
		Start()
}
