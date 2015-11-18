package util

import (
	"net/http"
	"strings"
)

type Request struct {
	Request *http.Request
	Url     string
}

func NewRequest(url string) *Request {
	if strings.HasPrefix(url, "http://") {
		req, _ := http.NewRequest("GET", url, nil)
		return &Request{
			Request: req,
			Url:     url,
		}
	}

	panic("Unimplemented protocol for handling url: " + url)
}
