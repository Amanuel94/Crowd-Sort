package dispatcher

import "network/modules/selector"

// Batches from Selector
// Takes Process + Pair
// Calls PrepareNeighbours
// Maintains Leaderboard

type Dispatcher[T any] struct {
	s *selector.Selector[T]
}

func NewDispatcher[T any](s *selector.Selector[T]) *Dispatcher[T] {
	return &Dispatcher[T]{
		s: s,
	}
}

// func (d *Dispatcher[T]) Assign() {
