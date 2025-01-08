package dispatcher

import (
	"context"
	"fmt"
	"network/interfaces"
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
	s        *selector.Selector[T]
	lb       [](any)
	n        int
	cpw      int         //capacity per worker
	pool     *pq[T]      // should contain already defined workers/processes
	tcounter int         // number of assigned tasks
	rank     map[any]int // maps id to rank
	channel  chan *shared.Pair[T]
}

func New[T any](cfg *DispatcherConfig[T]) *Dispatcher[T] {
	return &Dispatcher[T]{
		s:        cfg.s,
		lb:       cfg.lb,
		n:        cfg.n,
		cpw:      cfg.cpw,
		pool:     cfg.pool,
		tcounter: cfg.tcounter,
		rank:     cfg.rank,
		channel:  cfg.channel,
	}
}

func (d *Dispatcher[T]) assign(wg *sync.WaitGroup, process *interfaces.Comparator[T], pair *shared.Pair[T]) {
	defer wg.Done()
	val, err := (*process).CompareEntries(&pair.F, &pair.S)

	argue(err == nil, err.Error())
	switch val {
	case 1: // F > S
		pair.Order = shared.GT
	case -1: // F < S
		pair.Order = shared.LT
	case 0: // F == S
		pair.Order = shared.EQ
	}

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
	wg.Add(d.n)

	for d.tcounter < d.n {
		pair, err := backoff.Retry(context.TODO(), get_ready_result)
		if err != nil {
			break // TODO: handle this
		}
		worker := d.pool.Pop()
		argue((*worker).TaskCount() <= d.cpw, "Worker is overloaded")
		go d.assign(&wg, worker, pair)
		(*worker).Assigned()
		d.tcounter++
		d.pool.Push(*worker)

	}
	wg.Wait()
	close(d.channel)

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
