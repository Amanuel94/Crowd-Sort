package dispatcher

import (
	"context"
	"fmt"
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
	rank           map[any]int // maps id to rank
	channel        chan *shared.Pair[T]
}

func NewDispatcher[T any](s *selector.Selector[T]) *Dispatcher[T] {
	return &Dispatcher[T]{
		s:  s,
		lb: make([]any, 0),
	}
}

func (d *Dispatcher[T]) assign(wg *sync.WaitGroup, process *IProcess[T], pair *shared.Pair[T]) {
	defer wg.Done()
	err := (*process).CompareEntries(pair)

	argue(err == nil, err.Error())

	d.channel <- pair
	d.s.PrepareNeighbours(pair.Id)

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
			break // TODO: handle this
		}
		worker := d.worker_pool.Pop()
		argue((*worker).TaskCount() <= d.task_limit, "Worker is overloaded")
		go d.assign(&wg, worker, pair)
		(*worker).Assigned()
		d.tasks_assigned++
		d.worker_pool.Push(worker)

	}
	wg.Wait()

	fmt.Println("Every task has been assigned")

}

func (d *Dispatcher[T]) UpdateLeaderboard() {

	for pair := range d.channel {
		res := (*pair).Order
		argue(res != shared.NA, "Found Incomparable pairs")

		if res == shared.GT {
			fRank, sRank := d.rank[(*pair).F.GetIndex()], d.rank[(*pair).S.GetIndex()]
			d.lb[fRank], d.lb[sRank] = d.lb[sRank], d.lb[fRank]
			d.rank[(*pair).F.GetIndex()] = sRank
			d.rank[(*pair).S.GetIndex()] = fRank
		}
	}
}
