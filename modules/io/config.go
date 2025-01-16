package io

import (
	"iter"

	"github.com/Amanuel94/crowdsort/interfaces"
)

type Config[T any] struct {
	items       iter.Seq[*interfaces.Comparable[T]]
	comparators iter.Seq[func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)]
}

func NewConfig[T any](items iter.Seq[*interfaces.Comparable[T]], comparators iter.Seq[func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)]) *Config[T] {

	return &Config[T]{
		items:       items,
		comparators: comparators,
	}
}
