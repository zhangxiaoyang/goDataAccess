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

func (this *ConsolePipeline) Pipe(items []*common.Item, merge bool) {
	if merge {
		merged := []map[string]string{}
		for _, item := range items {
			merged = append(merged, item.GetAll())
		}
		if json, err := json.Marshal(merged); err == nil {
			fmt.Println(string(json))
		}
	} else {
		for _, item := range items {
			if json, err := json.Marshal(item.GetAll()); err == nil {
				fmt.Println(string(json))
			}
		}
	}
}
