package scheduler

import (
	"container/list"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"sync"
)

type Scheduler struct {
	lock  *sync.Mutex
	queue *list.List
}

func NewScheduler() *Scheduler {
	s := &Scheduler{queue: list.New()}
	s.lock = &sync.Mutex{}
	return s
}

func (this *Scheduler) Push(req *common.Request) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.queue.PushBack(req)
}

func (this *Scheduler) Poll() *common.Request {
	this.lock.Lock()
	defer this.lock.Unlock()
	e := this.queue.Front()
	if e != nil {
		this.queue.Remove(e)
		return e.Value.(*common.Request)
	}
	return nil
}

func (this *Scheduler) Count() int {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.queue.Len()
}
