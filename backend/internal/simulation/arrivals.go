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
func arrivalsToTime(arrivals []ArrivalGroup) []time.Duration {
	times := make([]time.Duration, 0, len(arrivals))

	for _, a := range arrivals {
		// Calculate arrival times for all passengers, assuming uniform distribution.
		for i := range a.Amount {
			t := a.Duration / time.Duration(a.Amount) * time.Duration(i)
			times = append(times, a.Start+t)
		}
	}

	return times
}
