package scheduler

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
)

type BaseScheduler interface {
	Push(*common.Request)
	Poll() *common.Request
	Count() int
}
