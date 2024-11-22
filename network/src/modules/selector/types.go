package selector

import (
	"network/shared/interfaces"

	"github.com/google/uuid"
)

// enum for score
type ord = int

const (
	NA ord = iota
	LT
	EQ
	GT
)

type pair[T any] struct {
	id    uuid.UUID
	f     interfaces.Comparable[T]
	s     interfaces.Comparable[T]
	order ord
}

func NewPair[T any](f interfaces.Comparable[T], s interfaces.Comparable[T]) *pair[T] {
	return &pair[T]{id: uuid.New(), f: f, s: s, order: NA}
}
