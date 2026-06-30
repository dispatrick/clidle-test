package main

import (
	"strings"
	"testing"
)

func word(s string) [_numChars]byte {
	var w [_numChars]byte
	copy(w[:], s)
	return w
}

func TestRowStates(t *testing.T) {
	tests := []struct {
		name   string
		answer string
		guess  string
		want   [_numChars]keyState
	}{
		{
			name:   "all correct",
			answer: "ROBOT",
			guess:  "ROBOT",
			want:   [_numChars]keyState{_keyStateCorrect, _keyStateCorrect, _keyStateCorrect, _keyStateCorrect, _keyStateCorrect},
		},
		{
			name:   "present and absent",
			answer: "ROBOT",
			guess:  "OTHER",
			want:   [_numChars]keyState{_keyStatePresent, _keyStatePresent, _keyStateAbsent, _keyStateAbsent, _keyStatePresent},
		},
		{
			name:   "duplicate letters matched at most once",
			answer: "ABBEY",
			guess:  "BABES",
			want:   [_numChars]keyState{_keyStatePresent, _keyStatePresent, _keyStateCorrect, _keyStateCorrect, _keyStateAbsent},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &model{answer: word(tt.answer)}
			got := m.rowStates(word(tt.guess))
			if got != tt.want {
				t.Errorf("rowStates(%q) against %q = %v, want %v", tt.guess, tt.answer, got, tt.want)
			}
		})
	}
}

func TestViewShare(t *testing.T) {
	m := &model{answer: word("ROBOT"), gameOver: true, gridRow: 2}
	m.grid[0] = word("OTHER")
	m.grid[1] = word("ROBOT")

	got := m.viewShare()

	for _, want := range []string{
		"clidle 2/6",
		"🟨🟨⬛⬛🟨",
		"🟩🟩🟩🟩🟩",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("viewShare() = %q, want it to contain %q", got, want)
		}
	}

	// Share text must be copyable verbatim, so no line may carry alignment padding.
	for _, line := range strings.Split(got, "\n") {
		if line != strings.TrimLeft(line, " ") {
			t.Errorf("viewShare() line %q has leading space padding", line)
		}
	}
}

func TestViewShareLoss(t *testing.T) {
	m := &model{answer: word("ROBOT"), gameOver: true, gridRow: 1}
	m.grid[0] = word("OTHER")

	if got := m.viewShare(); !strings.Contains(got, "clidle X/6") {
		t.Errorf("viewShare() for a loss = %q, want it to contain %q", got, "clidle X/6")
	}
}
