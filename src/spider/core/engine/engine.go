package engine

import (
	"crypto/md5"
	"io/ioutil"
	"log"
	"reflect"
	"spider/common"
	"spider/core/downloader"
	"spider/core/pipeline"
	"spider/core/processer"
	"spider/core/scheduler"
	"spider/plugin"
	"time"
)

type Engine struct {
	taskName        string
	scheduler       scheduler.BaseScheduler
	downloader      downloader.BaseDownloader
	processer       processer.BaseProcesser
	pipelines       []pipeline.BasePipeline
	plugins         []plugin.BasePlugin
	config          *common.Config
	resourceManager *common.ResourceManager
	retryCache      map[[md5.Size]byte]int
}

func NewEngine(taskName string) *Engine {
	e := &Engine{taskName: taskName}
	e.config = common.NewConfig()

	e.resourceManager = common.NewResourceManager(e.config.GetConcurrency())
	e.retryCache = make(map[[md5.Size]byte]int)

	e.scheduler = scheduler.NewScheduler()
	e.downloader = downloader.NewHttpDownloader()
	e.processer = processer.NewLazyProcesser()
	return e
}

func (this *Engine) Start() {
	log.Printf("started\n")
	log.Printf("config: %+v", this.config)
	log.Printf("scheduler: %s", reflect.TypeOf(this.scheduler).Elem().Name())
	log.Printf("downloader: %s", reflect.TypeOf(this.downloader).Elem().Name())
	log.Printf("processer: %s", reflect.TypeOf(this.processer).Elem().Name())
	for _, p := range this.pipelines {
		log.Printf("pipeline: %s", reflect.TypeOf(p).Elem().Name())
	}
	for _, p := range this.plugins {
		log.Printf("plugins: %s", reflect.TypeOf(p).Elem().Name())
	}

	for {
		if this.isDone() {
			log.Printf("finished\n")
			break
		} else {
			log.Printf("%d goroutines are running\n", this.resourceManager.Count())
			time.Sleep(this.config.GetWaitTime())
		}

		if this.isEmpty() {
			continue
		}

		if ok := this.resourceManager.Alloc(); !ok {
			log.Printf("blocked because no more resource\n")
			time.Sleep(this.config.GetPollingTime())
			continue
		}

		req := this.scheduler.Poll()
		go func(req *common.Request) {
			this.process(req)
			this.resourceManager.Free()
		}(req)
	}
}

func (this *Engine) process(req *common.Request) {
	this.hook(plugin.PreDownloaderType, req)
	resp, err := this.downloader.Download(req, this.config)

	if err != nil && this.config.GetMaxRetryTimes() > 0 {
		this.retry(req)
		return
	}
	log.Printf("downloaded ok %s\n", req.Url)

	var y = common.NewYield()
	this.hook(plugin.PreProcesserType, req)
	this.processer.Process(resp, y)

	log.Printf("generated %d requests from %s\n", len(y.GetAllRequests()), req.Url)
	for _, r := range y.GetAllRequests() {
		this.hook(plugin.PreSchedulerType, r)
		this.scheduler.Push(r)
	}

	log.Printf("generated %d items from %s\n", len(y.GetAllItems()), req.Url)
	for _, p := range this.pipelines {
		items := y.GetAllItems()
		this.hook(plugin.PrePipelineType, items)
		p.Pipe(items, y.GetMerge())
	}
}

func (this *Engine) retry(req *common.Request) {
	h := md5.Sum([]byte(req.Url))
	if _, ok := this.retryCache[h]; ok {
		this.retryCache[h]++
	} else {
		this.retryCache[h] = 1
	}
	if this.retryCache[h] <= this.config.GetMaxRetryTimes() {
		log.Printf("retry(%d) %s\n", this.retryCache[h], req.Url)
		this.hook(plugin.PreSchedulerType, req)
		this.scheduler.Push(req)
	} else {
		delete(this.retryCache, h)
	}
}

func (this *Engine) isDone() bool {
	return this.scheduler.Count() == 0 && this.resourceManager.Count() == 0
}

func (this *Engine) isEmpty() bool {
	if this.scheduler.Count() == 0 {
		return true
	}
	return false
}

func (this *Engine) hook(pluginType plugin.PluginType, args ...interface{}) {
	for _, p := range this.plugins {
		p.Do(pluginType, args...)
	}
}

func (this *Engine) SetStartUrl(url string) *Engine {
	r := common.NewRequest(url)
	this.hook(plugin.PreSchedulerType, r)
	this.scheduler.Push(r)
	return this
}

func (this *Engine) SetStartUrls(urls []string) *Engine {
	for _, url := range urls {
		r := common.NewRequest(url)
		this.hook(plugin.PreSchedulerType, r)
		this.scheduler.Push(r)
	}
	return this
}

func (this *Engine) SetScheduler(scheduler scheduler.BaseScheduler) *Engine {
	this.scheduler = scheduler
	return this
}

func (this *Engine) SetDownloader(downloader downloader.BaseDownloader) *Engine {
	this.downloader = downloader
	return this
}

func (this *Engine) SetProcesser(processer processer.BaseProcesser) *Engine {
	this.processer = processer
	return this
}

func (this *Engine) AddPipeline(pipeline pipeline.BasePipeline) *Engine {
	this.pipelines = append(this.pipelines, pipeline)
	return this
}

func (this *Engine) AddPlugin(plugin plugin.BasePlugin) *Engine {
	this.plugins = append(this.plugins, plugin)
	return this
}

func (this *Engine) SetConfig(config *common.Config) *Engine {
	this.config = config
	this.resourceManager = common.NewResourceManager(config.GetConcurrency())
	if !this.config.GetLogging() {
		log.SetOutput(ioutil.Discard)
	}
	return this
}
