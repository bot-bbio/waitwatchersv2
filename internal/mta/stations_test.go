package mta

import (
	"slices"
	"testing"
)

func TestResolveStation(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"Times Sq-42 St", "127"},
		{"72 St", "123"},
		{"96 St", "120"},
		{"14 St-Union Sq", "635"},
		{"Court Sq G", "G22"},
	}

	for _, tt := range tests {
		got, err := ResolveStation(tt.name)
		if err != nil {
			t.Errorf("ResolveStation(%s) returned error: %v", tt.name, err)
			continue
		}
		if !slices.Contains(got, tt.expected) {
			t.Errorf("ResolveStation(%s) = %v; expected to contain %s", tt.name, got, tt.expected)
		}
	}

	_, err := ResolveStation("NonExistentStation")
	if err == nil {
		t.Error("expected error for non-existent station, got nil")
	}
}
