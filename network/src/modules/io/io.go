// Handles side-effects
package io

import (
	"context"
	"time"
)

type IO[T comparable] struct{
	ctx context.Context
	canc context.CancelFunc
	key IOKey
}

// Initalizes the IO module

func Init[T comparable]() *IO[T] {
	newIO := &IO[T]{}
	newIO.createContext(NewConfig("io"))
	return newIO
}

func InitWithTimeOut[T comparable](timeOut time.Duration) *IO[T] {
	newIO := &IO[T]{}
	cfg := NewConfig("io")
	cfg.WithTimeout(timeOut)
	newIO.createContext(cfg)
	return newIO
}

func (io *IO[T]) createContext(cfg *Config) {

	io.key = *NewIOKey(cfg.key)
	if cfg.withTimeout {
		io.ctx, io.canc = context.WithTimeout(context.Background(), cfg.timeOut)
	} else {
		io.ctx = context.Background()
	}
}

func (io *IO[T]) Read() T {
	return io.ctx.Value(io.key).(T)
}

func (io *IO[T]) Write(key IOKey, value []T) {
	if (io.ctx.Value(key) == nil) {
		io.ctx = context.WithValue(io.ctx, io.key, value)
	} else {
		curr := io.ctx.Value(io.key).([]T)
		io.ctx = context.WithValue(io.ctx, io.key, append(curr, value...))
	}
}

func (io *IO[T]) Close() {
	io.canc()
}



