package util

type Yield struct {
	items    []*Item
	requests []*Request
}

func NewYield() *Yield {
	return &Yield{}
}

func (this *Yield) AddItem(item *Item) {
	this.items = append(this.items, item)
}

func (this *Yield) AddRequest(req *Request) {
	this.requests = append(this.requests, req)
}

func (this *Yield) GetAllItems() []*Item {
	return this.items
}

func (this *Yield) GetAllRequests() []*Request {
	return this.requests
}
