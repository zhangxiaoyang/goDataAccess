package engine

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"spider/common"
	"spider/core/extractor"
	"spider/core/pipeline"
	"strings"
)

type QuickEngine struct {
	engineFileName string
}

func NewQuickEngine(engineFileName string) *QuickEngine {
	return &QuickEngine{engineFileName: engineFileName}
}

func (this *QuickEngine) Start() {
	config := NewQuickEngineConfig(this.engineFileName)
	file, _ := os.Create(config.OutputFile)
	defer file.Close()

	NewEngine(config.TaskName).
		AddPipeline(pipeline.NewFilePipeline(file)).
		SetProcesser(NewQuickEngineProcesser(config)).
		SetStartUrls(config.StartUrls).
		SetConfig(common.NewConfig().SetConcurrency(config.Concurrency)).
		Start()
}

type QuickEngineConfig struct {
	TaskName    string   `json:"task_name"`
	BaseUrl     string   `json:"base_url"`
	MaxDepth    int      `json:"max_depth"`
	StartUrls   []string `json:"start_urls"`
	Items       Rules    `json:"items"`
	Urls        Rules    `json:"urls"`
	Merge       bool     `json:"merge"`
	OutputFile  string   `json:"output_file"`
	Concurrency int      `json:"concurrency"`
}

type Rules struct {
	ScopeRule string            `json:"scope_rule"`
	Rules     map[string]string `json:"rules"`
	TrimFunc  string            `json:"trim_func"`
}

func NewQuickEngineConfig(fileName string) *QuickEngineConfig {
	config := &QuickEngineConfig{}
	file, _ := ioutil.ReadFile(fileName)
	json.Unmarshal(file, config)
	return config
}

type QuickEngineProcesser struct {
	config *QuickEngineConfig
	depth  int
}

func NewQuickEngineProcesser(config *QuickEngineConfig) *QuickEngineProcesser {
	return &QuickEngineProcesser{config: config, depth: 0}
}

func (this *QuickEngineProcesser) processItems(resp *common.Response, y *common.Yield) {
	var trimFunc extractor.TrimFunc
	switch this.config.Items.TrimFunc {
	case "trim_html_tags":
		trimFunc = extractor.TrimHtmlTags
	case "trim_blank":
		trimFunc = extractor.TrimBlank
	}

	items := extractor.NewExtractor().
		SetScopeRule(this.config.Items.ScopeRule).
		SetRules(this.config.Items.Rules).
		SetTrimFunc(trimFunc).
		Extract(resp.Body)
	for _, item := range items {
		y.AddItem(item)
	}
}

func (this *QuickEngineProcesser) processRequests(resp *common.Response, y *common.Yield) {
	var trimFunc extractor.TrimFunc
	switch this.config.Urls.TrimFunc {
	case "trim_html_tags":
		trimFunc = extractor.TrimHtmlTags
	case "trim_blank":
		trimFunc = extractor.TrimBlank
	}

	items := extractor.NewExtractor().
		SetScopeRule(this.config.Urls.ScopeRule).
		SetRules(this.config.Urls.Rules).
		SetTrimFunc(trimFunc).
		Extract(resp.Body)
	for _, item := range items {
		for _, url := range item.GetAll() {
			if strings.HasPrefix(url, "http://") {
				y.AddRequest(common.NewRequest(url))
			} else {
				y.AddRequest(common.NewRequest(this.config.BaseUrl + url))
			}
		}
	}
}

func (this *QuickEngineProcesser) Process(resp *common.Response, y *common.Yield) {
	if this.depth > this.config.MaxDepth {
		return
	}

	this.depth++
	if this.config.Items.ScopeRule != "" {
		this.processItems(resp, y)
	}
	if this.config.Urls.ScopeRule != "" {
		this.processRequests(resp, y)
	}
	y.SetMerge(this.config.Merge)
}
