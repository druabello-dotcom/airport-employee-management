package simulation

import (
	"time"
)

type ArrivalGroup struct {
	Start    time.Duration
	Duration time.Duration
	Amount   int
}

// @param arrivals Slice of arrivalGroup, sorted by start.
func ArrivalsToEvents(arrivals []ArrivalGroup) []event {
	events := make([]event, 0, len(arrivals))

	for _, a := range arrivals {
		// Calculate arrival times for all passengers, assuming uniform distribution.
		for i := range a.Amount {
			t := a.Duration.Nanoseconds() / int64(a.Amount) * int64(i)
			events = append(events, event{
				time: a.Start + time.Duration(t),
			})
		}
	}

	return events
}
