package shared

import (
	"fmt"

	"github.com/Amanuel94/crowdsort/interfaces"
	"golang.org/x/exp/constraints"
)

type Ord = int                 // order: <, >, =, =?
type Status = string           // status of wire containing comparable items
type ComparatorStatus = string // status of comparator module
type MessageType = int         // type of ping message
type CmpFunc[T any] func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)

// Creates dynamically assigned status for wire.
func Assigned(assignee string) Status {
	return Status(fmt.Sprintf("ASSIGNED TO: %s", assignee))
}

const (
	NA Ord = iota
	LT
	EQ
	GT
)

const (
	PENDING   Status = "PENDING"
	COMPLETED Status = "COMPLETED"
)

const (
	ComparatorStatusIdle     ComparatorStatus = "IDLE"
	ComparatorStatusBusy     ComparatorStatus = "BUSY"
	ComparatorStatusDone     ComparatorStatus = "DONE"
	ComparatorStatusOverflow ComparatorStatus = "OVERFLOW"
)

const (
	TaskStatusUpdate MessageType = iota
	LeaderboardUpdate
	ComparatorStatusUpdate
)

// Wire is a struct that contains a comparable item.
type Wire[T any] struct {
	index  string
	value  interfaces.Comparable[T]
	status Status
}

// A comparable wrapper for constrained types (int, float, etc.).
type OrderedType[T constraints.Ordered] struct {
	index any
	value T
}

// A comparator wrapper for comparator
type ComparatorModule[T any] struct {
	pid      string
	cmp      CmpFunc[T]
	task_cnt int
	status   ComparatorStatus
}

// Contains  two comparables identified by the wires F and S
type Connector[T any] struct {
	Id          string
	F           string
	S           string
	Order       Ord
	AssignieeId string
}

// Not all fields are filled
type PingMessage struct {
	Type         MessageType
	F            string // required when Type == LeaderboardUpdate
	S            string // required when Type == LeaderboardUpdate
	AssignieeId  string // required when Type == LeaderboardUpdate
	WireId       string // required when Type == TaskStatusUpdate
	WireStatus   Status // required when Type == TaskStatusUpdate
	ComparatorId string // required when Type == ComparattorStatusUpdate
}
