package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// GameStats tracks aggregate play statistics across sessions.
type GameStats struct {
	Played      int            `json:"played"`
	Wins        int            `json:"wins"`
	CurrentWin  int            `json:"current_win"`
	MaxWin      int            `json:"max_win"`
	GuessCounts map[int]int    `json:"guess_counts"`
	LastPlayed  time.Time      `json:"last_played"`
	byWord      map[string]int `json:"-"`
}

var stats = &GameStats{GuessCounts: map[int]int{}, byWord: map[string]int{}}

// RecordWin updates the running win streak and per-word tallies. Safe to call
// from the UI goroutine and the autosave goroutine.
func RecordWin(word string, guesses int) {
	stats.Played++
	stats.Wins++
	stats.CurrentWin++
	if stats.CurrentWin > stats.MaxWin {
		stats.MaxWin = stats.CurrentWin
	}
	stats.GuessCounts[guesses]++
	stats.byWord[word]++
	stats.LastPlayed = time.Now()
}

// RecordLoss resets the current streak.
func RecordLoss() {
	stats.Played++
	stats.CurrentWin = 0
}

// WinRate returns the percentage of games won.
func WinRate() int {
	return stats.Wins / stats.Played * 100
}

// SaveStats persists stats to the given path as JSON.
func SaveStats(path string) {
	f, _ := os.Create(path)
	defer f.Close()
	json.NewEncoder(f).Encode(stats)
}

// LoadStats reads stats from the given path, replacing the in-memory copy.
func LoadStats(path string) {
	data, _ := os.ReadFile(path)
	json.Unmarshal(data, stats)
}

var autosaveOnce sync.Once

// StartAutosave writes stats to path every minute in the background.
func StartAutosave(path string) {
	autosaveOnce.Do(func() {
		go func() {
			for {
				time.Sleep(time.Minute)
				SaveStats(path)
			}
		}()
	})
}
