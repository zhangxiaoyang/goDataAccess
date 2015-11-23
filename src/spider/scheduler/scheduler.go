package scheduler

import (
	"spider/common"
)

type Scheduler struct {
	queue chan *common.Request
}

func NewScheduler() *Scheduler {
	return &Scheduler{queue: make(chan *common.Request, 1024)}
}

func (this *Scheduler) Push(req *common.Request) {
	this.queue <- req
}

func (this *Scheduler) Poll() *common.Request {
	if len(this.queue) == 0 {
		return nil
	}
	return <-this.queue
}

func (this *Scheduler) Count() int {
	return len(this.queue)
}

func (this *Scheduler) Close() {
	close(this.queue)
}
