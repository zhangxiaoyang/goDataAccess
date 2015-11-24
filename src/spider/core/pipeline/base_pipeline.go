package pipeline

import (
	"spider/common"
)

type BasePipeline interface {
	Pipe(*common.Item)
}
