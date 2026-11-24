package editor

import (
	"strings"
	"testing"
)

var wrapped = strings.Repeat("word ", 12) + "end"

const (
	secondRowStart = 55
	wrappedLength  = 63
)

func assertPos(t *testing.T, e *Editor, wantRow, wantCol int) {
	t.Helper()
	row, col := e.CursorPos()
	if row != wantRow || col != wantCol {
		t.Errorf("cursor is on row %d, column %d, want row %d, column %d", row, col, wantRow, wantCol)
	}
}

func TestMoveLeftStepsOverAnEmojiWhole(t *testing.T) {
	e := withText("a🎵b", 3)
	e.Move(MoveLeft, false)
	e.Move(MoveLeft, false)
	assertCaret(t, e, 1)
}

func TestMoveStopsAtTheEnds(t *testing.T) {
	e := withText("ab", 0)
	e.Move(MoveLeft, false)
	assertCaret(t, e, 0)
	e = withText("ab", 2)
	e.Move(MoveRight, false)
	assertCaret(t, e, 2)
}

func TestMoveLeftAcrossASoftWrapTakesOnePress(t *testing.T) {
	e := withText(wrapped, secondRowStart)
	assertPos(t, e, 1, 0)
	e.Move(MoveLeft, false)
	assertCaret(t, e, secondRowStart-1)
	assertPos(t, e, 0, 54)
}

func TestMoveRightAcrossASoftWrapTakesOnePress(t *testing.T) {
	e := withText(wrapped, secondRowStart-1)
	e.Move(MoveRight, false)
	assertCaret(t, e, secondRowStart)
	assertPos(t, e, 1, 0)
}

func TestMoveByWordSkipsSpacesBeforeTheWord(t *testing.T) {
	e := withText("run the tests now", 17)
	e.Move(MoveWordLeft, false)
	e.Move(MoveWordLeft, false)
	assertCaret(t, e, 8)
	e.Move(MoveWordRight, false)
	assertCaret(t, e, 13)
}

func TestMoveByWordSkipsPunctuation(t *testing.T) {
	e := withText("one,  two...", 0)
	e.Move(MoveWordRight, false)
	assertCaret(t, e, 3)
	e.Move(MoveWordRight, false)
	assertCaret(t, e, 9)
	e.Move(MoveWordLeft, false)
	assertCaret(t, e, 6)
}

func TestWordsHoldLettersDigitsAndUnderscores(t *testing.T) {
	e := withText("foo_bar9 x", 0)
	e.Move(MoveWordRight, false)
	assertCaret(t, e, 8)
}

func TestAnEmojiIsAWordOfItsOwn(t *testing.T) {
	e := withText("a🎵b", 3)
	e.Move(MoveWordLeft, false)
	assertCaret(t, e, 2)
	e.Move(MoveWordLeft, false)
	assertCaret(t, e, 1)
	e.Move(MoveWordLeft, false)
	assertCaret(t, e, 0)
}

func TestMoveByWordCrossesNewLines(t *testing.T) {
	e := withText("one\n\ntwo", 5)
	e.Move(MoveWordLeft, false)
	assertCaret(t, e, 0)
}

func TestDownKeepsTheColumnAcrossHardWraps(t *testing.T) {
	e := withText("abcdef\nab\nabcdef", 5)
	e.Move(MoveDown, false)
	assertCaret(t, e, 9)
	e.Move(MoveDown, false)
	assertCaret(t, e, 15)
	e.Move(MoveUp, false)
	e.Move(MoveUp, false)
	assertCaret(t, e, 5)
}

func TestUpMovesAcrossASoftWrap(t *testing.T) {
	e := withText(wrapped, wrappedLength)
	assertPos(t, e, 1, 8)
	e.Move(MoveUp, false)
	assertPos(t, e, 0, 8)
	assertCaret(t, e, 8)
	e.Move(MoveDown, false)
	assertCaret(t, e, wrappedLength)
}

func TestDownIntoASoftRowStopsBeforeItsWrapPoint(t *testing.T) {
	e := withText(strings.Repeat("x", 57)+"\n"+wrapped, 57)
	e.Move(MoveDown, false)
	assertPos(t, e, 1, 54)
}

func TestUpOnTheFirstRowAndDownOnTheLastDoNothing(t *testing.T) {
	e := withText("one\ntwo", 2)
	e.Move(MoveUp, false)
	assertCaret(t, e, 2)
	e = withText("one\ntwo", 5)
	e.Move(MoveDown, false)
	assertCaret(t, e, 5)
}

func TestAnotherKeyForgetsTheRememberedColumn(t *testing.T) {
	e := withText("abcdef\nab\nabcdef", 5)
	e.Move(MoveDown, false)
	e.Move(MoveLeft, false)
	e.Move(MoveDown, false)
	assertCaret(t, e, 10+1)
}

func TestMoveToTheStartOfARowOnASoftWrap(t *testing.T) {
	e := withText(wrapped, 60)
	e.Move(MoveRowStart, false)
	assertCaret(t, e, secondRowStart)
	e.Move(MoveLeft, false)
	e.Move(MoveRowStart, false)
	assertCaret(t, e, 0)
}

func TestMoveToTheEndOfARowSetsTheEndOfRowFlag(t *testing.T) {
	e := withText(wrapped, 3)
	e.Move(MoveRowEnd, false)
	assertCaret(t, e, secondRowStart)
	if !e.EndOfRow() {
		t.Error("end-of-row flag is not set after moving to the end of a row")
	}
	assertPos(t, e, 0, 55)
}

func TestWithoutTheFlagTheInsertionPointAtAWrapIsOnTheNextRow(t *testing.T) {
	e := withText(wrapped, secondRowStart)
	assertPos(t, e, 1, 0)
}

func TestMoveToTheEndOfAFullRowThenRight(t *testing.T) {
	e := withText(strings.Repeat("a", 57)+"bbb", 0)
	e.Move(MoveRowEnd, false)
	assertPos(t, e, 0, 57)
	e.Move(MoveRight, false)
	assertPos(t, e, 1, 1)
}

func TestAnyOtherKeyClearsTheEndOfRowFlag(t *testing.T) {
	for name, act := range map[string]func(*Editor){
		"left":   func(e *Editor) { e.Move(MoveLeft, false) },
		"right":  func(e *Editor) { e.Move(MoveRight, false) },
		"start":  func(e *Editor) { e.Move(MoveRowStart, false) },
		"insert": func(e *Editor) { e.Insert("x") },
	} {
		e := withText(wrapped, 3)
		e.Move(MoveRowEnd, false)
		act(e)
		if e.EndOfRow() {
			t.Errorf("%s left the end-of-row flag set", name)
		}
	}
}

func TestMoveToTheEndOfARowOnAHardLine(t *testing.T) {
	e := withText("abc\ndef", 1)
	e.Move(MoveRowEnd, false)
	assertCaret(t, e, 3)
	e.Move(MoveRowStart, false)
	assertCaret(t, e, 0)
}

func TestMoveToTheStartAndEndOfThePrompt(t *testing.T) {
	e := withText("one\ntwo\nthree", 6)
	e.Move(MoveStart, false)
	assertCaret(t, e, 0)
	e.Move(MoveEnd, false)
	assertCaret(t, e, 13)
}

var manyWords = strings.Repeat("word ", 30)

func TestHomeAfterEndGoesToTheStartOfTheSameRow(t *testing.T) {
	e := withText(wrapped, 3)
	e.Move(MoveRowEnd, false)
	e.Move(MoveRowStart, false)
	assertPos(t, e, 0, 0)
	assertCaret(t, e, 0)
}

func TestEndTwiceStaysAtTheEndOfTheSameRow(t *testing.T) {
	e := withText(wrapped, 3)
	e.Move(MoveRowEnd, false)
	e.Move(MoveRowEnd, false)
	assertPos(t, e, 0, 55)
	assertCaret(t, e, secondRowStart)
}

func TestDownAfterEndMovesToTheRowBelowTheDrawnRow(t *testing.T) {
	e := withText(manyWords, 3)
	e.Move(MoveRowEnd, false)
	e.Move(MoveDown, false)
	assertPos(t, e, 1, 54)
}

func TestShiftHomeAfterShiftEndSelectsFromTheDrawnRow(t *testing.T) {
	e := withText(wrapped, 3)
	e.Move(MoveRowEnd, true)
	e.Move(MoveRowStart, true)
	assertSelection(t, e, 0, 3)
	assertCaret(t, e, 0)
}

func TestShiftDownAfterShiftEndExtendsFromTheDrawnRow(t *testing.T) {
	e := withText(manyWords, 3)
	e.Move(MoveRowEnd, true)
	e.Move(MoveDown, true)
	assertSelection(t, e, 3, secondRowStart+54)
}
