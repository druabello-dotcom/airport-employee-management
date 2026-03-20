package simulation

import (
	"container/heap"
	"errors"
	"sort"
	"time"
)

var ErrNoEvents = errors.New("no queued events")

type sim struct {
	maxWait          time.Duration
	timePerPassenger time.Duration

	arrivalQueue   timeHeap
	freeQueue      timeHeap
	passengerQueue []time.Duration // Arrival times of currently queued passengers.
}

type Result struct {
	Time    time.Duration `json:"time"`
	MinOpen int           `json:"minOpen"`
}

func New(maxWait, timePerPassenger time.Duration, arrivals []ArrivalGroup) *sim {
	s := &sim{
		maxWait:          maxWait,
		timePerPassenger: timePerPassenger,
		freeQueue:        make(timeHeap, 0),
	}

	s.arrivalQueue = ArrivalsToTime(arrivals)
	heap.Init(&s.arrivalQueue)

	return s
}

// Finds the minimum number of open checkpoints at the time of each arrival.
// @param arrivals The arrival times of the passengers, must be sorted in ascending order.
func (s *sim) Run(arrivals []time.Duration) []Result {
	results := make([]Result, len(arrivals))

	checkpoints := make([]time.Duration, 0)
	for i, arrivalT := range arrivals {
		// Remove all idle checkpoints.
		for len(checkpoints) > 0 && checkpoints[0] < arrivalT {
			checkpoints = checkpoints[1:]
		}

		deadline := arrivalT + s.maxWait
		canServeInTime := len(checkpoints) > 0 && checkpoints[0] <= deadline
		if canServeInTime {
			t := checkpoints[0]
			checkpoints = checkpoints[1:]
			checkpoints = insertSorted(checkpoints, max(t, arrivalT)+s.timePerPassenger)
		} else {
			// Need to open new checkpoint.
			checkpoints = insertSorted(checkpoints, arrivalT+s.timePerPassenger)
		}

		results[i].Time = arrivalT
		results[i].MinOpen = len(checkpoints)
	}

	return results
}

func insertSorted(checkpoints []time.Duration, t time.Duration) []time.Duration {
	i := sort.Search(len(checkpoints), func(i int) bool {
		return checkpoints[i] >= t
	})

	checkpoints = append(checkpoints, 0)     // Increase size by one to be able to shift.
	copy(checkpoints[i+1:], checkpoints[i:]) // Shift slice after i.
	checkpoints[i] = t

	return checkpoints
}

// @return The minimum possible open checkpoints after the simulated event, and the time of the event.
func (s *sim) SimulateNextEvent() (Result, error) {
	var a, f *time.Duration
	if len(s.arrivalQueue) > 0 {
		a = &s.arrivalQueue[0]
	}
	if len(s.freeQueue) > 0 {
		f = &s.freeQueue[0]
	}

	if a == nil && f == nil {
		return Result{}, ErrNoEvents
	}

	var t time.Duration
	if f == nil {
		t = s.simulateArrival()
	} else if a == nil || *f < *a {
		t = s.simulateFree()
	} else {
		t = s.simulateArrival()
	}

	return Result{
		Time:    t,
		MinOpen: len(s.freeQueue),
	}, nil
}

func (s *sim) simulateArrival() (t time.Duration) {
	t = heap.Pop(&s.arrivalQueue).(time.Duration)

	s.passengerQueue = append(s.passengerQueue, t)

	if len(s.freeQueue) == 0 {
		s.queueFree(t + s.timePerPassenger)
	}

	return
}

func (s *sim) simulateFree() (t time.Duration) {
	t = heap.Pop(&s.freeQueue).(time.Duration)

	if len(s.passengerQueue) == 0 {
		return
	}

	s.passengerQueue = s.passengerQueue[1:]

	if len(s.passengerQueue) == 0 {
		return
	}

	// Assume the checkpoint will close, then we can correct if wrong.
	// This avoids reopening and then closing the same checkpoint.
	if !s.sufficientCheckpoints() {
		// We must reopen this one to avoid exceeding maxWait.
		s.queueFree(t + s.timePerPassenger)
	} else if len(s.freeQueue) == 0 {
		return
	}

	earliest := s.freeQueue[0]
	deadline := s.passengerQueue[0] + s.maxWait

	if earliest > deadline {
		// We cannot afford to have this passenger wait for next available checkpoint.
		s.queueFree(t)
	}

	return
}

// Simulates all passengers in the queue to check if the current amount of checkpoints is
// sufficient to not keep anyone waiting more than maxWait.
func (s *sim) sufficientCheckpoints() bool {
	if len(s.freeQueue) == 0 {
		return false
	}

	checkpoints := make(timeHeap, len(s.freeQueue))
	copy(checkpoints, s.freeQueue)

	for _, t := range s.passengerQueue {
		deadline := t + s.maxWait
		earliest := heap.Pop(&checkpoints).(time.Duration)

		if earliest > deadline {
			return false
		}

		heap.Push(&checkpoints, earliest+s.timePerPassenger)
	}

	return true
}

func (s *sim) queueFree(t time.Duration) {
	heap.Push(&s.freeQueue, t)
}

func (s *sim) IsFinished() bool {
	return len(s.arrivalQueue) == 0 && len(s.freeQueue) == 0
}
