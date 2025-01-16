package shared

import (
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
	F     uuid.UUID
	S     uuid.UUID
	Order Ord
}
