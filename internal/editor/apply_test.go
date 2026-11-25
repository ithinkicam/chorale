package editor

import "testing"

func TestApplyDoesEachAction(t *testing.T) {
	typed := func(s string) *Editor {
		e := newEditor()
		e.Insert(s)
		return e
	}
	tests := []struct {
		name     string
		editor   *Editor
		action   Action
		wantText string
		wantFrom int
		wantTo   int
	}{
		{"move", withText("ab", 2), Action{Kind: ActionMove, Move: MoveLeft}, "ab", 1, 1},
		{"extend", withText("ab", 2), Action{Kind: ActionMove, Move: MoveLeft, Extend: true}, "ab", 1, 2},
		{"select all", withText("ab", 1), Action{Kind: ActionSelectAll}, "ab", 0, 2},
		{"backspace", withText("ab", 2), Action{Kind: ActionBackspace}, "a", 1, 1},
		{"delete", withText("ab", 0), Action{Kind: ActionDelete}, "b", 0, 0},
		{"delete word back", withText("run the", 7), Action{Kind: ActionDeleteWordBack}, "run ", 4, 4},
		{"delete word forward", withText("run the", 0), Action{Kind: ActionDeleteWordForward}, " the", 0, 0},
		{"delete to row start", withText("ab cd", 5), Action{Kind: ActionDeleteToRowStart}, "", 0, 0},
		{"delete to row end", withText("ab cd", 2), Action{Kind: ActionDeleteToRowEnd}, "ab", 2, 2},
		{"undo", typed("ab"), Action{Kind: ActionUndo}, "", 0, 0},
		{"new line", typed("a"), Action{Kind: ActionNewLine}, "a\n", 2, 2},
		{"copy", withText("ab", 1), Action{Kind: ActionCopy}, "ab", 1, 1},
		{"cut", withText("ab", 1), Action{Kind: ActionCut}, "ab", 1, 1},
		{"interrupt", withText("ab", 1), Action{Kind: ActionInterrupt}, "ab", 1, 1},
		{"send", withText("ab", 1), Action{Kind: ActionSend}, "ab", 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.Apply(tt.action)
			assertText(t, tt.editor, tt.wantText)
			assertSelection(t, tt.editor, tt.wantFrom, tt.wantTo)
		})
	}
}

func TestApplyRedoRepeatsWhatUndoReverted(t *testing.T) {
	e := newEditor()
	e.Insert("ab")
	e.Apply(Action{Kind: ActionUndo})
	assertText(t, e, "")
	e.Apply(Action{Kind: ActionRedo})
	assertText(t, e, "ab")
}
