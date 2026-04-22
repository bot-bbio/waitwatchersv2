package mta

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/molus/mach/internal/models"
	"google.golang.org/protobuf/proto"
)

// Client handles communication with the MTA GTFS-RT API.
type Client struct {
	URLs []string
}

// NewClient creates a new MTA client for multiple endpoints.
func NewClient(urls ...string) *Client {
	return &Client{URLs: urls}
}

// Fetch retrieves and decodes the GTFS-RT feed from a specific URL.
func (c *Client) Fetch(ctx context.Context, url string) (*gtfs.FeedMessage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	feed := &gtfs.FeedMessage{}
	if err := proto.Unmarshal(body, feed); err != nil {
		return nil, fmt.Errorf("failed to unmarshal GTFS-RT feed: %w", err)
	}

	return feed, nil
}

// GetPredictions fetches from all URLs concurrently and converts them to model predictions.
func (c *Client) GetPredictions(ctx context.Context) ([]models.Prediction, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allPredictions []models.Prediction
	var firstErr error

	for _, url := range c.URLs {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			
			feed, err := c.Fetch(ctx, u)
			if err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}

			var preds []models.Prediction
			for _, entity := range feed.Entity {
				if entity.TripUpdate == nil {
					continue
				}

				tripID := *entity.TripUpdate.Trip.TripId
				routeID := *entity.TripUpdate.Trip.RouteId
				for _, update := range entity.TripUpdate.StopTimeUpdate {
					if update.StopId == nil || (update.Arrival == nil && update.Departure == nil) {
						continue
					}

					p := models.Prediction{
						TrainID:   tripID,
						StationID: *update.StopId,
						Line:      routeID,
					}

					if update.Arrival != nil && update.Arrival.Time != nil {
						p.ArrivalTime = time.Unix(*update.Arrival.Time, 0)
					}
					if update.Departure != nil && update.Departure.Time != nil {
						p.DepartureTime = time.Unix(*update.Departure.Time, 0)
					}

					preds = append(preds, p)
				}
			}

			mu.Lock()
			allPredictions = append(allPredictions, preds...)
			mu.Unlock()
		}(url)
	}

	wg.Wait()

	if firstErr != nil && len(allPredictions) == 0 {
		return nil, firstErr
	}

	return allPredictions, nil
}
