package engine

import (
	"spider/common"
	"spider/core/pipeline"
	"spider/core/processer"
	"spider/core/scheduler"
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
