package selector

type selector[T comparable] struct {
	g *Graph
	q *queue[T]
}

func NewSelector[T comparable]() *selector[T] {
	return &selector[T]{
		g: NewGraph(),
		q: NewQueue[T](),
	}
}

