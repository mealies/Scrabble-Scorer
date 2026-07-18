package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/mealies/wasmScrabbleScorer/internal/models"
	"github.com/mealies/wasmScrabbleScorer/internal/scoring"
)

func HandleStart(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	if r.Method != "POST" {
		fsthttp.Error(w, "Method not allowed", fsthttp.StatusMethodNotAllowed)
		return
	}

	var names []string
	if err := json.NewDecoder(r.Body).Decode(&names); err != nil {
		fsthttp.Error(w, err.Error(), fsthttp.StatusBadRequest)
		return
	}

	if len(names) < 2 || len(names) > 4 {
		fsthttp.Error(w, "Scrabble requires 2-4 players", fsthttp.StatusBadRequest)
		return
	}

	players := make([]*models.Player, len(names))
	for i, name := range names {
		players[i] = &models.Player{
			Name:              name,
			Score:             0,
			LastWord:          "",
			WordHistory:       []models.RoundRecord{},
			CurrentRoundWords: []string{},
		}
	}

	game := &models.Game{Players: players}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(fsthttp.StatusOK)
	json.NewEncoder(w).Encode(game)
}

func HandleScore(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	if r.Method != "POST" {
		fsthttp.Error(w, "Method not allowed", fsthttp.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Game        *models.Game `json:"game"`
		PlayerIndex int          `json:"playerIndex"`
		Input       string       `json:"input"`
		Multipliers []string     `json:"multipliers"`
		IsBingo     bool         `json:"isBingo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fsthttp.Error(w, err.Error(), fsthttp.StatusBadRequest)
		return
	}

	if req.Game == nil {
		fsthttp.Error(w, "No game in progress", fsthttp.StatusBadRequest)
		return
	}

	game := req.Game
	if req.PlayerIndex < 0 || req.PlayerIndex >= len(game.Players) {
		fsthttp.Error(w, "Invalid player index", fsthttp.StatusBadRequest)
		return
	}

	var score int
	val, err := strconv.Atoi(req.Input)
	if err == nil {
		score = val
		game.Players[req.PlayerIndex].LastWord = req.Input
	} else {
		score = scoring.CalculateWordScore(req.Input, req.Multipliers, req.IsBingo)
		game.Players[req.PlayerIndex].LastWord = req.Input
	}

	game.Players[req.PlayerIndex].CurrentRoundScore += score
	game.Players[req.PlayerIndex].CurrentRoundWords = append(game.Players[req.PlayerIndex].CurrentRoundWords, req.Input)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(fsthttp.StatusOK)
	json.NewEncoder(w).Encode(game)
}

func HandleEndRound(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	if r.Method != "POST" {
		fsthttp.Error(w, "Method not allowed", fsthttp.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Game        *models.Game `json:"game"`
		PlayerIndex int          `json:"playerIndex"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fsthttp.Error(w, err.Error(), fsthttp.StatusBadRequest)
		return
	}

	if req.Game == nil {
		fsthttp.Error(w, "No game in progress", fsthttp.StatusBadRequest)
		return
	}

	game := req.Game
	if req.PlayerIndex < 0 || req.PlayerIndex >= len(game.Players) {
		fsthttp.Error(w, "Invalid player index", fsthttp.StatusBadRequest)
		return
	}

	p := game.Players[req.PlayerIndex]
	p.IncreaseScore(p.CurrentRoundScore, p.CurrentRoundWords)
	p.CurrentRoundScore = 0
	p.CurrentRoundWords = []string{}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(fsthttp.StatusOK)
	json.NewEncoder(w).Encode(game)
}

func HandleFinish(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	if r.Method != "POST" {
		fsthttp.Error(w, "Method not allowed", fsthttp.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Game        *models.Game `json:"game"`
		WinnerIndex int          `json:"winnerIndex"`
		Leftovers   []int        `json:"leftovers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fsthttp.Error(w, err.Error(), fsthttp.StatusBadRequest)
		return
	}

	if req.Game == nil {
		fsthttp.Error(w, "No game in progress", fsthttp.StatusBadRequest)
		return
	}

	game := req.Game
	if req.WinnerIndex < 0 || req.WinnerIndex >= len(game.Players) {
		fsthttp.Error(w, "Invalid winner index", fsthttp.StatusBadRequest)
		return
	}

	for i := range game.Players {
		if i == req.WinnerIndex {
			continue
		}
		points := 0
		if i < len(req.Leftovers) {
			points = req.Leftovers[i]
		}
		game.Players[i].Score -= points
		game.Players[req.WinnerIndex].Score += points
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(fsthttp.StatusOK)
	json.NewEncoder(w).Encode(game)
}

func HandleStatus(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("null"))
}
