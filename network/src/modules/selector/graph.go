// Dependecy graph for the network

package selector

import (
	"github.com/google/uuid"
)

type Node[T any] struct {
	Value      pair[T]
	Neighbours []Node[T]
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

func (g *Graph[T]) AddNode(u pair[T]) {
	n := &Node[T]{
		Value:      u,
		Neighbours: []Node[T]{},
	}
	g.Nodes = append(g.Nodes, n)
	g.m[u.id] = n
}

func (g *Graph[T]) AddEdge(u pair[T], v pair[T]) {

	nu, oku := g.m[u.id]
	nv, okv := g.m[v.id]

	if !oku {
		g.AddNode(u)
	}
	if !okv {
		g.AddNode(v)
	}

	for _, neighbour := range nu.Neighbours {
		if neighbour.Value.id == v.id {
			return
		}
	}
	nu.Neighbours = append(nu.Neighbours, *nv)
	nv.Neighbours = append(nv.Neighbours, *nu)
	nu.Adj++
	nv.Adj++
}
