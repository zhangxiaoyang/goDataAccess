package scheduler

import (
	"spider/common"
)

type BaseScheduler interface {
	Push(*common.Request)
	Poll() *common.Request
	Count() int
	Close()
}
