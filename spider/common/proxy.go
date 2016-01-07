package common

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Proxy struct{}

func NewProxy() *Proxy {
	return &Proxy{}
}

type JsonResp struct {
	Level  int    `json:"level"`
	Num    int    `json:"num"`
	Result string `json:"result"`
}

func (this *Proxy) GetOneProxy(url string) (string, error) {
	resp, err := http.Get("http://127.0.0.1:1234/getOneProxy")
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
	return jsonResp.Result, nil
}
