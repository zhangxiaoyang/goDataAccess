package downloader

import (
	"spider/util"
)

type BaseDownloader interface {
	Download(*util.Request, *util.Config) (*util.Response, error)
}
