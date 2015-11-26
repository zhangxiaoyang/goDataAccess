package engine

import (
	"spider/common"
	"spider/core/downloader"
	"spider/core/pipeline"
	"spider/core/processer"
	"spider/core/scheduler"
	"spider/plugin"
)

type BaseEngine interface {
	Start()
	SetStartUrl(string) *BaseEngine
	SetStartUrls([]string) *BaseEngine
	SetScheduler(scheduler.BaseScheduler) *BaseEngine
	SetDownloader(downloader.BaseDownloader) *BaseEngine
	SetProcesser(processer.BaseProcesser) *BaseEngine
	AddPipeline(pipeline.BasePipeline) *BaseEngine
	AddPlugin(plugin.BasePlugin) *BaseEngine
	SetConfig(*common.Config) *BaseEngine
}
