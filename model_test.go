package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// letters for CRANE: C(1) R(2) A(3) N(4) E(5).
func TestHardModeViolation(t *testing.T) {
	tests := []struct {
		name    string
		answer  string
		guesses []string // rows already played
		guess   string   // the candidate guess
		want    string
	}{
		{
			name:   "no previous guesses accepts anything",
			answer: "CRANE",
			guess:  "STOMP",
			want:   "",
		},
		{
			name:    "green letter must stay in position",
			answer:  "CRANE",
			guesses: []string{"CLOSE"}, // C green (pos 1), E green (pos 5)
			guess:   "TRACE",           // C moved off position 1
			want:    "Guess must use C in position 1.",
		},
		{
			name:    "greens kept in position are allowed",
			answer:  "CRANE",
			guesses: []string{"CLOSE"},
			guess:   "CRANE",
			want:    "",
		},
		{
			name:    "present letter must be reused somewhere",
			answer:  "CRANE",
			guesses: []string{"AGENT"}, // A present, E present, N green (pos 4)
			guess:   "SHINE",           // keeps N (pos 4) and E, but drops A
			want:    "Guess must contain A.",
		},
		{
			name:    "present letter reused in a different position is allowed",
			answer:  "CRANE",
			guesses: []string{"AGENT"}, // A present, E present, N green (pos 4)
			guess:   "PLANE",           // A reused at pos 3, N kept, E kept
			want:    "",
		},
		{
			name:    "absent letters impose no constraint",
			answer:  "CRANE",
			guesses: []string{"MUDDY"}, // none present
			guess:   "STOMP",
			want:    "",
		},
		{
			name:    "constraints accumulate across rows",
			answer:  "CRANE",
			guesses: []string{"MUDDY", "AGENT"}, // A/E present, N green from row 2
			guess:   "SHINE",                    // drops A
			want:    "Guess must contain A.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &model{hardMode: true}
			copy(m.answer[:], tt.answer)
			for _, g := range tt.guesses {
				copy(m.grid[m.gridRow][:], g)
				m.gridRow++
			}
			var next [_numChars]byte
			copy(next[:], tt.guess)
			assert.Equal(t, tt.want, m.hardModeViolation(next))
		})
	}
}

func TestHardModeViolationDisabledWhenOff(t *testing.T) {
	m := &model{hardMode: false}
	copy(m.answer[:], "CRANE")
	copy(m.grid[0][:], "CLOSE")
	m.gridRow = 1

	var next [_numChars]byte
	copy(next[:], "STOMP") // violates green C, but hard mode is off
	assert.Equal(t, "", m.hardModeViolation(next))
}
