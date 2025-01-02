package dispatcher

import (
	"network/modules/selector"
	"network/shared"
	"sync"
)

// Batches from Selector
// Takes Process + Pair
// Calls PrepareNeighbours
// Maintains Leaderboard

type Dispatcher[T any] struct {
	s           *selector.Selector[T]
	lb          []any
	n_workers   int
	task_limit  int
	worker_pool []Process[T]
	task_pq     []int
}

func NewDispatcher[T any](s *selector.Selector[T]) *Dispatcher[T] {
	return &Dispatcher[T]{
		s:  s,
		lb: make([]any, 0),
	}
}

func (d *Dispatcher[T]) worker(wg *sync.WaitGroup, process Process[T], pair shared.Pair[T]) {
	defer wg.Done()
	process.CompareEntries(pair.F, pair.S)
	d.s.PrepareNeighbours(pair.Id)
}

func (d *Dispatcher[T]) Dispatch() {

	wg := sync.WaitGroup{}
	for {
		pair, ok := d.s.Next()
		if !ok {
			break
		}
		wg.Add(1)
		go d.worker(&wg, process, pair)
	}
	wg.Wait()

}
