// Handles side-effects
package io

import (
	"context"
	"fmt"
	"iter"

	"github.com/amanuel94/crowdsort/network/interfaces"
	"github.com/amanuel94/crowdsort/network/modules/dispatcher"
	"github.com/amanuel94/crowdsort/network/shared"
	"github.com/amanuel94/crowdsort/network/utils"
)

// Initalizes the IO module

func Init[T any](cfg *Config[T]) *IO[T] {
	fmt.Println("Initializing IO")
	newIO := &IO[T]{}
	newIO.ctx = *cfg.ctx
	newIO.canc = *cfg.canc
	items := utils.Map(func(v *interfaces.Comparable[T]) *shared.IndexedItem[T] {
		return shared.NewIndexedItem[T](*v).(*shared.IndexedItem[T])
	}, cfg.items)

	comparators := utils.Map(func(v func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)) *shared.IndexedComparator[T] {
		return shared.NewComparator[T](v).(*shared.IndexedComparator[T])
	}, cfg.comparators)

	dcfg := dispatcher.IndexedDispatcherConfig[T](items, comparators)
	newIO.d = dispatcher.New(dcfg)
	return newIO
}
func (io *IO[T]) ReadItems(key IOKey) iter.Seq[shared.IndexedItem[T]] {
	fmt.Println("Reading Values from IO-Context")
	return io.ctx.Value(key).(iter.Seq[shared.IndexedItem[T]])
}

func (io *IO[T]) WriteFromList(values []interfaces.Comparable[T], key IOKey) {
	asSeq := utils.SliceToSeq(values)
	io.WriteFromSeq(asSeq, key)
}

func (io *IO[T]) WriteFromSeq(values iter.Seq[interfaces.Comparable[T]], key IOKey) {
	indexedValues := utils.Map(func(v interfaces.Comparable[T]) shared.IndexedItem[T] {
		return shared.NewIndexedItem[T](v).(shared.IndexedItem[T])
	}, values)
	if io.ctx.Value(key) == nil {
		fmt.Println("Writing to empty IO")
		io.ctx = context.WithValue(io.ctx, key, indexedValues)
	} else {
		fmt.Println("Writing to non-empty IO")
		curr := io.ctx.Value(key).(iter.Seq[shared.IndexedItem[T]])
		io.ctx = context.WithValue(io.ctx, key, utils.Concat(curr, indexedValues))
	}
	fmt.Println("Done writing")
}

func (io *IO[T]) Clear(key IOKey) {
	io.ctx = context.WithValue(io.ctx, key, nil)
}

func (io *IO[T]) collectDispatcherMessages() {
	for msg := range io.d.MSG {
		io.msgBuffer = append(io.msgBuffer, msg)
	}
}

func (io *IO[T]) Close() {
	fmt.Println("Closing IO")
	io.canc()
}

func WriteInt(i *IO[int64], values []int64, key IOKey) {
	asSeq := utils.SliceToSeq(values)
	asComparable := utils.Map(func(v int64) interfaces.Comparable[int64] { return shared.NewInt(v) }, asSeq)
	i.WriteFromSeq(asComparable, key)
}

func (io *IO[T]) StartDispatcher() {
	go io.collectDispatcherMessages()
	io.msgBuffer = append(io.msgBuffer, "Starting Dispatcher")
	io.d.Dispatch()
}

func (io *IO[T]) showCollectedMessages() {
	for _, msg := range io.msgBuffer {
		fmt.Println(msg)
	}
}

func (io *IO[T]) ShowLeaderboard() {
	io.msgBuffer = append(io.msgBuffer, "Live Leaderboard")
	for range io.d.Ping {
		clearTable()
		io.showCollectedMessages()
		printTable([]string{"Index", "Value"}, io.d.GetLeaderboard())
		printProgressBar(io.d.GetTaskCount(), io.d.GetTotalTasks())

	}
	close(io.d.MSG)

}
