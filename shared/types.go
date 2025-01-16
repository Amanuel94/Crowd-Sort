package shared

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
	Id    string
	F     string
	S     string
	Order Ord
}
