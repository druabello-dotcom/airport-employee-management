package simulation

import "time"

// @param data Slice of Result, is modified such that it will be empty after the function call.
func FindIntervalMaximums(data []Result, interval time.Duration) []Result {
	res := make([]Result, 0)

	for i := 0; len(data) > 0; i++ {
		start := interval * time.Duration(i)
		r, j := findIntervalMaximum(data, start, interval)

		r.Time = start
		res = append(res, r)

		data = data[j:]
	}

	return res
}

// @return The maximum in the interval, and the first index outside of the interval.
func findIntervalMaximum(data []Result, start, interval time.Duration) (Result, int) {
	var mxOpen int
	var mxWait time.Duration

	for i, r := range data {
		if r.Time-start > interval {
			return Result{
				TimeWaited: mxWait,
				MinOpen:    mxOpen,
			}, i
		}

		mxOpen = max(mxOpen, r.MinOpen)
		mxWait = max(mxWait, r.TimeWaited)
	}

	return Result{
		TimeWaited: mxWait,
		MinOpen:    mxOpen,
	}, len(data)
}
