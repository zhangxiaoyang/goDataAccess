package core

import (
	"bufio"
	"container/ring"
	"encoding/json"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/agent/util"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type AgentServer struct {
	dbDir string
	rings map[string]*ring.Ring
}

func NewAgentServer(dbDir string) *AgentServer {
	server := &AgentServer{
		dbDir: dbDir,
		rings: map[string]*ring.Ring{},
	}
	server.readAllProxyToRing()
	return server
}

func (this *AgentServer) GetOneProxy(url *string, proxy *string) error {
	domain := util.ExtractDomain(*url)
	if _, ok := this.rings[domain]; !ok {
		domain = "*"
	}
	*proxy = fmt.Sprintf("%s", this.rings[domain].Value)
	this.rings[domain] = this.rings[domain].Next()
	log.Printf("handle url %s using rule: %s, proxy: %s\n", *url, domain, *proxy)
	return nil
}

func (this *AgentServer) readAllProxyToRing() {
	if fileInfos, err := ioutil.ReadDir(this.dbDir); err == nil {
		for _, f := range fileInfos {
			validProxyPath := path.Join(this.dbDir, f.Name())
			var domain string
			if this.isValidProxy(f.Name()) {
				domain = strings.Trim(strings.Trim(f.Name(), "validjson"), ".")

			} else if this.isCandidateProxy(f.Name()) {
				domain = "*"

			} else {
				continue
			}

			file, err := os.Open(validProxyPath)
			if err != nil {
				log.Printf("error %s", err)
				return
			}
			defer file.Close()

			r := bufio.NewReader(file)
			proxies := []string{}
			for {
				line, err := r.ReadString('\n')
				if err != nil || err == io.EOF {
					break
				}
				proxies = append(proxies, strings.TrimSpace(line))
			}
			this.rings[domain] = ring.New(len(proxies))
			for _, proxy := range proxies {
				addr := util.NewAddr()
				json.Unmarshal([]byte(proxy), addr)
				this.rings[domain].Value = addr.Serialize()
				this.rings[domain] = this.rings[domain].Next()
			}
		}
	}
}

func (this *AgentServer) isValidProxy(fileName string) bool {
	if strings.HasPrefix(fileName, "valid.") {
		return true
	}
	return false
}

func (this *AgentServer) isCandidateProxy(fileName string) bool {
	if fileName == "candidate.json" {
		return true
	}
	return false
}
