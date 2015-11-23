package pipeline

import (
	"spider/util"
)

type BasePipeline interface {
	Pipe(*util.Item)
}
