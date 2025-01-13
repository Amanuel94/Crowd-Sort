package dispatcher

import (
	"context"
	"fmt"
	"sync"

	"github.com/amanuel94/crowdsort/interfaces"
	"github.com/amanuel94/crowdsort/modules/selector"
	"github.com/amanuel94/crowdsort/shared"

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
	Ping     chan int
	MSG      chan interface{}
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
		Ping:     make(chan int),
		MSG:      make(chan interface{}),
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

func (d *Dispatcher[T]) collectSelectorMessages(s *selector.Selector[T]) {
	for msg := range s.MSG {
		(*d).MSG <- msg
	}
}

func (d *Dispatcher[T]) Dispatch() {

	go d.collectSelectorMessages(d.s)
	deferPanic(&d.MSG)

	get_ready_result := func() (*shared.Pair[T], error) {
		res, ok := d.s.Next()
		if !ok {
			d.MSG <- "DISPATCHER WARNING: No ready tasks"
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
	close(d.s.MSG)

	d.MSG <- fmt.Sprintf("DISPATCHER INFO: Finished dispatching %d tasks", d.tcounter)

}

func (d *Dispatcher[T]) UpdateLeaderboard() {

	deferPanic(&d.MSG)
	for pair := range d.channel {
		res := (*pair).Order
		argue(res != shared.NA, "Found Incomparable pairs")

		if res == shared.GT {
			fRank, sRank := d.rank[(*pair).F.GetIndex()], d.rank[(*pair).S.GetIndex()]
			d.lb[fRank], d.lb[sRank] = d.lb[sRank], d.lb[fRank]
			d.rank[(*pair).F.GetIndex()] = sRank
			d.rank[(*pair).S.GetIndex()] = fRank
		}

		d.Ping <- 1
	}
	close(d.Ping)
}

func (d *Dispatcher[T]) GetLeaderboard() [](shared.IndexedItem[T]) {
	lb := make([]shared.IndexedItem[T], len(d.lb))
	for i, item := range d.lb {
		lb[i] = item.(shared.IndexedItem[T])
	}
	return lb
}

func (d *Dispatcher[T]) GetTaskCount() int {
	return d.tcounter
}

func (d *Dispatcher[T]) GetTotalTasks() int {
	return d.n
}
