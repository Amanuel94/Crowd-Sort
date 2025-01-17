package interfaces

// Comparable is an interface that defines a method to compare two values
type Comparable[T any] interface {
	GetIndex() any
	GetValue() T
	SetValue(T)
	Compare(other Comparable[T]) int
}

// wrapper for  comparator modules

type Comparator[T any] interface {
	GetID() any
	CompareEntries(*Comparable[T], *Comparable[T]) (int, error)
	Assigned()
	TaskCount() int
}

type Keyed interface {
	GetKey() string
}
