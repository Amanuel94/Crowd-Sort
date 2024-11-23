package shared

import (
	"network/shared/interfaces"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

// Wrapper for indexing items
type IndexedItem[T any] struct {
	Index uuid.UUID
	Value interfaces.Comparable[T]
}

func (item IndexedItem[T]) GetIndex() any {
	return item.Index
}

func (item IndexedItem[T]) GetValue() T {
	return item.Value.GetValue()
}

func (item IndexedItem[T]) Compare(other interfaces.Comparable[T]) int {
	return item.Value.Compare(other)
}

func NewIndexedItem[T any](value interfaces.Comparable[T]) interfaces.Comparable[T] {
	return IndexedItem[T]{
		Index: uuid.New(),
		Value: value,
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
