package downloader

import (
	"errors"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type HttpDownloader struct{}

func NewHttpDownloader() *HttpDownloader {
	return &HttpDownloader{}
}

func (this *HttpDownloader) Download(req *common.Request, config *common.Config) (*common.Response, error) {
	for key, value := range config.GetHeaders() {
		req.Request.Header.Set(key, value)
	}

	transport := &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, config.GetConnectionTimeout())
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		ResponseHeaderTimeout: config.GetDownloadTimeout(),
		MaxIdleConnsPerHost:   config.GetMaxIdleConnsPerHost(),
	}
	client := &http.Client{
		Timeout:   2 * config.GetDownloadTimeout(),
		Transport: transport,
	}
	if req.ProxyUrl != "" {
		transport.Proxy = http.ProxyURL(&url.URL{Host: req.ProxyUrl})
	}
	if req.Jar != nil {
		client.Jar = req.Jar
	}
	if req.Error != nil {
		return nil, req.Error
	}

	resp, err := common.NewCurl(client, req).Do()
	if err != nil {
		fmt.Printf("curl %s error %s\n", req.Url, err)
		return nil, err
	}

	if config.GetSucc() != "" && !strings.Contains(resp.Body, config.GetSucc()) {
		return nil, errors.New(fmt.Sprintf("Invalid response body(succ: %s)", config.GetSucc()))
	}
	return resp, nil
}
