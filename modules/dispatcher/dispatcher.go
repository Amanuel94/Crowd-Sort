package dispatcher

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/modules/selector"
	"github.com/Amanuel94/crowdsort/shared"

	"github.com/cenkalti/backoff/v5"
)

type Dispatcher[T any] struct {
	s         *selector.Selector[T]
	lb        [](any)
	n         int
	cpw       int //capacity per worker
	workers   []*interfaces.Comparator[T]
	pool      *pq[T]      // should contain already defined workers/processes
	tcounter  int         // number of assigned tasks
	rank      map[any]int // maps id to rank
	channel   chan *shared.Connector[T]
	id2Item   map[any]*interfaces.Comparable[T]
	id2Worker map[any]*interfaces.Comparator[T]
	Ping      chan shared.PingMessage
	MSG       chan interface{}
}

func New[T any](cfg *DispatcherConfig[T]) *Dispatcher[T] {
	d := &Dispatcher[T]{
		s:         cfg.s,
		lb:        cfg.lb,
		cpw:       cfg.cpw,
		pool:      cfg.pool,
		workers:   cfg.workers,
		tcounter:  cfg.tcounter,
		rank:      cfg.rank,
		channel:   cfg.channel,
		id2Item:   make(map[any]*interfaces.Comparable[T]),
		id2Worker: make(map[any]*interfaces.Comparator[T]),
		Ping:      make(chan shared.PingMessage),
		MSG:       make(chan interface{}),
	}

	items := []interfaces.Comparable[T]{}
	for i, item := range d.lb {
		idx := i
		itemComparable := item.(interfaces.Comparable[T])
		items = append(items, itemComparable)
		itemWire := item.(*shared.Wire[T])
		itemIndex := itemWire.GetIndex()
		d.id2Item[itemIndex] = &itemComparable
		d.rank[itemIndex] = idx
	}

	for _, worker := range d.workers {
		d.id2Worker[(*worker).GetID()] = worker
	}
	d.s.CreateGraph(items)
	d.n = d.s.NPairs()

	return d
}

func (d *Dispatcher[T]) GetComparatorsFromPool() []*interfaces.Comparator[T] {

	comparators := append([]*interfaces.Comparator[T]{}, d.workers...)

	sort.Slice(comparators, func(i, j int) bool {
		return (*comparators[i]).TaskCount() > (*comparators[j]).TaskCount()
	})
	return comparators

}

func (d *Dispatcher[T]) assign(wg *sync.WaitGroup, worker *interfaces.Comparator[T], pair *shared.Connector[T]) {
	defer deferPanic(&d.MSG)
	defer wg.Done()

	pf, ps := d.id2Item[pair.F], d.id2Item[pair.S]
	pfi, psi := (*pf).(*shared.Wire[T]), (*ps).(*shared.Wire[T])

	workeri := (*worker).(*shared.ComparatorModule[T])

	statusMsg := shared.Assigned((*worker).GetID().(string))

	pfi.SetStatus(statusMsg)
	psi.SetStatus(statusMsg)
	workeri.SetStatus(shared.ComparatorStatusBusy)

	msgf := *shared.NewTaskStatusUpdate(pfi.GetIndex().(string))
	d.Ping <- msgf
	msgs := *shared.NewTaskStatusUpdate(psi.GetIndex().(string))
	d.Ping <- msgs
	csmsg := *shared.NewComparatorStatusUpdate((*worker).GetID().(string))
	d.Ping <- csmsg

	val, err := (*worker).CompareEntries(pf, ps)
	argue(err == nil, "Error in comparing")
	switch val {
	case 1: // F > S
		pair.Order = shared.GT
	case -1: // F < S
		pair.Order = shared.LT
	case 0: // F == S
		pair.Order = shared.EQ
	}

	d.MSG <- fmt.Sprintf("[INFO]: Comparator %s submitted.", (*worker).GetID())

	d.channel <- pair

}

func (d *Dispatcher[T]) collectSelectorMessages() {
	for msg := range d.s.MSG {
		(*d).MSG <- msg
	}
}

func (d *Dispatcher[T]) get_ready_result(attr func() (*shared.Connector[T], bool), worker *interfaces.Comparator[T]) func() (*shared.Connector[T], error) {
	get_ready_result := func() (*shared.Connector[T], error) {
		res, ok := attr()
		d.MSG <- fmt.Sprintf("[INFO]: Awaiting for new tasks to assign %v...", (*worker).GetID())
		if !ok {
			return nil, backoffError(ok, "No ready tasks")
		}
		return res, nil
	}

	return get_ready_result
}

func (d *Dispatcher[T]) Dispatch() {
	go d.collectSelectorMessages()
	defer deferPanic(&d.MSG)

	wg := sync.WaitGroup{}
	defer wg.Done()
	wg.Add(d.n)
	for d.tcounter < d.n {
		d.pool.mu.Lock()
		worker := d.pool.Pop()
		for len(d.pool.pq) > 0 && (*worker).TaskCount() >= d.cpw {
			workeri := (*worker).(*shared.ComparatorModule[T])
			workeri.SetStatus(shared.ComparatorStatusDone)
			d.Ping <- *shared.NewComparatorStatusUpdate((*worker).GetID().(string))

			worker = d.pool.Pop()
		}

		argue(len(d.pool.pq) > 0, "No available workers")

		bo := backoff.NewExponentialBackOff()
		bo.InitialInterval = 2 * time.Second
		opts := backoff.WithBackOff(bo)
		pair, err := backoff.Retry(context.TODO(), d.get_ready_result(d.s.Next, worker), opts)
		if err != nil {
			panic(err)
		}
		go d.assign(&wg, worker, pair)
		pair.AssignieeId = (*worker).GetID().(string)
		(*worker).Assigned()
		d.pool.Push(*worker)
		d.pool.mu.Unlock()
		d.tcounter++

	}
	wg.Wait()
	d.MSG <- fmt.Sprintf("[INFO]: Finished dispatching %d tasks", d.tcounter)
	close(d.channel)

}

func (d *Dispatcher[T]) UpdateLeaderboard() {
	count := 0
	defer deferPanic(&d.MSG)
	for pair := range d.channel {
		d.s.PrepareNeighbours(pair.Id)

		pf, ps := d.id2Item[pair.F], d.id2Item[pair.S]
		pfi, psi := (*pf).(*shared.Wire[T]), (*ps).(*shared.Wire[T])
		workeri := (*d.id2Worker[pair.AssignieeId]).(*shared.ComparatorModule[T])
		if d.s.GetRemainingComparision(pair.F) == 0 {
			pfi.SetStatus(shared.COMPLETED)
		} else {
			pfi.SetStatus(shared.PENDING)
		}
		d.Ping <- *shared.NewTaskStatusUpdate(pfi.GetIndex().(string))

		if d.s.GetRemainingComparision(pair.S) == 0 {
			psi.SetStatus(shared.COMPLETED)
		} else {
			psi.SetStatus(shared.PENDING)
		}
		d.Ping <- *shared.NewTaskStatusUpdate(psi.GetIndex().(string))

		workeri.SetStatus(shared.ComparatorStatusIdle)
		d.Ping <- *shared.NewComparatorStatusUpdate(workeri.GetID().(string))

		res := (*pair).Order
		argue(res != shared.NA, "Found Incomparable pairs")
		if res == shared.GT {
			pfv := (*pf).GetValue()
			psv := (*ps).GetValue()
			(*pf).SetValue(psv)
			(*ps).SetValue(pfv)
		}
		count += 1
		d.Ping <- *shared.NewLeaderboardUpdate((*pair).F, (*pair).S, (*pair).AssignieeId)
	}

	d.MSG <- fmt.Sprintf("[INFO] : %d tasks completed", count)
	close(d.Ping)
}

func (d *Dispatcher[T]) GetLeaderboard() [](shared.Wire[T]) {
	lb := make([]shared.Wire[T], len(d.lb))
	for i, item := range d.lb {
		lb[i] = *(item.(*shared.Wire[T]))
	}

	return lb
}

func (d *Dispatcher[T]) GetTaskCount() int {
	return d.tcounter
}

func (d *Dispatcher[T]) GetTotalTasks() int {
	return d.n
}
