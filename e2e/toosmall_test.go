package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/ithinkicam/chorale/internal/termtest"
)

const tooSmallWords = "Chorale needs at least 60 by 20."

func TestTooSmallMessageWithColourOff(t *testing.T) {
	term := termtest.Start(t, chorale, termtest.Options{Width: 40, Height: 20, NoColor: true})
	term.WaitForText(tooSmallWords)
	term.WaitForText("q quit")
	term.WaitQuiet(200 * time.Millisecond)

	if codes := termtest.ColourCodes(term.Raw()); len(codes) != 0 {
		t.Errorf("chorale sent %d colour codes with NO_COLOR=1, want 0; the first is %q", len(codes), codes[0].Sequence)
	}
}

func TestQuittingFromTheTooSmallScreen(t *testing.T) {
	term := termtest.Start(t, chorale, termtest.Options{Width: 40, Height: 20})
	screen := term.Screen()
	term.WaitForText(tooSmallWords)

	term.Send(ctrlCTwice)

	if code := term.Wait(); code != 0 {
		t.Errorf("exit status %d, want 0", code)
	}
	if left := screen.ModesLeftOn(); len(left) != 0 {
		t.Errorf("modes left on after quitting: %v, want none", left)
	}
	if _, _, visible := screen.Cursor(); !visible {
		t.Error("the cursor is hidden after quitting")
	}
}

func TestQDoesNothingOnTheTooSmallScreen(t *testing.T) {
	term := termtest.Start(t, chorale, termtest.Options{Width: 40, Height: 20})
	screen := term.Screen()
	term.WaitForText(tooSmallWords)

	term.Send("q")
	term.WaitQuiet(200 * time.Millisecond)

	if !screen.ModeOn(termtest.ModeAltScreen) {
		t.Errorf("chorale left the alternate screen after q, so it has ended. The screen shows:\n%s", screen.Text())
	}
	if !strings.Contains(screen.Text(), tooSmallWords) {
		t.Errorf("the screen no longer shows %q:\n%s", tooSmallWords, screen.Text())
	}
	term.Resize(60, 20)
	term.WaitForText("No sessions yet.")
}
