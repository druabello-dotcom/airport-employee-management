package simulation

import (
	"fmt"
	"strconv"
	"time"
)

type ArrivalGroup struct {
	Start    time.Duration
	Duration time.Duration
	Amount   int
}

func (ag *ArrivalGroup) ParseFromCSV(s []string) error {
	start, err := time.ParseDuration(s[0])
	if err != nil {
		return fmt.Errorf("parsing start time: %w", err)
	}

	amount, err := strconv.Atoi(s[1])
	if err != nil {
		return fmt.Errorf("parsing amount: %w", err)
	}

	ag.Start = start
	ag.Amount = amount

	return nil
}

// @param arrivals Slice of arrivalGroup, sorted by start.
func arrivalsToEvents(arrivals []ArrivalGroup) []event {
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
