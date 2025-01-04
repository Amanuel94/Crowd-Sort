package dispatcher

import (
	"network/shared"
)

type IProcess[T any] interface {
	GetIndex() any
	CompareEntries(*shared.Pair[T]) error
	Assigned()
	TaskCount() int
}

// for live standings
// TODO: Implement live prodcast
