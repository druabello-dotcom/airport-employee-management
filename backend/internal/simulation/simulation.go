package simulation

import "container/heap"

type sim struct {
	q eventQueue
}

func New(events []event) *sim {
	s := &sim{
		q: make(eventQueue, len(events)),
	}

	// Initialize event queue with events.
	for i := range events {
		s.q[i] = &events[i]
	}
	heap.Init(&s.q)

	return s
}

func (s *sim) SimulateNextEvent() {
	e := s.q.Pop().(*event)

	switch e.kind {
	case eventArrival:
		s.simulateArrival()

	case eventFree:
		s.simulateFree()
	}
}

func (s *sim) simulateArrival() {

}

func (s *sim) simulateFree() {

}

func (s *sim) IsFinished() bool {
	return len(s.q) == 0
}
