package pipeline

import (
	"spider/common"
)

type NullPipeline struct{}

func NewNullPipeline() *NullPipeline {
	return &NullPipeline{}
}

func (this *NullPipeline) Pipe(item *common.Item) {
}
