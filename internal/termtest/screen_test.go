package termtest

import (
	"slices"
	"testing"
)

// Each case feeds the escape sequence that turns a mode on, then the one
// that turns it off, and checks the screen reports both.
var modeCases = []struct {
	mode Mode
	on   string
	off  string
}{
	{ModeAltScreen, "\x1b[?1049h", "\x1b[?1049l"},
	{ModeAltScreen, "\x1b[?1047h", "\x1b[?1047l"},
	{ModeCursorHidden, "\x1b[?25l", "\x1b[?25h"},
	{ModeMouse, "\x1b[?1000h", "\x1b[?1000l"},
	{ModeMouse, "\x1b[?1002h", "\x1b[?1002l"},
	{ModeMouse, "\x1b[?1003h", "\x1b[?1003l"},
	{ModeMouseSGR, "\x1b[?1006h", "\x1b[?1006l"},
	{ModeBracketedPaste, "\x1b[?2004h", "\x1b[?2004l"},
	{ModeFocusEvents, "\x1b[?1004h", "\x1b[?1004l"},
	{ModeKeyboardEnhancements, "\x1b[>1u", "\x1b[<u"},
	{ModeKeyboardEnhancements, "\x1b[=3;1u", "\x1b[=0;1u"},
	{ModeKeyboardEnhancements, "\x1b[>4;1m", "\x1b[>4;0m"},
	{ModeKeyboardEnhancements, "\x1b[>4;2m", "\x1b[>4m"},
}

func TestScreenTracksModes(t *testing.T) {
	for _, c := range modeCases {
		t.Run(string(c.mode)+" "+quote(c.on), func(t *testing.T) {
			s := NewScreen(60, 20)
			if s.ModeOn(c.mode) {
				t.Fatalf("%s is on before anything was written", c.mode)
			}
			write(t, s, c.on)
			if !s.ModeOn(c.mode) {
				t.Fatalf("%s is off after %s", c.mode, quote(c.on))
			}
			write(t, s, c.off)
			if s.ModeOn(c.mode) {
				t.Fatalf("%s is still on after %s", c.mode, quote(c.off))
			}
			if !slices.Contains(s.ModesTurnedOn(), c.mode) {
				t.Fatalf("ModesTurnedOn() = %v, want it to include %s", s.ModesTurnedOn(), c.mode)
			}
			if left := s.ModesLeftOn(); len(left) != 0 {
				t.Fatalf("ModesLeftOn() = %v, want none", left)
			}
		})
	}
}

func TestScreenReportsModesLeftOn(t *testing.T) {
	s := NewScreen(60, 20)
	write(t, s, "\x1b[?1049h\x1b[?25l\x1b[?2004h\x1b[>1u")
	write(t, s, "\x1b[?1049l")

	want := []Mode{ModeBracketedPaste, ModeCursorHidden, ModeKeyboardEnhancements}
	if got := s.ModesLeftOn(); !slices.Equal(got, want) {
		t.Fatalf("ModesLeftOn() = %v, want %v", got, want)
	}
}

func TestKittyKeyboardStackPopsOneLevelAtATime(t *testing.T) {
	s := NewScreen(60, 20)
	write(t, s, "\x1b[>1u\x1b[>3u\x1b[<u")
	if !s.ModeOn(ModeKeyboardEnhancements) {
		t.Fatal("one push is left on the stack, but keyboard enhancements are off")
	}
	write(t, s, "\x1b[<u")
	if s.ModeOn(ModeKeyboardEnhancements) {
		t.Fatal("the stack is empty, but keyboard enhancements are on")
	}
}

func TestScreenLinesAreTrimmedPlainText(t *testing.T) {
	s := NewScreen(20, 3)
	write(t, s, "\x1b[1;31mred\x1b[0m word\r\n\r\n  last")

	want := []string{"red word", "", "  last"}
	if got := s.Lines(); !slices.Equal(got, want) {
		t.Fatalf("Lines() = %q, want %q", got, want)
	}
}

func TestScreenCursor(t *testing.T) {
	s := NewScreen(20, 5)
	write(t, s, "\x1b[3;7H")
	x, y, visible := s.Cursor()
	if x != 6 || y != 2 || !visible {
		t.Fatalf("Cursor() = %d, %d, %v; want 6, 2, true", x, y, visible)
	}
	write(t, s, "\x1b[?25l")
	if _, _, visible := s.Cursor(); visible {
		t.Fatal("cursor reported visible after it was hidden")
	}
}

func TestScreenAnswersQueriesWithoutBlocking(t *testing.T) {
	// A program that asks the terminal a question, such as the primary
	// device attributes, must not stall the screen.
	s := NewScreen(20, 5)
	write(t, s, "\x1b[c\x1b[c\x1b[cdone")
	if got := s.Lines()[0]; got != "done" {
		t.Fatalf("first line = %q, want %q", got, "done")
	}
}

func write(t *testing.T, s *Screen, text string) {
	t.Helper()
	if _, err := s.Write([]byte(text)); err != nil {
		t.Fatalf("Write(%s): %v", quote(text), err)
	}
}

func TestOnlcrStartsEachLineAtTheLeftEdge(t *testing.T) {
	s := NewScreen(20, 3)
	write(t, s, string(onlcr([]byte("first\n  second"))))
	if got := s.Lines()[1]; got != "  second" {
		t.Fatalf("second line = %q, want %q", got, "  second")
	}
}
