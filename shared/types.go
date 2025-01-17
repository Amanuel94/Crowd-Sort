package shared

import "github.com/Amanuel94/crowdsort/interfaces"

// enum for score
type Ord = int

const (
	NA Ord = iota
	LT
	EQ
	GT
)

type Connector[T any] struct {
	Id    string
	F     string
	S     string
	Order Ord
}

type CmpFunc[T any] func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)
