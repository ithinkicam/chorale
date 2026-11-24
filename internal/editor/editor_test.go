package editor

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
)

const textWidth = 57

func newEditor() *Editor {
	return New(textWidth, ansi.GraphemeWidth)
}

func withText(text string, caret int) *Editor {
	e := newEditor()
	e.gs = Graphemes(text)
	e.caret = caret
	e.anchor = caret
	return e
}

func assertText(t *testing.T, e *Editor, want string) {
	t.Helper()
	if got := e.Text(); got != want {
		t.Errorf("text is %q, want %q", got, want)
	}
}

func assertCaret(t *testing.T, e *Editor, want int) {
	t.Helper()
	if got := e.Caret(); got != want {
		t.Errorf("insertion point is %d, want %d", got, want)
	}
}

func TestInsertPutsLettersAtTheEnd(t *testing.T) {
	e := newEditor()
	e.Insert("hello")
	assertText(t, e, "hello")
	assertCaret(t, e, 5)
}

func TestInsertGoesInTheMiddle(t *testing.T) {
	e := withText("helo", 3)
	e.Insert("l")
	assertText(t, e, "hello")
	assertCaret(t, e, 4)
}

func TestInsertAcceptsWideCharactersAndEmoji(t *testing.T) {
	e := newEditor()
	e.Insert("合唱 🎵 é")
	assertText(t, e, "合唱 🎵 é")
	assertCaret(t, e, 7)
}

func TestInsertJoinsACombiningAccent(t *testing.T) {
	e := newEditor()
	e.Insert("e")
	e.Insert("́")
	assertText(t, e, "é")
	assertCaret(t, e, 1)
	if got := e.Len(); got != 1 {
		t.Errorf("text has %d characters, want 1", got)
	}
}

func TestInsertJoinsAnEmojiBuiltFromSeveralKeystrokes(t *testing.T) {
	e := newEditor()
	for _, part := range []string{"👨", "‍", "👩"} {
		e.Insert(part)
	}
	if got := e.Len(); got != 1 {
		t.Errorf("text has %d characters, want 1", got)
	}
	assertCaret(t, e, 1)
}

func TestInsertAfterAnAccentedLetterKeepsTheInsertionPointAfterIt(t *testing.T) {
	e := withText("éb", 1)
	e.Insert("x")
	assertText(t, e, "éxb")
	assertCaret(t, e, 2)
}

func TestGraphemesKeepsAnEmojiWhole(t *testing.T) {
	got := Graphemes("a🎵é")
	if len(got) != 3 {
		t.Errorf("got %d graphemes %q, want 3", len(got), got)
	}
}
