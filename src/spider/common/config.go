package common

import (
	"time"
)

type Config struct {
	concurrency         int
	pollingTime         time.Duration
	waitTime            time.Duration
	downloadTimeout     time.Duration
	connectionTimeout   time.Duration
	maxIdleConnsPerHost int
	maxRetryTimes       int
	logging             bool
}

func NewConfig() *Config {
	return &Config{
		concurrency:         2,
		pollingTime:         200 * time.Millisecond,
		waitTime:            200 * time.Millisecond,
		downloadTimeout:     2 * time.Minute,
		connectionTimeout:   2 * time.Second,
		maxIdleConnsPerHost: 10,
		maxRetryTimes:       2,
		logging:             true,
	}
}

func (this *Config) SetConcurrency(concurrency int) *Config {
	this.concurrency = concurrency
	return this
}

func (this *Config) SetPollingTime(pollingTime time.Duration) *Config {
	this.pollingTime = pollingTime
	return this
}

func (this *Config) SetWaitTime(waitTime time.Duration) *Config {
	this.waitTime = waitTime
	return this
}

func (this *Config) SetDownloadTimeout(downloadTimeout time.Duration) *Config {
	this.downloadTimeout = downloadTimeout
	return this
}

func (this *Config) SetConnectionTimeout(connectionTimeout time.Duration) *Config {
	this.connectionTimeout = connectionTimeout
	return this
}

func (this *Config) SetMaxIdleConnsPerHost(maxIdleConnsPerHost int) *Config {
	this.maxIdleConnsPerHost = maxIdleConnsPerHost
	return this
}

func (this *Config) SetMaxRetryTimes(maxRetryTimes int) *Config {
	this.maxRetryTimes = maxRetryTimes
	return this
}

func (this *Config) SetLogging(logging bool) *Config {
	this.logging = logging
	return this
}

func (this *Config) GetConcurrency() int {
	return this.concurrency
}

func (this *Config) GetPollingTime() time.Duration {
	return this.pollingTime
}

func (this *Config) GetWaitTime() time.Duration {
	return this.waitTime
}

func (this *Config) GetDownloadTimeout() time.Duration {
	return this.downloadTimeout
}

func (this *Config) GetConnectionTimeout() time.Duration {
	return this.connectionTimeout
}

func (this *Config) GetMaxIdleConnsPerHost() int {
	return this.maxIdleConnsPerHost
}

func (this *Config) GetMaxRetryTimes() int {
	return this.maxRetryTimes
}

func (this *Config) GetLogging() bool {
	return this.logging
}
