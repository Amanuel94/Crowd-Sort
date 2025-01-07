package io

import (
	"context"
	"network/modules/dispatcher"
)

type Config[T any] struct {
	ctx  *context.Context
	canc *context.CancelFunc
	key  string
	d    *dispatcher.Dispatcher[T]
}

func NewConfig[T any](ctx *context.Context) *Config[T] {

	canc := context.CancelFunc(func() {})
	return &Config[T]{
		ctx:  ctx,
		canc: &canc,
	}
}

func (cfg *Config[T]) WithCancelFunc(canc *context.CancelFunc) {
	cfg.canc = canc
}

// func (cfg *Config[T]) With(timeOut time.Duration) {
// 	cfg.withTimeout = true
// 	cfg.timeOut = timeOut
// }
