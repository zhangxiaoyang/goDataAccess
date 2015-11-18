package engine

import (
	"crypto/md5"
	"spider/downloader"
	"spider/pipeline"
	"spider/processer"
	"spider/scheduler"
	"spider/util"
	"time"
)

type Engine struct {
	taskName   string
	processer  processer.BaseProcesser
	downloader downloader.BaseDownloader
	pipelines  []pipeline.BasePipeline
	scheduler  scheduler.BaseScheduler
	config     *util.Config
	count      int
	retryCache map[[md5.Size]byte]int
}

func NewEngine(taskName string) *Engine {
	e := &Engine{taskName: taskName}
	e.config = util.NewConfig()
	e.config.Concurrency = 2
	e.config.PollingTime = 200 * time.Millisecond
	e.config.WaitTime = 200 * time.Millisecond
	e.config.DownloadTimeout = 2 * time.Minute
	e.config.ConnectionTimeout = 2 * time.Second
	e.config.MaxIdleConnsPerHost = 10
	e.config.MaxRetryTimes = 2

	e.count = 0
	e.retryCache = make(map[[md5.Size]byte]int)

	e.scheduler = scheduler.NewScheduler()
	e.downloader = downloader.NewHttpDownloader()
	e.processer = processer.NewLazyProcesser()
	e.pipelines = append(e.pipelines, pipeline.NewConsolePipeline("\t"))
	return e
}

func (this *Engine) SetStartUrls(urls []string) *Engine {
	for _, url := range urls {
		this.scheduler.Push(util.NewRequest(url))
	}
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

func (this *Engine) SetScheduler(scheduler scheduler.BaseScheduler) *Engine {
	this.scheduler = scheduler
	return this
}

func (this *Engine) SetConfig(config *util.Config) *Engine {
	this.config = config
	return this
}

func (this *Engine) Start() {
	for {
		if this.isDone() {
			break
		} else {
			time.Sleep(this.config.WaitTime)
		}

		if this.isFull() {
			time.Sleep(this.config.PollingTime)
			continue
		}

		if this.isEmpty() {
			continue
		}

		req := this.next()
		go func(req *util.Request) {
			this.process(req)
		}(req)
	}
}

func (this *Engine) process(req *util.Request) {
	for _, pipe := range this.pipelines {
		resp, err := this.downloader.Download(req, this.config)

		if err != nil && this.config.MaxRetryTimes > 0 {
			this.retry(req)
			continue
		}

		items := this.processer.Process(resp)
		if items != nil {
			pipe.Pipe(items)
		}
	}

	this.count--
}

func (this *Engine) retry(req *util.Request) {
	h := md5.Sum([]byte(req.Url))
	if _, ok := this.retryCache[h]; ok {
		this.retryCache[h]++
	} else {
		this.retryCache[h] = 1
	}
	if this.retryCache[h] <= this.config.MaxRetryTimes {
		this.scheduler.Push(req)
	} else {
		delete(this.retryCache, h)
	}
}

func (this *Engine) isDone() bool {
	return this.scheduler.Count() == 0 && this.count == 0
}

func (this *Engine) isFull() bool {
	if this.count < this.config.Concurrency {
		return false
	}
	return true
}

func (this *Engine) isEmpty() bool {
	if this.scheduler.Count() == 0 {
		return true
	}
	return false
}

func (this *Engine) next() *util.Request {
	this.count++
	return this.scheduler.Poll()
}
