package models

import (
	"testing"
)

func TestPlayerIncreaseScore(t *testing.T) {
	p := &Player{
		Name:  "Test",
		Score: 10,
	}

	p.IncreaseScore(5, []string{"WORD"})

	if p.Score != 15 {
		t.Errorf("Expected score 15, got %d", p.Score)
	}

	if len(p.WordHistory) != 1 {
		t.Errorf("Expected word history length 1, got %d", len(p.WordHistory))
	}

	if p.WordHistory[0].Score != 5 {
		t.Errorf("Expected round score 5, got %d", p.WordHistory[0].Score)
	}

	if p.WordHistory[0].Words[0] != "WORD" {
		t.Errorf("Expected word WORD, got %s", p.WordHistory[0].Words[0])
	}
}
