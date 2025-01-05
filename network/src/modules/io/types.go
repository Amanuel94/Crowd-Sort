package io

import (
	"context"
	"network/modules/dispatcher"
)

type IO[T any] struct {
	ctx  context.Context
	canc context.CancelFunc
	key  IOKey
	d    dispatcher.Dispatcher[T]
}

// identifies an io stream
type IOKey struct {
	key string
}

func NewIOKey(key string) *IOKey {
	return &IOKey{
		key: key,
	}
}
