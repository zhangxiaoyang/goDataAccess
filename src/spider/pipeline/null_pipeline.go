package pipeline

import (
	"spider/util"
)

type NullPipeline struct{}

func NewNullPipeline() *NullPipeline {
	return &NullPipeline{}
}

func (this *NullPipeline) Pipe(item *util.Item) {
}
