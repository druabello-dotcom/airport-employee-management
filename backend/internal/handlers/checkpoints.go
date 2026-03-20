package handlers

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/druabello/airport-employee-management/internal/simulation"
)

const defaultDuration = 30 * time.Minute

type checkpointsReq struct {
	MaxWait          duration `json:"maxWait"`
	ResultInterval   duration `json:"resultInterval"`
	TimePerPassenger duration `json:"timePerPassenger"`
}

// Wrapper to allow unmarshalling json into a duration.
type duration struct {
	time.Duration
}

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

	jsonStr := r.FormValue("config")
	var config checkpointsReq
	if err := json.Unmarshal([]byte(jsonStr), &config); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
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

	arrivals := make([]simulation.ArrivalGroup, len(records)-1)
	if err = arrivals[0].ParseFromCSV(records[1]); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse arrival group: %v", err), http.StatusBadRequest)
		return
	}
	for i := 1; i < len(records)-1; i++ {
		if err = arrivals[i].ParseFromCSV(records[i+1]); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse arrival group: %v", err), http.StatusBadRequest)
			return
		}

		// Update duration of previous.
		arrivals[i-1].Duration = arrivals[i].Start - arrivals[i-1].Start
	}
	arrivals[len(arrivals)-1].Duration = defaultDuration

	sim := simulation.New(config.MaxWait.Duration, config.TimePerPassenger.Duration, arrivals)

	arrivalTimes := simulation.ArrivalsToTime(arrivals)
	simRes := sim.Run(arrivalTimes)

	resp := simulation.FindIntervalMaximums(simRes, config.ResultInterval.Duration)

	j, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failure creating response", http.StatusInternalServerError)
		log.Printf("Failed creating checkpoints resp: %v\n", err)
		return
	}

	if _, err := w.Write(j); err != nil {
		log.Printf("Error writing checkpoints resp: %v\n", err)
	}
}

func (d *duration) UnmarshalJSON(b []byte) error {
	var j any
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	switch value := j.(type) {
	case float64:
		d.Duration = time.Duration(value)
	case string:
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid duration: %#v", j)
	}

	return nil
}
