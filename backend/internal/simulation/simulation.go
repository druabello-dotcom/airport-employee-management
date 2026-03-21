package simulation

import (
	"container/heap"
	"time"
)

type sim struct {
	maxWait          time.Duration
	timePerPassenger time.Duration
}

type Result struct {
	Time    time.Duration `json:"time"`
	MinOpen int           `json:"minOpen"`
}

func New(maxWait, timePerPassenger time.Duration) *sim {
	s := &sim{
		maxWait:          maxWait,
		timePerPassenger: timePerPassenger,
	}

	return s
}

// Finds the minimum number of open checkpoints at the time of each arrival.
// @param arrivals The arrival times of the passengers, must be sorted in ascending order.
func (s *sim) Run(arrivals []time.Duration) []Result {
	results := make([]Result, len(arrivals))

	checkpoints := make(timeHeap, 0)
	for i, arrivalT := range arrivals {
		// Remove all idle checkpoints.
		for len(checkpoints) > 0 && checkpoints[0] < arrivalT {
			heap.Pop(&checkpoints)
		}

		deadline := arrivalT + s.maxWait
		canServeInTime := len(checkpoints) > 0 && checkpoints[0] <= deadline
		if canServeInTime {
			t := heap.Pop(&checkpoints).(time.Duration)
			heap.Push(&checkpoints, t+s.timePerPassenger)
		} else {
			// Need to open new checkpoint.
			heap.Push(&checkpoints, arrivalT+s.timePerPassenger)
		}

		results[i].Time = arrivalT
		results[i].MinOpen = len(checkpoints)
	}

	return results
}
