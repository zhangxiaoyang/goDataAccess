package pipeline

import (
	"fmt"
	"spider/util"
)

type ConsolePipeline struct {
	splitter string
}

func NewConsolePipeline(splitter string) *ConsolePipeline {
	return &ConsolePipeline{splitter: splitter}
}

func (this *ConsolePipeline) Pipe(item *util.Item) {
	for k, v := range item.GetAll() {
		fmt.Println(k + this.splitter + v)
	}
}
