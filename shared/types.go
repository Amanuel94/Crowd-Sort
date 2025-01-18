package shared

import (
	"fmt"

	"github.com/Amanuel94/crowdsort/interfaces"
	"golang.org/x/exp/constraints"
)

// enum for score
type Ord = int

const (
	NA Ord = iota
	LT
	EQ
	GT
)

type Status = string
type MessageType = int
type ComparatorStatus = string

func Assigned(assignee string) Status {
	return Status(fmt.Sprintf("ASSIGNED TO: %s", assignee))
}

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

// Wrapper for indexing items
type Wire[T any] struct {
	index  string
	value  interfaces.Comparable[T]
	status Status
}
type OrderedType[T constraints.Ordered] struct {
	index any
	value T
}
type ComparatorModule[T any] struct {
	pid      string
	cmp      CmpFunc[T]
	task_cnt int
	status   ComparatorStatus
}

type Connector[T any] struct {
	Id          string
	F           string
	S           string
	Order       Ord
	AssignieeId string
}

type CmpFunc[T any] func(*interfaces.Comparable[T], *interfaces.Comparable[T]) (int, error)

type PingMessage struct {
	Type         MessageType
	F            string // when the message is about leaderboard update
	S            string // when the message is about leaderboard update
	AssignieeId  string // when the message is about leaderboard update
	WireId       string // when the message is about status update
	WireStatus   Status // when the message is about status update
	ComparatorId string // when the message is about status update
}
