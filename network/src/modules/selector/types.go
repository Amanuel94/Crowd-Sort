package selector

type indexer[T comparable] struct{
	index int
	value T
}

type pair[T comparable] struct{
	f T	
	s T
}