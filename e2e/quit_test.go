package e2e

import (
	"testing"

	"github.com/ithinkicam/chorale/internal/termtest"
)

// Scenario: Quitting restores the terminal.
func TestQuittingRestoresTheTerminal(t *testing.T) {
	term := termtest.Start(t, chorale, termtest.Options{Width: 60, Height: 20})
	screen := term.Screen()

	// GIVEN Chorale is showing the empty session list on the alternate screen.
	term.WaitForText("No sessions yet.")
	if !screen.ModeOn(termtest.ModeAltScreen) {
		t.Fatal("the session list is showing, but not on the alternate screen")
	}

	term.Send("q")

	if code := term.Wait(); code != 0 {
		t.Errorf("exit status %d, want 0", code)
	}
	if left := screen.ModesLeftOn(); len(left) != 0 {
		t.Errorf("modes left on after quitting: %v (turned on during the run: %v)",
			left, screen.ModesTurnedOn())
	}
	if _, _, visible := screen.Cursor(); !visible {
		t.Error("the cursor is hidden after quitting")
	}
}
