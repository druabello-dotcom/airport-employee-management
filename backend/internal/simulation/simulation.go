package simulation

import (
	"container/heap"
	"time"
)

type sim struct {
	maxWait          time.Duration
	timePerPassenger time.Duration

	arrivalQueue   eventQueue
	freeQueue      eventQueue
	passengerQueue []time.Duration // Arrival times of currently queued passengers.
}

func New(maxWait time.Duration, arrivals []arrivalGroup) *sim {
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

func (s *sim) SimulateNextEvent() {
	a := s.arrivalQueue.front()
	f := s.freeQueue.front()

	if a.time < f.time {
		s.simulateArrival()
	} else {
		s.simulateFree()
	}
}

func (s *sim) simulateArrival() {

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

	if len(s.freeQueue) == 0 {
		// There are still passengers queued, but no checkpoint open.
		heap.Push(&s.freeQueue, &event{
			time: e.time,
			kind: eventFree,
		})
		return
	}

	nextFree := s.freeQueue.front().time
	remWait := s.maxWait - (e.time - s.passengerQueue[0])

	if remWait < nextFree-e.time {
		// We cannot afford to have this passenger wait.
		heap.Push(&s.freeQueue, &event{
			time: e.time,
			kind: eventFree,
		})
	}
}

func (s *sim) IsFinished() bool {
	return len(s.arrivalQueue) == 0 && len(s.freeQueue) == 0
}
