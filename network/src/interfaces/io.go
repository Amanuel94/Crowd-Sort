package contracts

// interace for generic io with a clients
type IO[T comparable] struct {
    value T
}
