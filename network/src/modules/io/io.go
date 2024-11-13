// Handles side-effects
package io

import (
	"context"
	"fmt"
	"iter"
	"network/interfaces"
	"network/utils"
	"time"
)

type IO[T any] struct{
	ctx context.Context
	canc context.CancelFunc
	key IOKey 
}

// Initalizes the IO module

func Init[T any]() *IO[T] {
	fmt.Println("Initializing IO")
	newIO := &IO[T]{}
	newIO.createContext(NewConfig("io"))
	return newIO
}

func InitWithTimeOut[T any](timeOut time.Duration) *IO[T] {
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

func (io *IO[T]) WriteFromList(values []interfaces.Comparable[T]) {
	asSeq := utils.SliceToSeq(values)
	io.WriteFromSeq(asSeq)
}

func (io *IO[T]) WriteFromSeq(values iter.Seq[interfaces.Comparable[T]]) {
	indexedValues := utils.Map[interfaces.Comparable[T], IndexedItem[T]](func(v interfaces.Comparable[T]) IndexedItem[T] {return *NewIndexedItem[T](v)}, values)
	if (io.ctx.Value(io.key) == nil) {
		io.ctx = context.WithValue(io.ctx, io.key, indexedValues)
	} else {
		curr := io.ctx.Value(io.key).(iter.Seq[IndexedItem[T]])
		io.ctx = context.WithValue(io.ctx, io.key, utils.Concat(curr, indexedValues))
	}
}

func (io *IO[int64]) WriteInt(values []int64) {
	asSeq := utils.SliceToSeq(values)
	asComparable := utils.Map[int64, interfaces.Comparable[int64]](func (v int64) interfaces.Comparable[int64] {return NewInt64(v)}, asSeq)
	io.WriteFromSeq(asComparable)
}

func (io *IO[T]) Clear() {
	io.ctx = context.WithValue(io.ctx, io.key, nil)
}

func (io *IO[T]) Close() {
	fmt.Println("Closing IO")
	io.canc()
}



