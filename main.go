package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fastly/compute-sdk-go/fsthttp"
)

//go:embed static/index.html
var indexPage []byte

var letterValues = map[rune]int{
	'_': 0,
	'A': 1, 'E': 1, 'I': 1, 'O': 1, 'U': 1, 'L': 1, 'N': 1, 'R': 1, 'S': 1, 'T': 1,
	'D': 2, 'G': 2,
	'B': 3, 'C': 3, 'M': 3, 'P': 3,
	'F': 4, 'H': 4, 'V': 4, 'W': 4, 'Y': 4,
	'K': 5,
	'J': 8, 'X': 8,
	'Q': 10, 'Z': 10,
}

func CalculateWordScore(word string, multipliers []string) int {
	totalScore := 0
	wordMultiplier := 1

	runes := []rune(strings.ToUpper(word))
	for i, char := range runes {
		letterValue, ok := letterValues[char]
		if !ok {
			continue
		}

		mult := "none"
		if i < len(multipliers) {
			mult = multipliers[i]
		}

		switch strings.ToLower(mult) {
		case "dl":
			totalScore += letterValue * 2
		case "tl":
			totalScore += letterValue * 3
		case "dw":
			totalScore += letterValue
			wordMultiplier *= 2
		case "tw":
			totalScore += letterValue
			wordMultiplier *= 3
		default:
			totalScore += letterValue
		}
	}

	return totalScore * wordMultiplier
}

type Player struct {
	Name              string `json:"Name"`
	Score             int    `json:"Score"`
	LastWord          string `json:"LastWord"`
	WordHistory       []int  `json:"wordHistory"`
	CurrentRoundScore int    `json:"currentRoundScore"`
}

func (p *Player) IncreaseScore(s int) {
	p.Score += s
	p.WordHistory = append(p.WordHistory, s)
}

type Game struct {
	Players []*Player `json:"Players"`
}

func handleStart(w fsthttp.ResponseWriter, r *fsthttp.Request) {
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

	players := make([]*Player, len(names))
	for i, name := range names {
		players[i] = &Player{Name: name, Score: 0, LastWord: "", WordHistory: []int{}}
	}

	game := &Game{Players: players}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(fsthttp.StatusOK)
	json.NewEncoder(w).Encode(game)
}

func handleScore(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	if r.Method != "POST" {
		fsthttp.Error(w, "Method not allowed", fsthttp.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Game        *Game    `json:"game"`
		PlayerIndex int      `json:"playerIndex"`
		Input       string   `json:"input"`
		Multipliers []string `json:"multipliers"`
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
		score = CalculateWordScore(req.Input, req.Multipliers)
		game.Players[req.PlayerIndex].LastWord = req.Input
	}

	game.Players[req.PlayerIndex].CurrentRoundScore += score
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(fsthttp.StatusOK)
	json.NewEncoder(w).Encode(game)
}

func handleEndRound(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	if r.Method != "POST" {
		fsthttp.Error(w, "Method not allowed", fsthttp.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Game        *Game `json:"game"`
		PlayerIndex int   `json:"playerIndex"`
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
	p.IncreaseScore(p.CurrentRoundScore)
	p.CurrentRoundScore = 0

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(fsthttp.StatusOK)
	json.NewEncoder(w).Encode(game)
}

func handleFinish(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	if r.Method != "POST" {
		fsthttp.Error(w, "Method not allowed", fsthttp.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Game        *Game `json:"game"`
		WinnerIndex int   `json:"winnerIndex"`
		Leftovers   []int `json:"leftovers"`
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

func handleStatus(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("null"))
}

func main() {
	// Log service version
	fmt.Println("FASTLY_SERVICE_VERSION:", os.Getenv("FASTLY_SERVICE_VERSION"))

	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		switch r.URL.Path {
		case "/api/start":
			handleStart(w, r)
		case "/api/score":
			handleScore(w, r)
		case "/api/end-round":
			handleEndRound(w, r)
		case "/api/finish":
			handleFinish(w, r)
		case "/api/status":
			handleStatus(w, r)
		case "/":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			io.Copy(w, bytes.NewReader(indexPage))
		default:
			// Catch all other requests and return a 404.
			w.WriteHeader(fsthttp.StatusNotFound)
			fmt.Fprintf(w, "The page you requested could not be found\n")
		}
	})
}
