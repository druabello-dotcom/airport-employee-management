package simulation

import (
	"time"
)

type arrivalGroup struct {
	start    time.Duration
	duration time.Duration
	amount   int
}

// @param arrivals Slice of arrivalGroup, sorted by start.
func ArrivalsToEvents(arrivals []arrivalGroup) []event {
	events := make([]event, 0, len(arrivals))

	for _, a := range arrivals {
		// Calculate arrival times for all passengers, assuming uniform distribution.
		for i := range a.amount {
			t := a.duration.Nanoseconds() / int64(a.amount) * int64(i)
			events = append(events, event{
				time: a.start + time.Duration(t),
			})
		}
	}

	return events
}
