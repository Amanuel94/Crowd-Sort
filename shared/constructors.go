package shared

import (
	"github.com/Amanuel94/crowdsort/interfaces"
	"github.com/Amanuel94/crowdsort/utils"

	"golang.org/x/exp/constraints"
)

func NewWire[T any](value interfaces.Comparable[T]) interfaces.Comparable[T] {
	return Wire[T]{
		index:  utils.NewIdentifier(),
		value:  value,
		status: PENDING,
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

func (item Wire[T]) GetStatus() Status {
	return item.status
}

func (item *Wire[T]) SetStatus(status Status) {
	item.status = status
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

func NewComparator[T any](cmp CmpFunc[T]) interfaces.Comparator[T] {
	return &ComparatorModule[T]{
		pid:      utils.NewIdentifier(),
		cmp:      cmp,
		task_cnt: 0,
		status:   ComparatorStatusIdle,
	}
}

func NewConnector[T any](f string, s string) *Connector[T] {
	return &Connector[T]{Id: utils.NewIdentifier(), F: f, S: s, Order: NA}
}

func (c *ComparatorModule[T]) SetStatus(s ComparatorStatus) {
	c.status = s
}

func (c *ComparatorModule[T]) GetStatus() ComparatorStatus {
	return c.status
}

func (c *Connector[T]) GetKey() string {
	return c.Id
}

func NewLeaderboardUpdate(f string, s string, assigneeId string) *PingMessage {
	return &PingMessage{
		Type:        LeaderboardUpdate,
		F:           f,
		S:           s,
		AssignieeId: assigneeId,
	}
}

func NewTaskStatusUpdate(wid string) *PingMessage {
	return &PingMessage{
		Type:   TaskStatusUpdate,
		WireId: wid,
	}
}

func NewComparatorStatusUpdate(id string) *PingMessage {
	return &PingMessage{
		Type:         ComparatorStatusUpdate,
		ComparatorId: id,
	}
}

func AsModule[T any](s *interfaces.Comparator[T]) *ComparatorModule[T] {
	return (*s).(*ComparatorModule[T])
}

func AsWire[T any](s *interfaces.Comparable[T]) *Wire[T] {
	return (*s).(*Wire[T])
}
