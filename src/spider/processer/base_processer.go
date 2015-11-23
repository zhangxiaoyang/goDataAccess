package processer

import (
	"spider/common"
)

type BaseProcesser interface {
	Process(*common.Response, *common.Yield)
}
