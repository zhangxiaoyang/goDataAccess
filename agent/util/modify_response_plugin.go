package util

import (
	"encoding/json"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/plugin"
	"log"
	"regexp"
)

type ModifyResponsePlugin struct{}

func NewModifyResponsePlugin() *ModifyResponsePlugin {
	return &ModifyResponsePlugin{}
}

func (this *ModifyResponsePlugin) Do(pluginType plugin.PluginType, args ...interface{}) {
	if pluginType == plugin.AfterDownloaderType {
		resp := args[0].(*common.Response)
		req := args[2].(*common.Request)

		meta := map[string]string{
			"url":    req.Url,
			"domain": this.extractDomain(req.Url),
			"proxy":  req.ProxyUrl,
		}
		metaStr, err := json.Marshal(meta)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body = fmt.Sprintf("%s\n%s", metaStr, resp.Body)
	}
}

func (this *ModifyResponsePlugin) extractDomain(url string) string {
	return regexp.MustCompile(`http[s]?://([\w\-\.]+)`).FindStringSubmatch(url)[1]
}
