package downloader

import (
	"io/ioutil"
	"net"
	"net/http"
	"spider/common"
)

type HttpDownloader struct{}

func NewHttpDownloader() *HttpDownloader {
	return &HttpDownloader{}
}

func (this *HttpDownloader) Download(req *common.Request, config *common.Config) (*common.Response, error) {
	client := &http.Client{
		Timeout: 2 * config.DownloadTimeout,
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, config.ConnectionTimeout)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			ResponseHeaderTimeout: config.DownloadTimeout,
			MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
		},
	}

	resp, err := client.Do(req.Request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return common.NewResponse(resp, req.Url, body), nil
}
