package downloader

import (
	"spider/common"
)

type BaseDownloader interface {
	Download(*common.Request, *common.Config) (*common.Response, error)
}
