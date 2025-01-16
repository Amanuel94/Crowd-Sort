// Dependecy graph for the network

package selector

import (
	"github.com/Amanuel94/crowdsort/shared"
)

type Node[T any] struct {
	Value      *shared.Pair[T]
	Neighbours []*Node[T]
	Adj        int
}

type Graph[T any] struct {
	Nodes []*Node[T]
	m     map[string]*Node[T]
}

func NewGraph[T any]() *Graph[T] {
	return &Graph[T]{
		Nodes: []*Node[T]{},
		m:     make(map[string]*Node[T]),
	}
}

func (g *Graph[T]) AddNode(u *shared.Pair[T]) *Node[T] {
	n := &Node[T]{
		Value:      u,
		Neighbours: []*Node[T]{},
	}
	g.Nodes = append(g.Nodes, n)
	g.m[u.Id] = n
	return n
}

func (g *Graph[T]) AddEdge(src_id string, dest_id string, msg *chan interface{}) {

	nsrc, oku := g.m[src_id]
	ndest, okv := g.m[dest_id]

	deferPanic(msg)
	argue(oku && okv, "Nodes not found")
	for _, neighbour := range nsrc.Neighbours {
		if neighbour.Value.Id == dest_id {
			return
		}
	}
	nsrc.Neighbours = append(nsrc.Neighbours, ndest)
	ndest.Adj++
}
