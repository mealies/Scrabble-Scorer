package scoring

import (
	"testing"
)

func TestCalculateWordScore(t *testing.T) {
	tests := []struct {
		name        string
		word        string
		multipliers []string
		isBingo     bool
		want        int
	}{
		{
			name:        "Simple word",
			word:        "HELLO",
			multipliers: []string{"none", "none", "none", "none", "none"},
			isBingo:     false,
			want:        8, // H(4)+E(1)+L(1)+L(1)+O(1) = 8
		},
		{
			name:        "Double Letter",
			word:        "HELLO",
			multipliers: []string{"dl", "none", "none", "none", "none"},
			isBingo:     false,
			want:        12, // H(4*2)+E(1)+L(1)+L(1)+O(1) = 12
		},
		{
			name:        "Triple Word",
			word:        "HELLO",
			multipliers: []string{"none", "none", "none", "none", "tw"},
			isBingo:     false,
			want:        24, // (H(4)+E(1)+L(1)+L(1)+O(1)) * 3 = 24
		},
		{
			name:        "Bingo only",
			word:        "SCRABBLE",
			multipliers: nil,
			isBingo:     true,
			want:        64, // S(1)+C(3)+R(1)+A(1)+B(3)+B(3)+L(1)+E(1) = 14 + 50 = 64
		},
		{
			name:        "Multipliers and Bingo",
			word:        "QUIZZES",
			multipliers: []string{"none", "none", "none", "dl", "none", "none", "none"},
			isBingo:     true,
			want:        94, // Q(10)+U(1)+I(1)+Z(10*2)+Z(10)+E(1)+S(1) = 44; 44 + 50 = 94
		},
		{
			name:        "Empty word",
			word:        "",
			multipliers: nil,
			isBingo:     false,
			want:        0,
		},
		{
			name:        "Case insensitive word",
			word:        "hello",
			multipliers: nil,
			isBingo:     false,
			want:        8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateWordScore(tt.word, tt.multipliers, tt.isBingo); got != tt.want {
				t.Errorf("CalculateWordScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateWordScoreComplex(t *testing.T) {
	// "QUIZZES" with DL on first Z
	// Q(10), U(1), I(1), Z(10), Z(10), E(1), S(1)
	// Base: 10+1+1+10+10+1+1 = 34
	// With DL on first Z (index 3): 10+1+1+(10*2)+10+1+1 = 44
	// With Bingo: 44 + 50 = 94
	got := CalculateWordScore("QUIZZES", []string{"none", "none", "none", "dl", "none", "none", "none"}, true)
	if got != 94 {
		t.Errorf("CalculateWordScore(QUIZZES, dl on Z, bingo) = %d, want 94", got)
	}

	// "AX" with TW
	// A(1), X(8) = 9. 9 * 3 = 27
	got = CalculateWordScore("AX", []string{"none", "tw"}, false)
	if got != 27 {
		t.Errorf("CalculateWordScore(AX, tw on X) = %d, want 27", got)
	}
}
