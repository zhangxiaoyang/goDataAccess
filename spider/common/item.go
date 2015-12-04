package common

type Item struct {
	item map[string]string
}

func NewItem() *Item {
	return &Item{item: make(map[string]string)}
}

func (this *Item) Set(key string, value string) {
	this.item[key] = value
}

func (this *Item) Get(key string, defaultVal string) string {
	if val, ok := this.item[key]; ok {
		return val
	}
	return defaultVal
}

func (this *Item) GetAll() map[string]string {
	return this.item
}
