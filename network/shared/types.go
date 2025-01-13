package shared

import (
	"github.com/amanuel94/crowdsort/interfaces"

	"github.com/google/uuid"
)

// wrapper for pair
// ****************************************************

// enum for score
type Ord = int

const (
	NA Ord = iota
	LT
	EQ
	GT
)

type Pair[T any] struct {
	Id    uuid.UUID
	F     interfaces.Comparable[T]
	S     interfaces.Comparable[T]
	Order Ord
}

func NewPair[T any](f interfaces.Comparable[T], s interfaces.Comparable[T]) *Pair[T] {
	return &Pair[T]{Id: uuid.New(), F: f, S: s, Order: NA}
}
