package common

import (
	"net/http"
)

type Response struct {
	Response *http.Response
	Url      string
	Body     string
}

func NewResponse(resp *http.Response, url string, body string) *Response {
	return &Response{
		Response: resp,
		Url:      url,
		Body:     body,
	}
}
