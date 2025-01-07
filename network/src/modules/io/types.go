package io

import (
	"context"
)

type IO[T any] struct {
	ctx  context.Context
	canc context.CancelFunc
}

type IOKey struct {
	key string
}
