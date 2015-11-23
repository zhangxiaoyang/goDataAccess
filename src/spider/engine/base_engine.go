package engine

import (
	"spider/pipeline"
	"spider/processer"
	"spider/scheduler"
	"spider/common"
)

type BaseEngine interface {
	Start()
	SetStartUrls([]string) *BaseEngine
	SetProcesser(processer.BaseProcesser) *BaseEngine
	SetPipeline(pipeline.BasePipeline) *BaseEngine
	SetPipelines([]pipeline.BasePipeline) *BaseEngine
	SetScheduler(scheduler.BaseScheduler) *BaseEngine
	SetConfig(*common.Config) *BaseEngine
}
