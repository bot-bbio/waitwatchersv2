package models

import "testing"

func TestNewStation(t *testing.T) {
	s := Station{ID: "123", Name: "Times Sq - 42 St"}
	if s.ID != "123" {
		t.Errorf("expected ID 123, got %s", s.ID)
	}
	if s.Name != "Times Sq - 42 St" {
		t.Errorf("expected Name Times Sq - 42 St, got %s", s.Name)
	}
}

func TestNewTrain(t *testing.T) {
	tr := Train{ID: "T1", Line: "1", IsExpress: false}
	if tr.ID != "T1" {
		t.Errorf("expected ID T1, got %s", tr.ID)
	}
	if tr.Line != "1" {
		t.Errorf("expected Line 1, got %s", tr.Line)
	}
	if tr.IsExpress != false {
		t.Errorf("expected IsExpress false, got %v", tr.IsExpress)
	}
}
