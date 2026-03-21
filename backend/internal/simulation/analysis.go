package simulation

import "time"

// @param data Slice of Result, is modified such that it will be empty after the function call.
func FindIntervalMaximums(data []Result, interval time.Duration) []Result {
	res := make([]Result, 0)

	for i := 0; len(data) > 0; i++ {
		start := interval * time.Duration(i)
		mx, j := findIntervalMaximum(data, start, interval)

		res = append(res, Result{
			Time:    start,
			MinOpen: mx,
		})

		data = data[j:]
	}

	return res
}

// @return The maximum in the interval, and the first index outside of the interval.
func findIntervalMaximum(data []Result, start, interval time.Duration) (int, int) {
	var mx int

	for i, r := range data {
		if r.Time-start > interval {
			return mx, i
		}

		mx = max(mx, r.MinOpen)
	}

	return mx, len(data)
}
