package simulation

import "time"

const (
	eventArrival eventType = iota
	eventFree
)

type eventType int

type event struct {
	time time.Duration
	kind eventType
}

type eventQueue []*event

func (q eventQueue) Len() int {
	return len(q)
}

func (q eventQueue) Less(i, j int) bool {
	return q[i].time < q[j].time
}

func (q eventQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *eventQueue) Push(x any) {
	e := x.(*event)
	*q = append(*q, e)
}

func (q *eventQueue) Pop() any {
	n := len(*q)
	e := (*q)[n-1]
	*q = (*q)[0 : n-1]
	return e
}
