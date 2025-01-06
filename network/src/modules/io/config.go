package io

import "time"

type Config struct {
	withTimeout bool
	timeOut time.Duration
	key string
}


func NewConfig(key string) *Config {
	return &Config{
		key: key,
	}
}

func (cfg *Config) WithTimeout(timeOut time.Duration) {
	cfg.withTimeout = true
	cfg.timeOut = timeOut
}
