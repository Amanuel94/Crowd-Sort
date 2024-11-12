package io

import (
	"fmt"
	"iter"

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
type IndexedItem[T constraints.Ordered] struct{
	Index uuid.UUID
	Value T
}

// Implements the comparable interface
func (item IndexedItem[T]) Compare(other IndexedItem[T]) int {
	if item.Value > other.Value {
		return 1
	} else if item.Value < other.Value {
		return -1
	}
	return 0
}



func NewIndexedItem[T constraints.Ordered](value T) *IndexedItem[T] {
	return &IndexedItem[T]{
		Index: uuid.New(),
		Value: value,
	}
}

// prints indexed items
func PrintIndexedItem[T constraints.Ordered](items iter.Seq[IndexedItem[T]]) {
	buff := []IndexedItem[T]{}
	for item := range items {
		fmt.Println(item.Index, item.Value)
		buff = append(buff, item)
	}
}