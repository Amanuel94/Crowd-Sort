package selector

import (
	"network/shared"
	"network/shared/interfaces"
	"reflect"

	"github.com/google/uuid"
)

type selector[T any] struct {
	g *Graph[T]
	q *queue[pair[T]]
	// leaves  *queue[Node[T]]
	alg     string
	batched bool
}

// TODO: Implement bitonic sort and shell sort
func NewSelector[T any](cfg Config) *selector[T] {
	argue(cfg.alg == "BEMS", "Invalid algorithm name")
	return &selector[T]{
		g:       NewGraph[T](),
		q:       NewQueue[pair[T]](),
		alg:     cfg.alg,
		batched: false,
	}
}

func (s *selector[T]) CreateGraph(u [](interfaces.Comparable[T])) {

	argue(len(u) > 0, "Empty input")
	dummyInstance := shared.IndexedItem[T]{}
	argue(reflect.TypeOf(u[0]) == reflect.TypeOf(dummyInstance), "Invalid type")

	n_nodes := len(u)
	pair_indices := BEMS_pairs_generator(n_nodes, 1, 0)

	pmap := make(map[uuid.UUID]pair[T])
	for _, pi := range pair_indices {
		i, j := pi[0], pi[1]
		if i >= n_nodes || j >= n_nodes {
			continue
		}
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

func (s *selector[T]) Batch() (pair[T], bool) {

	if !s.batched {
		s.firstBatch()
		s.batched = true
	}
	if s.q.size == 0 {
		return pair[T]{}, false
	}
	return s.q.Dequeue(), true
}

func (s *selector[T]) PrepareNeighbours(id uuid.UUID) {
	node, ok := s.g.m[id]
	argue(ok, "Node not found")
	for _, neighbour := range node.Neighbours {
		neighbour.Adj--
		if neighbour.Adj == 1 {
			s.q.Enqueue(neighbour.Value)
		}
	}
}

func (s *selector[T]) firstBatch() {

	for _, node := range s.g.Nodes {
		if node.Adj == 1 {
			s.q.Enqueue(node.Value)
		}
	}

}
