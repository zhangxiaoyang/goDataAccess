package pipeline

import (
	"os"
	"spider/util"
)

type FilePipeline struct {
	path string
	file *os.File
}

func NewFilePipeline(path string) *FilePipeline {
	file, err := os.Create(path)
	if err != nil {
		return nil
	}

	p := &FilePipeline{path: path, file: file}
	//runtime.SetFinalizer(this.file, this.file.Close)//How to deal with this?
	return p
}

func (this *FilePipeline) Pipe(items *util.Items) {
	for k, v := range items.GetAll() {
		this.file.WriteString(k + "\t" + v + "\n")
	}
}
