package main

import (
	_ "embed"
	"sort"
	"strings"
)

// Spanish word lists. Both lists are lowercase, one word per line, and only
// contain the letters A-Z plus Ñ (accents are stripped, see normalizeES).
var (
	//go:embed words_es_answers.txt
	wordsESAnswers string

	//go:embed words_es_allowed.txt
	wordsESAllowed string
)

// _defaultLanguage is the language code used when none is requested, or when
// an unknown one is requested.
const _defaultLanguage = "en"

// uiStrings holds every user-visible string of the game.
type uiStrings struct {
	// score is a format string taking the current score.
	score string
	// guessTooShort is shown when an incomplete guess is submitted.
	guessTooShort string
	// invalidWord is shown when the guess is not in the dictionary.
	invalidWord string
	// win is shown when the answer has been guessed.
	win string
	// loss is a format string taking the answer, shown when the guesses run out.
	loss string
	// controlQuit and controlRestart label the controls hint.
	controlQuit    string
	controlRestart string
}

// keyboardRow is a single row of the on-screen keyboard. Keys that are a
// single rune are letters, anything longer is a label (such as "ENTER").
type keyboardRow struct {
	keys    []string
	padding int
}

// language bundles everything that varies between languages: the word lists,
// the on-screen keyboard, the UI strings, and the letter folding rules.
type language struct {
	code       string
	name       string
	dictionary Dictionary
	keyboard   []keyboardRow
	strings    uiStrings
	// fold maps an input rune to its canonical (uppercase, unaccented) form.
	fold func(rune) rune
}

// normalize folds an entire word to its canonical form.
func (l language) normalize(word string) string {
	return strings.Map(l.fold, word)
}

// isLetter reports whether the given rune is a letter of this language's
// alphabet, as shown on the on-screen keyboard.
func (l language) isLetter(r rune) bool {
	for _, row := range l.keyboard {
		for _, key := range row.keys {
			if key == string(r) {
				return true
			}
		}
	}
	return false
}

// languages is the registry of all supported languages.
var languages = map[string]language{
	"en": {
		code:       "en",
		name:       "English",
		dictionary: EnglishDictionary,
		keyboard: []keyboardRow{
			{keys: []string{"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P"}, padding: 2},
			{keys: []string{"A", "S", "D", "F", "G", "H", "J", "K", "L"}, padding: 4},
			{keys: []string{"ENTER", "Z", "X", "C", "V", "B", "N", "M", "DELETE"}, padding: 0},
		},
		strings: uiStrings{
			score:          "Score: %d",
			guessTooShort:  "Your guess must be a 5-letter word.",
			invalidWord:    "That's not a valid word.",
			win:            "You win!",
			loss:           "The word was %s. Better luck next time!",
			controlQuit:    "quit",
			controlRestart: "restart",
		},
		fold: toAsciiUpper,
	},
	"es": {
		code:       "es",
		name:       "Español",
		dictionary: spanishDictionary(),
		keyboard: []keyboardRow{
			{keys: []string{"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P"}, padding: 3},
			{keys: []string{"A", "S", "D", "F", "G", "H", "J", "K", "L", "Ñ"}, padding: 3},
			{keys: []string{"ENVIAR", "Z", "X", "C", "V", "B", "N", "M", "BORRAR"}, padding: 0},
		},
		strings: uiStrings{
			score:          "Puntuación: %d",
			guessTooShort:  "Tu intento debe ser una palabra de 5 letras.",
			invalidWord:    "No está en la lista de palabras.",
			win:            "¡Has ganado!",
			loss:           "La palabra era %s. ¡Más suerte la próxima vez!",
			controlQuit:    "salir",
			controlRestart: "reiniciar",
		},
		fold: foldES,
	},
}

// lookupLanguage returns the language for the given code. If the code is
// unknown or empty, it returns the default language and false.
func lookupLanguage(code string) (language, bool) {
	code = strings.ToLower(strings.TrimSpace(code))
	if lang, ok := languages[code]; ok {
		return lang, true
	}
	return languages[_defaultLanguage], false
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

// spanishDictionary builds the Spanish dictionary from the embedded word
// lists. The allowed list is a superset of the answer list.
func spanishDictionary() Dictionary {
	answers := parseWordList(wordsESAnswers, foldES)
	allowed := parseWordList(wordsESAllowed, foldES)

	allWords := make(map[string]struct{}, len(allowed)+len(answers))
	for _, word := range allowed {
		allWords[word] = struct{}{}
	}
	for _, word := range answers {
		allWords[word] = struct{}{}
	}

	return Dictionary{commonWords: answers, allWords: allWords}
}

// parseWordList splits an embedded word list into normalized words, skipping
// blank lines and comments.
func parseWordList(data string, fold func(rune) rune) []string {
	lines := strings.Split(data, "\n")
	words := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, strings.Map(fold, line))
	}
	return words
}

// normalizeES folds a Spanish word to its canonical form: uppercase, with
// vowel diacritics stripped. Ñ is kept as a distinct letter.
func normalizeES(word string) string {
	return strings.Map(foldES, word)
}

// foldES folds a single Spanish rune to uppercase, stripping vowel diacritics.
func foldES(r rune) rune {
	switch r {
	case 'á', 'à', 'ä', 'â', 'Á', 'À', 'Ä', 'Â':
		return 'A'
	case 'é', 'è', 'ë', 'ê', 'É', 'È', 'Ë', 'Ê':
		return 'E'
	case 'í', 'ì', 'ï', 'î', 'Í', 'Ì', 'Ï', 'Î':
		return 'I'
	case 'ó', 'ò', 'ö', 'ô', 'Ó', 'Ò', 'Ö', 'Ô':
		return 'O'
	case 'ú', 'ù', 'ü', 'û', 'Ú', 'Ù', 'Ü', 'Û':
		return 'U'
	case 'ñ', 'Ñ':
		return 'Ñ'
	default:
		return toAsciiUpper(r)
	}
}
