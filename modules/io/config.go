package io

import (
	"iter"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/shared"
)

// TODO: Add options for other I/O interfaces
type Config[T any] struct {
	items       iter.Seq[*interfaces.Comparable[T]]
	comparators iter.Seq[shared.CmpFunc[T]]
	verbose     int
	bufferSize  int
	cpw         int
}

func NewConfig[T any](items iter.Seq[*interfaces.Comparable[T]], comparators iter.Seq[shared.CmpFunc[T]], verbose int) *Config[T] {

	return &Config[T]{
		items:       items,
		comparators: comparators,
		verbose:     verbose,
		bufferSize:  15,
		cpw:         1000,
	}
}

func (cfg *Config[T]) WithBufferSize(bufferSize int) *Config[T] {
	cfg.bufferSize = bufferSize
	return cfg
}

func (cfg *Config[T]) WithTaskLimit(cpw int) *Config[T] {
	cfg.cpw = cpw
	return cfg
}
