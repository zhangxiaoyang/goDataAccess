package common

import (
	"log"
	"net/rpc"
)

type Proxy struct {
	address string
}

func NewProxy() *Proxy {
	return &Proxy{address: "127.0.0.1:1234"}
}

func (this *Proxy) GetOneProxy(url string) string {
	client, err := rpc.DialHTTP("tcp", this.address)
	if err != nil {
		log.Printf("dialing error %s\n", err)
		return ""
	}
	defer client.Close()

	var proxy string
	err = client.Call("AgentServer.GetOneProxy", &url, &proxy)
	if err != nil {
		log.Printf("proxy error %s\n", err)
		return ""
	}
	return proxy
}
