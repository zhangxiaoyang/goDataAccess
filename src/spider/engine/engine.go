package engine

import (
	"spider/downloader"
	"spider/pipeline"
	"spider/processer"
	"spider/scheduler"
	"spider/util"
	"time"
)

type Engine struct {
	TaskName   string
	Processer  processer.BaseProcesser
	Downloader downloader.BaseDownloader
	Pipelines  []pipeline.BasePipeline
	Scheduler  scheduler.BaseScheduler
	Config     *util.Config
	count      int
}

func NewEngine(TaskName string) *Engine {
	return (&Engine{TaskName: TaskName}).init()
}

func (this *Engine) SetStartUrls(urls []string) *Engine {
	for _, url := range urls {
		this.Scheduler.Push(util.NewRequest(url))
	}

	return this
}

func (this *Engine) Start() {
	for {
		if this.isDone() {
			break
		} else {
			time.Sleep(this.Config.WaitTime)
		}

		if this.isFull() {
			time.Sleep(this.Config.PollingTime)
			continue
		}

		req := this.next()
		if req == nil {
			continue
		} else {
			this.count++
		}

		that := this
		go func(req *util.Request) {
			for _, p := range this.Pipelines {
				p.Pipe(this.Processer.Process(this.Downloader.Download(req, this.Config)))
			}
			that.count--
		}(req)
	}
}

func (this *Engine) isFull() bool {
	if this.count < this.Config.Concurrency {
		return false
	}
	return true
}

func (this *Engine) isDone() bool {
	return this.Scheduler.Count() == 0 && this.count == 0
}

func (this *Engine) next() *util.Request {
	return this.Scheduler.Poll()
}

func (this *Engine) init() *Engine {
	this.Config = util.NewConfig()
	this.Config.Concurrency = 2
	this.Config.PollingTime = 200 * time.Millisecond
	this.Config.WaitTime = 200 * time.Millisecond
	this.Config.DownloadTimeout = 2 * time.Minute
	this.Config.ConnectionTimeout = 2 * time.Second
	this.Config.MaxIdleConnsPerHost = 10

	this.count = 0

	this.Scheduler = scheduler.NewScheduler()
	this.Downloader = downloader.NewHttpDownloader()
	this.Processer = processer.NewLazyProcesser()
	this.Pipelines = append(this.Pipelines, pipeline.NewConsolePipeline())

	return this
}
