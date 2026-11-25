package editor

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var pasteCleaner = strings.NewReplacer("\r\n", "\n", "\r", "\n")

// cleanPaste returns pasted text as it is inserted.
func cleanPaste(s string) string {
	s = stripEscapes(pasteCleaner.Replace(s))
	s = strings.ReplaceAll(s, "\t", "    ")
	return strings.Map(func(r rune) rune {
		if r != '\n' && unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}

// Paste inserts cleaned s at the insertion point as one edit, replacing the
// selection if there is one.
func (e *Editor) Paste(s string) {
	s = cleanPaste(s)
	if s == "" {
		e.EndGroup()
		return
	}
	start, end := e.Selection()
	e.apply(start, end, s, editOther)
}

const (
	c1CSI = '\u009b'
	c1OSC = '\u009d'
	c1DCS = '\u0090'
	c1SOS = '\u0098'
	c1PM  = '\u009e'
	c1APC = '\u009f'
)

// stripEscapes removes each complete escape sequence. A sequence that never
// ends loses only its introducer.
func stripEscapes(s string) string {
	out, _ := scan(s)
	return out
}

// scanner counts the bytes it examines, so a test can check the scan stays linear.
type scanner struct {
	examined int
}

func scan(s string) (string, int) {
	var sc scanner
	var b strings.Builder
	for i := 0; i < len(s); {
		sc.examined++
		end := sc.escapeEnd(s, i)
		if end == i {
			_, size := utf8.DecodeRuneInString(s[i:])
			b.WriteString(s[i : i+size])
			end = i + size
		}
		i = end
	}
	return b.String(), sc.examined
}

// escapeEnd returns the index after the escape sequence or introducer that
// starts at i, or i when there is none.
func (sc *scanner) escapeEnd(s string, i int) int {
	if s[i] == 0x1b {
		return sc.escapeFormEnd(s, i)
	}
	r, size := utf8.DecodeRuneInString(s[i:])
	switch r {
	case c1CSI:
		return sc.sequenceEnd(s, i+size)
	case c1OSC:
		return sc.stringEnd(s, i+size, true)
	case c1DCS, c1SOS, c1PM, c1APC:
		return sc.stringEnd(s, i+size, false)
	}
	return i
}

func (sc *scanner) escapeFormEnd(s string, i int) int {
	if i+1 >= len(s) {
		return i + 1
	}
	switch c := s[i+1]; {
	case c == '[':
		return sc.sequenceEnd(s, i+2)
	case c == ']':
		return sc.stringEnd(s, i+2, true)
	case c == 'P' || c == 'X' || c == '^' || c == '_':
		return sc.stringEnd(s, i+2, false)
	case c >= 0x21 && c <= 0x2f:
		j := i + 1
		for j < len(s) && s[j] >= 0x21 && s[j] <= 0x2f {
			sc.examined++
			j++
		}
		if j < len(s) && s[j] >= 0x30 && s[j] <= 0x7e {
			return j + 1
		}
		return i + 1
	case c >= 0x30 && c <= 0x7e:
		return i + 2
	}
	return i + 1
}

// sequenceEnd returns the index after the control sequence whose parameters
// start at start, or start when it has no final byte.
func (sc *scanner) sequenceEnd(s string, start int) int {
	for j := start; j < len(s); j++ {
		sc.examined++
		switch b := s[j]; {
		case b >= 0x20 && b <= 0x3f:
		case b >= 0x40 && b <= 0x7e:
			return j + 1
		default:
			return start
		}
	}
	return start
}

// stringEnd returns the index after the control string whose text starts at
// start, or start when it has no terminator.
func (sc *scanner) stringEnd(s string, start int, bell bool) int {
	for j := start; j < len(s); j++ {
		sc.examined++
		switch b := s[j]; {
		case bell && b == 0x07:
			return j + 1
		case b == 0x1b && j+1 < len(s) && s[j+1] == '\\':
			return j + 2
		case b == 0xc2 && j+1 < len(s) && s[j+1] == 0x9c:
			return j + 2
		case b == 0xc2 && j+1 < len(s) && s[j+1] >= 0x80 && s[j+1] <= 0x9f:
			return start
		case b < 0x20 || b == 0x7f:
			return start
		}
	}
	return start
}
