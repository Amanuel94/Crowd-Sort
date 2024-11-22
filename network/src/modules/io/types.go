package io

import (
	"context"
	"fmt"
	"iter"
	"network/shared"
	"network/shared/interfaces"

	"golang.org/x/exp/constraints"
)

type IO[T any] struct {
	ctx  context.Context
	canc context.CancelFunc
	key  IOKey
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

// Wrapper for constrained types
type OrderedType[T constraints.Ordered] struct {
	Value T
}

func (o *OrderedType[T]) GetIndex() any {
	return o.Value
}

func (o *OrderedType[T]) GetValue() T {
	return o.Value
}

func (o *OrderedType[T]) Compare(other interfaces.Comparable[T]) int {
	if o.Value < other.GetValue() {
		return -1
	} else if o.Value > other.GetValue() {
		return 1
	}
	return 0
}

func NewInt[T constraints.Integer](value T) interfaces.Comparable[T] {
	return &OrderedType[T]{
		Value: value,
	}
}

// prints indexed items
func PrintIndexedItem[T any](items iter.Seq[shared.IndexedItem[T]]) {
	buff := []shared.IndexedItem[T]{}
	for item := range items {
		fmt.Println(item.Index, item.Value)
		buff = append(buff, item)
	}
}
