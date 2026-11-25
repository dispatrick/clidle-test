package main

import (
	"math/rand"
	"testing"
)

const _numDraws = 10

func drawWords(rng *rand.Rand, n int) []string {
	words := make([]string, n)
	for i := range words {
		words[i] = EnglishDictionary.GetRandomCommonWord(rng)
	}
	return words
}

func TestGetRandomCommonWordSeededIsDeterministic(t *testing.T) {
	a := drawWords(rand.New(rand.NewSource(42)), _numDraws)
	b := drawWords(rand.New(rand.NewSource(42)), _numDraws)
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("draw %d differs for the same seed: %q != %q", i, a[i], b[i])
		}
	}
}

func TestGetRandomCommonWordDifferentSeeds(t *testing.T) {
	a := drawWords(rand.New(rand.NewSource(1)), _numDraws)
	b := drawWords(rand.New(rand.NewSource(2)), _numDraws)
	for i := range a {
		if a[i] != b[i] {
			return
		}
	}
	t.Fatalf("seeds 1 and 2 produced the same %d words: %v", _numDraws, a)
}

func TestGetRandomCommonWordNilRng(t *testing.T) {
	word := EnglishDictionary.GetRandomCommonWord(nil)
	if !EnglishDictionary.IsWord(word) {
		t.Fatalf("nil rng returned %q, which is not a dictionary word", word)
	}
	if len(word) != _numChars {
		t.Fatalf("nil rng returned %q, want %d characters", word, _numChars)
	}
}
