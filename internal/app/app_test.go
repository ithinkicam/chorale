package app_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/ithinkicam/chorale/internal/app"
	"github.com/ithinkicam/chorale/internal/termtest"
)

const (
	titleWords  = "♪ chorale"
	emptyWords  = "No sessions yet."
	labelWords  = "▸ Prompt"
	sessionKeys = "Enter send · Ctrl-C quit"
)

func TestEmptySessionListAt60By20WithColour(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)

	lines := mt.Lines()
	if !strings.HasPrefix(lines[0], titleWords) {
		t.Errorf("first line is %q, want it to start with %q", lines[0], titleWords)
	}
	if label := labelRow(t, lines); label <= rowOf(t, lines, emptyWords) {
		t.Errorf("the label %q is on row %d, want it below %q", labelWords, label+1, emptyWords)
	}
	if last := lines[len(lines)-1]; !strings.Contains(last, sessionKeys) {
		t.Errorf("last line is %q, want it to show %q", last, sessionKeys)
	}
	mt.RequireGolden()
	assertViewFits(t, mt, 60, 20)
}

func TestEmptySessionListAt120By30WithColour(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 120, Height: 30})
	mt.WaitForText(emptyWords)

	lines := mt.Lines()
	if len(lines) != 30 {
		t.Fatalf("screen has %d rows, want 30", len(lines))
	}
	if !strings.HasPrefix(lines[0], titleWords) {
		t.Errorf("first line is %q, want it to start with %q", lines[0], titleWords)
	}
	labelRow(t, lines)
	if !strings.Contains(lines[29], sessionKeys) {
		t.Errorf("row 30 is %q, want it to show %q", lines[29], sessionKeys)
	}
	mt.RequireGolden()
	assertViewFits(t, mt, 120, 30)
}

func TestEmptySessionListSnapshotWithColourOff(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20, NoColor: true})
	mt.WaitForText(emptyWords)
	for _, words := range []string{titleWords, emptyWords, labelWords, sessionKeys} {
		if !strings.Contains(mt.Text(), words) {
			t.Errorf("the screen doesn't show %q:\n%s", words, mt.Text())
		}
	}
	mt.RequireGolden()
}

func TestCursorInTheEmptyPromptBox(t *testing.T) {
	s := startPromptSession(t)

	label := labelRow(t, s.frame())
	x, y, visible := s.mt.Screen().Cursor()
	if !visible {
		t.Error("the cursor is hidden")
	}
	if x != 2 || y != label+1 {
		t.Errorf("cursor at column %d, row %d; want column 3, row %d, below the label", x+1, y+1, label+2)
	}
}

func TestCursorInAnEmptyPromptBox(t *testing.T) {
	s := startPromptSession(t)

	lines := s.frame()
	screen := s.mt.Screen()
	x, y, visible := screen.Cursor()
	if !visible {
		t.Error("the cursor is hidden")
	}
	if style, set := screen.CursorStyle(); !set || screen.CursorBlinks() {
		t.Errorf("cursor style is %d, set %v; want a steady style that chorale set", style, set)
	}
	if want := labelRow(t, lines) + 1; x != 2 || y != want {
		t.Errorf("cursor at column %d, row %d; want column 3, row %d, the prompt box's first text row", x+1, y+1, want+1)
	}
}

func labelRow(t *testing.T, lines []string) int {
	t.Helper()
	return rowOf(t, lines, labelWords)
}

func rowOf(t *testing.T, lines []string, words string) int {
	t.Helper()
	for i, line := range lines {
		if strings.Contains(line, words) {
			return i
		}
	}
	t.Fatalf("the screen doesn't show %q:\n%s", words, strings.Join(lines, "\n"))
	return 0
}

func assertViewFits(t *testing.T, mt *termtest.ModelTerm, width, height int) {
	t.Helper()
	rows := strings.Split(mt.FinalView(), "\n")
	if len(rows) > height {
		t.Errorf("the view has %d rows, more than the window's %d", len(rows), height)
	}
	for i, row := range rows {
		if w := ansi.StringWidth(row); w > width {
			t.Errorf("view row %d is %d columns wide, wider than the window's %d: %q", i+1, w, width, row)
		}
	}
}
