// heap implementation for (interfaces.Comparator management
package dispatcher

import (
	"iter"

	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/shared"
)

type pq[T any] struct {
	pq []*(interfaces.Comparator[T])
}

func NewPQ[T any]() *pq[T] {
	return &pq[T]{
		pq: make([]*(interfaces.Comparator[T]), 0),
	}
}

func FromList[T any](processes []*(interfaces.Comparator[T])) *pq[T] {

	return &pq[T]{
		pq: processes,
	}

}

func FromSeq[T any](processes iter.Seq[*shared.IndexedComparator[T]]) *pq[T] {
	pq := NewPQ[T]()
	for process := range processes {
		pq.Push(process)
	}
	return pq
}

func (p *pq[T]) Push(item interfaces.Comparator[T]) {
	p.pq = append(p.pq, &item)
	for i := len(p.pq) - 1; i > 0; {

		parent := (i - 1) / 2
		if (*p.pq[parent]).TaskCount() < (*p.pq[i]).TaskCount() {
			break
		}
		p.pq[parent], p.pq[i] = p.pq[i], p.pq[parent]
		i = parent
	}
}

func (p *pq[T]) Pop() *(interfaces.Comparator[T]) {
	if len(p.pq) == 0 {
		return nil
	}
	item := p.pq[0]
	p.pq[0] = p.pq[len(p.pq)-1]
	p.pq = p.pq[:len(p.pq)-1]
	for i := 0; i < len(p.pq); {
		left := 2*i + 1
		right := 2*i + 2
		if left >= len(p.pq) {
			break
		}
		min := left
		if right < len(p.pq) && (*p.pq[right]).TaskCount() < (*p.pq[left]).TaskCount() {
			min = right
		}
		if (*p.pq[i]).TaskCount() < (*p.pq[min]).TaskCount() {
			break
		}
		p.pq[i], p.pq[min] = p.pq[min], p.pq[i]
		i = min
	}
	return item
}
