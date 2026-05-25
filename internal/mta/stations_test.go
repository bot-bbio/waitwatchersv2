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
		// Lookup with brackets
		{"42 St-Times Sq / Port Authority [1,2,3,7,A,C,E,N,Q,R,W,S]", "127"},
		{"42 St-Times Sq / Port Authority [1,2,3,7,A,C,E,N,Q,R,W,S]", "A27"},
		// Lookup without brackets (CLI style)
		{"42 St-Times Sq / Port Authority", "127"},
		{"72 St", "123"},
		{"72 St [1,2,3]", "123"},
		{"72 St Q", "Q03"},
		{"72 St Q [Q]", "Q03"},
		{"96 St", "120"},
		{"96 St Q", "Q05"},
		{"14 St-Union Sq", "635"},
		{"Court Sq", "G22"},
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
