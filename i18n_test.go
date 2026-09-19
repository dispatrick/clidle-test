package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

// TestWordLists checks the integrity of every registered language's word
// lists.
func TestWordLists(t *testing.T) {
	for code, lang := range languages {
		code, lang := code, lang
		t.Run(code, func(t *testing.T) {
			if len(lang.Answers) < 50 {
				t.Errorf("answer list is too short: %d words", len(lang.Answers))
			}
			if len(lang.Allowed) < len(lang.Answers) {
				t.Errorf("allowed list (%d) is smaller than answer list (%d)",
					len(lang.Allowed), len(lang.Answers))
			}

			checkWords(t, "answers", lang.Answers)
			checkWords(t, "allowed", lang.Allowed)

			for _, answer := range lang.Answers {
				if !lang.dictionary().IsWord(answer) {
					t.Errorf("answer %q is not an accepted guess", answer)
				}
			}
		})
	}
}

// checkWords verifies that a word list is sorted, deduplicated, and made up of
// correctly sized uppercase A-Z words.
func checkWords(t *testing.T, name string, words []string) {
	t.Helper()

	if len(words) == 0 {
		t.Errorf("%s list is empty", name)
		return
	}
	if !sort.StringsAreSorted(words) {
		t.Errorf("%s list is not sorted", name)
	}

	seen := make(map[string]struct{}, len(words))
	for _, word := range words {
		if len(word) != _numChars {
			t.Errorf("%s: %q is not %d characters long", name, word, _numChars)
		}
		for _, r := range word {
			if !isAsciiUpper(r) {
				t.Errorf("%s: %q contains a non A-Z character %q", name, word, r)
				break
			}
		}
		if _, ok := seen[word]; ok {
			t.Errorf("%s: %q is duplicated", name, word)
		}
		seen[word] = struct{}{}
	}
}

// TestStringsComplete checks that every UI string resolves to a non-empty
// value in every language, either directly or via the English fallback.
func TestStringsComplete(t *testing.T) {
	english := languages[_langDefault]
	if english == nil {
		t.Fatalf("default language %q is not registered", _langDefault)
	}

	for _, key := range _stringKeys {
		if english.Strings[key] == "" {
			t.Errorf("English is missing a value for key %q", key)
		}
	}

	for code, lang := range languages {
		for _, key := range _stringKeys {
			if lang.s(key) == "" {
				t.Errorf("language %q resolves key %q to an empty string", code, key)
			}
		}
	}
}

// TestStringsFormatVerbs checks that localized format strings keep the same
// verbs as the English original, so formatting them cannot go wrong.
func TestStringsFormatVerbs(t *testing.T) {
	if got := fmt.Sprintf(languages["cy"].s(keyScore), 42); !strings.Contains(got, "42") {
		t.Errorf("Welsh score string did not format the score: %q", got)
	}
	if got := fmt.Sprintf(languages["cy"].s(keyLoss), "BWRDD"); !strings.Contains(got, "BWRDD") {
		t.Errorf("Welsh loss string did not format the answer: %q", got)
	}
}

// TestStringsFallBackToEnglish checks that a missing translation falls back to
// the English string instead of rendering an empty label.
func TestStringsFallBackToEnglish(t *testing.T) {
	partial := &language{Code: "xx", Name: "Test", Strings: map[string]string{}}
	if got, want := partial.s(keyWin), languages[_langDefault].Strings[keyWin]; got != want {
		t.Errorf("fallback returned %q, want %q", got, want)
	}
}

func TestLookupLanguage(t *testing.T) {
	tests := []struct {
		code  string
		want  string
		known bool
	}{
		{"en", "en", true},
		{"cy", "cy", true},
		{"CY", "cy", true},
		{" cy ", "cy", true},
		{"fr", "en", false},
		{"", "en", false},
		{"not-a-language", "en", false},
	}

	for _, test := range tests {
		if got := lookupLanguage(test.code).Code; got != test.want {
			t.Errorf("lookupLanguage(%q) = %q, want %q", test.code, got, test.want)
		}
		if got := isLanguage(test.code); got != test.known {
			t.Errorf("isLanguage(%q) = %v, want %v", test.code, got, test.known)
		}
	}
}

// TestWelshDigraphs checks that digraphs are treated as two ordinary letters.
func TestWelshDigraphs(t *testing.T) {
	welsh := languages["cy"]
	for _, word := range []string{"BWRDD", "LLAWN"} {
		if !welsh.dictionary().IsWord(word) {
			t.Errorf("%q should be an accepted Welsh guess", word)
		}
	}
}
