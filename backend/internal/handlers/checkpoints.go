package handlers

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/druabello/airport-employee-management/internal/simulation"
)

const defaultDuration = 30 * time.Minute

func HandleCheckpoints(w http.ResponseWriter, r *http.Request) {
	const (
		maxRequestSize = 1 << 20 // 1 MiB
		maxRequestBody = 1 << 24 // 16 MiB
	)

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	defer r.Body.Close() // nolint:errcheck

	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Invalid gzip", http.StatusBadRequest)
			return
		}
		defer gz.Close() // nolint:errcheck

		r.Body = http.MaxBytesReader(w, gz, maxRequestBody)
	}

	if err := r.ParseMultipartForm(maxRequestBody); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not retrive file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close() // nolint:errcheck

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read csv: %v", err), http.StatusBadRequest)
		return
	}

	arrivals := make([]simulation.ArrivalGroup, 0, len(records))
	ag, err := csvToArrivalGroup(records[0])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse arrival group: %v", err), http.StatusBadRequest)
		return
	}
	arrivals = append(arrivals, ag)
	for i := 1; i < len(records); i++ {
		ag, err = csvToArrivalGroup(records[i])
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse arrival group: %v", err), http.StatusBadRequest)
			return
		}

		arrivals[i-1].Duration = ag.Start - arrivals[i-1].Start // Update duration of previous.

		arrivals = append(arrivals, ag)
	}

	arrivals[len(arrivals)-1].Duration = defaultDuration
}

func csvToArrivalGroup(s []string) (simulation.ArrivalGroup, error) {
	start, err := time.ParseDuration(s[0])
	if err != nil {
		return simulation.ArrivalGroup{}, fmt.Errorf("parsing start time: %w", err)
	}

	amount, err := strconv.Atoi(s[1])
	if err != nil {
		return simulation.ArrivalGroup{}, fmt.Errorf("parsing amount: %w", err)
	}

	return simulation.ArrivalGroup{
		Start:  start,
		Amount: amount,
	}, nil
}
