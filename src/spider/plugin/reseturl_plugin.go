package plugin

import (
	"spider/common"
)

type ResetUrlPlugin struct{}

func NewResetUrlPlugin() *ResetUrlPlugin {
	return &ResetUrlPlugin{}
}

func (this *ResetUrlPlugin) Do(pluginType PluginType, params ...interface{}) {
	if pluginType == PreDownloaderType {
		req := params[0].(*common.Request)
		*req = *common.NewRequest("http://wx.qq.com")
	}
}
