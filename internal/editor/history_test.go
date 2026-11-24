package editor

import (
	"fmt"
	"strings"
	"testing"
)

func TestUndoRemovesTypingAWordAtATime(t *testing.T) {
	e := newEditor()
	e.Insert("run the tests")
	e.Undo()
	assertText(t, e, "run the ")
	e.Undo()
	assertText(t, e, "run ")
	e.Undo()
	assertText(t, e, "")
}

func TestTypingOneKeyAtATimeIsUndoneLikeTypingItAllAtOnce(t *testing.T) {
	e := newEditor()
	for _, g := range Graphemes("run the tests") {
		e.Insert(g)
	}
	e.Undo()
	assertText(t, e, "run the ")
}

func TestANewLineEndsATypingGroup(t *testing.T) {
	e := newEditor()
	e.Insert("one\ntwo")
	e.Undo()
	assertText(t, e, "one\n")
}

func TestUndoAtTheStartOfTheHistoryDoesNothing(t *testing.T) {
	e := newEditor()
	if e.Undo() {
		t.Error("undo reported a change with nothing to undo")
	}
	e.Insert("a")
	e.Undo()
	if e.Undo() {
		t.Error("undo reported a change after reaching the start")
	}
}

func TestRedoRestoresTheUndoneEditAndTheInsertionPoint(t *testing.T) {
	e := newEditor()
	e.Insert("run the tests")
	e.DeleteWordBack()
	e.Undo()
	assertText(t, e, "run the tests")
	e.Redo()
	assertText(t, e, "run the ")
	assertCaret(t, e, 8)
}

func TestRedoWithNothingUndoneDoesNothing(t *testing.T) {
	e := newEditor()
	e.Insert("a")
	if e.Redo() {
		t.Error("redo reported a change with nothing to redo")
	}
}

func TestUndoRestoresTheSelection(t *testing.T) {
	e := newEditor()
	e.Insert("run the tests")
	e.Move(MoveWordLeft, true)
	e.Insert("build")
	assertText(t, e, "run the build")
	e.Undo()
	assertText(t, e, "run the tests")
	assertSelection(t, e, 8, 13)
	assertCaret(t, e, 8)
}

func TestUndoRestoresTheInsertionPointFromBeforeTheEdit(t *testing.T) {
	e := newEditor()
	e.Insert("ab")
	e.Move(MoveStart, false)
	e.Delete()
	assertText(t, e, "b")
	e.Undo()
	assertText(t, e, "ab")
	assertCaret(t, e, 0)
}

func TestMovingStartsANewTypingGroup(t *testing.T) {
	e := newEditor()
	e.Insert("ab")
	e.Move(MoveLeft, false)
	e.Insert("c")
	assertText(t, e, "acb")
	e.Undo()
	assertText(t, e, "ab")
}

func TestEachDeletionIsOneUndo(t *testing.T) {
	e := newEditor()
	e.Insert("one two")
	e.DeleteWordBack()
	e.DeleteWordBack()
	assertText(t, e, "")
	e.Undo()
	assertText(t, e, "one ")
	e.Undo()
	assertText(t, e, "one two")
}

func TestClearIsOneUndo(t *testing.T) {
	e := newEditor()
	e.Insert("half a thought")
	e.Clear()
	assertText(t, e, "")
	e.Undo()
	assertText(t, e, "half a thought")
}

func TestAnEditThatChangesNothingIsNotRecorded(t *testing.T) {
	e := newEditor()
	e.Insert("a")
	before := e.HistoryLen()
	e.Move(MoveStart, false)
	e.Backspace()
	if got := e.HistoryLen(); got != before {
		t.Errorf("history has %d snapshots after an edit that changed nothing, want %d", got, before)
	}
}

func TestANewEditAfterUndoDiscardsRedo(t *testing.T) {
	e := newEditor()
	e.Insert("run the tests")
	e.Undo()
	e.Insert("x")
	assertText(t, e, "run the x")
	if e.Redo() {
		t.Error("redo reported a change after a new edit")
	}
	assertText(t, e, "run the x")
}

func TestHistoryKeepsAtMostTwoHundredSnapshots(t *testing.T) {
	e := newEditor()
	for range 300 {
		e.Move(MoveStart, false)
		e.Insert("x")
	}
	if got := e.HistoryLen(); got != 200 {
		t.Fatalf("history has %d snapshots, want 200", got)
	}
	undone := 0
	for e.Undo() {
		undone++
	}
	if undone != 199 {
		t.Errorf("undid %d edits, want 199", undone)
	}
}

func TestHistoryTextStaysUnderSixteenMegabytes(t *testing.T) {
	var prompt strings.Builder
	for i := range 10000 {
		fmt.Fprintf(&prompt, "line number %d of a long prompt\n", i)
	}
	e := newEditor()
	e.replace(0, 0, prompt.String())
	e.hist = newHistory(e.Text())
	for range 300 {
		e.Move(MoveStart, false)
		e.Insert("x")
	}

	if got := e.HistorySize(); got > 16<<20 {
		t.Errorf("history stores %d bytes of text, want at most %d", got, 16<<20)
	}
	if got := e.HistoryLen(); got < 10 || got >= 200 {
		t.Errorf("history has %d snapshots, want between 10 and 199", got)
	}
	undone := 0
	for e.Undo() {
		undone++
	}
	if undone != e.HistoryLen()-1 {
		t.Errorf("undid %d edits, want %d", undone, e.HistoryLen()-1)
	}
}
