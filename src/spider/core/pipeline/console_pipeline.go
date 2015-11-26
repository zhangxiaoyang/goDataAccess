package pipeline

import (
	"fmt"
	"spider/common"
)

type ConsolePipeline struct {
	splitter string
}

func NewConsolePipeline(splitter string) *ConsolePipeline {
	return &ConsolePipeline{splitter: splitter}
}

func (this *ConsolePipeline) Pipe(items []*common.Item) {
	for _, item := range items {
		for k, v := range item.GetAll() {
			fmt.Println(k + this.splitter + v)
		}
	}
}
