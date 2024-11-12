package selector

type queue[T any] struct{
	q chan T	
} 

func NewQueue[T any]() *queue[T] {
	return &queue[T]{q: make(chan T)}
}

func (q *queue[T]) Enqueue(value T) {
	q.q <- value
}

func (q *queue[T]) Dequeue() T {
	return <-q.q
}