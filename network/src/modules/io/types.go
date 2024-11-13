package io

import (
	"fmt"
	"iter"
	"network/interfaces"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

// identifies an io stream
type IOKey struct {	
	key string
}

func NewIOKey(key string) *IOKey {
	return &IOKey{
		key: key,
	}
}

// Wrapper for indexing items
type IndexedItem[T any] struct {
	Index uuid.UUID
	Value interfaces.Comparable[T]
}

func  (item *IndexedItem[T]) GetValue() T {
	return item.Value.GetValue()
}

func (item *IndexedItem[T]) Compare(other *IndexedItem[T]) int {
	return item.Value.Compare(other.Value)
}

func NewIndexedItem[T any](value interfaces.Comparable[T]) *IndexedItem[T] {
	return &IndexedItem[T]{
		Index: uuid.New(),
		Value: value,
	}
}

// Wrapper for constriants types
type OrderedType[T constraints.Ordered] struct {
	Value T
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
	
func NewInt64(value int64) interfaces.Comparable[int64] {
	return &OrderedType[int64]{
		Value: value,
	}
}

// prints indexed items
func PrintIndexedItem[T any](items iter.Seq[IndexedItem[T]]) {
	buff := []IndexedItem[T]{}
	for item := range items {
		fmt.Println(item.Index, item.Value)
		buff = append(buff, item)
	}
}