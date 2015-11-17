package pipeline

import (
	"spider/util"
)

type ConsolePipeline struct{}

func NewConsolePipeline() *ConsolePipeline {
	return &ConsolePipeline{}
}

func (this *ConsolePipeline) Pipe(items *util.Items) {
	innerItems := items.GetAll()
	for k, v := range innerItems {
		println(k + "\t" + v)
	}
}
