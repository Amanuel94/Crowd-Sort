package interfaces

// Comparable is an interface that defines a method to compare two values
type Comparable[T any] interface {
	GetIndex() any
	GetValue() T
	Compare(other Comparable[T]) int
}

// wrapper for  comparator modules

type Comparator[T any] interface {
	GetIndex() any
	CompareEntries(*Comparable[T], *Comparable[T]) (int, error)
	Assigned()
	TaskCount() int
}
