package util

import (
	"strings"
)

type Addr struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func NewAddr() *Addr {
	return &Addr{}
}

func (this *Addr) Serialize() string {
	return this.IP + ":" + this.Port
}

func (this *Addr) Deserialize(s string) *Addr {
	tmp := strings.Split(s, ":")[:2]
	this.IP, this.Port = tmp[0], tmp[1]
	return this
}
