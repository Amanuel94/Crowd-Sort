package shared

import (
	"github.com/Amanuel94/crowdsort/interfaces"

	"github.com/google/uuid"
	"golang.org/x/exp/constraints"
)

// Wrapper for indexing items
type IndexedItem[T any] struct {
	Index uuid.UUID // primary key
	Value interfaces.Comparable[T]
}

func NewIndexedItem[T any](value interfaces.Comparable[T]) interfaces.Comparable[T] {
	return IndexedItem[T]{
		Index: uuid.New(),
		Value: value,
	}
}

func (item IndexedItem[T]) GetIndex() any {
	return item.Index
}

func (item IndexedItem[T]) GetValue() T {
	return item.Value.GetValue()
}

func (item IndexedItem[T]) Compare(other interfaces.Comparable[T]) int {
	return item.Value.Compare(other)
}

func (item IndexedItem[T]) SetValue(val T) {
	item.Value.SetValue(val)
}

// Wrapper for constrained types
type OrderedType[T constraints.Ordered] struct {
	Index any
	Value T
}

func (o *OrderedType[T]) GetIndex() any {
	return o.Index
}

func (o *OrderedType[T]) GetValue() T {
	return o.Value
}

func (o *OrderedType[T]) Compare(other interfaces.Comparable[T]) int {
	if o.Value < other.GetValue() {
		return -1
	} else if o.Value > other.GetValue() {
		return 1
	}
	return 0
}

func (o *OrderedType[T]) SetValue(val T) {
	o.Value = val
}

func NewInt[T constraints.Integer](value T) interfaces.Comparable[T] {
	return &OrderedType[T]{
		Index: uuid.New(),
		Value: value,
	}
}

type IndexedComparator[T any] struct {
	index    uuid.UUID
	cmp      func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)
	task_cnt int
}

func (ic IndexedComparator[T]) GetIndex() any {
	return ic.index
}

func (ic IndexedComparator[T]) CompareEntries(f *interfaces.Comparable[T], s *interfaces.Comparable[T]) (int, error) {
	return ic.cmp(f, s)
}

func (ic *IndexedComparator[T]) Assigned() {
	(ic).task_cnt++
}

func (ic IndexedComparator[T]) TaskCount() int {
	return ic.task_cnt
}

// Constructor for Creating Comparator Modules

func NewComparator[T any](cmp func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)) interfaces.Comparator[T] {
	return &IndexedComparator[T]{
		index:    uuid.New(),
		cmp:      cmp,
		task_cnt: 0,
	}
}

func NewPair[T any](f uuid.UUID, s uuid.UUID) *Pair[T] {
	return &Pair[T]{Id: uuid.New(), F: f, S: s, Order: NA}
}
