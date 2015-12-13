package main

import (
	"crypto/md5"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/da/util"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"github.com/zhangxiaoyang/goDataAccess/spider/plugin"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

var gAuth = &Auth{Jar: map[string]*cookiejar.Jar{}, IsAuthed: map[string]bool{}}
var gConfig *common.Config

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

	e := engine.
		NewQuickEngine(configFilePath).
		SetOutputFile(outFile).
		GetEngine()
	gConfig = e.GetConfig()
	e.SetStartUrls(urls).
		AddPlugin(plugin.NewProxyPlugin()).
		AddPlugin(plugin.NewCookiePlugin(GetCookieFunc))
	e.Start()
}

type Auth struct {
	Jar      map[string]*cookiejar.Jar
	IsAuthed map[string]bool
}

func GetCookieFunc(req *common.Request) (*cookiejar.Jar, error) {
	if _, ok := gAuth.IsAuthed[req.ProxyUrl]; ok {
		log.Printf("have authed %+v\n", gAuth.Jar[req.ProxyUrl])
		return gAuth.Jar[req.ProxyUrl], nil
	}

	baseUrl := "http://bgp.he.net"
	transport := &http.Transport{
		Proxy: http.ProxyURL(&url.URL{Host: req.ProxyUrl}),
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, gConfig.GetConnectionTimeout())
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		ResponseHeaderTimeout: gConfig.GetDownloadTimeout(),
		MaxIdleConnsPerHost:   gConfig.GetMaxIdleConnsPerHost(),
	}
	gAuth.Jar[req.ProxyUrl], _ = cookiejar.New(nil)
	client := &http.Client{
		Jar:       gAuth.Jar[req.ProxyUrl],
		Timeout:   2 * gConfig.GetDownloadTimeout(),
		Transport: transport,
	}

	var p string
	var i string
	{
		u := baseUrl + "/i"
		resp, err := common.NewCurl(client, common.NewRequest(u)).Do()
		if err != nil {
			log.Printf("1. auth failed(%s) %s\n", u, err)
			return nil, err
		}
		i = strings.Trim(resp.Response.Header.Get("ETag"), "\"")
	}
	{
		u := baseUrl + "/dns/qq.com"
		_, err := common.NewCurl(client, common.NewRequest(u)).Do()
		if err != nil {
			log.Printf("2. auth failed(%s) %s\n", u, err)
			return nil, err
		}
		path := ""
		for _, c := range gAuth.Jar[req.ProxyUrl].Cookies(req.Request.URL) {
			if c.Name == "path" {
				path = c.Value
				break
			}
		}
		decodedPath, _ := url.QueryUnescape(path)
		p = fmt.Sprintf("%x", md5.Sum([]byte(decodedPath)))
	}
	{
		u := baseUrl + "/cc"
		_, err := common.NewCurl(client, common.NewRequest(u)).Do()
		if err != nil {
			log.Printf("3. auth failed(%s) %s\n", u, err)
			return nil, err
		}
	}
	{
		u := baseUrl + "/jc"
		form := url.Values{}
		form.Add("p", p)
		form.Add("i", i)
		r := common.NewRequest(u)
		r.Request, _ = http.NewRequest("POST", u, strings.NewReader(form.Encode()))
		_, err := common.NewCurl(client, r).Do()
		if err != nil {
			log.Printf("4.auth failed(%s) %s\n", u, err)
			return nil, err
		}
	}
	gAuth.IsAuthed[req.ProxyUrl] = true
	log.Printf("auth succeed %+v\n", gAuth.Jar[req.ProxyUrl])
	return gAuth.Jar[req.ProxyUrl], nil
}
