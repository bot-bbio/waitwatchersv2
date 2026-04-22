package engine

import (
	"testing"
	"time"

	"github.com/molus/mach/internal/models"
)

func TestCalculateWaitDelta(t *testing.T) {
	now := time.Now()

	// Scenario: Comparing Line 1 vs Line 2
	preds := []models.Prediction{
		{TrainID: "L1", StationID: "127", ArrivalTime: now.Add(2 * time.Minute), Line: "1"},
		{TrainID: "L1", StationID: "137", ArrivalTime: now.Add(12 * time.Minute), Line: "1"},
		{TrainID: "E1", StationID: "127", ArrivalTime: now.Add(5 * time.Minute), Line: "2"},
		{TrainID: "E1", StationID: "137", ArrivalTime: now.Add(10 * time.Minute), Line: "2"},
	}

	res, err := CalculateWaitDelta(preds, []string{"127"}, []string{"137"})
	if err != nil {
		t.Fatalf("CalculateWaitDelta failed: %v", err)
	}

	if len(res.Options) < 2 {
		t.Fatalf("expected at least 2 options, got %d", len(res.Options))
	}

	// Option 0 (Line 2): 10 mins
	// Option 1 (Line 1): 12 mins
	// Delta: 10 - 12 = -2 mins
	expected := -2 * time.Minute
	if res.WaitDelta != expected {
		t.Errorf("expected delta %v, got %v", expected, res.WaitDelta)
	}
	
	if res.Options[0].Line != "2" || res.Options[1].Line != "1" {
		t.Errorf("expected lines 2 and 1, got %s and %s", res.Options[0].Line, res.Options[1].Line)
	}
}

func TestStressTestTimesSqTo168St(t *testing.T) {
	now := time.Now()
	// Times Sq (127/R16/725) -> 168 St (112/A09)
	// Comparing A train vs 1 train.
	preds := []models.Prediction{
		{TrainID: "A1", StationID: "A27", ArrivalTime: now.Add(2 * time.Minute), Line: "A"}, // Times Sq A
		{TrainID: "A1", StationID: "A09", ArrivalTime: now.Add(15 * time.Minute), Line: "A"}, // 168 St A
		{TrainID: "1_1", StationID: "127", ArrivalTime: now.Add(1 * time.Minute), Line: "1"}, // Times Sq 1
		{TrainID: "1_1", StationID: "112", ArrivalTime: now.Add(20 * time.Minute), Line: "1"}, // 168 St 1
	}

	res, err := CalculateWaitDelta(preds, []string{"127", "A27", "R16", "725"}, []string{"112", "A09"})
	if err != nil {
		t.Fatalf("CalculateWaitDelta failed: %v", err)
	}

	if res.Options[0].Line != "A" {
		t.Errorf("expected A train to be fastest, got %s", res.Options[0].Line)
	}
	if res.WaitDelta != -5*time.Minute {
		t.Errorf("expected -5m delta, got %v", res.WaitDelta)
	}
}

func TestStressTestGrandCentralToYankeeStadium(t *testing.T) {
	now := time.Now()
	// Grand Central (631/723) -> 161 St-Yankee Stadium (414/D11)
	// Comparing 4 train vs D train.
	preds := []models.Prediction{
		{TrainID: "4_1", StationID: "631", ArrivalTime: now.Add(3 * time.Minute), Line: "4"}, // GC 4
		{TrainID: "4_1", StationID: "414", ArrivalTime: now.Add(18 * time.Minute), Line: "4"}, // 161 St 4
		{TrainID: "D1", StationID: "D17", ArrivalTime: now.Add(5 * time.Minute), Line: "D"}, // 34 St Herald Sq (proxy for GC user walking/transferring)
		// Wait, the stress test says "waiting at Grand Central".
		// Actually D train doesn't stop at Grand Central. User must transfer or we only consider GC lines.
		// Stress test says: "is it faster to take the 4 or the B/D train".
		// This implies the user is at a location where both are accessible, or we compare routes from GC.
		// For the purpose of the test, let's assume GC complexes include the transfer logic or we use the resolved IDs.
		{TrainID: "D1", StationID: "D15", ArrivalTime: now.Add(6 * time.Minute), Line: "D"}, // 47-50 Sts
		{TrainID: "D1", StationID: "D11", ArrivalTime: now.Add(20 * time.Minute), Line: "D"}, // 161 St D
	}

	// We'll pass the station IDs for Grand Central and 161 St.
	res, err := CalculateWaitDelta(preds, []string{"631", "723", "D15"}, []string{"414", "D11"})
	if err != nil {
		t.Fatalf("CalculateWaitDelta failed: %v", err)
	}

	if res.Options[0].Line != "4" {
		t.Errorf("expected 4 train to be fastest, got %s", res.Options[0].Line)
	}
}

func TestCalculateWaitDeltaOnlyOneRoute(t *testing.T) {
	now := time.Now()
	preds := []models.Prediction{
		{TrainID: "1_TR", StationID: "127", ArrivalTime: now.Add(1 * time.Minute), Line: "1"},
		{TrainID: "1_TR", StationID: "112", ArrivalTime: now.Add(20 * time.Minute), Line: "1"},
	}

	res, err := CalculateWaitDelta(preds, []string{"127"}, []string{"112"})
	if err != nil {
		t.Fatalf("CalculateWaitDelta failed: %v", err)
	}

	if res.Options[0].Line != "1" {
		t.Errorf("expected line 1, got %s", res.Options[0].Line)
	}
	if res.WaitDelta != 0 {
		t.Errorf("expected 0 delta, got %v", res.WaitDelta)
	}
}

func TestCalculateWaitDeltaNoRoutes(t *testing.T) {
	_, err := CalculateWaitDelta(nil, []string{"127"}, []string{"137"})
	if err == nil {
		t.Error("expected error for no routes, got nil")
	}
}
