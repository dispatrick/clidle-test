package language

// messagesEnglish contains the English user-facing strings. These are the
// canonical strings, and are used as a fallback for unsupported languages.
var messagesEnglish = Messages{
	Score:          "Score: %d",
	GuessTooShort:  "Your guess must be a 5-letter word.",
	GuessNotAWord:  "That's not a valid word.",
	Win:            "You win!",
	Loss:           "The word was %s. Better luck next time!",
	KeyEnter:       "ENTER",
	KeyDelete:      "DELETE",
	ControlQuit:    "quit",
	ControlRestart: "restart",
}
