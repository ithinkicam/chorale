package editor

import "testing"

func TestCleanPasteTurnsCarriageReturnsIntoNewLines(t *testing.T) {
	if got, want := CleanPaste("a\r\nb\rc\nd"), "a\nb\nc\nd"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteRemovesAColourCodeWhole(t *testing.T) {
	if got, want := CleanPaste("safe\x1b[31mtext"), "safetext"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteRemovesAnOperatingSystemCommandWhole(t *testing.T) {
	tests := map[string]string{
		"ended by bell":            "a\x1b]52;c;aGk=\x07b",
		"ended by string terminal": "a\x1b]0;a title\x1b\\b",
	}
	for name, in := range tests {
		if got, want := CleanPaste(in), "ab"; got != want {
			t.Errorf("%s: cleaned paste is %q, want %q", name, got, want)
		}
	}
}

func TestCleanPasteRemovesADeviceControlStringWhole(t *testing.T) {
	if got, want := CleanPaste("a\x1bP1$rabc\x1b\\b"), "ab"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteRemovesALoneEscapeByte(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{"at the end", "a\x1b", "a"},
		{"at the end of a line", "a\x1b\n", "a\n"},
	}
	for _, tt := range tests {
		if got := CleanPaste(tt.in); got != tt.want {
			t.Errorf("%s: cleaned paste is %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestCleanPasteTurnsATabIntoFourSpaces(t *testing.T) {
	if got, want := CleanPaste("a\tb"), "a    b"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteRemovesC0C1AndDeleteCharacters(t *testing.T) {
	if got, want := CleanPaste("a\x00b\x07c\x08d\u0085e\u009bf\x7fg"), "abcdefg"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteKeepsTextEmojiAndNewLines(t *testing.T) {
	in := "合唱 🎵 é\n\nend"
	if got := CleanPaste(in); got != in {
		t.Errorf("cleaned paste is %q, want it unchanged", got)
	}
}

func TestPasteInsertsSeveralLinesAtTheInsertionPoint(t *testing.T) {
	e := newEditor()
	e.Paste("first\r\nsecond\nthird")
	assertText(t, e, "first\nsecond\nthird")
	assertCaret(t, e, 18)
}

func TestPasteReplacesTheSelection(t *testing.T) {
	e := newEditor()
	e.Insert("run the tests")
	e.Move(MoveWordLeft, true)
	e.Paste("build")
	assertText(t, e, "run the build")
	if e.HasSelection() {
		t.Error("selection remains after pasting over it")
	}
}

func TestPasteIsOneUndoGroup(t *testing.T) {
	e := newEditor()
	e.Paste("first\r\nsecond\nthird")
	e.Undo()
	assertText(t, e, "")
}

func TestUndoingAPasteOverASelectionRestoresIt(t *testing.T) {
	e := newEditor()
	e.Insert("run the tests")
	e.Move(MoveWordLeft, true)
	e.Paste("build")
	e.Undo()
	assertText(t, e, "run the tests")
	assertSelection(t, e, 8, 13)
}

func TestPasteDoesNotJoinTheTypingBeforeIt(t *testing.T) {
	e := newEditor()
	e.Insert("ab")
	e.Paste("cd")
	e.Insert("ef")
	e.Undo()
	assertText(t, e, "abcd")
	e.Undo()
	assertText(t, e, "ab")
}

func TestPastingOnlyEscapeSequencesChangesNothing(t *testing.T) {
	e := newEditor()
	e.Insert("a")
	before := e.HistoryLen()
	e.Paste("\x1b[31m")
	assertText(t, e, "a")
	if got := e.HistoryLen(); got != before {
		t.Errorf("history has %d snapshots, want %d", got, before)
	}
}
