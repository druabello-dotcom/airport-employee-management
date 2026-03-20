package simulation

import (
	"container/heap"
	"errors"
	"fmt"
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

	s.arrivalQueue = arrivalsToTime(arrivals)
	heap.Init(&s.arrivalQueue)

	return s
}

// Runs the entire simulation and collects the results.
func (s *sim) Run() ([]Result, error) {
	results := make([]Result, 0)

	for !s.IsFinished() {
		res, err := s.SimulateNextEvent()
		if err != nil {
			return nil, fmt.Errorf("simulating event: %w", err)
		}

		results = append(results, res)
	}

	return results, nil
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
