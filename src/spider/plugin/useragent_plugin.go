package plugin

import (
	"spider/common"
)

type UserAgentPlugin struct{}

func NewUserAgentPlugin() *UserAgentPlugin {
	return &UserAgentPlugin{}
}

func (this *UserAgentPlugin) Do(pluginType PluginType, params ...interface{}) {
	if pluginType == PreDownloaderType {
		req := params[0].(*common.Request)
		req.Request.Header.Set("User-Agent", "golang spider")
	}
}
