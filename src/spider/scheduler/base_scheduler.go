package scheduler

import (
	"spider/util"
)

type BaseScheduler interface {
	Push(*util.Request)
	Poll() *util.Request
	Count() int
	Close()
}
