package dispatcher

import (
	"iter"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/shared"

	"github.com/Amanuel94/crowdsort/modules/selector"
)

type DispatcherConfig[T any] struct {
	s        *selector.Selector[T]
	lb       [](any)
	n        int
	workers  []*interfaces.Comparator[T]
	cpw      int         //capacity per worker
	pool     *pq[T]      // should contain already defined workers/comparators
	tcounter int         // number of assigned tasks
	rank     map[any]int // maps id to rank
	channel  chan *shared.Connector[T]
}

func NewDispatcherConfig[T any](items iter.Seq[*shared.Wire[T]], comparators iter.Seq[*shared.ComparatorModule[T]]) *DispatcherConfig[T] {

	lb := []any{}
	rank := make(map[any]int)
	idx := 0
	for item := range items {
		lb = append(lb, item)
		rank[(*item).GetIndex()] = idx
	}

	pq := FromSeq(comparators)
	scfg := selector.NewConfig()

	return &DispatcherConfig[T]{
		s:        selector.NewSelector[T](*scfg),
		lb:       lb,
		n:        len(lb),
		cpw:      1000,
		pool:     pq,
		workers:  append([]*(interfaces.Comparator[T]){}, pq.pq...),
		tcounter: 0,
		rank:     rank,
		channel:  make(chan *shared.Connector[T]),
	}

}

func WithTaskLimit[T any](cfg *DispatcherConfig[T], limit int) *DispatcherConfig[T] {
	cfg.cpw = limit
	return cfg
}
