package models

type RoundRecord struct {
	Words []string `json:"words"`
	Score int      `json:"score"`
}

type Player struct {
	Name              string        `json:"Name"`
	Score             int           `json:"Score"`
	LastWord          string        `json:"LastWord"`
	WordHistory       []RoundRecord `json:"wordHistory"`
	CurrentRoundScore int           `json:"currentRoundScore"`
	CurrentRoundWords []string      `json:"currentRoundWords"`
}

func (p *Player) IncreaseScore(s int, words []string) {
	p.Score += s
	p.WordHistory = append(p.WordHistory, RoundRecord{Words: words, Score: s})
}

type Game struct {
	Players []*Player `json:"Players"`
}
