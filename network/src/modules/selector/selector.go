package selector

import (
	"network/shared/interfaces"

	"github.com/google/uuid"
)

type selector[T interfaces.Comparable[T]] struct {
	g   *Graph[T]
	q   *queue[T]
	alg string
}

// TODO: Implement bitonic sort and shell sort
func NewSelector[T interfaces.Comparable[T]](cfg Config) *selector[T] {
	argue(cfg.alg == "BEMS", "Invalid algorithm name")
	return &selector[T]{
		g:   NewGraph[T](),
		q:   NewQueue[T](),
		alg: cfg.alg,
	}
}

func (s *selector[T]) CreateGraph(u []interfaces.Comparable[T]) {
	n_nodes := len(u)
	pair_indices := BEMS_pairs_generator(n_nodes, 1, 0)

	pmap := make(map[uuid.UUID]pair[T])
	for _, pi := range pair_indices {
		i, j := pi[0], pi[1]
		pair := NewPair(u[i], u[j])
		fprev, fok := pmap[pair.f.GetIndex().(uuid.UUID)]
		sprev, sok := pmap[pair.s.GetIndex().(uuid.UUID)]
		if fok {
			s.g.AddEdge(*pair, fprev)
		}
		if sok {
			s.g.AddEdge(*pair, sprev)
		}
		pmap[pair.f.GetIndex().(uuid.UUID)] = *pair
		pmap[pair.s.GetIndex().(uuid.UUID)] = *pair
	}
}

func (s *selector[T]) Batch() {
	for _, node := range s.g.Nodes {
		if node.Adj == 1 {
			s.q.Enqueue(node.Value)
			node.Adj-- // avoid re-adding the same node
		}
	}
}

func (s *selector[T]) PrepareNeighbours(id uuid.UUID) {
	node, ok := s.g.m[id]
	argue(ok, "Node not found")
	for _, neighbour := range node.Neighbours {
		neighbour.Adj--
	}

}
