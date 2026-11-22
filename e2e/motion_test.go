package e2e

import (
	"strings"
	"testing"

	"github.com/ithinkicam/chorale/internal/termtest"
)

// Scenario: Idle screen stays still.
func TestIdleScreenStaysStill(t *testing.T) {
	term := startSessionList(t, termtest.Options{Width: 60, Height: 20})

	if n := term.IdleBytes(); n != 0 {
		t.Fatalf("with no input, chorale wrote %d bytes in %v", n, termtest.IdleWindow)
	}
}

// Scenario: Startup draws the session list directly.
func TestStartupDrawsTheSessionListDirectly(t *testing.T) {
	term := startSessionList(t, termtest.Options{Width: 60, Height: 20})

	// The session list never contains "needs", so finding the word
	// anywhere in the output means the too-small message was drawn. A
	// single word is checked because a renderer may move the cursor
	// between words instead of writing spaces.
	if strings.Contains(string(term.Raw()), "needs") {
		t.Fatalf("chorale drew the too-small message before the session list: %q", term.Raw())
	}
}
