package extractor

import (
	"regexp"
	"spider/common"
)

type TrimFunc func(string) string

type Extractor struct {
	scopeRule string
	rules     map[string]string
	trimFunc  TrimFunc
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (this *Extractor) Extract(body string) []*common.Item {
	items := []*common.Item{}
	scopes := regexp.MustCompile(this.scopeRule).FindAllString(body, -1)
	for _, scope := range scopes {
		item := common.NewItem()
		for key, rule := range this.rules {
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

func (this *Extractor) SetRules(rules map[string]string) *Extractor {
	this.rules = rules
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
