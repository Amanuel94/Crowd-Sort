// Handles side-effects
package io

import (
	"context"
	"fmt"
	"iter"
	"network/utils"
	"time"

	"golang.org/x/exp/constraints"
)

type IO[T constraints.Ordered] struct{
	ctx context.Context
	canc context.CancelFunc
	key IOKey 
}

// Initalizes the IO module

func Init[T constraints.Ordered]() *IO[T] {
	fmt.Println("Initializing IO")
	newIO := &IO[T]{}
	newIO.createContext(NewConfig("io"))
	return newIO
}

func InitWithTimeOut[T constraints.Ordered](timeOut time.Duration) *IO[T] {
	fmt.Println("Initializing IO with timeout")
	newIO := &IO[T]{}
	cfg := NewConfig("io")
	cfg.WithTimeout(timeOut)
	newIO.createContext(cfg)
	return newIO
}

func (io *IO[T]) createContext(cfg *Config) {

	io.key = *NewIOKey(cfg.key)
	if cfg.withTimeout {
		io.ctx, io.canc = context.WithTimeout(context.Background(), cfg.timeOut)
	} else {
		io.ctx = context.Background()
		io.canc = func() {}
	}
}

func (io *IO[T]) Read() iter.Seq[IndexedItem[T]] {
	return io.ctx.Value(io.key).(iter.Seq[IndexedItem[T]])
}

func (io *IO[T]) Write(values []T) {
	asSeq := utils.SliceToSeq(values)
	indexedValues := utils.Map[T, IndexedItem[T]](func(v T) IndexedItem[T] {return *NewIndexedItem[T](v)}, asSeq)
	if (io.ctx.Value(io.key) == nil) {
		io.ctx = context.WithValue(io.ctx, io.key, indexedValues)
	} else {
		curr := io.ctx.Value(io.key).(iter.Seq[IndexedItem[T]])
		io.ctx = context.WithValue(io.ctx, io.key, utils.Concat(curr, indexedValues))
	}
}

func (io *IO[T]) Clear() {
	io.ctx = context.WithValue(io.ctx, io.key, nil)
}


func (io *IO[T]) Close() {
	fmt.Println("Closing IO")
	io.canc()
}



