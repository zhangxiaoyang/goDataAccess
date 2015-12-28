package extractor

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"regexp"
)

type TrimFunc func(string) string

type Extractor struct {
	scopeRule string
	kvRule    map[string]string
	trimFunc  TrimFunc
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (this *Extractor) Extract(resp *common.Response) []*common.Item {
	items := []*common.Item{}
	scopes := regexp.MustCompile(this.scopeRule).FindAllString(resp.Body, -1)
	for _, scope := range scopes {
		item := common.NewItem()
		for key, rule := range this.kvRule {
			if key == "_URL_" {
				item.Set(key, resp.Url)
				continue
			}
			value := regexp.MustCompile(rule).FindStringSubmatch(scope)[1]
			if this.trimFunc != nil {
				item.Set(key, this.trimFunc(value))
			} else {
				item.Set(key, value)
			}
		}
		items = append(items, item)
	}
	return items
}

func (this *Extractor) SetScopeRule(scopeRule string) *Extractor {
	this.scopeRule = scopeRule
	return this
}

func (this *Extractor) SetRules(kvRule map[string]string) *Extractor {
	this.kvRule = kvRule
	return this
}

func (this *Extractor) SetTrimFunc(trimFunc TrimFunc) *Extractor {
	this.trimFunc = trimFunc
	return this
}

func TrimBlank(s string) string {
	return regexp.MustCompile(`[\s]`).ReplaceAllString(s, "")
}

func TrimHtmlTags(s string) string {
	return regexp.MustCompile(`(<.*?>)|(&nbsp;)|([\s])`).ReplaceAllString(s, "")
}
