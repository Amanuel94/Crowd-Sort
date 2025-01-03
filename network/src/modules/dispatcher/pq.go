// heap implementation for Iprocess management

package dispatcher

type pq[T any] struct {
	pq []*IProcess[T]
}

func NewPQ[T any]() *pq[T] {
	return &pq[T]{
		pq: make([]*IProcess[T], 0),
	}
}

func (p *pq[T]) Push(item *IProcess[T]) {
	p.pq = append(p.pq, item)
	for i := len(p.pq) - 1; i > 0; {

		parent := (i - 1) / 2
		if (*p.pq[parent]).TaskCount() < (*p.pq[i]).TaskCount() {
			break
		}
		p.pq[parent], p.pq[i] = p.pq[i], p.pq[parent]
		i = parent
	}
}

func (p *pq[T]) Pop() *IProcess[T] {
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
