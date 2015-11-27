package pipeline

import (
	"encoding/json"
	"os"
	"spider/common"
)

type FilePipeline struct {
	file *os.File
}

func NewFilePipeline(file *os.File) *FilePipeline {
	return &FilePipeline{file: file}
}

func (this *FilePipeline) Pipe(items []*common.Item) {
	for _, item := range items {
		if json, err := json.Marshal(item.GetAll()); err == nil {
			this.file.WriteString(string(json) + "\n")
		}
	}
}
