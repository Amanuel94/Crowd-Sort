package utils

import (
	"iter"
)

// Map applies a function to each element of a sequence and returns a new sequence with the results.
func Map[T1, T2 any](f func(T1) T2,  seq iter.Seq[T1]) iter.Seq[T2] {
	return func(yield func (T2) bool){
		for v := range seq {
			if !yield(f(v)){
				return
			}
		}
	}
}

// SliceToSeq converts a slice to a sequence.
func SliceToSeq[T any](slice []T) iter.Seq[T] {
	return func(yield func (T) bool){
		for _, v := range slice {
			if !yield(v){
				return
			}
		}
	}
}

// Concat concatenates multiple sequences into a single sequence.
func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func (T) bool){
		for _, seq := range seqs {
			for v := range seq {
				if !yield(v){
					return
				}
			}
		}
	}
}