package e2e

import (
	"strings"
	"testing"

	"github.com/ithinkicam/chorale/internal/termtest"
)

// Scenario: Version flag.
func TestVersionFlag(t *testing.T) {
	term := termtest.Start(t, chorale, termtest.Options{Width: 60, Height: 20, Args: []string{"--version"}})

	if code := term.Wait(); code != 0 {
		t.Errorf("exit status %d, want 0", code)
	}
	assertOneLine(t, term, "chorale "+testVersion)
	if fromModes := term.Screen().ModesTurnedOn(); len(fromModes) != 0 {
		t.Errorf("--version turned on terminal modes: %v", fromModes)
	}
}

// Scenario: Development build version.
func TestDevelopmentBuildVersion(t *testing.T) {
	term := termtest.Start(t, choraleDev, termtest.Options{Width: 60, Height: 20, Args: []string{"--version"}})

	if code := term.Wait(); code != 0 {
		t.Errorf("exit status %d, want 0", code)
	}
	assertOneLine(t, term, "chorale dev")
}

// assertOneLine checks the program printed exactly want and a line ending.
func assertOneLine(t *testing.T, term *termtest.Term, want string) {
	t.Helper()
	// The terminal turns "\n" into "\r\n".
	got := strings.ReplaceAll(string(term.Raw()), "\r\n", "\n")
	if got != want+"\n" {
		t.Errorf("printed %q, want exactly %q", got, want+"\n")
	}
}
