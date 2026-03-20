package simulation

import "time"

type timeHeap []time.Duration

func (q timeHeap) Len() int {
	return len(q)
}

func (q timeHeap) Less(i, j int) bool {
	return q[i] < q[j]
}

func (q timeHeap) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *timeHeap) Push(x any) {
	e := x.(time.Duration)
	*q = append(*q, e)
}

func (q *timeHeap) Pop() any {
	n := len(*q)
	e := (*q)[n-1]
	*q = (*q)[0 : n-1]
	return e
}
