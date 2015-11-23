package processer

import (
	"spider/util"
)

type LazyProcesser struct{}

func NewLazyProcesser() *LazyProcesser {
	return &LazyProcesser{}
}

func (this *LazyProcesser) Process(resp *util.Response, y *util.Yield) {
	y.AddItem(func() *util.Item {
		item := util.NewItem()
		item.Set("html", resp.Body)
		return item
	}())
}
