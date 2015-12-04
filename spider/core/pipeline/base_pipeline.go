package pipeline

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
)

type BasePipeline interface {
	Pipe([]*common.Item, bool)
}
