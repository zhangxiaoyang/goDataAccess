package engine

import (
	"crypto/md5"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/downloader"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/pipeline"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/processer"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/scheduler"
	"github.com/zhangxiaoyang/goDataAccess/spider/plugin"
	"io/ioutil"
	"log"
	"reflect"
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
	startedTime := time.Now()
	defer func() {
		log.Printf("[engine.go] took %s\n", time.Since(startedTime))
	}()

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
			log.Printf("%d goroutines are running, %d tasks left\n", this.resourceManager.Count(), this.scheduler.Count())
			time.Sleep(this.config.GetWaitTime())
		}

		if this.isEmpty() {
			continue
		}

		if ok := this.resourceManager.Alloc(); !ok {
			log.Printf("waiting for resource\n")
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
	if req.Depth > this.config.GetMaxDepth() {
		log.Printf("Skip %s, because depth(%d) > %d\n", req.Url, req.Depth, this.config.GetMaxDepth())
		return
	}

	this.hook(plugin.BeforeDownloaderType, req, this.config)
	resp, err := this.downloader.Download(req, this.config)
	if err != nil {
		log.Printf("downloaded failed(%s)\n", err)
		if this.config.GetMaxRetryTimes() > 0 {
			this.retry(req, resp, err)
		} else {
			this.hook(plugin.AfterDownloaderType, resp, err)
			log.Printf("downloaded failed(retried %d times) %s\n", this.config.GetMaxRetryTimes(), req.Url)
		}
		return
	} else {
		this.hook(plugin.AfterDownloaderType, resp, err)
	}
	log.Printf("downloaded ok %s\n", req.Url)

	var y = common.NewYield()
	this.hook(plugin.BeforeProcesserType, resp, y)
	this.processer.Process(resp, y)
	this.hook(plugin.AfterProcesserType)

	log.Printf("generated %d requests from %s\n", len(y.GetAllRequests()), req.Url)
	for _, r := range y.GetAllRequests() {
		r.Depth = req.Depth + 1
		this.hook(plugin.BeforeSchedulerType, r)
		this.scheduler.Push(r)
		this.hook(plugin.AfterSchedulerType)
	}

	if y.GetMerge() {
		log.Printf("generated %d items(merged) from %s\n", len(y.GetAllItems()), req.Url)
	} else {
		log.Printf("generated %d items from %s\n", len(y.GetAllItems()), req.Url)
	}
	for _, p := range this.pipelines {
		items := y.GetAllItems()
		this.hook(plugin.BeforePipelineType, items, y.GetMerge())
		p.Pipe(items, y.GetMerge())
		this.hook(plugin.AfterPipelineType)
	}
}

func (this *Engine) retry(req *common.Request, resp *common.Response, err error) {
	h := md5.Sum([]byte(req.Url))
	if _, ok := this.retryCache[h]; ok {
		this.retryCache[h]++
	} else {
		this.retryCache[h] = 1
	}
	if this.retryCache[h] <= this.config.GetMaxRetryTimes() {
		log.Printf("retry(%d) %s\n", this.retryCache[h], req.Url)
		this.hook(plugin.BeforeSchedulerType, req)
		this.scheduler.Push(req)
		this.hook(plugin.AfterSchedulerType)
	} else {
		this.hook(plugin.AfterDownloaderType, resp, err)
		delete(this.retryCache, h)
		log.Printf("downloaded failed(retried %d times) %s\n", this.config.GetMaxRetryTimes(), req.Url)
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
	this.hook(plugin.BeforeSchedulerType, r)
	this.scheduler.Push(r)
	this.hook(plugin.AfterSchedulerType)
	return this
}

func (this *Engine) SetStartUrls(urls []string) *Engine {
	for _, url := range urls {
		r := common.NewRequest(url)
		this.hook(plugin.BeforeSchedulerType, r)
		this.scheduler.Push(r)
		this.hook(plugin.AfterSchedulerType)
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

func (this *Engine) SetPipeline(pipeline pipeline.BasePipeline) *Engine {
	this.pipelines = this.pipelines[:0]
	this.pipelines = append(this.pipelines, pipeline)
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

func (this *Engine) GetConfig() *common.Config {
	return this.config
}
