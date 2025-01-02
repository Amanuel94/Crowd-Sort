package dispatcher

import (
	"context"
	"network/modules/selector"
	"network/shared"
	"sync"

	"github.com/cenkalti/backoff/v5"
	"github.com/google/uuid"
)

// Batches from Selector
// Takes Process + Pair
// Calls PrepareNeighbours
// Maintains Leaderboard

type Dispatcher[T any] struct {
	s              *selector.Selector[T]
	lb             []any
	n_workers      int
	n_tasks        int
	task_limit     int
	worker_pool    *pq[T] // should contain already defined workers/processes
	tasks_assigned int
	failed         []uuid.UUID
}

func NewDispatcher[T any](s *selector.Selector[T]) *Dispatcher[T] {
	return &Dispatcher[T]{
		s:  s,
		lb: make([]any, 0),
	}
}

func (d *Dispatcher[T]) assign(wg *sync.WaitGroup, process IProcess[T], pair shared.Pair[T]) shared.Ord {
	defer wg.Done()
	res, err := process.CompareEntries(pair.F, pair.S)

	if err != nil {
		return shared.NA
	}
	d.s.PrepareNeighbours(pair.Id)
	return res
}

func (d *Dispatcher[T]) Dispatch() {

	get_ready_result := func() (*shared.Pair[T], error) {
		res, ok := d.s.Next()
		if !ok {
			return nil, backoffError(ok, "No ready tasks")
		}
		return res, nil
	}

	wg := sync.WaitGroup{}
	wg.Add(d.n_tasks)

	for d.tasks_assigned < d.n_tasks {
		pair, err := backoff.Retry(context.TODO(), get_ready_result)
		if err != nil {
			break
		}
		worker := d.worker_pool.Pop()
		argue((*worker).TaskCount() <= d.task_limit, "Worker is overloaded")
		go d.assign(&wg, *worker, *pair)
		(*worker).Assigned()
		d.tasks_assigned++
		d.worker_pool.Push(worker)

	}
	wg.Wait()

}
