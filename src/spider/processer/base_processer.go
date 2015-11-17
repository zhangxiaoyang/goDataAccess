package processer

import (
	"spider/util"
)

type BaseProcesser interface {
	Process(*util.Response) *util.Items
}
