package dispatcher

import (
	"network/shared"
	"network/shared/interfaces"
)

type IProcess[T any] interface {
	GetIndex() any
	CompareEntries(interfaces.Comparable[T], interfaces.Comparable[T]) (shared.Ord, error)
	Assigned()
	TaskCount() int
}

// for live standings
// TODO: Implement live prodcast
