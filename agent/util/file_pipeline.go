package util

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"os"
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

	for _, item := range items {
		if html := item.Get("html", ""); html != "" {
			this.file.WriteString(html + "\n")
		}
	}
}
