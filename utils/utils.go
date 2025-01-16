// utility functions

package utils

import (
	"iter"
	"sync"

	"golang.org/x/exp/rand"
)

// Map applies a function to each element of a sequence and returns a new sequence with the results.
func Map[T1, T2 any](f func(T1) T2, seq iter.Seq[T1]) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// SliceToSeq converts a slice to a sequence.
func SliceToSeq[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range slice {
			if !yield(v) {
				return
			}
		}
	}
}

// Concat concatenates multiple sequences into a single sequence.
func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// bit count
func Bit_count(x int) int {
	count := 0
	for x != 0 {
		x &= x - 1
		count++
	}
	return count
}

func Bit_Len(x int) int {
	count := 0
	for x != 0 {
		x >>= 1
		count++
	}
	return count
}

// next power of 2
func NextPower(x int) int {
	return 1 << (Bit_Len(x))
}

// randint generates a random integer in the range [a, b].
func RandInt(a, b int) int {
	return rand.Intn(b-a+1) + a
}

// new wait group with counter
func NewWaitGroup(n int) *sync.WaitGroup {
	wg := sync.WaitGroup{}
	wg.Add(n)
	return &wg
}
