package common

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
	return &Config{
		Concurrency:         2,
		PollingTime:         200 * time.Millisecond,
		WaitTime:            200 * time.Millisecond,
		DownloadTimeout:     2 * time.Minute,
		ConnectionTimeout:   2 * time.Second,
		MaxIdleConnsPerHost: 10,
		MaxRetryTimes:       2,
	}
}
