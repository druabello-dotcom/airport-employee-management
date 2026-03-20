package simulation

import (
	"container/heap"
	"errors"
	"time"
)

var ErrNoEvents = errors.New("no queued events")

type sim struct {
	maxWait          time.Duration
	timePerPassenger time.Duration

	arrivalQueue   eventQueue
	freeQueue      eventQueue
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
		freeQueue:        make(eventQueue, 0),
	}

	arrivalEvents := arrivalsToEvents(arrivals)

	s.arrivalQueue = make(eventQueue, len(arrivalEvents))

	// Initialize event queue with events.
	for i := range arrivalEvents {
		s.arrivalQueue[i] = &arrivalEvents[i]
	}
	heap.Init(&s.arrivalQueue)

	return s
}

// @return The minimum possible open checkpoints after the simulated event, and the time of the event.
func (s *sim) SimulateNextEvent() (Result, error) {
	var a, f *event
	if len(s.arrivalQueue) > 0 {
		a = s.arrivalQueue.front()
	}
	if len(s.freeQueue) > 0 {
		f = s.freeQueue.front()
	}

	if a == nil && f == nil {
		return Result{}, ErrNoEvents
	}

	var t time.Duration
	if f == nil {
		t = s.simulateArrival()
	} else if a == nil || f.time < a.time {
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
	e := heap.Pop(&s.arrivalQueue).(*event)
	t = e.time

	s.passengerQueue = append(s.passengerQueue, e.time)

	if len(s.freeQueue) == 0 {
		s.queueFree(e.time + s.timePerPassenger)
	}

	return
}

func (s *sim) simulateFree() (t time.Duration) {
	e := heap.Pop(&s.freeQueue).(*event)
	t = e.time

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
		s.queueFree(e.time + s.timePerPassenger)
	} else if len(s.freeQueue) == 0 {
		return
	}

	earliest := s.freeQueue.front().time
	deadline := s.passengerQueue[0] + s.maxWait

	if earliest > deadline {
		// We cannot afford to have this passenger wait for next available checkpoint.
		s.queueFree(e.time)
	}

	return
}

// Simulates all passengers in the queue to check if the current amount of checkpoints is
// sufficient to not keep anyone waiting more than maxWait.
func (s *sim) sufficientCheckpoints() bool {
	if len(s.freeQueue) == 0 {
		return false
	}

	checkpoints := make(eventQueue, len(s.freeQueue))
	copy(checkpoints, s.freeQueue)

	for _, t := range s.passengerQueue {
		deadline := t + s.maxWait
		earliest := heap.Pop(&checkpoints).(*event)

		if earliest.time > deadline {
			return false
		}

		heap.Push(&checkpoints, &event{
			time: earliest.time + s.timePerPassenger,
		})
	}

	return true
}

func (s *sim) queueFree(t time.Duration) {
	heap.Push(&s.freeQueue, &event{
		time: t,
	})
}

func (s *sim) IsFinished() bool {
	return len(s.arrivalQueue) == 0 && len(s.freeQueue) == 0
}
