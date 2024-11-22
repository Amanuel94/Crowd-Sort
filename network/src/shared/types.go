package shared

import (
	"network/shared/interfaces"

	"github.com/google/uuid"
)

// Wrapper for indexing items
type IndexedItem[T any] struct {
	Index uuid.UUID
	Value interfaces.Comparable[T]
}

func (item *IndexedItem[T]) GetIndex() uuid.UUID {
	return item.Index
}

func (item *IndexedItem[T]) GetValue() T {
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
