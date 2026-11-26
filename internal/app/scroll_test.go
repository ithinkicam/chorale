package app_test

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/ithinkicam/chorale/internal/termtest"
)

var newLine = named('j', tea.ModCtrl)

func numberedLines(n int) string {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i+1)
	}
	return strings.Join(lines, "\n")
}

func (s *session) typeNumberedLines(n int) {
	s.t.Helper()
	for i := 1; i <= n; i++ {
		if i > 1 {
			s.press(newLine)
		}
		s.typeText(fmt.Sprintf("line %d", i))
	}
}

func TestALongPromptScrollsAt60By20(t *testing.T) {
	s := startPromptSession(t)

	s.typeNumberedLines(20)

	rows := s.boxRows()
	if len(rows) != 8 {
		t.Errorf("the prompt box shows %d text rows, %q, want 8", len(rows), rows)
	} else if rows[7] != "line 20" {
		t.Errorf("the last text row is %q, want %q", rows[7], "line 20")
	}
	s.wantScreenWith("12 lines above")
	s.wantCursorAt(8, 10)
	s.wantViewFits(60, 20)
	s.mt.RequireGolden()
}

func TestScrollingFollowsTheInsertionPoint(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd(numberedLines(20))

	s.press(optionUp)

	s.wantBoxRows("line 1")
	s.wantScreenWith("12 lines below")
	s.wantScreenWithout("lines above")
}

func TestResizeKeepsTheInsertionPoint(t *testing.T) {
	s := startSession(t, termtest.Options{Width: 120, Height: 30})
	s.given(words(18), 69)

	s.mt.Resize(60, 20)
	s.settle()

	s.wantText(words(18))
	s.wantCaret(69)
	s.wantBoxRows(words(11), words(7))
	s.wantCursorAt(2, 17)
}

func TestTheBoxGrowsWithItsText(t *testing.T) {
	tests := []struct {
		name   string
		height int
		lines  int
		want   int
	}{
		{"empty", 20, 0, 1},
		{"three lines", 20, 3, 3},
		{"eight lines at 20 rows", 20, 8, 8},
		{"nine lines at 20 rows", 20, 9, 8},
		{"twelve lines at 24 rows", 24, 12, 12},
		{"thirteen lines at 24 rows", 24, 13, 12},
		{"twenty lines at 40 rows", 40, 20, 12},
		{"twenty lines at 22 rows", 22, 20, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := startSession(t, termtest.Options{Width: 60, Height: tt.height})
			if tt.lines > 0 {
				s.givenAtEnd(numberedLines(tt.lines))
			}

			rows := s.boxRows()

			if len(rows) != tt.want {
				t.Errorf("the prompt box shows %d text rows, %q, want %d", len(rows), rows, tt.want)
			}
		})
	}
}

func TestScrollingMovesByTheLeastAmount(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd(numberedLines(20))
	s.wantBoxRows("line 13")
	s.wantScreenWith("12 lines above")

	s.press(up, up, up)
	s.wantBoxRows("line 13")
	s.wantScreenWith("12 lines above")

	s.press(up, up, up, up, up)
	s.wantBoxRows("line 12")
	s.wantScreenWith("11 lines above · 1 line below")
	s.wantCursorAt(1, 10)
}

func TestTheCountsSayHowManyRowsAreOutOfSightOnBothSides(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd(numberedLines(20))

	s.press(optionUp)
	s.press(down, down, down, down, down, down, down, down, down)

	s.wantBoxRows("line 3")
	s.wantScreenWith("▸ Prompt · 2 lines above · 10 lines below")
}

func TestAWrappedPromptSnapshotAt60By20(t *testing.T) {
	s := startPromptSession(t)

	s.typeText(words(15))

	s.wantBoxRows(words(11), words(4))
	s.wantCursorAt(2, 22)
	s.mt.RequireGolden()
}

func TestASelectionSnapshotAt60By20(t *testing.T) {
	s := startPromptSession(t)
	s.givenSelected(words(15), 24)

	s.wantBoxRows(words(11), words(4))
	s.wantSelected("word word word word word")
	s.mt.RequireGolden()
}

func TestDeletingLinesKeepsTheBoxFull(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd(numberedLines(20))

	for range len("\nline 20") {
		s.press(backspace)
	}

	rows := s.boxRows()
	if len(rows) != 8 {
		t.Errorf("the prompt box shows %d text rows, %q, want 8", len(rows), rows)
	} else if rows[7] != "line 19" {
		t.Errorf("the last text row is %q, want %q", rows[7], "line 19")
	}
	s.wantScreenWith("11 lines above")
	s.wantCursorAt(8, 10)
}

func TestGrowingTheWindowShowsMoreOfAScrolledPrompt(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd(numberedLines(20))

	s.mt.Resize(60, 24)
	s.settle()

	rows := s.boxRows()
	if len(rows) != 12 {
		t.Errorf("the prompt box shows %d text rows, %q, want 12", len(rows), rows)
	} else if rows[11] != "line 20" {
		t.Errorf("the last text row is %q, want %q", rows[11], "line 20")
	}
	s.wantScreenWith("8 lines above")
	s.wantCursorAt(12, 10)
}

func TestShrinkingTheWindowKeepsTheInsertionPointInView(t *testing.T) {
	s := startSession(t, termtest.Options{Width: 60, Height: 24})
	s.givenAtEnd(numberedLines(20))

	s.mt.Resize(60, 20)
	s.settle()

	rows := s.boxRows()
	if len(rows) != 8 {
		t.Errorf("the prompt box shows %d text rows, %q, want 8", len(rows), rows)
	} else if rows[7] != "line 20" {
		t.Errorf("the last text row is %q, want %q", rows[7], "line 20")
	}
	s.wantScreenWith("12 lines above")
	s.wantCursorAt(8, 10)
}

func TestAWrappedLineOutOfSightCountsByItsRows(t *testing.T) {
	s := startPromptSession(t)

	s.paste(strings.Repeat("x", 150) + "\n" + strings.TrimPrefix(numberedLines(9), "line 1\n"))

	rows := s.boxRows()
	if len(rows) != 8 {
		t.Errorf("the prompt box shows %d text rows, %q, want 8", len(rows), rows)
	} else if rows[0] != "line 2" || rows[7] != "line 9" {
		t.Errorf("the text rows run from %q to %q, want %q to %q", rows[0], rows[7], "line 2", "line 9")
	}
	s.wantScreenWith("3 lines above")
	s.wantScreenWithout("1 line above")
}
