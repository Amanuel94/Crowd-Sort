// Generic queue implementation
package selector

type listNode[T any] struct {
	value T
	next  *listNode[T]
}

type queue[T any] struct {
	head *listNode[T]
	tail *listNode[T]
	size int
}

func NewQueue[T any]() *queue[T] {
	head := &listNode[T]{}
	tail := &listNode[T]{next: head}
	head.next = tail
	return &queue[T]{
		head: head,
		tail: tail,
	}
}

func (q *queue[T]) Enqueue(value T) {
	node := &listNode[T]{value: value}
	q.tail.next.next = node
	q.tail.next = node
	q.size++
}

func (q *queue[T]) Dequeue() T {
	argue(q.size > 0, "Empty queue")
	node := q.head.next
	q.head.next = node.next
	return node.value
}
