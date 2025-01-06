package globals

// Comparable is an interface that defines a method to compare two values

type Comparable[T any] interface {
	GetIndex() any
	GetValue() T
	Compare(other Comparable[T]) int
}

// IProcess is an interface that defines a method to compare two values

type IProcess[T any] interface {
	GetIndex() any
	CompareEntries(Comparable[T], Comparable[T]) error
	Assigned()
	TaskCount() int
}
