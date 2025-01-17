package io

import (
	"sync"

	"github.com/Amanuel94/crowdsort/modules/dispatcher"
)

type IO[T any] struct {
	d         *dispatcher.Dispatcher[T]
	wg        *sync.WaitGroup
	msgBuffer []interface{}
}
