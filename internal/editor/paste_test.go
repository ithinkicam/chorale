package editor

import (
	"strings"
	"testing"
	"time"
)

func TestCleanPasteTurnsCarriageReturnsIntoNewLines(t *testing.T) {
	if got, want := cleanPaste("a\r\nb\rc\nd"), "a\nb\nc\nd"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteRemovesAColourCodeWhole(t *testing.T) {
	if got, want := cleanPaste("safe\x1b[31mtext"), "safetext"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteRemovesAnOperatingSystemCommandWhole(t *testing.T) {
	tests := map[string]string{
		"ended by bell":            "a\x1b]52;c;aGk=\x07b",
		"ended by string terminal": "a\x1b]0;a title\x1b\\b",
	}
	for name, in := range tests {
		if got, want := cleanPaste(in), "ab"; got != want {
			t.Errorf("%s: cleaned paste is %q, want %q", name, got, want)
		}
	}
}

func TestCleanPasteRemovesADeviceControlStringWhole(t *testing.T) {
	if got, want := cleanPaste("a\x1bP1$rabc\x1b\\b"), "ab"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteRemovesALoneEscapeByte(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{"at the end", "a\x1b", "a"},
		{"at the end of a line", "a\x1b\n", "a\n"},
		{"before a new line and more text", "a\x1b\nb", "a\nb"},
	}
	for _, tt := range tests {
		if got := cleanPaste(tt.in); got != tt.want {
			t.Errorf("%s: cleaned paste is %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestCleanPasteTurnsATabIntoFourSpaces(t *testing.T) {
	if got, want := cleanPaste("a\tb"), "a    b"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteRemovesC0C1AndDeleteCharacters(t *testing.T) {
	if got, want := cleanPaste("a\x00b\x07c\x08d\u0085e\x7ff"), "abcdef"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteKeepsTextEmojiAndNewLines(t *testing.T) {
	in := "合唱 🎵 é\n\nend"
	if got := cleanPaste(in); got != in {
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
	before := e.historyLen()
	e.Paste("\x1b[31m")
	assertText(t, e, "a")
	if got := e.historyLen(); got != before {
		t.Errorf("history has %d snapshots, want %d", got, before)
	}
}

func TestCleanPasteKeepsTextAfterAControlStringThatNeverEnds(t *testing.T) {
	for _, intro := range []string{"]", "P", "X", "^", "_"} {
		in := "a\x1b" + intro + "\nsecond\nthird"
		if got, want := cleanPaste(in), "a\nsecond\nthird"; got != want {
			t.Errorf("escape %s: cleaned paste is %q, want %q", intro, got, want)
		}
	}
}

func TestCleanPasteKeepsTextAfterAControlStringThatEndsWithTheText(t *testing.T) {
	if got, want := cleanPaste("a\x1b]0;title"), "a0;title"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteKeepsTextAfterAControlSequenceWithNoFinalByte(t *testing.T) {
	tests := []struct{ in, want string }{
		{"a\x1b[\nsecond\nthird", "a\nsecond\nthird"},
		{"a\x1b[31\nsecond", "a31\nsecond"},
		{"a\x1b[31", "a31"},
	}
	for _, tt := range tests {
		if got := cleanPaste(tt.in); got != tt.want {
			t.Errorf("cleaned paste of %q is %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCleanPasteTreatsC1IntroducersLikeTheirEscapeForms(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{"control sequence", "\u009b31mtext", "text"},
		{"operating system command", "a\u009d0;title\x07b", "ab"},
		{"device control string", "a\u0090qabc\u009cb", "ab"},
		{"string ended by an escape and a backslash", "a\u009d0;title\x1b\\b", "ab"},
		{"unfinished", "a\u009d0;x\nsecond", "a0;x\nsecond"},
	}
	for _, tt := range tests {
		if got := cleanPaste(tt.in); got != tt.want {
			t.Errorf("%s: cleaned paste is %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestCleanPasteTakesLinearTimeForC1Strings(t *testing.T) {
	timeOf := func(in string) time.Duration {
		start := time.Now()
		cleanPaste(in)
		return time.Since(start)
	}
	const repeats = 20000
	escForm := timeOf(strings.Repeat("\x1b]ab", repeats))
	c1Form := timeOf(strings.Repeat("\u009dab", repeats))

	limit := 20 * max(escForm, 2*time.Millisecond)
	if c1Form > limit {
		t.Errorf("cleaning %d C1 strings took %v, want at most %v, 20 times the %v the escape form took", repeats, c1Form, limit, escForm)
	}
}

func TestCleanPasteKeepsTheWordAfterAnEscapeAndASpace(t *testing.T) {
	if got, want := cleanPaste("a\x1b word"), "a word"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}
