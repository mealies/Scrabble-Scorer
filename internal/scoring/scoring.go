package scoring

import (
	"strings"
)

var LetterValues = map[rune]int{
	'_': 0,
	'A': 1, 'E': 1, 'I': 1, 'O': 1, 'U': 1, 'L': 1, 'N': 1, 'R': 1, 'S': 1, 'T': 1,
	'D': 2, 'G': 2,
	'B': 3, 'C': 3, 'M': 3, 'P': 3,
	'F': 4, 'H': 4, 'V': 4, 'W': 4, 'Y': 4,
	'K': 5,
	'J': 8, 'X': 8,
	'Q': 10, 'Z': 10,
}

func CalculateWordScore(word string, multipliers []string, isBingo bool) int {
	totalScore := 0
	wordMultiplier := 1

	runes := []rune(strings.ToUpper(word))
	for i, char := range runes {
		letterValue, ok := LetterValues[char]
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

	score := totalScore * wordMultiplier
	if isBingo {
		score += 50
	}
	return score
}
