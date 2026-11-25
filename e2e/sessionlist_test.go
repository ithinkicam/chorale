package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/ithinkicam/chorale/internal/termtest"
)

var sessionListWords = []string{"♪ chorale", "No sessions yet.", "▸ Prompt", "Enter send · Ctrl-C quit"}

func startSessionList(t *testing.T, opts termtest.Options) *termtest.Term {
	t.Helper()
	term := termtest.Start(t, chorale, opts)
	for _, words := range sessionListWords {
		term.WaitForText(words)
	}
	term.WaitQuiet(200 * time.Millisecond)
	return term
}

func TestEmptySessionListWithColourOff(t *testing.T) {
	term := startSessionList(t, termtest.Options{Width: 60, Height: 20, NoColor: true})

	if codes := termtest.ColourCodes(term.Raw()); len(codes) != 0 {
		t.Errorf("chorale sent %d colour codes with NO_COLOR=1, want 0; the first is %q", len(codes), codes[0].Sequence)
	}
}

func TestColourOnUsesTheTerminalPalette(t *testing.T) {
	term := startSessionList(t, termtest.Options{Width: 60, Height: 20})

	codes := termtest.ColourCodes(term.Raw())
	if len(codes) == 0 {
		t.Fatal("with colour on, chorale sent no colour codes")
	}
	for _, c := range codes {
		if !c.Palette {
			t.Errorf("colour code %q isn't one of the terminal's 16 palette colours", c.Sequence)
		}
	}
}

func TestCursorSitsInTheEmptyPromptBox(t *testing.T) {
	term := startSessionList(t, termtest.Options{Width: 60, Height: 20})

	labelX, labelY := find(t, term.Screen().Lines(), "▸ Prompt")
	x, y, visible := term.Screen().Cursor()
	if !visible {
		t.Error("the cursor is hidden")
	}
	if x != labelX+2 || y != labelY+1 {
		t.Errorf("cursor at column %d, row %d; want column %d, row %d, below the %q label",
			x+1, y+1, labelX+3, labelY+2, "▸ Prompt")
	}
}

// find assumes nothing wide comes before text, since columns count cells.
func find(t *testing.T, lines []string, text string) (x, y int) {
	t.Helper()
	for row, line := range lines {
		if col := strings.Index(line, text); col >= 0 {
			return len([]rune(line[:col])), row
		}
	}
	t.Fatalf("the screen doesn't show %q", text)
	return 0, 0
}
