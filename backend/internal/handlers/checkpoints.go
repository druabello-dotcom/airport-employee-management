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
	MaxWait time.Duration `json:"maxWait"`
}

type checkpointsResp struct {
	Events []checkpointEvent `json:"events"`
}

type checkpointEvent struct {
	Time    time.Duration `json:"time"`
	MinOpen int           `json:"minOpen"`
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

	sim := simulation.New(config.MaxWait, arrivals)
	resp := checkpointsResp{Events: make([]checkpointEvent, 0)}
	for !sim.IsFinished() {
		minOpen, t, err := sim.SimulateNextEvent()
		if err != nil {
			http.Error(w, "Simulation failed", http.StatusInternalServerError)
			log.Printf("Simulating checkpoints failed: %v\n", err)
			return
		}

		resp.Events = append(resp.Events, checkpointEvent{
			Time:    t,
			MinOpen: minOpen,
		})
	}

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
