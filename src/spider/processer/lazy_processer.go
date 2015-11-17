package processer

import (
	"spider/util"
)

type LazyProcesser struct{}

func NewLazyProcesser() *LazyProcesser {
	return &LazyProcesser{}
}

func (this *LazyProcesser) Process(resp *util.Response) *util.Items {
	items := util.NewItems()
	items.Set("html", resp.Body)
	return items
}
