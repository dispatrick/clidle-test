package language

import (
	"reflect"
	"testing"
)

func TestLookup(t *testing.T) {
	tests := []struct {
		code     string
		wantCode string
		wantOK   bool
	}{
		{"en", "en", true},
		{"EN", "en", true},
		{"ca", "ca", true},
		{"CA", "ca", true},
		{"ca-ES", "ca", true},
		{"ca_ES.UTF-8", "ca", true},
		{"catalan", "ca", true},
		{" ca ", "ca", true},
		{"", "en", false},
		{"xx", "en", false},
	}

	for _, test := range tests {
		lang, ok := Lookup(test.code)
		if lang.Code != test.wantCode || ok != test.wantOK {
			t.Errorf("Lookup(%q) = (%q, %v), want (%q, %v)", test.code, lang.Code, ok, test.wantCode, test.wantOK)
		}
	}
}

// TestMessagesComplete ensures that every language translates every message.
func TestMessagesComplete(t *testing.T) {
	for code, lang := range Languages {
		messages := reflect.ValueOf(lang.Messages)
		for i := 0; i < messages.NumField(); i++ {
			if messages.Field(i).String() == "" {
				t.Errorf("language %q is missing message %q", code, messages.Type().Field(i).Name)
			}
		}
	}
}
