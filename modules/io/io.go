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
	msgBuffer := make([]interface{}, 0)
	RegisterMessage("[INFO]: Initializing IO", cfg.verbose, &msgBuffer)

	items := utils.Map(func(v *interfaces.Comparable[T]) *shared.Wire[T] {
		item := shared.NewWire(*v).(shared.Wire[T])
		return &item
	}, cfg.items)
	comparators := utils.Map(func(v shared.CmpFunc[T]) *shared.ComparatorModule[T] {
		return shared.NewComparator(v).(*shared.ComparatorModule[T])
	}, cfg.comparators)

	newIO := &IO[T]{}

	dcfg := dispatcher.NewDispatcherConfig(items, comparators)
	newIO.d = dispatcher.New(dcfg)
	newIO.msgBuffer = msgBuffer
	newIO.wg = utils.NewWaitGroup(2)
	newIO.verbose = cfg.verbose

	RegisterMessage("[INFO]: IO Initialized", newIO.verbose, &newIO.msgBuffer)

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

func (io *IO[T]) ShowLeaderboard() {
	cnt := 0
	go io.d.UpdateLeaderboard()
	for p := range io.d.Ping {
		clearTable()
		fmt.Printf("Live Leaderboard\n\n")
		printTable([]string{"Wire", "Value"}, io.d.GetLeaderboard(), p)
		fmt.Println()
		fmt.Println()
		tble := io.d.GetComparatorsFromPool()
		printWorkerStatusTable(tble, p)

		fmt.Println()
		fmt.Println()
		printProgressBar(io.d.GetTaskCount(), io.d.GetTotalTasks())
		fmt.Println()
		if p.Type == shared.LeaderboardUpdate {
			printUpdate(p)
		}
		if io.verbose > 1 {
			io.showCollectedMessages()
		}
		cnt++

	}
	fmt.Println("INFO: Final Leaderboard ", cnt)
	io.wg.Done()

}

func (io *IO[T]) showCollectedMessages() {
	for _, msg := range io.msgBuffer {
		fmt.Println(msg)
	}
}

func (io *IO[T]) Wait() {
	io.wg.Wait()

}

func RegisterMessage(msg string, verbose int, msgBuffer *[]interface{}) {
	if verbose > 0 {
		fmt.Println(msg)
		*msgBuffer = append(*msgBuffer, msg)
	}
}
