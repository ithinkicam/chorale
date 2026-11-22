package e2e

import (
	"testing"

	"github.com/ithinkicam/chorale/internal/termtest"
)

func TestQuittingRestoresTheTerminal(t *testing.T) {
	term := termtest.Start(t, chorale, termtest.Options{Width: 60, Height: 20})
	screen := term.Screen()

	term.WaitForText("No sessions yet.")
	if !screen.ModeOn(termtest.ModeAltScreen) {
		t.Fatal("the session list is showing, but not on the alternate screen")
	}

	term.Send("q")
	defer func() { t.Logf("modes turned on during the run: %v", screen.ModesTurnedOn()) }()

	if code := term.Wait(); code != 0 {
		t.Errorf("exit status %d, want 0", code)
	}
	if left := screen.ModesLeftOn(); len(left) != 0 {
		t.Errorf("modes left on after quitting: %v, want none (turned on during the run: %v)",
			left, screen.ModesTurnedOn())
	}
	if _, _, visible := screen.Cursor(); !visible {
		t.Error("the cursor is hidden after quitting")
	}
}
