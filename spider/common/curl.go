package common

import (
	"compress/gzip"
	"errors"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"net/http"
)

type Curl struct {
	client *http.Client
	req    *Request
}

func NewCurl(client *http.Client, req *Request) *Curl {
	return &Curl{client: client, req: req}
}

func (this *Curl) Do() (*Response, error) {
	resp, err := this.client.Do(this.req.Request)
	if err != nil {
		return NewResponse(nil, this.req.Url, ""), err
	}
	defer resp.Body.Close()

	var body string
	if resp.StatusCode == 200 {
		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, _ := gzip.NewReader(resp.Body)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)
				if err != nil && err != io.EOF {
					return NewResponse(nil, this.req.Url, ""), err
				}
				if n == 0 {
					break
				}
				body += string(buf)
			}
		} else {
			contentType := resp.Header.Get("Content-Type")
			newBody, err := charset.NewReader(resp.Body, contentType)
			if err != nil {
				return NewResponse(nil, this.req.Url, ""), err
			}

			bodyByte, err := ioutil.ReadAll(newBody)
			if err != nil {
				return NewResponse(nil, this.req.Url, ""), err
			}
			body = string(bodyByte)
		}
	} else {
		return NewResponse(nil, this.req.Url, ""), errors.New(fmt.Sprintf("Response StatusCode: %d", resp.StatusCode))
	}
	return NewResponse(resp, this.req.Url, body), nil
}
