package processer

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
)

type BaseProcesser interface {
	Process(*common.Response, *common.Yield)
}
