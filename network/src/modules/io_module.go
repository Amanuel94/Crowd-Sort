package io_module

type io[T comparable] struct {
    value T
}

func (i io[T]) get() T {
	return i.value
}