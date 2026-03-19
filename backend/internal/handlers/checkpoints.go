package handlers

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"net/http"
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

	arrivals := make([]simulation.ArrivalGroup, len(records))
	if err = arrivals[0].ParseFromCSV(records[0]); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse arrival group: %v", err), http.StatusBadRequest)
		return
	}
	for i := 1; i < len(records); i++ {
		if err = arrivals[i].ParseFromCSV(records[i]); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse arrival group: %v", err), http.StatusBadRequest)
			return
		}

		// Update duration of previous.
		arrivals[i-1].Duration = arrivals[i].Start - arrivals[i-1].Start
	}

	arrivals[len(arrivals)-1].Duration = defaultDuration
}
