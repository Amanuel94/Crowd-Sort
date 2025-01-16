// Handles side-effects
package io

import (
	"fmt"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/modules/dispatcher"
	"github.com/Amanuel94/crowdsort/shared"
	"github.com/Amanuel94/crowdsort/utils"
)

// Initalizes the IO module

func New[T any](cfg *Config[T]) *IO[T] {
	fmt.Println("INFO: Initializing IO")
	newIO := &IO[T]{}
	items := utils.Map(func(v *interfaces.Comparable[T]) *shared.IndexedItem[T] {
		item := shared.NewIndexedItem[T](*v).(shared.IndexedItem[T])
		return &item
	}, cfg.items)

	comparators := utils.Map(func(v func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)) *shared.IndexedComparator[T] {
		return shared.NewComparator[T](v).(*shared.IndexedComparator[T])
	}, cfg.comparators)

	dcfg := dispatcher.IndexedDispatcherConfig[T](items, comparators)
	newIO.d = dispatcher.New(dcfg)
	newIO.msgBuffer = make([]interface{}, 0)
	newIO.wg = utils.NewWaitGroup(2)

	fmt.Println("INFO: IO Initialized")

	return newIO
}

func (io *IO[T]) collectDispatcherMessages() {
	for msg := range io.d.MSG {
		fmt.Println(msg)
		io.msgBuffer = append(io.msgBuffer, msg)
	}
}

func (io *IO[T]) StartDispatcher() {
	go io.collectDispatcherMessages()
	fmt.Println("INFO: Starting Dispatcher")
	io.d.Dispatch()
	io.wg.Done()
}

// func (io *IO[T]) showCollectedMessages() {
// 	for _, msg := range io.msgBuffer {
// 		fmt.Println(msg)
// 	}
// }

func (io *IO[T]) ShowLeaderboard() {
	cnt := 0
	go io.d.UpdateLeaderboard()
	for range io.d.Ping {
		clearTable()
		fmt.Printf("Live Leaderboard\n")
		printTable([]string{"Index", "Value"}, io.d.GetLeaderboard())
		fmt.Println()
		printProgressBar(io.d.GetTaskCount(), io.d.GetTotalTasks())
		fmt.Println()
		// io.showCollectedMessages()
		cnt++

	}
	fmt.Println("INFO: Leaderboard Updated", cnt)
	io.wg.Done()
	// close(io.d.MSG)

}

func (io *IO[T]) Wait() {
	io.wg.Wait()

}
