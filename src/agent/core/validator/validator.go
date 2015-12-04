package validator

import (
	"agent/util"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Validator struct {
	validateUrl string
	succ        string
	userAgent   string
	timeout     time.Duration
}

func NewValidator() *Validator {
	return &Validator{
		validateUrl: "http://www.baidu.com",
		succ:        "baidu",
		userAgent:   "Mozilla/5.0 (X11; Linux i686) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.35 Safari/537.36",
		timeout:     10 * time.Second,
	}
}

func (this *Validator) Validate(addr *util.Addr) bool {
	req, _ := http.NewRequest("GET", this.validateUrl, nil)
	req.Header.Set("User-Agent", this.userAgent)
	agentUrl := &url.URL{Host: addr.Serialize()}

	client := &http.Client{
		Timeout: this.timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(agentUrl),
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, this.timeout)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			ResponseHeaderTimeout: this.timeout,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
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
					return false
				}
				if n == 0 {
					break
				}
				body += string(buf)
			}
		default:
			bodyByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return false
			}
			body = string(bodyByte)
		}
	} else {
		return false
	}

	if this.succ != "" && !strings.Contains(body, this.succ) {
		return false
	}
	return true
}

func (this *Validator) SetValidateUrl(validateUrl string) *Validator {
	this.validateUrl = validateUrl
	return this
}

func (this *Validator) SetSucc(succ string) *Validator {
	this.succ = succ
	return this
}
