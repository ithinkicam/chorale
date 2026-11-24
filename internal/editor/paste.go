package editor

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/x/ansi"
)

var pasteCleaner = strings.NewReplacer("\r\n", "\n", "\r", "\n")

// CleanPaste returns pasted text as it is inserted: with new lines for
// carriage returns, escape sequences removed whole, tabs as four spaces, and
// every other control character removed.
func CleanPaste(s string) string {
	s = ansi.Strip(pasteCleaner.Replace(s))
	s = strings.ReplaceAll(s, "\t", "    ")
	return strings.Map(func(r rune) rune {
		if r != '\n' && unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}

// Paste inserts cleaned s at the insertion point as one edit, replacing the
// selection if there is one.
func (e *Editor) Paste(s string) {
	s = CleanPaste(s)
	if s == "" {
		return
	}
	start, end := e.Selection()
	e.apply(start, end, s, false)
}
