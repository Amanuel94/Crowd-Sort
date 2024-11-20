// Dependecy graph for the network

package selector

import (
	"github.com/google/uuid"
)

type Node struct {
	Index uuid.UUID
	Neighbours []uuid.UUID
}

type Graph struct {
	Nodes []*Node
	m map[uuid.UUID]*Node
}

func NewGraph() *Graph {
	return &Graph{
		Nodes: []*Node{},
		m: make(map[uuid.UUID]*Node),
	}
}

func (g* Graph) AddNode(u uuid.UUID) {
	n := &Node{
		Index: u,
		Neighbours: []uuid.UUID{},
	}
	g.Nodes = append(g.Nodes, n)
	g.m[u] = n
}

func (g* Graph) AddEdge(u uuid.UUID, v uuid.UUID) {

	nu, oku := g.m[u]
	nv, okv := g.m[v]

	argue(!oku, "Node u not found")
	argue(!okv, "Node v not found")

	nu.Neighbours = append(nu.Neighbours, u)
	nv.Neighbours = append(nv.Neighbours, v)
}