package dispatcher

import "network/shared/interfaces"

type Process[T any] interface {
	GetIndex() any
	CompareEntries(interfaces.Comparable[T], interfaces.Comparable[T])
	TaskCount() int
}

// for live standings
// TODO: Implement live prodcast
