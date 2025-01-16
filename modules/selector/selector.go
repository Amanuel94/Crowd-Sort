package selector

import (
	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/shared"

	"github.com/google/uuid"
)

type Selector[T any] struct {
	g       *Graph[T]
	q       *queue[*shared.Pair[T]]
	alg     string
	batched bool
	MSG     chan interface{}
	Rank    map[uuid.UUID]int
}

// TODO: Implement bitonic sort and shell sort
func NewSelector[T any](cfg Config) *Selector[T] {
	argue(cfg.alg == "BEMS", "Invalid algorithm name")
	return &Selector[T]{
		g:       NewGraph[T](),
		q:       NewQueue[*shared.Pair[T]](),
		alg:     cfg.alg,
		batched: false,
		Rank:    make(map[uuid.UUID]int),
	}
}

func (s *Selector[T]) NPairs() int {
	return len(s.g.Nodes)
}

func (s *Selector[T]) RegisterItems(u [](interfaces.Comparable[T])) {
	for i, item := range u {
		s.Rank[item.GetIndex().(uuid.UUID)] = i
	}
}

func (s *Selector[T]) CreateGraph(u [](interfaces.Comparable[T])) {

	deferPanic(&s.MSG)
	argue(len(u) > 0, "Empty input")

	n_nodes := len(u)
	pair_indices := BEMS_pairs_generator(n_nodes, 1, 0, &s.MSG)
	pmap := make(map[uuid.UUID]uuid.UUID)
	for _, pi := range pair_indices {
		i, j := pi[0], pi[1]
		if i >= n_nodes || j >= n_nodes {
			continue
		}
		pair := shared.NewPair[T](u[i].GetIndex().(uuid.UUID), u[j].GetIndex().(uuid.UUID))
		s.g.AddNode(pair)

		fprev, fok := pmap[pair.F]
		sprev, sok := pmap[pair.S]
		if fok {
			s.g.AddEdge(fprev, pair.Id, &s.MSG)
		}
		if sok {
			s.g.AddEdge(sprev, pair.Id, &s.MSG)
		}
		pmap[pair.F] = pair.Id
		pmap[pair.S] = pair.Id
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

	return s.q.Dequeue(&s.MSG), true
}

// Enqueue pairs with 0 dependencies
func (s *Selector[T]) PrepareNeighbours(id uuid.UUID) {
	node, ok := s.g.m[id]
	deferPanic(&s.MSG)
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

// func PrintGraph[T any](g *Graph[T]) {
// 	for _, node := range g.Nodes {
// 		fmt.Print(node.Value.F.GetValue(), node.Value.S.GetValue())
// 		fmt.Println(":")
// 		for _, neighbour := range node.Neighbours {
// 			fmt.Print("(", neighbour.Value.F.GetValue(), neighbour.Value.S.GetValue(), ")")
// 		}
// 		fmt.Println()
// 	}
// }
