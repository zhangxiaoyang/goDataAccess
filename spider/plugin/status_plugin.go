package plugin

import (
	"container/list"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"os"
	"sync"
)

type StatusPlugin struct {
	file      *os.File
	lock      *sync.Mutex
	succQueue *list.List
	failQueue *list.List
}

func NewStatusPlugin(file *os.File) *StatusPlugin {
	return &StatusPlugin{
		file:      file,
		succQueue: list.New(),
		failQueue: list.New(),
		lock:      &sync.Mutex{},
	}
}

func (this *StatusPlugin) Do(pluginType PluginType, args ...interface{}) {
	if pluginType == AfterDownloaderType {
		resp := args[0].(*common.Response)
		err, _ := args[1].(error)
		status := ""
		if err != nil {
			this.failQueue.PushBack(resp.Url)
			status = "fail"
		} else {
			status = "succ"
			this.succQueue.PushBack(resp.Url)
		}
		this.file.WriteString(fmt.Sprintf("%s\t%s\n", status, resp.Url))
	}
}
