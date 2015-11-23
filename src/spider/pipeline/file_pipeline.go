package pipeline

import (
	"os"
	"spider/common"
)

type FilePipeline struct {
	file     *os.File
	splitter string
}

func NewFilePipeline(file *os.File, splitter string) *FilePipeline {
	return &FilePipeline{file: file, splitter: splitter}
}

func (this *FilePipeline) Pipe(item *common.Item) {
	for k, v := range item.GetAll() {
		this.file.WriteString(k + this.splitter + v + "\n")
	}
}
