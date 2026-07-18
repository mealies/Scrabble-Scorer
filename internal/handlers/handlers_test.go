package handlers

import (
	"encoding/json"
	"testing"

	"github.com/mealies/wasmScrabbleScorer/internal/models"
)

func TestSerialization(t *testing.T) {
	// Verify that models serialize as expected for the frontend
	p := &models.Player{
		Name:  "Alice",
		Score: 100,
		WordHistory: []models.RoundRecord{
			{Words: []string{"HELLO"}, Score: 8},
		},
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	// Check if key names match what the JS expects (some are capitalized, some are not)
	if _, ok := m["Name"]; !ok {
		t.Errorf("Expected key 'Name', got %v", m)
	}
	if _, ok := m["wordHistory"]; !ok {
		t.Errorf("Expected key 'wordHistory', got %v", m)
	}
}
