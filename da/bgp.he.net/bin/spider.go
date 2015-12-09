package main

import (
	"compress/gzip"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/da/util"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/rpc"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		log.Printf("lost argument")
		return
	}

	configFilePath, inFilePath, outFilePath := os.Args[1], os.Args[2], os.Args[3]
	outFile, _ := os.Create(outFilePath)
	defer outFile.Close()

	log.Printf("load urls from %s", inFilePath)
	domains := util.LoadUrlsFromFile(inFilePath)

	urls := []string{}
	for _, domain := range domains {
		urls = append(urls, "http://bgp.he.net/dns/"+domain)
	}
	engine.
		NewQuickEngine(configFilePath).
		SetOutputFile(outFile).
		GetEngine().
		SetStartUrls(urls).
		SetDownloader(NewHeDownloader()).
		Start()
}

type HeDownloader struct {
	jar      map[string]*cookiejar.Jar
	isAuthed map[string]bool
}

func NewHeDownloader() *HeDownloader {
	return &HeDownloader{jar: map[string]*cookiejar.Jar{}, isAuthed: map[string]bool{}}
}

func (this *HeDownloader) Download(req *common.Request, config *common.Config) (*common.Response, error) {
	proxyUrl := this.getOneProxy(req.Url)
	if proxyUrl == "" {
		return nil, errors.New(fmt.Sprintf("get proxy failed"))
	}
	log.Printf("use proxy %s\n", proxyUrl)

	this.auth(proxyUrl, config)
	return this.send(proxyUrl, req, config)
}

func (this *HeDownloader) auth(proxyUrl string, config *common.Config) bool {
	if _, ok := this.isAuthed[proxyUrl]; ok {
		log.Printf("have authed %+v\n", this.jar[proxyUrl])
		return true
	}

	var p string
	var i string
	this.jar[proxyUrl], _ = cookiejar.New(nil)
	{
		u := "http://bgp.he.net/i"
		resp, err := this.send(proxyUrl, common.NewRequest(u), config)
		if err != nil {
			log.Printf("auth failed(%s) %s\n", u, err)
			return false
		}
		i = strings.Trim(resp.Response.Header.Get("ETag"), "\"")
	}
	{
		u := "http://bgp.he.net/dns/qq.com"
		req := common.NewRequest(u)
		_, err := this.send(proxyUrl, req, config)
		if err != nil {
			log.Printf("auth failed(%s) %s\n", u, err)
			return false
		}
		path := ""
		for _, c := range this.jar[proxyUrl].Cookies(req.Request.URL) {
			if c.Name == "path" {
				path = c.Value
				break
			}
		}
		decodedPath, _ := url.QueryUnescape(path)
		p = fmt.Sprintf("%x", md5.Sum([]byte(decodedPath)))
	}
	{
		u := "http://bgp.he.net/cc"
		_, err := this.send(proxyUrl, common.NewRequest(u), config)
		if err != nil {
			log.Printf("auth failed(%s) %s\n", u, err)
			return false
		}
	}
	{
		u := "http://bgp.he.net/jc"
		form := url.Values{}
		form.Add("p", p)
		form.Add("i", i)
		req := common.NewRequest(u)
		req.Request, _ = http.NewRequest("POST", u, strings.NewReader(form.Encode()))
		_, err := this.send(proxyUrl, req, config)
		if err != nil {
			log.Printf("auth failed(%s) %s\n", u, err)
			return false
		}
	}
	this.isAuthed[proxyUrl] = true
	log.Printf("auth succeed %+v\n", this.jar[proxyUrl])
	return true
}

func (this *HeDownloader) send(proxyUrl string, req *common.Request, config *common.Config) (*common.Response, error) {
	for key, value := range config.GetHeaders() {
		req.Request.Header.Set(key, value)
	}

	client := &http.Client{
		Jar:     this.jar[proxyUrl],
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

	if _, ok := this.isAuthed[proxyUrl]; ok {
		if config.GetSucc() != "" && !strings.Contains(body, config.GetSucc()) {
			return nil, errors.New(fmt.Sprintf("Invalid response body(succ: %s)", config.GetSucc()))
		}
	}
	return common.NewResponse(resp, req.Url, body), nil
}

func (this *HeDownloader) getOneProxy(url string) string {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	if err != nil {
		log.Printf("dialing error %s\n", err)
		return ""
	}
	defer client.Close()

	var proxy string
	err = client.Call("AgentServer.GetOneProxy", &url, &proxy)
	if err != nil {
		log.Printf("error %s\n", err)
		return ""
	}
	return proxy
}
