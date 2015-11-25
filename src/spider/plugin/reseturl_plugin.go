package plugin

import (
	"spider/common"
)

type ResetUrlPlugin struct {
	pluginType PluginType
}

func NewResetUrlPlugin() *ResetUrlPlugin {
	return &ResetUrlPlugin{pluginType: PreDownloaderType}
}

func (this *ResetUrlPlugin) Do(params ...interface{}) {
	req := params[0].(*common.Request)
	*req = *common.NewRequest("http://wx.qq.com")
}

func (this *ResetUrlPlugin) GetPluginType() PluginType {
	return this.pluginType
}
