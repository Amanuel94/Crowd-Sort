package io

import (
	"context"
	"iter"

	"github.com/amanuel94/crowdsort/interfaces"
)

type Config[T any] struct {
	ctx         *context.Context
	canc        *context.CancelFunc
	items       iter.Seq[*interfaces.Comparable[T]]
	comparators iter.Seq[func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)]
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
