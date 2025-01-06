package dispatcher

import (
	"network/modules/selector"
	"network/shared"
)

type DispatcherConfig[T any] struct {
	s        *selector.Selector[T]
	lb       [](any)
	n        int
	cpw      int         //capacity per worker
	pool     *pq[T]      // should contain already defined workers/processes
	tcounter int         // number of assigned tasks
	rank     map[any]int // maps id to rank
	channel  chan *shared.Pair[T]
}

func IndexedDispatcherConfig[T any](s *selector.Selector[T], items []*shared.IndexedItem[T], tasks_limit int, processes []*IProcess[T]) *DispatcherConfig[T] {

	lb := []any{}
	rank := make(map[any]int)
	for idx, item := range items {
		lb = append(lb, item)
		rank[(*item).GetIndex()] = idx
	}

	pq := FromList(processes)

	return &DispatcherConfig[T]{
		s:        s,
		lb:       lb,
		n:        len(lb),
		cpw:      tasks_limit,
		pool:     pq,
		tcounter: 0,
		rank:     rank,
		channel:  make(chan *shared.Pair[T]),
	}

}
