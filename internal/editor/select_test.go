package editor

import "testing"

func assertSelection(t *testing.T, e *Editor, wantStart, wantEnd int) {
	t.Helper()
	start, end := e.Selection()
	if start != wantStart || end != wantEnd {
		t.Errorf("selection is %d to %d, want %d to %d", start, end, wantStart, wantEnd)
	}
}

func TestShiftMovementExtendsTheSelectionFromWhereTheInsertionPointWas(t *testing.T) {
	e := withText("run the tests", 13)
	e.Move(MoveWordLeft, true)
	assertSelection(t, e, 8, 13)
	assertCaret(t, e, 8)
	if got := e.SelectedText(); got != "tests" {
		t.Errorf("selected text is %q, want %q", got, "tests")
	}
}

func TestExtendingBackTowardTheAnchorShrinksTheSelection(t *testing.T) {
	e := withText("run the tests", 8)
	e.Move(MoveRight, true)
	e.Move(MoveRight, true)
	e.Move(MoveLeft, true)
	assertSelection(t, e, 8, 9)
	e.Move(MoveLeft, true)
	if e.HasSelection() {
		t.Error("selection remains after extending back to the anchor")
	}
}

func TestExtendingLeftOfTheAnchorSelectsBackwards(t *testing.T) {
	e := withText("abcdef", 3)
	e.Move(MoveLeft, true)
	e.Move(MoveLeft, true)
	assertSelection(t, e, 1, 3)
	assertCaret(t, e, 1)
}

func TestMovingWithoutShiftClearsTheSelection(t *testing.T) {
	e := withText("run the tests", 13)
	e.Move(MoveWordLeft, true)
	e.Move(MoveLeft, false)
	if e.HasSelection() {
		t.Error("selection remains after moving without Shift")
	}
	assertCaret(t, e, 7)
}

func TestTypingReplacesTheSelection(t *testing.T) {
	e := withText("run the tests", 13)
	e.Move(MoveWordLeft, true)
	e.Insert("build")
	assertText(t, e, "run the build")
	assertCaret(t, e, 13)
	if e.HasSelection() {
		t.Error("selection remains after typing over it")
	}
}

func TestSelectAllSelectsTheWholePrompt(t *testing.T) {
	e := withText("one\ntwo", 2)
	e.SelectAll()
	assertSelection(t, e, 0, 7)
	if got := e.SelectedText(); got != "one\ntwo" {
		t.Errorf("selected text is %q, want %q", got, "one\ntwo")
	}
}

func TestSelectedCountCountsANewLineAndAnEmojiAsOneCharacterEach(t *testing.T) {
	e := withText("a🎵\nb", 0)
	e.SelectAll()
	if got := e.SelectedCount(); got != 4 {
		t.Errorf("selected count is %d, want 4", got)
	}
}

func TestSelectedCountIsZeroWithNothingSelected(t *testing.T) {
	e := withText("abc", 2)
	if got := e.SelectedCount(); got != 0 {
		t.Errorf("selected count is %d, want 0", got)
	}
}
