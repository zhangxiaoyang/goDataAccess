package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"spider/common"
	"spider/core/engine"
	"spider/core/extractor"
	"spider/core/pipeline"
)

type Spider struct {
	TaskName      string            `json:"task_name"`
	StartUrls     []string          `json:"start_urls"`
	ItemScopeRule string            `json:"item_scope_rule"`
	ItemRules     map[string]string `json:"item_rules"`
	TrimFunc      string            `json:"trim_func"`
	OutputFile    string            `json:"output_file"`
}

type MyProcesser struct {
	spider *Spider
}

func NewMyProcesser(spider *Spider) *MyProcesser {
	return &MyProcesser{spider: spider}
}

func (this *MyProcesser) Process(resp *common.Response, y *common.Yield) {
	var trimFunc extractor.TrimFunc
	switch this.spider.TrimFunc {
	case "trim_html_tags":
		trimFunc = extractor.TrimHtmlTags
	case "trim_blank":
		trimFunc = extractor.TrimBlank
	}

	items := extractor.NewExtractor().
		SetItemScopeRule(this.spider.ItemScopeRule).
		SetItemRules(this.spider.ItemRules).
		SetTrimFunc(trimFunc).
		Extract(resp.Body)

	for _, item := range items {
		y.AddItem(item)
	}
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
