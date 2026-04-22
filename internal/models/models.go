package models

import "time"

// Station represents a subway station.
type Station struct {
	ID   string
	Name string
}

// Train represents a specific train and its status.
type Train struct {
	ID        string
	Line      string // e.g., "1", "2", "3", "A", "C", "E"
	IsExpress bool
}

// Prediction represents a predicted arrival/departure at a station.
type Prediction struct {
	TrainID       string
	StationID     string
	Line          string // e.g., "1", "2"
	IsExpress     bool
	ArrivalTime   time.Time
	DepartureTime time.Time
}

