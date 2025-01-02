// binary heap implementation

package dispatcher

type pq[T any] struct {
	pq []Process[T]
}

func NewPQ[T any]() *pq[T] {
	return &pq[T]{
		pq: make([]Process[T], 0),
	}
}

func (p *pq[T]) Push(item Process[T]) {
	p.pq = append(p.pq, item)
	for i := len(p.pq) - 1; i > 0; {

		parent := (i - 1) / 2
		if p.pq[parent].Compare(p.pq[i]) < 0 {
			break
		}
		p.pq[parent], p.pq[i] = p.pq[i], p.pq[parent]
		i = parent
	}
}
