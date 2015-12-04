package downloader

import (
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
)

type BaseDownloader interface {
	Download(*common.Request, *common.Config) (*common.Response, error)
}
