package pipeline

import (
	"encoding/json"
	"os"
	"spider/common"
	"sync"
)

type FilePipeline struct {
	file *os.File
	lock *sync.Mutex
}

func NewFilePipeline(file *os.File) *FilePipeline {
	return &FilePipeline{file: file, lock: &sync.Mutex{}}
}

func (this *FilePipeline) Pipe(items []*common.Item, merge bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if merge {
		merged := []map[string]string{}
		for _, item := range items {
			merged = append(merged, item.GetAll())
		}
		if json, err := json.Marshal(merged); err == nil {
			this.file.WriteString(string(json) + "\n")
		}
	} else {
		for _, item := range items {
			if json, err := json.Marshal(item.GetAll()); err == nil {
				this.file.WriteString(string(json) + "\n")
			}
		}
	}
}
