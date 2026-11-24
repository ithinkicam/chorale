package editor

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestSetWidthKeepsTheInsertionPointOnTheSameCharacter(t *testing.T) {
	e := New(117, ansi.GraphemeWidth)
	e.Insert(strings.Repeat("word ", 17) + "word")
	e.caret, e.anchor = 69, 69
	assertPos(t, e, 0, 69)

	e.SetWidth(57, ansi.GraphemeWidth)
	assertCaret(t, e, 69)
	assertPos(t, e, 1, 14)
	if got := e.Text(); got != strings.Repeat("word ", 17)+"word" {
		t.Errorf("text changed to %q on resize", got)
	}
}

func TestCursorPosAndMoveMeasureByCodePointsUnderWcWidth(t *testing.T) {
	heart := "❤️"
	text := "a" + heart + "b"

	byCodePoints := New(textWidth, ansi.WcWidth)
	byCodePoints.Insert(text)
	assertPos(t, byCodePoints, 0, 3)
	byCodePoints.Move(MoveLeft, false)
	assertPos(t, byCodePoints, 0, 2)

	byGraphemes := New(textWidth, ansi.GraphemeWidth)
	byGraphemes.Insert(text)
	assertPos(t, byGraphemes, 0, 4)
	byGraphemes.Move(MoveLeft, false)
	assertPos(t, byGraphemes, 0, 3)
}

func TestMoveByRowUnderWcWidthWrapsByCodePoints(t *testing.T) {
	e := New(textWidth, ansi.WcWidth)
	e.Insert(strings.Repeat("❤️", 60))
	assertPos(t, e, 1, 3)
	e.Move(MoveUp, false)
	assertPos(t, e, 0, 3)
	assertCaret(t, e, 3)
}

func TestHistorySizeIsTheTextOfEverySnapshot(t *testing.T) {
	e := newEditor()
	check := func(step string) {
		t.Helper()
		sum := 0
		for _, s := range e.hist.snaps {
			sum += len(s.text)
		}
		if got := e.HistorySize(); got != sum {
			t.Errorf("after %s history size is %d, want %d", step, got, sum)
		}
	}

	e.Insert("run the tests")
	check("typing")
	e.DeleteWordBack()
	check("deleting a word")
	e.Undo()
	check("undo")
	e.Redo()
	check("redo")
	e.Undo()
	e.Insert("x")
	check("an edit after undo")
	e.Paste("pasted\nlines")
	check("paste")
	e.Clear()
	check("clear")
	for range 300 {
		e.Move(MoveStart, false)
		e.Insert("y")
	}
	check("300 edits")
}

func TestSetWidthClearsTheEndOfRowFlag(t *testing.T) {
	e := withText(wrapped, 3)
	e.Move(MoveRowEnd, false)
	e.SetWidth(textWidth-1, ansi.GraphemeWidth)
	if e.endOfRow {
		t.Error("end-of-row flag is still set after SetWidth")
	}
}

func TestSetWidthForgetsTheRememberedColumn(t *testing.T) {
	e := withText("abcdef\nab\nabcdef", 5)
	e.Move(MoveDown, false)
	e.SetWidth(textWidth-1, ansi.GraphemeWidth)
	e.Move(MoveDown, false)
	assertCaret(t, e, 10+2)
}
