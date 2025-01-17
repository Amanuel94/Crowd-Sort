package selector

import (
	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/shared"
)

type Selector[T any] struct {
	g       *graph[*shared.Connector[T]]
	q       *queue[*shared.Connector[T]]
	alg     string
	batched bool
	MSG     chan interface{}
}

// TODO: Implement bitonic sort and shell sort
func NewSelector[T any](cfg Config) *Selector[T] {

	return &Selector[T]{
		g:       NewGraph[*shared.Connector[T]](),
		q:       NewQueue[*shared.Connector[T]](),
		alg:     cfg.alg,
		batched: false,
		MSG:     make(chan interface{}),
	}
}

func (s *Selector[T]) NPairs() int {
	return len(s.g.nodes)
}

func (s *Selector[T]) CreateGraph(u [](interfaces.Comparable[T])) {

	defer deferPanic(&s.MSG)
	argue(len(u) > 0, "Empty input")

	n_nodes := len(u)
	pair_indices := BEMS_pairs_generator(n_nodes, 1, 0, &s.MSG)
	pmap := make(map[string]string)
	for _, pi := range pair_indices {
		i, j := pi[0], pi[1]
		if i >= n_nodes || j >= n_nodes {
			continue
		}
		pair := shared.NewConnector[T](u[i].GetIndex().(string), u[j].GetIndex().(string))
		s.g.addNode(&pair)

		fprev, fok := pmap[pair.F]
		sprev, sok := pmap[pair.S]
		if fok {
			s.g.addEdge(fprev, pair.Id, &s.MSG)
		}
		if sok {
			s.g.addEdge(sprev, pair.Id, &s.MSG)
		}
		pmap[pair.F] = pair.Id
		pmap[pair.S] = pair.Id
	}

}

func (s *Selector[T]) Next() (*shared.Connector[T], bool) {

	if !s.batched {
		s.firstBatch()
		s.batched = true
	}
	if s.q.size == 0 {
		return &shared.Connector[T]{}, false
	}

	return s.q.dequeue(&s.MSG), true
}

// Enqueue pairs with 0 dependencies
func (s *Selector[T]) PrepareNeighbours(id string) {
	node, ok := s.g.m[id]
	defer deferPanic(&s.MSG)
	argue(ok, "Node not found")
	for _, neighbour := range node.neighbours {
		neighbour.adj--
		if neighbour.adj == 0 {
			s.q.enqueue(*neighbour.value)
		}
	}

}

func (s *Selector[T]) firstBatch() {

	for _, node := range s.g.nodes {
		if node.adj == 0 {
			s.q.enqueue(*node.value)
		}
	}

}
