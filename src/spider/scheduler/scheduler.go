package scheduler

import (
	"spider/util"
)

type Scheduler struct {
	queue chan *util.Request
}

func NewScheduler() *Scheduler {
	s := &Scheduler{queue: make(chan *util.Request, 1024)}
	return s
}

func (this *Scheduler) Push(req *util.Request) {
	this.queue <- req
}

func (this *Scheduler) Poll() *util.Request {
	if len(this.queue) == 0 {
		return nil
	}
	return <-this.queue
}

func (this *Scheduler) Count() int {
	return len(this.queue)
}
