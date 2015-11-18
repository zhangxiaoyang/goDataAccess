package util

import (
	"time"
)

type Config struct {
	Concurrency         int
	PollingTime         time.Duration
	WaitTime            time.Duration
	DownloadTimeout     time.Duration
	ConnectionTimeout   time.Duration
	MaxIdleConnsPerHost int
	MaxRetryTimes       int
}

func NewConfig() *Config {
	return &Config{}
}
