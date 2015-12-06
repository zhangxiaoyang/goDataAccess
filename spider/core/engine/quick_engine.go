package engine

import (
	"encoding/json"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/extractor"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/pipeline"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

type QuickEngine struct {
	quickEngineConfigPath string
	file                  *os.File
}

func NewQuickEngine(quickEngineConfigPath string) *QuickEngine {
	return &QuickEngine{quickEngineConfigPath: quickEngineConfigPath}
}

func (this *QuickEngine) GetEngine() *Engine {
	c := NewQuickEngineConfig(this.quickEngineConfigPath)
	if c.OutputFile != "" {
		this.file, _ = os.Create(c.OutputFile)
		defer this.file.Close()
	}

	return NewEngine(c.TaskName).
		AddPipeline(pipeline.NewFilePipeline(this.file)).
		SetProcesser(NewQuickEngineProcesser(c)).
		SetStartUrls(c.StartUrls).
		SetConfig(c.ToCommonConfig())
}

func (this *QuickEngine) Start() {
	this.GetEngine().Start()
}

func (this *QuickEngine) SetOutputFile(file *os.File) *QuickEngine {
	this.file = file
	return this
}

type QuickEngineConfig struct {
	TaskName    string   `json:"task_name"`
	BaseUrl     string   `json:"base_url"`
	StartUrls   []string `json:"start_urls"`
	ItemRule    _Rules   `json:"item_rule"`
	RequestRule _Rules   `json:"request_rule"`
	Merge       bool     `json:"merge"`
	OutputFile  string   `json:"output_file"`
	Config      _Config  `json:"config"`
}

type _Rules struct {
	ScopeRule string            `json:"scope_rule"`
	KVRule    map[string]string `json:"kv_rule"`
	TrimFunc  string            `json:"trim_func"`
}

type _Config struct {
	Concurrency         int    `json:"concurrency"`
	PollingTime         string `json:"polling_time"`
	WaitTime            string `json:"wait_time"`
	DownloadTimeout     string `json:"download_timeout"`
	ConnectionTimeout   string `json:"connection_timeout"`
	MaxIdleConnsPerHost int    `json:"max_idle_conns_per_host"`
	MaxRetryTimes       int    `json:"max_retry_times"`
	Logging             bool   `json:"logging"`
	UserAgent           string `json:"user_agent"`
	Succ                string `json:"succ"`
}

func NewQuickEngineConfig(fileName string) *QuickEngineConfig {
	c := &QuickEngineConfig{}
	file, _ := ioutil.ReadFile(fileName)
	json.Unmarshal(file, c)

	t := reflect.TypeOf(&c.Config)
	v := reflect.ValueOf(&c.Config)
	config := common.NewConfig()
	for i := 0; i < t.Elem().NumField(); i++ {
		field := t.Elem().Field(i)
		value := v.Elem().FieldByName(field.Name)
		funcName := "Get" + field.Name
		switch value.Type().Kind() {
		case reflect.Int:
			if value.Int() == 0 {
				defaultValue := reflect.ValueOf(config).MethodByName(funcName).Call([]reflect.Value{})
				value.SetInt(defaultValue[0].Int())
			}
		case reflect.String:
			if value.String() == "" {
				defaultValue := reflect.ValueOf(config).MethodByName(funcName).Call([]reflect.Value{})
				value.SetString(fmt.Sprintf("%s", defaultValue[0].Interface()))
				//value.SetString(defaultValue[0].Interface().(time.Duration).String())
			}
		case reflect.Bool:
			if value.Bool() == false {
				defaultValue := reflect.ValueOf(config).MethodByName(funcName).Call([]reflect.Value{})
				value.SetBool(defaultValue[0].Bool())
			}
		}
	}
	return c
}

func (this *QuickEngineConfig) stringToDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func (this *QuickEngineConfig) ToCommonConfig() *common.Config {
	e := common.NewConfig().SetConcurrency(this.Config.Concurrency).
		SetPollingTime(this.stringToDuration(this.Config.PollingTime)).
		SetWaitTime(this.stringToDuration(this.Config.WaitTime)).
		SetDownloadTimeout(this.stringToDuration(this.Config.DownloadTimeout)).
		SetConnectionTimeout(this.stringToDuration(this.Config.ConnectionTimeout)).
		SetMaxIdleConnsPerHost(this.Config.MaxIdleConnsPerHost).
		SetMaxRetryTimes(this.Config.MaxRetryTimes).
		SetLogging(this.Config.Logging).
		SetUserAgent(this.Config.UserAgent).
		SetSucc(this.Config.Succ)

	return e
}

type QuickEngineProcesser struct {
	config *QuickEngineConfig
}

func NewQuickEngineProcesser(config *QuickEngineConfig) *QuickEngineProcesser {
	return &QuickEngineProcesser{config: config}
}

func (this *QuickEngineProcesser) processItems(resp *common.Response, y *common.Yield) {
	var TrimFunc extractor.TrimFunc
	switch this.config.ItemRule.TrimFunc {
	case "trim_html_tags":
		TrimFunc = extractor.TrimHtmlTags
	case "trim_blank":
		TrimFunc = extractor.TrimBlank
	}

	items := extractor.NewExtractor().
		SetScopeRule(this.config.ItemRule.ScopeRule).
		SetRules(this.config.ItemRule.KVRule).
		SetTrimFunc(TrimFunc).
		Extract(resp.Body)
	for _, item := range items {
		y.AddItem(item)
	}
}

func (this *QuickEngineProcesser) processRequests(resp *common.Response, y *common.Yield) {
	var TrimFunc extractor.TrimFunc
	switch this.config.RequestRule.TrimFunc {
	case "trim_html_tags":
		TrimFunc = extractor.TrimHtmlTags
	case "trim_blank":
		TrimFunc = extractor.TrimBlank
	}

	items := extractor.NewExtractor().
		SetScopeRule(this.config.RequestRule.ScopeRule).
		SetRules(this.config.RequestRule.KVRule).
		SetTrimFunc(TrimFunc).
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
	common.Try(func() {
		if this.config.ItemRule.ScopeRule != "" {
			this.processItems(resp, y)
		}
		if this.config.RequestRule.ScopeRule != "" {
			this.processRequests(resp, y)
		}
		y.SetMerge(this.config.Merge)
	}, func(e interface{}) {
		log.Printf("pannic %s\n", e)
	})
}
