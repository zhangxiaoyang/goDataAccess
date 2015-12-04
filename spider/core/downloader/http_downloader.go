package downloader

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type HttpDownloader struct{}

func NewHttpDownloader() *HttpDownloader {
	return &HttpDownloader{}
}

func (this *HttpDownloader) Download(req *common.Request, config *common.Config) (*common.Response, error) {
	req.Request.Header.Set("User-Agent", config.GetUserAgent())

	client := &http.Client{
		Timeout: 2 * config.GetDownloadTimeout(),
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, config.GetConnectionTimeout())
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			ResponseHeaderTimeout: config.GetDownloadTimeout(),
			MaxIdleConnsPerHost:   config.GetMaxIdleConnsPerHost(),
		},
	}

	resp, err := client.Do(req.Request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body string
	if resp.StatusCode == 200 {
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			reader, _ := gzip.NewReader(resp.Body)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)
				if err != nil && err != io.EOF {
					return nil, err
				}
				if n == 0 {
					break
				}
				body += string(buf)
			}
		default:
			bodyByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			body = string(bodyByte)
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Response StatusCode: %d", resp.StatusCode))
	}

	if config.GetSucc() != "" && !strings.Contains(body, config.GetSucc()) {
		return nil, errors.New(fmt.Sprintf("Invalid response body(succ: %s)", config.GetSucc()))
	}
	return common.NewResponse(resp, req.Url, body), nil
}
