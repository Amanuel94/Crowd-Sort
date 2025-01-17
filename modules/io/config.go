package io

import (
	"iter"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/shared"
)

// TODO: Add options for  other I/O interfaces
type Config[T any] struct {
	items       iter.Seq[*interfaces.Comparable[T]]
	comparators iter.Seq[shared.CmpFunc[T]]
}

func NewConfig[T any](items iter.Seq[*interfaces.Comparable[T]], comparators iter.Seq[shared.CmpFunc[T]]) *Config[T] {

	return &Config[T]{
		items:       items,
		comparators: comparators,
	}
}
