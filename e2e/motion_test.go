package e2e

import (
	"strings"
	"testing"

	"github.com/ithinkicam/chorale/internal/termtest"
)

func TestIdleScreenStaysStill(t *testing.T) {
	term := startSessionList(t, termtest.Options{Width: 60, Height: 20})

	if n := term.IdleBytes(); n != 0 {
		t.Fatalf("with no input, chorale wrote %d bytes in %v", n, termtest.IdleWindow)
	}
}

func TestStartupDrawsTheSessionListDirectly(t *testing.T) {
	term := startSessionList(t, termtest.Options{Width: 60, Height: 20})

	// Check one word, because a renderer may move the cursor between words instead of writing spaces.
	if strings.Contains(string(term.Raw()), "needs") {
		t.Fatalf("chorale drew the too-small message before the session list: %q", term.Raw())
	}
}
