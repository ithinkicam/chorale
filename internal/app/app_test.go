package app_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/ithinkicam/chorale/internal/app"
	"github.com/ithinkicam/chorale/internal/termtest"
)

// The words on the empty session list, from the app-shell spec.
const (
	titleWords = "♪ chorale"
	emptyWords = "No sessions yet."
	keyWords   = "q quit"
)

// Scenario: Empty session list at 60 by 20 with colour.
func TestEmptySessionListAt60By20WithColour(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)

	lines := mt.Lines()
	if !strings.HasPrefix(lines[0], titleWords) {
		t.Errorf("first line is %q, want it to start with %q", lines[0], titleWords)
	}
	if last := lines[len(lines)-1]; !strings.Contains(last, keyWords) {
		t.Errorf("last line is %q, want it to show %q", last, keyWords)
	}
	assertViewFits(t, 60, 20)
}

// Scenario: Empty session list at 120 by 30 with colour.
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
	if !strings.Contains(lines[29], keyWords) {
		t.Errorf("row 30 is %q, want it to show %q", lines[29], keyWords)
	}
	assertViewFits(t, 120, 30)
}

// assertViewFits checks the model's own view, before any renderer could
// clip it, has no line wider than the window and no more rows than it has.
func assertViewFits(t *testing.T, width, height int) {
	t.Helper()
	var m tea.Model = app.New()
	m, _ = m.Update(tea.WindowSizeMsg{Width: width, Height: height})
	rows := strings.Split(m.View().Content, "\n")
	if len(rows) > height {
		t.Errorf("the view has %d rows, more than the window's %d", len(rows), height)
	}
	for i, row := range rows {
		if w := ansi.StringWidth(row); w > width {
			t.Errorf("view row %d is %d columns wide, wider than the window's %d: %q", i+1, w, width, row)
		}
	}
}
