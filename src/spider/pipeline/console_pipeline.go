package pipeline

import (
	"fmt"
	"spider/util"
)

type ConsolePipeline struct{}

func NewConsolePipeline() *ConsolePipeline {
	return &ConsolePipeline{}
}

func (this *ConsolePipeline) Pipe(items *util.Items) {
	for k, v := range items.GetAll() {
		fmt.Println(k + "\t" + v)
	}
}
