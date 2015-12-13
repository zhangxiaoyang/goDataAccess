package plugin

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
)

type ProxyPlugin struct{}

func NewProxyPlugin() *ProxyPlugin {
	return &ProxyPlugin{}
}

func (this *ProxyPlugin) Do(pluginType PluginType, args ...interface{}) {
	if pluginType == PreDownloaderType {
		req := args[0].(*common.Request)
		req.ProxyUrl = common.NewProxy().GetOneProxy(req.Url)
	}
}
