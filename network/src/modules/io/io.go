// Handles side-effects
package io

import (
	"context"
	"fmt"
	"iter"
	"network/shared"
	"network/shared/interfaces"
	"network/utils"
	"time"
)

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

func (io *IO[T]) Read() iter.Seq[shared.IndexedItem[T]] {
	fmt.Println("Reading from IO")
	return io.ctx.Value(io.key).(iter.Seq[shared.IndexedItem[T]])
}

func (io *IO[T]) WriteFromList(values []interfaces.Comparable[T]) {
	asSeq := utils.SliceToSeq(values)
	io.WriteFromSeq(asSeq)
}

func (io *IO[T]) WriteFromSeq(values iter.Seq[interfaces.Comparable[T]]) {
	indexedValues := utils.Map(func(v interfaces.Comparable[T]) shared.IndexedItem[T] { return *shared.NewIndexedItem(v) }, values)
	if io.ctx.Value(io.key) == nil {
		fmt.Println("Writing to empty IO")
		io.ctx = context.WithValue(io.ctx, io.key, indexedValues)
	} else {
		fmt.Println("Writing to non-empty IO")
		curr := io.ctx.Value(io.key).(iter.Seq[shared.IndexedItem[T]])
		io.ctx = context.WithValue(io.ctx, io.key, utils.Concat(curr, indexedValues))
	}
	fmt.Println("Done writing")
}

func (io *IO[T]) Clear() {
	io.ctx = context.WithValue(io.ctx, io.key, nil)
}

func (io *IO[T]) Close() {
	fmt.Println("Closing IO")
	io.canc()
}

func WriteInt(i *IO[int64], values []int64) {
	asSeq := utils.SliceToSeq(values)
	asComparable := utils.Map(func(v int64) interfaces.Comparable[int64] { return NewInt(v) }, asSeq)
	i.WriteFromSeq(asComparable)
}
