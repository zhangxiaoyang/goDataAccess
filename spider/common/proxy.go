package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Proxy struct{}

func NewProxy() *Proxy {
	return &Proxy{}
}

type JsonResp struct {
	Num    int    `json:"num"`
	Result string `json:"result"`
}

func (this *Proxy) GetOneProxy(u string) (string, error) {
	resp, err := http.Get("http://127.0.0.1:1234/getOneProxy?url=" + url.QueryEscape(u))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var jsonResp JsonResp
	json.Unmarshal(body, &jsonResp)
	if jsonResp.Num <= 0 {
		return "", errors.New("Agent donot have enough proxies!")
	}
	return jsonResp.Result, nil
}
