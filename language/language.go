// Package language contains the user-facing strings for every language that
// clidle supports. English is the default, and is used as a fallback whenever
// an unknown language is requested.
package language

import "strings"

// Messages holds every user-facing string in the game. Using a struct instead
// of a map means a missing translation is a compile-time error.
type Messages struct {
	// Score is a format string for the default status line. It takes the
	// current score as its only argument.
	Score string
	// GuessTooShort is shown when the guess is not long enough.
	GuessTooShort string
	// GuessNotAWord is shown when the guess is not in the dictionary.
	GuessNotAWord string
	// Win is shown when the player guesses the answer.
	Win string
	// Loss is a format string shown when the player runs out of guesses. It
	// takes the answer as its only argument.
	Loss string
	// KeyEnter is the label of the "enter" key on the keyboard.
	KeyEnter string
	// KeyDelete is the label of the "delete" key on the keyboard.
	KeyDelete string
	// ControlQuit is the label of the "ctrl+c" control.
	ControlQuit string
	// ControlRestart is the label of the "ctrl+r" control.
	ControlRestart string
}

// Language is a language supported by the game.
type Language struct {
	// Code is the ISO 639-1 code of the language.
	Code string
	// Name is the English name of the language.
	Name string
	// Messages are the user-facing strings of the language.
	Messages Messages
}

// English is the default language.
var English = Language{
	Code:     "en",
	Name:     "English",
	Messages: messagesEnglish,
}

// Catalan is the Catalan language.
var Catalan = Language{
	Code:     "ca",
	Name:     "Catalan",
	Messages: messagesCatalan,
}

// Languages contains every supported language, keyed by language code.
var Languages = map[string]Language{
	English.Code: English,
	Catalan.Code: Catalan,
}

// Lookup returns the language matching the given code. Matching is
// case-insensitive, and region subtags (such as the "ES" in "ca-ES") are
// ignored. The English name of the language ("catalan") is also accepted. If
// the language is not supported, English is returned along with false.
func Lookup(code string) (Language, bool) {
	code = strings.TrimSpace(code)
	if code == "" {
		return English, false
	}

	// Strip the region and encoding subtags, e.g. "ca_ES.UTF-8" -> "ca".
	if idx := strings.IndexAny(code, "-_."); idx != -1 {
		code = code[:idx]
	}
	code = strings.ToLower(code)

	if lang, ok := Languages[code]; ok {
		return lang, true
	}
	for _, lang := range Languages {
		if strings.ToLower(lang.Name) == code {
			return lang, true
		}
	}
	return English, false
}
