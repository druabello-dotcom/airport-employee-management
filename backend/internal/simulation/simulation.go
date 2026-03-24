package simulation

import (
	"container/heap"
	"time"
)

type sim struct {
	maxCheckpoints   int
	maxWait          time.Duration
	timePerPassenger time.Duration
}

type Result struct {
	Time       time.Duration `json:"time"`
	TimeWaited time.Duration `json:"wait"`
	MinOpen    int           `json:"minOpen"`
}

func New(maxCheckpoints int, maxWait, timePerPassenger time.Duration) *sim {
	s := &sim{
		maxCheckpoints:   maxCheckpoints,
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

		var serveTime time.Duration

		deadline := arrivalT + s.maxWait
		canServeInTime := len(checkpoints) > 0 && checkpoints[0] <= deadline
		if !canServeInTime && len(checkpoints) < s.maxCheckpoints {
			// Need to open new checkpoint.
			heap.Push(&checkpoints, arrivalT+s.timePerPassenger)
			serveTime = arrivalT
		} else {
			t := heap.Pop(&checkpoints).(time.Duration)
			heap.Push(&checkpoints, t+s.timePerPassenger)
			serveTime = t
		}

		results[i].Time = arrivalT
		results[i].MinOpen = len(checkpoints)
		results[i].TimeWaited = serveTime - arrivalT
	}

	return results
}

// Checks whether maxWait will be exceeded if using checkpointCnt checkpoints over the next
// maxWait time interval of arrivals.
func (s *sim) exceedsMaxWait(checkpointCnt int, checkpoints timeHeap, arrivals []time.Duration) bool {
	c := make(timeHeap, len(checkpoints))
	copy(c, checkpoints)

	for len(c) > checkpointCnt {
		heap.Pop(&c)
	}

	for _, a := range arrivals {
		if a > arrivals[0]+s.maxWait {
			return false
		}

		deadline := a + s.maxWait
		canServeInTime := len(c) > 0 && c[0] <= deadline
		if canServeInTime {
			heap.Pop(&c)
			heap.Push(&c, a+s.timePerPassenger)
		} else {
			return true
		}
	}

	return false
}
