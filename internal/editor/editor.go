// Package editor is the prompt editor's logic: grapheme text, the insertion
// point, selection, layout, movement, deletion, history and paste cleaning.
package editor

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Editor holds a prompt and the insertion point in it.
type Editor struct {
	gs     []string
	caret  int
	anchor int
	width  int
	method ansi.Method

	endOfRow   bool
	wantCol    int
	wantColSet bool
	rows       []Row
	rowsOK     bool

	hist history
	last editKind
}

type editKind int

const (
	editOther editKind = iota
	editTyping
	editBackspace
	editDelete
)

// New returns an empty editor that lays text out to width columns, measured by method.
func New(width int, method ansi.Method) *Editor {
	return &Editor{width: width, method: method, hist: newHistory("")}
}

// SetWidth changes the width the text is laid out to, and how characters are measured.
func (e *Editor) SetWidth(width int, method ansi.Method) {
	e.width = width
	e.method = method
	e.rowsOK = false
}

// Graphemes splits s into user-perceived characters.
func Graphemes(s string) []string {
	var out []string
	for s != "" {
		g, _ := ansi.FirstGraphemeCluster(s, ansi.GraphemeWidth)
		out = append(out, g)
		s = s[len(g):]
	}
	return out
}

// Text returns the whole prompt.
func (e *Editor) Text() string {
	return strings.Join(e.gs, "")
}

// Len returns the number of characters in the prompt.
func (e *Editor) Len() int {
	return len(e.gs)
}

// Caret returns the insertion point as a character index.
func (e *Editor) Caret() int {
	return e.caret
}

// Insert types s at the insertion point, replacing the selection if there is
// one. Undo removes what was typed a word at a time.
func (e *Editor) Insert(s string) {
	start, end := e.Selection()
	prev := ""
	if start > 0 {
		prev = e.gs[start-1]
	}

	var chunk strings.Builder
	for _, g := range Graphemes(s) {
		chunk.WriteString(g)
		closes := isSpace(g) && !isSpace(prev)
		prev = g
		if !closes {
			continue
		}
		e.apply(start, end, chunk.String(), editTyping)
		e.last = editOther
		chunk.Reset()
		start, end = e.caret, e.caret
	}
	if chunk.Len() > 0 {
		e.apply(start, end, chunk.String(), editTyping)
	}
}

func isSpace(g string) bool {
	return g == "" || g == " " || g == "\n"
}

// apply swaps the characters from start to end for s as one edit in the
// history. An edit of the same kind as the last one joins its group, except
// for editOther.
func (e *Editor) apply(start, end int, s string, kind editKind) {
	caret, anchor := e.caret, e.anchor
	if !e.replace(start, end, s) {
		return
	}
	e.record(caret, anchor, kind != editOther && kind == e.last)
	e.last = kind
}

// Selection returns the start and end of the selected characters, equal when
// nothing is selected.
func (e *Editor) Selection() (start, end int) {
	return min(e.anchor, e.caret), max(e.anchor, e.caret)
}

// HasSelection reports whether any characters are selected.
func (e *Editor) HasSelection() bool {
	return e.anchor != e.caret
}

// SelectedText returns the selected characters.
func (e *Editor) SelectedText() string {
	start, end := e.Selection()
	return strings.Join(e.gs[start:end], "")
}

// SelectedCount returns how many characters are selected.
func (e *Editor) SelectedCount() int {
	start, end := e.Selection()
	return end - start
}

// SelectAll selects the whole prompt, leaving the insertion point at its end.
func (e *Editor) SelectAll() {
	e.anchor = 0
	e.caret = len(e.gs)
	e.endOfRow = false
	e.wantColSet = false
	e.last = editOther
}

// replace swaps the characters from start to end for s, and puts the insertion
// point after s. The characters on either side of the change are segmented
// again together with s, so accents and joiners attach to their neighbours.
func (e *Editor) replace(start, end int, s string) bool {
	if start == end && s == "" {
		return false
	}
	from := max(start-1, 0)
	to := min(end+1, len(e.gs))
	left := strings.Join(e.gs[from:start], "")
	right := strings.Join(e.gs[end:to], "")

	joined := Graphemes(left + s + right)
	caret := from + len(Graphemes(left+s))

	gs := make([]string, 0, len(e.gs)-(to-from)+len(joined))
	gs = append(gs, e.gs[:from]...)
	gs = append(gs, joined...)
	gs = append(gs, e.gs[to:]...)

	e.gs = gs
	e.caret = min(caret, from+len(joined))
	e.anchor = e.caret
	e.endOfRow = false
	e.wantColSet = false
	e.rowsOK = false
	return true
}

// Chars returns the prompt's characters, indexed as Rows are. The caller must not change it.
func (e *Editor) Chars() []string {
	return e.gs
}
