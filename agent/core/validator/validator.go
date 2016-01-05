package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/pipeline"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

type Validator struct {
	dbPath   string
	rulePath string
}

func NewValidator(dbPath string, rulePath string) *Validator {
	return &Validator{
		dbPath:   dbPath,
		rulePath: rulePath,
	}
}

func (this *Validator) Start() {
	ruleFilePath := path.Join(this.rulePath, "validate.json")

	file, _ := os.Create("xxx.txt")
	defer file.Close()

	engine.
		NewQuickEngine(ruleFilePath).
		GetEngine().
		SetDownloader(NewValidateDownloader()).
		SetPipeline(pipeline.NewFilePipeline(file)).
		Start()
}

type ValidateDownloader struct {
	validateUrl string
	succ        string
}

func NewValidateDownloader() *ValidateDownloader {
	return &ValidateDownloader{}
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
	byteBody, _ := json.Marshal(NewAddr().Deserialize(proxyUrl))
	resp.Body = string(byteBody)
	return resp, nil
}
