// Handles side-effects
package io

import (
	"context"
)

type IO[T comparable] struct{
	ctx context.Context
	canc context.CancelFunc
	key IOKey
}

// Initalizes the IO module
func (io *IO[T]) Init(cfg *Config) {

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

func (io *IO[T]) Write(value T) {
	io.ctx = context.WithValue(io.ctx, io.key, value)
}

func (io *IO[T]) Close() {
	io.canc()
}



