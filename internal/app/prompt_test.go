package app_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

var (
	left       = named(tea.KeyLeft, 0)
	right      = named(tea.KeyRight, 0)
	up         = named(tea.KeyUp, 0)
	down       = named(tea.KeyDown, 0)
	end        = named(tea.KeyEnd, 0)
	backspace  = named(tea.KeyBackspace, 0)
	optionLeft = named(tea.KeyLeft, tea.ModAlt)
	optionUp   = named(tea.KeyUp, tea.ModAlt)
	undo       = named('z', tea.ModCtrl)
	redo       = named('y', tea.ModCtrl)
)

func TestTypingAWord(t *testing.T) {
	s := startPromptSession(t)

	s.typeText("hello")

	waitForPrompt(t, s.mt, "hello")
	s.wantBoxRow("hello")
	s.wantCursorColumn(8)
}

func TestTypingInTheMiddle(t *testing.T) {
	s := startPromptSession(t)
	s.given("helo", 3)

	s.typeText("l")

	s.wantText("hello")
	s.wantCaret(4)
}

func TestWideCharactersAndEmoji(t *testing.T) {
	s := startPromptSession(t)

	s.typeText("合唱 🎵 é")

	s.wantText("合唱 🎵 é")
	s.wantBoxRow("合唱 🎵 é")
	s.wantCursorColumn(12)
}

func TestAnAccentTypedOnItsOwnJoinsItsLetter(t *testing.T) {
	s := startPromptSession(t)
	s.given("e", 1)

	s.typeText("́")

	s.wantText("é")
	s.wantBoxRow("e\u0301")
	s.press(left)
	s.wantCaret(0)
}

func TestTabDoesNothing(t *testing.T) {
	s := startPromptSession(t)
	s.given("run", 3)

	s.pressUnbound(named(tea.KeyTab, 0))

	s.wantText("run")
	s.wantCaret(3)
}

func TestCursorAfterASoftWrap(t *testing.T) {
	s := startPromptSession(t)

	s.typeText(words(12) + " end")

	s.wantRows(words(11), "word end")
	if row := s.state().CaretRow; row != 1 {
		t.Errorf("the insertion point is on text row %d, want 2", row+1)
	}
	s.wantBoxRow("word end")
	s.wantCursorColumn(11)
}

func TestCursorAtTheEndOfAFullRow(t *testing.T) {
	s := startPromptSession(t)
	s.given(strings.Repeat("a", 57)+"bbb", 0)

	s.press(end)

	if row := s.state().CaretRow; row != 0 {
		t.Errorf("the insertion point is on text row %d, want 1", row+1)
	}
	s.wantBoxRow(strings.Repeat("a", 57))
	s.wantCursorColumn(60)
	s.press(right)
	if row := s.state().CaretRow; row != 1 {
		t.Errorf("after Right, the insertion point is on text row %d, want 2", row+1)
	}
	s.wantBoxRow("bbb")
	s.wantCursorColumn(4)
}

func TestMovingOverAnEmoji(t *testing.T) {
	s := startPromptSession(t)
	s.given("a🎵b", 3)

	s.press(left, left)

	s.wantCaret(1)
}

func TestMovingLeftAcrossASoftWrap(t *testing.T) {
	s := startPromptSession(t)
	s.given(words(12)+" end", 55)

	s.press(left)

	s.wantCaret(54)
	if row := s.state().CaretRow; row != 0 {
		t.Errorf("the insertion point is on text row %d, want 1", row+1)
	}
	s.wantBoxRow(words(11))
	s.wantCursorColumn(57)
}

func TestMovingByWord(t *testing.T) {
	s := startPromptSession(t)
	s.given("run the tests now", 17)

	s.press(optionLeft, optionLeft)

	s.wantCaret(8)
	s.press(named(tea.KeyRight, tea.ModAlt))
	s.wantCaret(13)
}

func TestDownKeepsTheColumn(t *testing.T) {
	s := startPromptSession(t)
	s.given("abcdef\nab\nabcdef", 5)

	s.press(down, down)

	s.wantCaret(15)
}

func TestUpMovesAcrossASoftWrap(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd(words(12) + " end")

	s.press(up)

	if row := s.state().CaretRow; row != 0 {
		t.Errorf("the insertion point is on text row %d, want 1", row+1)
	}
	s.wantBoxRow(words(11))
	s.wantCursorColumn(11)
	s.wantCaret(8)
}

func TestUpOnTheFirstRow(t *testing.T) {
	s := startPromptSession(t)
	s.given("one\ntwo", 2)

	s.press(up)

	s.wantCaret(0)
	s.press(down, down)
	s.wantCaret(7)
}

func TestStartAndEndOfThePrompt(t *testing.T) {
	s := startPromptSession(t)
	s.given("one\ntwo\nthree", 6)

	s.press(optionUp)

	s.wantCaret(0)
	s.press(named(tea.KeyDown, tea.ModAlt))
	s.wantCaret(13)
}

func TestSelectingAWord(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("run the tests")

	s.press(named(tea.KeyLeft, tea.ModAlt|tea.ModShift))

	s.wantSelected("tests")
	s.wantCursorColumn(11)
}

func TestTypingReplacesTheSelection(t *testing.T) {
	s := startPromptSession(t)
	s.givenSelected("run the tests", 5)
	s.wantSelected("tests")

	s.typeText("build")

	s.wantText("run the build")
	s.wantSelected("")
}

func TestRightMovesToTheEndOfASelection(t *testing.T) {
	s := startPromptSession(t)
	s.given("run the tests now", 13)
	s.press(named(tea.KeyLeft, tea.ModShift), named(tea.KeyLeft, tea.ModShift), named(tea.KeyLeft, tea.ModShift),
		named(tea.KeyLeft, tea.ModShift), named(tea.KeyLeft, tea.ModShift))
	s.wantSelected("tests")
	s.wantCaret(8)

	s.press(right)

	s.wantSelected("")
	s.wantCaret(13)
}

func TestSelectingEverything(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("one\ntwo")

	s.press(named('a', tea.ModCtrl|tea.ModShift))

	s.wantSelected("one\ntwo")
}

func TestBackspaceOverAnEmoji(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("go 🎵")

	s.press(backspace)

	s.wantText("go ")
}

func TestDeletingAWord(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("run the tests")

	s.press(named(tea.KeyBackspace, tea.ModAlt))

	s.wantText("run the ")
}

func TestDeletingAWordAfterAnEmoji(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("go 🎵")

	s.press(named(tea.KeyBackspace, tea.ModAlt))

	s.wantText("go ")
}

func TestDeletingToTheEndOfTheRow(t *testing.T) {
	s := startPromptSession(t)
	s.given("keep this\nand this", 4)

	s.press(named('k', tea.ModCtrl))

	s.wantText("keep\nand this")
}

func TestBackspaceAtTheStartOfAWrappedRow(t *testing.T) {
	s := startPromptSession(t)
	s.given(words(12)+" end", 55)

	s.press(backspace)

	s.wantText(words(10) + " wordword end")
}

func TestUndoTypingAWordAtATime(t *testing.T) {
	s := startPromptSession(t)
	s.typeText("run the tests")

	s.press(undo)

	s.wantText("run the ")
	s.press(undo, undo)
	s.wantText("")
}

func TestUndoingARunOfBackspace(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("run the tests")

	s.press(backspace, backspace, backspace, backspace, backspace)
	s.wantText("run the ")
	s.press(undo)

	s.wantText("run the tests")
}

func TestRedoAfterUndo(t *testing.T) {
	s := startPromptSession(t)
	s.typeText("run the tests")
	s.press(named(tea.KeyBackspace, tea.ModAlt))
	s.wantText("run the ")

	s.press(undo, redo)

	s.wantText("run the ")
	s.wantCaret(8)
}

func TestUndoRestoresTheSelection(t *testing.T) {
	s := startPromptSession(t)
	s.givenSelected("run the tests", 5)

	s.typeText("build")
	s.wantText("run the build")
	s.press(undo)

	s.wantText("run the tests")
	s.wantSelected("tests")
}

func TestANewEditDiscardsRedo(t *testing.T) {
	s := startPromptSession(t)
	s.typeText("run the tests")
	s.press(undo)

	s.typeText("x")
	s.press(redo)

	s.wantText("run the x")
}

func TestWrappingAtASpaceAt60Columns(t *testing.T) {
	s := startPromptSession(t)

	s.typeText(words(15))

	s.wantRowsNoWiderThan(57)
	for i, row := range s.state().Rows {
		for _, word := range strings.Fields(row) {
			if word != "word" {
				t.Errorf("text row %d holds %q, a split word", i+1, word)
			}
		}
	}
	if text := s.state().Text; strings.Contains(text, "\n") {
		t.Errorf("the prompt is %q, want no new lines in it", text)
	}
	s.wantText(words(15))
	s.wantBoxRow(words(4))
}

func TestALongWordIsBroken(t *testing.T) {
	s := startPromptSession(t)

	s.typeText(strings.Repeat("x", 100))

	s.wantRows(strings.Repeat("x", 57), strings.Repeat("x", 43))
	s.wantBoxRow(strings.Repeat("x", 43))
}

func TestWideCharactersWrapWhole(t *testing.T) {
	s := startPromptSession(t)

	s.typeText(strings.Repeat("合", 40))

	s.wantRows(strings.Repeat("合", 28), strings.Repeat("合", 12))
	s.wantRowsNoWiderThan(57)
	s.wantBoxRow(strings.Repeat("合", 12))
}

const clearedWords = "Prompt cleared. Ctrl-Z brings it back."

func TestCtrlCClearsThePrompt(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("half a thought")

	s.press(ctrlC())

	s.wantText("")
	s.mt.WaitForText(clearedWords)
	requireRunning(t, s.mt)
	s.press(undo)
	s.wantText("half a thought")
}

func TestCtrlCClearsAPromptOfOnlySpaces(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("  \n ")

	s.press(ctrlC())

	s.wantText("")
	s.mt.WaitForText(clearedWords)
}

func TestKeyLineWithText(t *testing.T) {
	s := startPromptSession(t)

	s.typeText("run")

	s.wantLastLine("Enter send · Ctrl-C clear")
}

func TestKeyLineWithASelection(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("run the tests")

	s.press(named(tea.KeyLeft, tea.ModAlt|tea.ModShift))

	s.wantLastLine("Enter send · Ctrl-C copy")
}

func TestNewLineKeys(t *testing.T) {
	s := startPromptSession(t)

	s.typeText("a")
	s.press(named(tea.KeyEnter, tea.ModShift))
	s.typeText("b")
	s.press(named(tea.KeyEnter, tea.ModAlt))
	s.typeText("c")
	s.press(named('j', tea.ModCtrl))
	s.typeText("d")

	s.wantText("a\nb\nc\nd")
}

const noAgentWords = "No agent is connected yet."

func TestSendingWithNoAgentConnected(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("run the tests")

	s.press(named(tea.KeyEnter, 0))

	s.mt.WaitForText(noAgentWords + " Your prompt is kept.")
	s.wantText("run the tests")
}

func TestEnterOnAnEmptyPrompt(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("  \n ")

	s.press(named(tea.KeyEnter, 0))

	s.wantViewWithout(noAgentWords)
	s.wantScreenWithout(noAgentWords)
	s.wantText("  \n ")
}

func TestAStatusMessageGoesAtTheNextKey(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("run the tests")
	s.press(named(tea.KeyEnter, 0))
	s.mt.WaitForText(noAgentWords)

	s.press(left)

	s.wantViewWithout(noAgentWords)
}

func TestAPasteClearsAStatusMessage(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("run the tests")
	s.press(named(tea.KeyEnter, 0))
	s.mt.WaitForText(noAgentWords)

	s.paste("x")

	s.wantScreenWithout(noAgentWords)
	s.wantText("run the testsx")
}

func TestCmdZUndoes(t *testing.T) {
	s := startPromptSession(t)
	s.typeText("run the tests")

	s.press(named('z', tea.ModSuper))

	s.wantText("run the ")
}

func TestCmdASelectsEverything(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("one\ntwo")

	s.press(named('a', tea.ModSuper))

	s.wantSelected("one\ntwo")
}

func TestShiftCmdLeftSelectsToTheStartOfTheRow(t *testing.T) {
	s := startPromptSession(t)
	s.givenAtEnd("run the tests")

	s.press(named(tea.KeyLeft, tea.ModShift|tea.ModSuper))

	s.wantSelected("run the tests")
	s.wantCursorColumn(3)
}
