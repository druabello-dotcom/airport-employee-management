package simulation

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrNoTimeColumn   = errors.New("no time column in csv")
	ErrNoPeopleColumn = errors.New("no people column in csv")
)

type ArrivalGroup struct {
	Start    time.Duration
	Duration time.Duration
	Amount   int
}

func (ag *ArrivalGroup) ParseFromCSV(s []string, columnToIdx map[string]int) error {
	timeIdx, ok := columnToIdx["time"]
	if !ok {
		return ErrNoTimeColumn
	}

	start, err := time.ParseDuration(s[timeIdx])
	if err != nil {
		return fmt.Errorf("parsing start time: %w", err)
	}

	peopleIdx, ok := columnToIdx["people"]
	if !ok {
		return ErrNoPeopleColumn
	}

	amount, err := strconv.Atoi(s[peopleIdx])
	if err != nil {
		return fmt.Errorf("parsing amount: %w", err)
	}

	ag.Start = start
	ag.Amount = amount

	return nil
}

// @param arrivals Slice of arrivalGroup, sorted by start.
func ArrivalsToTime(arrivals []ArrivalGroup) []time.Duration {
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
