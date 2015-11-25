package engine

import (
	"crypto/md5"
	"spider/common"
	"spider/core/downloader"
	"spider/core/pipeline"
	"spider/core/processer"
	"spider/core/scheduler"
	"time"
)

type Engine struct {
	taskName        string
	processer       processer.BaseProcesser
	downloader      downloader.BaseDownloader
	pipelines       []pipeline.BasePipeline
	scheduler       scheduler.BaseScheduler
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
	e.pipelines = append(e.pipelines, pipeline.NewConsolePipeline("\t"))
	return e
}

func (this *Engine) SetStartUrl(url string) *Engine {
	this.scheduler.Push(common.NewRequest(url))
	return this
}

func (this *Engine) SetStartUrls(urls []string) *Engine {
	for _, url := range urls {
		this.scheduler.Push(common.NewRequest(url))
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

func (this *Engine) SetPipelines(pipelines []pipeline.BasePipeline) *Engine {
	this.pipelines = pipelines
	return this
}

func (this *Engine) SetPipeline(pipeline pipeline.BasePipeline) *Engine {
	this.pipelines = this.pipelines[:0]
	this.pipelines = append(this.pipelines, pipeline)
	return this
}

func (this *Engine) SetConfig(config *common.Config) *Engine {
	this.config = config
	this.resourceManager = common.NewResourceManager(config.GetConcurrency())
	return this
}

func (this *Engine) Start() {
	for {
		if this.isDone() {
			break
		} else {
			time.Sleep(this.config.GetWaitTime())
		}

		if this.isEmpty() {
			continue
		}

		if ok := this.resourceManager.Alloc(); !ok {
			time.Sleep(this.config.GetPollingTime())
			continue
		}

		req := this.scheduler.Poll()
		go func(req *common.Request) {
			this.process(req)
			this.resourceManager.Dealloc()
		}(req)
	}
}

func (this *Engine) process(req *common.Request) {
	for _, pipe := range this.pipelines {
		resp, err := this.downloader.Download(req, this.config)

		if err != nil && this.config.GetMaxRetryTimes() > 0 {
			this.retry(req)
			continue
		}

		var y = common.NewYield()
		this.processer.Process(resp, y)
		for _, r := range y.GetAllRequests() {
			this.scheduler.Push(r)
		}
		for _, i := range y.GetAllItems() {
			pipe.Pipe(i)
		}
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
