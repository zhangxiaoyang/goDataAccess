package common

type Yield struct {
	items    []*Item
	requests []*Request
	merge    bool
}

func NewYield() *Yield {
	return &Yield{merge: false}
}

func (this *Yield) AddItem(item *Item) {
	this.items = append(this.items, item)
}

func (this *Yield) AddRequest(req *Request) {
	this.requests = append(this.requests, req)
}

func (this *Yield) SetMerge(merge bool) {
	this.merge = merge
}

func (this *Yield) GetAllItems() []*Item {
	return this.items
}

func (this *Yield) GetAllRequests() []*Request {
	return this.requests
}

func (this *Yield) GetMerge() bool {
	return this.merge
}
