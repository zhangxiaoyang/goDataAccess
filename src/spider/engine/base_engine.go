package engine

import (
	"spider/util"
)

type BaseEngine interface {
	Start()
	isFull() bool
	isDone() bool
	next() *util.Request
}
