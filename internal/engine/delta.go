package engine

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/molus/mach/internal/models"
)

// CalculationResult holds the arrival times and the calculated delta between lines.
type CalculationResult struct {
	Options   []RouteOption
	WaitDelta time.Duration
}

// RouteOption represents a specific line's best arrival time.
type RouteOption struct {
	Line    string
	Arrival time.Time
}

// CalculateWaitDelta identifies all viable routes and calculates the delta between the top two.
// It accepts a list of station ID slices (one for each complex) to allow cross-complex comparison.
func CalculateWaitDelta(preds []models.Prediction, stationSets ...[]string) (CalculationResult, error) {
	if len(stationSets) < 2 || len(stationSets)%2 != 0 {
		return CalculationResult{}, fmt.Errorf("invalid station sets: must provide at least one (origin set, dest set) pair")
	}

	now := time.Now()
	optionsMap := make(map[string]time.Time)

	// For each origin set/dest set pair provided, find the earliest arrival for each line.
	for i := 0; i < len(stationSets); i += 2 {
		origins := stationSets[i]
		dests := stationSets[i+1]

		lineArrivals := findLineArrivals(preds, origins, dests, now)
		for line, arrival := range lineArrivals {
			if existing, ok := optionsMap[line]; !ok || arrival.Before(existing) {
				optionsMap[line] = arrival
			}
		}
	}

	if len(optionsMap) == 0 {
		return CalculationResult{}, fmt.Errorf("no valid routes found for the given stations")
	}

	// Convert map to slice and sort by arrival time.
	options := make([]RouteOption, 0, len(optionsMap))
	for line, arrival := range optionsMap {
		options = append(options, RouteOption{Line: line, Arrival: arrival})
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].Arrival.Before(options[j].Arrival)
	})

	res := CalculationResult{
		Options: options,
	}

	// Calculate WaitDelta if we have at least two options.
	if len(options) >= 2 {
		res.WaitDelta = options[0].Arrival.Sub(options[1].Arrival)
	}

	return res, nil
}

func findLineArrivals(preds []models.Prediction, originIDs, destIDs []string, now time.Time) map[string]time.Time {
	// First, identify all train IDs that stop at any of the origin IDs in the FUTURE.
	trainsAtOrigin := make(map[string]time.Time)
	for _, p := range preds {
		for _, originID := range originIDs {
			if strings.HasPrefix(p.StationID, originID) && p.ArrivalTime.After(now) {
				trainsAtOrigin[p.TrainID] = p.ArrivalTime
				break
			}
		}
	}

	// Now find the earliest arrival time at any of the destination IDs for each LINE.
	bestLineArrivals := make(map[string]time.Time)
	for _, p := range preds {
		for _, destID := range destIDs {
			if strings.HasPrefix(p.StationID, destID) {
				originArrival, ok := trainsAtOrigin[p.TrainID]
				if !ok {
					continue
				}

				if !p.ArrivalTime.After(originArrival) {
					continue
				}

				if current, ok := bestLineArrivals[p.Line]; !ok || p.ArrivalTime.Before(current) {
					bestLineArrivals[p.Line] = p.ArrivalTime
				}
				break
			}
		}
	}

	return bestLineArrivals
}
