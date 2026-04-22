package mta

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

func TestClientFetchMock(t *testing.T) {
	// Create a sample FeedMessage
	feed := &gtfs.FeedMessage{
		Header: &gtfs.FeedHeader{
			GtfsRealtimeVersion: proto.String("2.0"),
			Timestamp:           proto.Uint64(uint64(time.Now().Unix())),
		},
		Entity: []*gtfs.FeedEntity{
			{
				Id: proto.String("1"),
				TripUpdate: &gtfs.TripUpdate{
					Trip: &gtfs.TripDescriptor{
						TripId: proto.String("trip1"),
						RouteId: proto.String("1"),
					},
					StopTimeUpdate: []*gtfs.TripUpdate_StopTimeUpdate{
						{
							StopId: proto.String("101N"),
							Arrival: &gtfs.TripUpdate_StopTimeEvent{
								Time: proto.Int64(time.Now().Add(5 * time.Minute).Unix()),
							},
						},
					},
				},
			},
		},
	}

	data, err := proto.Marshal(feed)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer ts.Close()

	c := NewClient(ts.URL)
	preds, err := c.GetPredictions(context.Background())
	if err != nil {
		t.Fatalf("GetPredictions failed: %v", err)
	}

	if len(preds) != 1 {
		t.Errorf("expected 1 prediction, got %d", len(preds))
	}

	p := preds[0]
	if p.TrainID != "trip1" {
		t.Errorf("expected trip1, got %s", p.TrainID)
	}
	if p.StationID != "101N" {
		t.Errorf("expected 101N, got %s", p.StationID)
	}
}

func TestGetPredictionsMultipleURLs(t *testing.T) {
	// Setup two mock servers
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := &gtfs.FeedMessage{
			Header: &gtfs.FeedHeader{GtfsRealtimeVersion: proto.String("2.0")},
			Entity: []*gtfs.FeedEntity{{
				Id: proto.String("1"),
				TripUpdate: &gtfs.TripUpdate{
					Trip: &gtfs.TripDescriptor{TripId: proto.String("T1"), RouteId: proto.String("A")},
					StopTimeUpdate: []*gtfs.TripUpdate_StopTimeUpdate{{
						StopId: proto.String("A01"),
						Arrival: &gtfs.TripUpdate_StopTimeEvent{Time: proto.Int64(time.Now().Unix())},
					}}},
			}},
		}
		data, _ := proto.Marshal(feed)
		w.Write(data)
	}))
	defer ts1.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		feed := &gtfs.FeedMessage{
			Header: &gtfs.FeedHeader{GtfsRealtimeVersion: proto.String("2.0")},
			Entity: []*gtfs.FeedEntity{{
				Id: proto.String("2"),
				TripUpdate: &gtfs.TripUpdate{
					Trip: &gtfs.TripDescriptor{TripId: proto.String("T2"), RouteId: proto.String("1")},
					StopTimeUpdate: []*gtfs.TripUpdate_StopTimeUpdate{{
						StopId: proto.String("101"),
						Arrival: &gtfs.TripUpdate_StopTimeEvent{Time: proto.Int64(time.Now().Unix())},
					}}},
			}},
		}
		data, _ := proto.Marshal(feed)
		w.Write(data)
	}))
	defer ts2.Close()

	c := NewClient(ts1.URL, ts2.URL)
	preds, err := c.GetPredictions(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(preds) != 2 {
		t.Errorf("expected 2 predictions, got %d", len(preds))
	}

	lines := make(map[string]bool)
	for _, p := range preds {
		lines[p.Line] = true
	}
	if !lines["A"] || !lines["1"] {
		t.Errorf("expected lines A and 1, got %v", lines)
	}
}

func TestClientFetchError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := NewClient(ts.URL)
	_, err := c.Fetch(context.Background(), ts.URL)
	if err == nil {
		t.Error("expected error for 500 status code, got nil")
	}
}

func TestClientFetchInvalidURL(t *testing.T) {
	c := NewClient(":")
	_, err := c.Fetch(context.Background(), ":")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}
