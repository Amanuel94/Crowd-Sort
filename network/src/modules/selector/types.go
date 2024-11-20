package selector

import "network/interfaces"

// enum for score
type ord = int
const (
	NA ord = iota
	LT 
	EQ 
	GT 
)

type pair[T any] struct{
	f interfaces.Comparable[T]	
	s interfaces.Comparable[T]
	order ord
	
	
}