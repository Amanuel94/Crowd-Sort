package selector

type queue[T any] struct{
	q chan pair[T]	
} 

func NewQueue[T any]() *queue[T] {
	return &queue[T]{q: make(chan pair[T])}
}

func (q *queue[T]) Enqueue(value pair[T]) {
	q.q <- value
}

func (q *queue[T]) Dequeue() pair[T] {
	return <-q.q
}