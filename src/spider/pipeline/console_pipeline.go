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

func (this *ConsolePipeline) Pipe(items *util.Items) {
	for k, v := range items.GetAll() {
		fmt.Println(k + this.splitter + v)
	}
}
