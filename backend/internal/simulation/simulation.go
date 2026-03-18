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

func New(maxWait time.Duration, arrivals []ArrivalGroup) *sim {
	s := &sim{
		maxWait:      maxWait,
		arrivalQueue: make(eventQueue, len(arrivals)),
		freeQueue:    make(eventQueue, 0),
	}

	arrivalEvents := ArrivalsToEvents(arrivals)

	// Initialize event queue with events.
	for i := range arrivalEvents {
		s.arrivalQueue[i] = &arrivalEvents[i]
	}
	heap.Init(&s.arrivalQueue)

	return s
}

// @return The minimum possible open checkpoints after the simulated event.
func (s *sim) SimulateNextEvent() (int, error) {
	var a, f *event
	if len(s.arrivalQueue) > 0 {
		a = s.arrivalQueue.front()
	}
	if len(s.freeQueue) > 0 {
		f = s.freeQueue.front()
	}

	if a == nil && f == nil {
		return 0, ErrNoEvents
	}

	if f == nil {
		s.simulateArrival()
	} else if a == nil || f.time < a.time {
		s.simulateFree()
	} else {
		s.simulateArrival()
	}

	return len(s.freeQueue), nil
}

func (s *sim) simulateArrival() {
	e := heap.Pop(&s.arrivalQueue).(*event)

	s.passengerQueue = append(s.passengerQueue, e.time)

	if len(s.freeQueue) == 0 {
		s.openCheckpoint(e.time)
		return
	}
}

func (s *sim) simulateFree() {
	e := heap.Pop(&s.freeQueue).(*event)

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
		s.openCheckpoint(e.time + s.timePerPassenger)
	} else if len(s.freeQueue) == 0 {
		return
	}

	earliest := s.freeQueue.front().time
	deadline := s.passengerQueue[0] + s.maxWait

	if earliest > deadline {
		// We cannot afford to have this passenger wait for next available checkpoint.
		s.openCheckpoint(e.time)
	}
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

func (s *sim) openCheckpoint(t time.Duration) {
	heap.Push(&s.freeQueue, &event{
		time: t,
	})
}

func (s *sim) IsFinished() bool {
	return len(s.arrivalQueue) == 0 && len(s.freeQueue) == 0
}
