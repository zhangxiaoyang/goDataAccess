package common

import (
	"sync"
)

type ResourceManager struct {
	lock     *sync.Mutex
	count    int
	capacity int
}

func NewResourceManager(capacity int) *ResourceManager {
	return &ResourceManager{lock: &sync.Mutex{}, count: 0, capacity: capacity}
}

func (this *ResourceManager) Alloc() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.count < this.capacity {
		this.count++
		return true
	}
	return false
}

func (this *ResourceManager) Free() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.count--
}

func (this *ResourceManager) Count() int {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.count
}
