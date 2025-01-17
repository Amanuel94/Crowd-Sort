package shared

import (
	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/lithammer/shortuuid"

	"golang.org/x/exp/constraints"
)

// Wrapper for indexing items
type Wire[T any] struct {
	index string
	value interfaces.Comparable[T]
}

func NewWire[T any](value interfaces.Comparable[T]) interfaces.Comparable[T] {
	return Wire[T]{
		index: shortuuid.New(),
		value: value,
	}
}

func (item Wire[T]) GetIndex() any {
	return item.index
}

func (item Wire[T]) GetValue() T {
	return item.value.GetValue()
}

func (item Wire[T]) Compare(other interfaces.Comparable[T]) int {
	return item.value.Compare(other)
}

func (item Wire[T]) SetValue(val T) {
	item.value.SetValue(val)
}

// Wrapper for constrained types
type OrderedType[T constraints.Ordered] struct {
	index any
	value T
}

func (o *OrderedType[T]) GetIndex() any {
	return o.index
}

func (o *OrderedType[T]) GetValue() T {
	return o.value
}

func (o *OrderedType[T]) Compare(other interfaces.Comparable[T]) int {
	if o.value < other.GetValue() {
		return -1
	} else if o.value > other.GetValue() {
		return 1
	}
	return 0
}

func (o *OrderedType[T]) SetValue(val T) {
	o.value = val
}

func NewInt[T constraints.Integer](value T) interfaces.Comparable[T] {
	return &OrderedType[T]{
		index: nil,
		value: value,
	}
}

type ComparatorModule[T any] struct {
	pid      string
	cmp      func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)
	task_cnt int
}

func (ic ComparatorModule[T]) GetID() any {
	return ic.pid
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
		pid:      shortuuid.New(),
		cmp:      cmp,
		task_cnt: 0,
	}
}

func NewPair[T any](f string, s string) *Connector[T] {
	return &Connector[T]{Id: shortuuid.New(), F: f, S: s, Order: NA}
}
