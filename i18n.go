package main

import (
	"sort"
	"strings"
)

// Keys for the localizable UI strings. Every key must be present in the
// English table; other languages may omit a key and fall back to English.
const (
	keyScore         = "score"
	keyGuessLength   = "guess_length"
	keyInvalidWord   = "invalid_word"
	keyWin           = "win"
	keyLoss          = "loss"
	keyControlQuit   = "control_quit"
	keyControlResart = "control_restart"
)

// _stringKeys lists every key that the UI may ask for.
var _stringKeys = []string{
	keyScore,
	keyGuessLength,
	keyInvalidWord,
	keyWin,
	keyLoss,
	keyControlQuit,
	keyControlResart,
}

// _langDefault is the language used when none is requested, or when an unknown
// language is requested.
const _langDefault = "en"

// language holds everything that varies between locales: the UI copy and the
// word lists.
type language struct {
	Code    string            // "en", "cy"
	Name    string            // "English", "Cymraeg"
	Strings map[string]string // UI copy, keyed by the constants above
	Answers []string          // puzzle solutions
	Allowed []string          // accepted guesses, a superset of Answers

	dict Dictionary
}

// s returns the string for the given key, falling back to the English string
// if this language does not define it.
func (l *language) s(key string) string {
	if l != nil {
		if str, ok := l.Strings[key]; ok && str != "" {
			return str
		}
	}
	if fallback := languages[_langDefault]; fallback != nil && fallback != l {
		return fallback.Strings[key]
	}
	return ""
}

// dictionary returns the word list for this language.
func (l *language) dictionary() Dictionary {
	if l == nil {
		return languages[_langDefault].dict
	}
	return l.dict
}

var (
	_englishAnswers = EnglishDictionary.commonWords
	_englishAllowed = sortedWords(EnglishDictionary.allWords)

	_welshAnswers = parseWordList(_cyAnswers)
	_welshAllowed = parseWordList(_cyAllowed)
)

// languages is the registry of supported languages, keyed by language code.
var languages = map[string]*language{
	"en": {
		Code: "en",
		Name: "English",
		Strings: map[string]string{
			keyScore:         "Score: %d",
			keyGuessLength:   "Your guess must be a 5-letter word.",
			keyInvalidWord:   "That's not a valid word.",
			keyWin:           "You win!",
			keyLoss:          "The word was %s. Better luck next time!",
			keyControlQuit:   "quit",
			keyControlResart: "restart",
		},
		Answers: _englishAnswers,
		Allowed: _englishAllowed,
		dict:    EnglishDictionary,
	},
	"cy": {
		Code: "cy",
		Name: "Cymraeg",
		Strings: map[string]string{
			keyScore:         "Sgôr: %d",
			keyGuessLength:   "Rhaid i'ch cynnig fod yn air 5 llythyren.",
			keyInvalidWord:   "Nid yw hwnnw'n air dilys.",
			keyWin:           "Rydych chi wedi ennill!",
			keyLoss:          "%s oedd y gair. Pob lwc y tro nesaf!",
			keyControlQuit:   "gadael",
			keyControlResart: "ailddechrau",
		},
		Answers: _welshAnswers,
		Allowed: _welshAllowed,
		dict:    newDictionary(_welshAnswers, _welshAllowed),
	},
}

// normalizeLangCode canonicalizes a language code supplied by the user.
func normalizeLangCode(code string) string {
	return strings.ToLower(strings.TrimSpace(code))
}

// lookupLanguage returns the language for the given code, falling back to
// English if the code is unknown.
func lookupLanguage(code string) *language {
	if lang, ok := languages[normalizeLangCode(code)]; ok {
		return lang
	}
	return languages[_langDefault]
}

// isLanguage reports whether the given code names a supported language.
func isLanguage(code string) bool {
	_, ok := languages[normalizeLangCode(code)]
	return ok
}

// languageCodes returns the supported language codes, sorted.
func languageCodes() []string {
	codes := make([]string, 0, len(languages))
	for code := range languages {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return codes
}
