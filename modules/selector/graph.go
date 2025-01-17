// Dependecy graph for the network

package selector

import (
	"github.com/Amanuel94/crowdsort/shared"
)

type node[T any] struct {
	value      *shared.Connector[T]
	neighbours []*node[T]
	adj        int
}

type graph[T any] struct {
	nodes []*node[T]
	m     map[string]*node[T]
}

func NewGraph[T any]() *graph[T] {
	return &graph[T]{
		nodes: []*node[T]{},
		m:     make(map[string]*node[T]),
	}
}

func (g *graph[T]) addNode(u *shared.Connector[T]) *node[T] {
	n := &node[T]{
		value:      u,
		neighbours: []*node[T]{},
	}
	g.nodes = append(g.nodes, n)
	g.m[u.Id] = n
	return n
}

func (g *graph[T]) addEdge(src_id string, dest_id string, msg *chan interface{}) {

	nsrc, oku := g.m[src_id]
	ndest, okv := g.m[dest_id]

	deferPanic(msg)
	argue(oku && okv, "Nodes not found")
	for _, neighbour := range nsrc.neighbours {
		if neighbour.value.Id == dest_id {
			return
		}
	}
	nsrc.neighbours = append(nsrc.neighbours, ndest)
	ndest.adj++
}
