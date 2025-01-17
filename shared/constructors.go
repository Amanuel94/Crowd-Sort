package shared

import (
	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/lithammer/shortuuid"

	"golang.org/x/exp/constraints"
)

// Wrapper for indexing items
type Wire[T any] struct {
	Index string // primary key
	Value interfaces.Comparable[T]
}

func NewWire[T any](value interfaces.Comparable[T]) interfaces.Comparable[T] {
	return Wire[T]{
		Index: shortuuid.New(),
		Value: value,
	}
}

func (item Wire[T]) GetIndex() any {
	return item.Index
}

func (item Wire[T]) GetValue() T {
	return item.Value.GetValue()
}

func (item Wire[T]) Compare(other interfaces.Comparable[T]) int {
	return item.Value.Compare(other)
}

func (item Wire[T]) SetValue(val T) {
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
		Index: shortuuid.New(),
		Value: value,
	}
}

type ComparatorModule[T any] struct {
	index    string
	cmp      func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)
	task_cnt int
}

func (ic ComparatorModule[T]) GetIndex() any {
	return ic.index
}

func (ic ComparatorModule[T]) CompareEntries(f *interfaces.Comparable[T], s *interfaces.Comparable[T]) (int, error) {
	return ic.cmp(f, s)
}

func (ic *ComparatorModule[T]) Assigned() {
	(ic).task_cnt++
}

func (ic ComparatorModule[T]) TaskCount() int {
	return ic.task_cnt
}

// Constructor for Creating Comparator Modules

func NewComparator[T any](cmp func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)) interfaces.Comparator[T] {
	return &ComparatorModule[T]{
		index:    shortuuid.New(),
		cmp:      cmp,
		task_cnt: 0,
	}
}

func NewPair[T any](f string, s string) *Connector[T] {
	return &Connector[T]{Id: shortuuid.New(), F: f, S: s, Order: NA}
}
