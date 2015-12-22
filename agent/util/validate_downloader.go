package util

import (
	"errors"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type ValidateDownloader struct {
	validateUrl string
	succ        string
}

func NewValidateDownloader(validateUrl string, succ string) *ValidateDownloader {
	return &ValidateDownloader{validateUrl: validateUrl, succ: succ}
}

func (this *ValidateDownloader) Download(req *common.Request, config *common.Config) (*common.Response, error) {
	proxyUrl := req.Url
	req.Request, _ = http.NewRequest("GET", this.validateUrl, nil)
	for key, value := range config.GetHeaders() {
		req.Request.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: 2 * config.GetDownloadTimeout(),
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{Host: proxyUrl}),
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

	resp, err := common.NewCurl(client, req).Do()
	if err != nil {
		fmt.Printf("curl %s error %s\n", req.Url, err)
		return nil, err
	}

	if config.GetSucc() != "" && !strings.Contains(resp.Body, config.GetSucc()) {
		return nil, errors.New(fmt.Sprintf("Invalid response body(succ: %s)", config.GetSucc()))
	}
	resp.Body = proxyUrl
	return resp, nil
}
