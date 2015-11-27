package pipeline

import (
	"encoding/json"
	"fmt"
	"spider/common"
)

type ConsolePipeline struct{}

func NewConsolePipeline() *ConsolePipeline {
	return &ConsolePipeline{}
}

func (this *ConsolePipeline) Pipe(items []*common.Item) {
	for _, item := range items {
		if json, err := json.Marshal(item.GetAll()); err == nil {
			fmt.Println(string(json))
		}
	}
}
