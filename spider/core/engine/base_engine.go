package engine

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/downloader"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/pipeline"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/processer"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/scheduler"
	"github.com/zhangxiaoyang/goDataAccess/spider/plugin"
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
