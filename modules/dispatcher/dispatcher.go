package dispatcher

import (
	"context"
	"fmt"
	"sync"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/modules/selector"
	"github.com/Amanuel94/crowdsort/shared"

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
	id2Item  map[any]*shared.IndexedItem[T]
	Ping     chan int
	MSG      chan interface{}
}

func New[T any](cfg *DispatcherConfig[T]) *Dispatcher[T] {
	d := &Dispatcher[T]{
		s:        cfg.s,
		lb:       cfg.lb,
		cpw:      cfg.cpw,
		pool:     cfg.pool,
		tcounter: cfg.tcounter,
		rank:     cfg.rank,
		channel:  cfg.channel,
		id2Item:  make(map[any]*shared.IndexedItem[T]),
		Ping:     make(chan int),
		MSG:      make(chan interface{}),
	}
	items := []interfaces.Comparable[T]{}
	for i, item := range d.lb {
		idx := i
		item_intf := item.(interfaces.Comparable[T])
		items = append(items, item_intf)
		item_indxd := item.(*shared.IndexedItem[T])
		d.id2Item[(*item_indxd).GetIndex()] = item.(*shared.IndexedItem[T])
		d.rank[(*item_indxd).GetIndex()] = idx
	}
	d.s.CreateGraph(items)
	d.n = d.s.NPairs()
	// fmt.Println(d.n)
	return d
}

func (d *Dispatcher[T]) assign(wg *sync.WaitGroup, worker *interfaces.Comparator[T], pair *shared.Pair[T]) {
	defer wg.Done()
	// d.MSG <- fmt.Sprintf("DISPATCHER INFO: Assigning %v to %v", pair.Id, (*worker).GetIndex())

	val, err := (*worker).CompareEntries(&pair.F, &pair.S)
	d.MSG <- fmt.Sprintf("DISPATCHER INFO: Comparing %v and %v", pair.F.GetValue(), pair.S.GetValue())
	argue(err == nil, "Error in comparing")
	switch val {
	case 1: // F > S
		pair.Order = shared.GT
	case -1: // F < S
		pair.Order = shared.LT
	case 0: // F == S
		pair.Order = shared.EQ
	}

	d.channel <- pair
	// d.s.PrepareNeighbours(pair.Id)

}

func (d *Dispatcher[T]) collectSelectorMessages() {
	for msg := range d.s.MSG {
		(*d).MSG <- msg
	}
}

func (d *Dispatcher[T]) get_ready_result(attr func() (*shared.Pair[T], bool)) func() (*shared.Pair[T], error) {
	get_ready_result := func() (*shared.Pair[T], error) {
		res, ok := attr()
		if !ok {
			// d.MSG <- "INFO: Waiting For Tasks to Dispatch "
			return nil, backoffError(ok, "No ready tasks")
		}
		return res, nil
	}

	return get_ready_result
}

func (d *Dispatcher[T]) Dispatch() {

	go d.collectSelectorMessages()
	deferPanic(&d.MSG)

	wg := sync.WaitGroup{}
	wg.Add(d.n)

	for d.tcounter < d.n {

		pair, err := backoff.Retry(context.TODO(), d.get_ready_result(d.s.Next))
		if err != nil {
			panic(err)
		}
		worker := d.pool.Pop()
		argue((*worker).TaskCount() <= d.cpw, "Worker is overloaded")
		go d.assign(&wg, worker, pair)
		(*worker).Assigned()
		d.pool.Push(*worker)
		d.tcounter++

	}
	wg.Wait()
	d.MSG <- fmt.Sprintf("DISPATCHER INFO: Finished dispatching %d tasks", d.tcounter)
	close(d.channel)
	// close(d.MSG)

}

func (d *Dispatcher[T]) UpdateLeaderboard() {
	count := 0
	deferPanic(&d.MSG)
	for pair := range d.channel {
		d.s.PrepareNeighbours(pair.Id)
		res := (*pair).Order
		argue(res != shared.NA, "Found Incomparable pairs")
		// d.MSG <- fmt.Sprintf("PAIRS::::: %v - %v", pair.F.GetValue(), pair.S.GetValue())
		fRank, sRank := d.rank[(*pair).F.GetIndex()], d.rank[(*pair).S.GetIndex()]
		if d.lb[fRank].(interfaces.Comparable[T]).Compare(d.lb[sRank].(interfaces.Comparable[T])) > 0 {
			// if res == shared.GT {
			d.lb[fRank], d.lb[sRank] = d.lb[sRank], d.lb[fRank]
		}
		count += 1
		d.Ping <- 1
	}

	d.MSG <- fmt.Sprintf("INFO : %d tasks completed", count)
	close(d.Ping)
}

func (d *Dispatcher[T]) GetLeaderboard() [](shared.IndexedItem[T]) {
	lb := make([]shared.IndexedItem[T], len(d.lb))

	for i, item := range d.lb {
		lb[i] = *(item.(*shared.IndexedItem[T]))
	}

	return lb
}

func (d *Dispatcher[T]) GetTaskCount() int {
	return d.tcounter
}

func (d *Dispatcher[T]) GetTotalTasks() int {
	return d.n
}
