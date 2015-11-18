package downloader

import (
	"io/ioutil"
	"net"
	"net/http"
	"spider/util"
)

type HttpDownloader struct{}

func NewHttpDownloader() *HttpDownloader {
	return &HttpDownloader{}
}

func (this *HttpDownloader) Download(req *util.Request, config *util.Config) (*util.Response, error) {
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

	return util.NewResponse(resp, req.Url, body), nil
}
