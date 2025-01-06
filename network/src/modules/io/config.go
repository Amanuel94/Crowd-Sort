package io

import (
	"network/modules/dispatcher"
	"time"
)

type Config[T any] struct {
	withTimeout bool
	timeOut     time.Duration
	key         string
	d           *dispatcher.Dispatcher[T]
}

func NewConfig[T any](key string) *Config[T] {

	return &Config[T]{
		key:         key,
		withTimeout: false,
		timeOut:     0,
		d:           nil,
	}
}

func (cfg *Config[T]) WithTimeout(timeOut time.Duration) {
	cfg.withTimeout = true
	cfg.timeOut = timeOut
}
