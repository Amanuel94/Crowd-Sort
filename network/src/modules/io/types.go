package io

import (
	"context"
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

// // prints indexed items
// func PrintIndexedItem[T any](items iter.Seq[shared.IndexedItem[T]]) {
// 	buff := []shared.IndexedItem[T]{}
// 	for item := range items {
// 		fmt.Println(item.Index, item.Value)
// 		buff = append(buff, item)
// 	}
// }
