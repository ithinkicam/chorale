package editor

import (
	"strings"
	"testing"
)

func TestBackspaceRemovesAnEmojiWhole(t *testing.T) {
	e := withText("go 🎵", 4)
	e.Backspace()
	assertText(t, e, "go ")
	assertCaret(t, e, 3)
}

func TestBackspaceRemovesALetterWithItsAccent(t *testing.T) {
	e := withText("aé", 2)
	e.Backspace()
	assertText(t, e, "a")
}

func TestBackspaceAtTheStartDoesNothing(t *testing.T) {
	e := withText("ab", 0)
	e.Backspace()
	assertText(t, e, "ab")
}

func TestDeleteRemovesTheCharacterAfterTheInsertionPoint(t *testing.T) {
	e := withText("a🎵b", 1)
	e.Delete()
	assertText(t, e, "ab")
	assertCaret(t, e, 1)
}

func TestDeleteAtTheEndDoesNothing(t *testing.T) {
	e := withText("ab", 2)
	e.Delete()
	assertText(t, e, "ab")
}

func TestBackspaceAndDeleteRemoveTheSelection(t *testing.T) {
	e := withText("run the tests", 13)
	e.Move(MoveWordLeft, true)
	e.Backspace()
	assertText(t, e, "run the ")
	assertCaret(t, e, 8)

	e = withText("run the tests", 8)
	e.Move(MoveWordRight, true)
	e.Delete()
	assertText(t, e, "run the ")
}

func TestBackspaceAtTheStartOfASoftWrappedRowRemovesTheSpaceBefore(t *testing.T) {
	e := withText(wrapped, secondRowStart)
	e.Backspace()
	assertText(t, e, strings.Repeat("word ", 10)+"wordword end")
}

func TestDeleteWordBackRemovesTheWordBeforeTheInsertionPoint(t *testing.T) {
	e := withText("run the tests", 13)
	e.DeleteWordBack()
	assertText(t, e, "run the ")
	assertCaret(t, e, 8)
}

func TestDeleteWordBackAfterAnEmojiRemovesOnlyTheEmoji(t *testing.T) {
	e := withText("go 🎵", 4)
	e.DeleteWordBack()
	assertText(t, e, "go ")
}

func TestDeleteWordBackSkipsSpacesBeforeTheWord(t *testing.T) {
	e := withText("run the   ", 10)
	e.DeleteWordBack()
	assertText(t, e, "run ")
}

func TestDeleteWordForwardRemovesTheWordAfterTheInsertionPoint(t *testing.T) {
	e := withText("run the tests", 4)
	e.DeleteWordForward()
	assertText(t, e, "run  tests")
	assertCaret(t, e, 4)
}

func TestDeleteWordForwardBeforeAnEmojiRemovesOnlyTheEmoji(t *testing.T) {
	e := withText("go 🎵now", 3)
	e.DeleteWordForward()
	assertText(t, e, "go now")
}

func TestDeleteWordsReplaceTheSelection(t *testing.T) {
	e := withText("one two three", 13)
	e.Move(MoveLeft, true)
	e.DeleteWordBack()
	assertText(t, e, "one two thre")
}

func TestDeleteToRowEndKeepsTheNewLine(t *testing.T) {
	e := withText("keep this\nand this", 4)
	e.DeleteToRowEnd()
	assertText(t, e, "keep\nand this")
	assertCaret(t, e, 4)
}

func TestDeleteToRowEndAtTheEndOfALineDoesNothing(t *testing.T) {
	e := withText("keep\nthis", 4)
	e.DeleteToRowEnd()
	assertText(t, e, "keep\nthis")
}

func TestDeleteToRowStartRemovesBackToTheStartOfTheLine(t *testing.T) {
	e := withText("abc\ndef", 6)
	e.DeleteToRowStart()
	assertText(t, e, "abc\nf")
	assertCaret(t, e, 4)
}

func TestDeleteToRowStartOnASoftWrappedRowStopsAtTheRow(t *testing.T) {
	e := withText(wrapped, 60)
	e.DeleteToRowStart()
	assertText(t, e, strings.Repeat("word ", 11)+"end")
}

func TestDeleteToRowEndOnASoftWrappedRowStopsAtTheRow(t *testing.T) {
	e := withText(wrapped, 3)
	e.DeleteToRowEnd()
	assertText(t, e, "wor"+"word end")
}

func TestClearEmptiesThePrompt(t *testing.T) {
	e := withText("one\ntwo", 4)
	e.Move(MoveRight, true)
	e.Clear()
	assertText(t, e, "")
	assertCaret(t, e, 0)
	if e.HasSelection() {
		t.Error("selection remains after clearing")
	}
}
