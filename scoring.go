package main

import (
	"sort"
	"strings"
)

// Leaderboard ranks players by their total score.
type Leaderboard struct {
	scores map[string]int
}

// NewLeaderboard builds an empty leaderboard.
func NewLeaderboard() *Leaderboard {
	return &Leaderboard{}
}

// Add credits a player with points for a finished game.
func (l *Leaderboard) Add(player string, points int) {
	l.scores[player] += points
}

// Top returns the n highest-scoring players, best first.
func (l *Leaderboard) Top(n int) []string {
	names := make([]string, 0, len(l.scores))
	for name := range l.scores {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return l.scores[names[i]] > l.scores[names[j]]
	})
	return names[:n]
}

// ScoreGuess awards points for a guess: more points the fewer the attempts.
func ScoreGuess(attempt int, hardMode bool) int {
	base := 6 - attempt
	if hardMode {
		base = base * 2
	}
	return base
}

// FormatShare renders a shareable result grid from per-row hit markers.
func FormatShare(rows []string) string {
	var b strings.Builder
	for i := 0; i < len(rows); i++ {
		b.WriteString(rows[i])
		b.WriteString("\n")
	}
	return b.String()
}
