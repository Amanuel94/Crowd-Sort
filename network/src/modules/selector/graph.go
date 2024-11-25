// Dependecy graph for the network

package selector

import (
	"github.com/google/uuid"
)

type Node[T any] struct {
	Value      *pair[T]
	Neighbours []*Node[T]
	Adj        int
}

type Graph[T any] struct {
	Nodes []*Node[T]
	m     map[uuid.UUID]*Node[T]
}

func NewGraph[T any]() *Graph[T] {
	return &Graph[T]{
		Nodes: []*Node[T]{},
		m:     make(map[uuid.UUID]*Node[T]),
	}
}

func (g *Graph[T]) AddNode(u *pair[T]) *Node[T] {
	n := &Node[T]{
		Value:      u,
		Neighbours: []*Node[T]{},
	}
	g.Nodes = append(g.Nodes, n)
	g.m[u.id] = n
	return n
}

func (g *Graph[T]) AddEdge(src_id uuid.UUID, dest_id uuid.UUID) {

	nsrc, oku := g.m[src_id]
	ndest, okv := g.m[dest_id]

	// if !oku {
	// 	nsrc = g.AddNode(src)
	// }
	// if !okv {
	// 	ndest = g.AddNode(dest)
	// }

	argue(oku && okv, "Nodes not found")
	for _, neighbour := range nsrc.Neighbours {
		if neighbour.Value.id == dest_id {
			return
		}
	}
	nsrc.Neighbours = append(nsrc.Neighbours, ndest)
	ndest.Adj++
}
