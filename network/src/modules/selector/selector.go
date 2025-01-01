package selector

import (
	"network/shared"
	"network/shared/interfaces"
	"reflect"

	"github.com/google/uuid"
)

type Selector[T any] struct {
	g *Graph[T]
	q *queue[*shared.Pair[T]]
	// leaves  *queue[Node[T]]
	alg     string
	batched bool
}

// TODO: Implement bitonic sort and shell sort
func NewSelector[T any](cfg Config) *Selector[T] {
	argue(cfg.alg == "BEMS", "Invalid algorithm name")
	return &Selector[T]{
		g:       NewGraph[T](),
		q:       NewQueue[*shared.Pair[T]](),
		alg:     cfg.alg,
		batched: false,
	}
}

func (s *Selector[T]) CreateGraph(u [](interfaces.Comparable[T])) {

	argue(len(u) > 0, "Empty input")
	dummyInstance := shared.IndexedItem[T]{}
	argue(reflect.TypeOf(u[0]) == reflect.TypeOf(dummyInstance), "Invalid type")

	n_nodes := len(u)
	pair_indices := BEMS_pairs_generator(n_nodes, 1, 0)

	pmap := make(map[uuid.UUID]uuid.UUID)
	for _, pi := range pair_indices {
		i, j := pi[0], pi[1]
		if i >= n_nodes || j >= n_nodes {
			continue
		}
		pair := shared.NewPair(u[i], u[j])
		s.g.AddNode(pair)

		fprev, fok := pmap[pair.F.GetIndex().(uuid.UUID)]
		sprev, sok := pmap[pair.S.GetIndex().(uuid.UUID)]
		if fok {
			s.g.AddEdge(fprev, pair.Id)
		}
		if sok {
			s.g.AddEdge(sprev, pair.Id)
		}
		pmap[pair.F.GetIndex().(uuid.UUID)] = pair.Id
		pmap[pair.S.GetIndex().(uuid.UUID)] = pair.Id
	}
}

func (s *Selector[T]) Next() (*shared.Pair[T], bool) {

	if !s.batched {
		s.firstBatch()
		s.batched = true
	}
	if s.q.size == 0 {
		return &shared.Pair[T]{}, false
	}
	return s.q.Dequeue(), true
}

// Enqueue pairs with 0 dependencies
func (s *Selector[T]) PrepareNeighbours(id uuid.UUID) {
	node, ok := s.g.m[id]
	argue(ok, "Node not found")
	for _, neighbour := range node.Neighbours {
		neighbour.Adj--
		if neighbour.Adj == 0 {
			s.q.Enqueue(neighbour.Value)
		}
	}
}

func (s *Selector[T]) firstBatch() {

	for _, node := range s.g.Nodes {
		if node.Adj == 0 {
			s.q.Enqueue(node.Value)
		}
	}

}
