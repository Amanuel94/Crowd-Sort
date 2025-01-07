package dispatcher

import (
	"network/interfaces"
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

func IndexedDispatcherConfig[T any](items []*shared.IndexedItem[T], processes []*(interfaces.Comparator[T])) *DispatcherConfig[T] {

	lb := []any{}
	rank := make(map[any]int)
	for idx, item := range items {
		lb = append(lb, item)
		rank[(*item).GetIndex()] = idx
	}

	pq := FromList(processes)
	scfg := selector.NewConfig()

	return &DispatcherConfig[T]{
		s:        selector.NewSelector[T](*scfg),
		lb:       lb,
		n:        len(lb),
		cpw:      len(lb) / len(processes),
		pool:     pq,
		tcounter: 0,
		rank:     rank,
		channel:  make(chan *shared.Pair[T]),
	}

}

func WithTaskLimit[T any](cfg *DispatcherConfig[T], limit int) {
	cfg.cpw = limit
}
