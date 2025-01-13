package io

import (
	"context"

	"github.com/amanuel94/crowdsort/modules/dispatcher"
)

type IO[T any] struct {
	ctx       context.Context
	canc      context.CancelFunc
	d         *dispatcher.Dispatcher[T]
	msgBuffer []interface{}
}

type IOKey struct {
	Key string
}
