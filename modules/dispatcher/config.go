package dispatcher

import (
	"iter"

	"github.com/Amanuel94/crowdsort/shared"

	"github.com/Amanuel94/crowdsort/modules/selector"
)

type DispatcherConfig[T any] struct {
	s           *selector.Selector[T]
	lb          [](any)
	n           int
	cpw         int         //capacity per worker
	pool        *pq[T]      // should contain already defined workers/processes
	tcounter    int         // number of assigned tasks
	rank        map[any]int // maps id to rank
	channel     chan *shared.Pair[T]
	refresh_cnt int
}

func IndexedDispatcherConfig[T any](items iter.Seq[*shared.IndexedItem[T]], processes iter.Seq[*shared.IndexedComparator[T]]) *DispatcherConfig[T] {

	lb := []any{}
	rank := make(map[any]int)
	idx := 0
	for item := range items {
		lb = append(lb, item)
		rank[(*item).GetIndex()] = idx
	}

	// pq := FromList(processes)
	pq := FromSeq(processes)
	scfg := selector.NewConfig()

	return &DispatcherConfig[T]{
		s:           selector.NewSelector[T](*scfg),
		lb:          lb,
		n:           len(lb),
		cpw:         len(lb)/len(pq.pq) + 1,
		pool:        pq,
		tcounter:    0,
		rank:        rank,
		channel:     make(chan *shared.Pair[T]),
		refresh_cnt: 0,
	}

}

func WithTaskLimit[T any](cfg *DispatcherConfig[T], limit int) {
	cfg.cpw = limit
}
