package dispatcher

import (
	"context"
	"network/modules/selector"
	"network/shared"
	"sync"

	"github.com/cenkalti/backoff/v5"
)

// Batches from Selector
// Takes Process + Pair
// Calls PrepareNeighbours
// Maintains Leaderboard

type Dispatcher[T any] struct {
	s               *selector.Selector[T]
	lb              []any
	n_workers       int
	n_tasks         int
	task_limit      int
	worker_pool     []*pq[T]
	tasks_completed int
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

	get_ready_result := func() (*shared.Pair[T], error) {
		res, ok := d.s.Next()
		if !ok {
			return nil, argue(ok, "No ready tasks")
		}
		return res, nil
	}

	wg := sync.WaitGroup{}
	wg.Add(d.n_tasks)
	task_cnt := 0
	for task_cnt < d.n_tasks {
		pair, err := backoff.Retry(context.TODO(), get_ready_result)
		if err != nil {
			break
		}
	}
	wg.Wait()

}
