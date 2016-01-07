package util

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/plugin"
	"strconv"
)

type AddLevelPlugin struct {
	level int
}

func NewAddLevelPlugin(level int) *AddLevelPlugin {
	return &AddLevelPlugin{level: level}
}

func (this *AddLevelPlugin) Do(pluginType plugin.PluginType, args ...interface{}) {
	if pluginType == plugin.BeforePipelineType {
		items := args[0].([]*common.Item)
		for _, item := range items {
			item.Set("level", strconv.Itoa(this.level))
		}
	}
}
