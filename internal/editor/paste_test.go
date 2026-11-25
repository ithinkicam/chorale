package editor

import (
	"strings"
	"testing"
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

func TestScanningC1StringsExaminesAConstantNumberOfBytesPerInputByte(t *testing.T) {
	for _, unit := range []string{"\u009dab", "\x1b]ab", "\u0090ab", "\x1b[1"} {
		in := strings.Repeat(unit, 20000)
		_, examined := scan(in)
		if limit := 4 * len(in); examined > limit {
			t.Errorf("scanning %q repeated examined %d bytes of %d, want at most %d", unit, examined, len(in), limit)
		}
	}
}

func TestCleanPasteKeepsTheWordAfterAnEscapeAndASpace(t *testing.T) {
	if got, want := cleanPaste("a\x1b word"), "a word"; got != want {
		t.Errorf("cleaned paste is %q, want %q", got, want)
	}
}

func TestCleanPasteHandlesOtherShapesOfEscapeSequence(t *testing.T) {
	tests := []struct{ name, in, want string }{
		{"control sequence with a private marker", "a\x1b[?2004hb", "ab"},
		{"escape form string ended by a C1 terminator", "a\x1bPqabc\u009cb", "ab"},
		{"operating system command ended by a C1 terminator", "a\x1b]0;t\u009cb", "ab"},
		{"C1 start of string", "a\u0098xyz\u009cb", "ab"},
		{"C1 privacy message", "a\u009exyz\u009cb", "ab"},
		{"C1 application program command", "a\u009fxyz\u009cb", "ab"},
		{"raw 8-bit control sequence introducer", "a\x9b31mb", "a�31mb"},
		{"raw 8-bit operating system command", "a\x9d0;t\x07b", "a�0;tb"},
		{"null inside a control sequence", "a\x1b[3\x001mb", "a31mb"},
		{"null inside a control string", "a\x1b]0;t\x00x\x07b", "a0;txb"},
		{"escape and intermediates with no final byte", "a\x1b(\nb", "a(\nb"},
		{"escape and intermediates at the end", "a\x1b((", "a(("},
	}
	for _, tt := range tests {
		if got := cleanPaste(tt.in); got != tt.want {
			t.Errorf("%s: cleaned paste is %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestInsertCleansTextLikePaste(t *testing.T) {
	for _, text := range []string{"a\u009b31mb", "a\u001b[31mb", "\u0007", "a\tb", "a\r\nb"} {
		t.Run(text, func(t *testing.T) {
			pasted := newEditor()
			pasted.Paste(text)
			typed := newEditor()

			typed.Insert(text)

			assertText(t, typed, pasted.Text())
		})
	}
}

func TestTextThatCleansToNothingStillCountsAsAKey(t *testing.T) {
	tests := []struct {
		name string
		do   func(e *Editor)
	}{
		{"insert", func(e *Editor) { e.Insert("\u0007") }},
		{"paste", func(e *Editor) { e.Paste("\x1b[31m") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := withText("abcdef\nab\nabcdef", 5)
			e.Move(MoveDown, false)

			tt.do(e)
			e.Move(MoveDown, false)

			assertCaret(t, e, 9+3)
		})
	}
}

func TestInsertThatCleansToNothingEndsTheUndoGroup(t *testing.T) {
	e := newEditor()
	e.Insert("ab")
	e.Insert("\u0007")
	e.Insert("c")

	e.Undo()
	assertText(t, e, "ab")
	e.Undo()
	assertText(t, e, "")
	if e.Undo() {
		t.Errorf("a third undo went back to %q, want only two undo groups", e.Text())
	}
}
