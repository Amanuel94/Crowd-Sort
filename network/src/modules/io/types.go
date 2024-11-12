package io

import (
	"fmt"
	"iter"

	"github.com/google/uuid"
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

// wrapper for indexing items

type IndexedItem[T comparable] struct{
	Index uuid.UUID
	Value T
}



func NewIndexedItem[T comparable](value T) *IndexedItem[T] {
	return &IndexedItem[T]{
		Index: uuid.New(),
		Value: value,
	}
}

// prints indexed items
func PrintIndexedItem[T comparable](items iter.Seq[IndexedItem[T]]) {
	for item := range items {
		fmt.Println(item.Index, item.Value)
	}
}