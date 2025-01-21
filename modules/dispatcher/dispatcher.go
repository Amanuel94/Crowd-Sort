package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/modules/selector"
	"github.com/Amanuel94/crowdsort/shared"
	"github.com/Amanuel94/crowdsort/utils"

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

	for i, item := range d.lb {
		itemComparable := item.(interfaces.Comparable[T])
		itemWire := item.(*shared.Wire[T])
		itemIndex := itemWire.GetIndex()

		d.id2Item[itemIndex] = &itemComparable
		d.rank[itemIndex] = i
	}

	for _, worker := range d.workers {
		d.id2Worker[(*worker).GetID()] = worker
	}

	comparables := make([]interfaces.Comparable[T], len(d.lb))
	for i, item := range d.lb {
		comparables[i] = item.(interfaces.Comparable[T])
	}

	d.s.CreateGraph(comparables)
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
	workeri := shared.AsModule(worker)

	d.MSG <- fmt.Sprintf("[INFO]: Assigning task to %s", (*worker).GetID())

	// Set statuses
	statusMsg := shared.Assigned((*worker).GetID().(string))
	pfi.SetStatus(statusMsg)
	psi.SetStatus(statusMsg)
	workeri.SetStatus(shared.ComparatorStatusBusy)

	// Send status updates
	d.Ping <- *shared.NewTaskStatusUpdate(pfi.GetIndex().(string))
	d.Ping <- *shared.NewTaskStatusUpdate(psi.GetIndex().(string))
	d.Ping <- *shared.NewComparatorStatusUpdate((*worker).GetID().(string))

	// Perform comparison

	start := time.Now()
	val, err := (*worker).CompareEntries(pf, ps)
	elapsed := time.Since(start)

	if err != nil {
		d.MSG <- fmt.Sprintf("[ERROR]: Error in comparing: %v", err)
		return
	}

	// Set order based on comparison result
	switch val {
	case 1:
		pair.Order = shared.GT
	case -1:
		pair.Order = shared.LT
	case 0:
		pair.Order = shared.EQ
	}

	d.MSG <- fmt.Sprintf("[INFO]: Comparator %s submitted. Time Elapsed: %v ms", (*worker).GetID(), elapsed.Milliseconds())

	updateStatus := func(wire *shared.Wire[T]) {
		if d.s.GetRemainingComparision((*wire).GetIndex().(string)) == 0 {
			wire.SetStatus(shared.COMPLETED)
		} else {
			wire.SetStatus(shared.PENDING)
		}
		d.Ping <- *shared.NewTaskStatusUpdate(wire.GetIndex().(string))
	}

	updateStatus(pfi)
	updateStatus(psi)

	if (*worker).TaskCount() < d.cpw {
		shared.AsModule(worker).SetStatus(shared.ComparatorStatusIdle)
	} else if (*worker).TaskCount() == d.cpw {
		shared.AsModule(worker).SetStatus(shared.ComparatorStatusDone)
	} else {
		shared.AsModule(worker).SetStatus(shared.ComparatorStatusOverflow)
	}
	d.Ping <- *shared.NewComparatorStatusUpdate((*worker).GetID().(string))

	// Send the pair to the channel
	d.channel <- pair

}

func (d *Dispatcher[T]) collectSelectorMessages() {
	for msg := range d.s.MSG {
		(*d).MSG <- msg
	}
}

func (d *Dispatcher[T]) getReadyResult(attr func() (*shared.Connector[T], bool), worker *interfaces.Comparator[T]) func() (*shared.Connector[T], error) {
	return func() (*shared.Connector[T], error) {
		d.MSG <- fmt.Sprintf("[INFO]: Awaiting new tasks to assign to %v...", (*worker).GetID())
		res, ok := attr()
		if !ok {
			return nil, backoffError(ok, "No ready tasks")
		}
		return res, nil
	}
}

func (d *Dispatcher[T]) getWorker() (*interfaces.Comparator[T], error) {
	var idleWorker *interfaces.Comparator[T]
	var nonIdleWorkers []*interfaces.Comparator[T]

	for len(d.pool.pq) > 0 {
		worker := d.pool.Pop()
		workeri := shared.AsModule(worker)
		if workeri.GetStatus() == shared.ComparatorStatusIdle {
			idleWorker = worker
			break
		}
		nonIdleWorkers = append(nonIdleWorkers, worker)
	}

	for _, worker := range nonIdleWorkers {
		workeri := shared.AsModule(worker)
		if workeri.GetStatus() == shared.ComparatorStatusBusy {
			d.pool.Push(workeri)
		}
	}

	if len(d.pool.pq) == 0 && idleWorker == nil {
		return nil, backoff.Permanent(errors.New("no workers available"))
	}

	if idleWorker != nil {
		return idleWorker, nil
	}

	d.MSG <- "[INFO]: All workers are busy. Awaiting for a worker to be available..."
	return nil, backoffError(false, "all workers are occupied")
}

func getBackoffRetryOptionWithInitialInterval(interval time.Duration) backoff.RetryOption {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = interval
	opts := backoff.WithBackOff(bo)
	return opts
}

func (d *Dispatcher[T]) Dispatch() {
	go d.collectSelectorMessages()
	defer deferPanic(&d.MSG)
	defer close(d.channel)
	wg := utils.NewWaitGroup(d.n)

	d.pool.mu.Lock()
	for d.tcounter < d.n {

		// TODO: Assumes every worker completes its tasks. Handle non-happy path
		opts := getBackoffRetryOptionWithInitialInterval(3 * time.Second)
		worker, err := backoff.Retry(context.TODO(), d.getWorker, opts)
		if err != nil {
			d.MSG <- fmt.Sprintf("[ERROR]: %v", err)
			panic(err)
		}

		argue(len(d.pool.pq) > 0, "All workers are overloaded.")

		opts = getBackoffRetryOptionWithInitialInterval(3 * time.Second)
		pair, err := backoff.Retry(context.TODO(), d.getReadyResult(d.s.Next, worker), opts)
		if err != nil {
			panic(err)
		}

		go d.assign(wg, worker, pair)
		pair.AssignieeId = (*worker).GetID().(string)
		(*worker).Assigned()
		d.pool.Push(*worker)
		d.tcounter++
	}

	d.pool.mu.Unlock()
	wg.Wait()
	d.MSG <- fmt.Sprintf("[INFO]: Finished dispatching %d tasks", d.tcounter)

}

func (d *Dispatcher[T]) UpdateLeaderboard() {
	count := 0
	defer deferPanic(&d.MSG)
	defer close(d.Ping)

	for pair := range d.channel {
		d.s.PrepareNeighbours(pair.Id)

		pf := d.id2Item[pair.F]
		ps := d.id2Item[pair.S]

		if (d.s.GetRemainingComparision(pair.F)) == 0 {
			shared.AsWire(pf).SetStatus(shared.COMPLETED)
			d.Ping <- *shared.NewTaskStatusUpdate(pair.F)
		}

		if (d.s.GetRemainingComparision(pair.S)) == 0 {
			shared.AsWire(ps).SetStatus(shared.COMPLETED)
			d.Ping <- *shared.NewTaskStatusUpdate(pair.S)
		}

		res := pair.Order
		argue(res != shared.NA, "Found Incomparable pairs")

		if res == shared.GT {
			pfValue := (*pf).GetValue()
			psValue := (*ps).GetValue()
			(*pf).SetValue(psValue)
			(*ps).SetValue(pfValue)
		}

		count++
		d.Ping <- *shared.NewLeaderboardUpdate(pair.F, pair.S, pair.AssignieeId)
	}

	d.MSG <- fmt.Sprintf("[INFO]: %d tasks completed", count)
}

func (d *Dispatcher[T]) GetLeaderboard() [](shared.Wire[T]) {
	lb := make([]shared.Wire[T], len(d.lb))
	for i, item := range d.lb {
		wireItem := item.(*shared.Wire[T])
		lb[i] = *wireItem
	}

	return lb
}

func (d *Dispatcher[T]) GetTaskCount() int {
	return d.tcounter
}

func (d *Dispatcher[T]) GetTotalTasks() int {
	return d.n
}
