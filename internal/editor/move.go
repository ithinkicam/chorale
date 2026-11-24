package editor

import (
	"sort"
	"unicode"
	"unicode/utf8"
)

// Movement is a way of moving the insertion point.
type Movement int

// The ways of moving the insertion point.
const (
	MoveLeft Movement = iota
	MoveRight
	MoveWordLeft
	MoveWordRight
	MoveUp
	MoveDown
	MoveRowStart
	MoveRowEnd
	MoveStart
	MoveEnd
)

// Move moves the insertion point. With extend it leaves the selection anchor
// where it was, so the selection grows or shrinks.
func (e *Editor) Move(m Movement, extend bool) {
	if m != MoveUp && m != MoveDown {
		e.wantColSet = false
	}
	row := e.caretRow()
	e.endOfRow = false
	e.last = editOther

	switch m {
	case MoveLeft:
		e.caret = e.stepLeft(extend)
	case MoveRight:
		e.caret = e.stepRight(extend)
	case MoveWordLeft:
		e.caret = e.wordLeft(e.caret)
	case MoveWordRight:
		e.caret = e.wordRight(e.caret)
	case MoveUp:
		e.moveVertically(row, -1)
	case MoveDown:
		e.moveVertically(row, 1)
	case MoveRowStart:
		e.caret = e.layout()[row].Start
	case MoveRowEnd:
		e.caret = e.layout()[row].End
		e.endOfRow = true
	case MoveStart:
		e.caret = 0
	case MoveEnd:
		e.caret = len(e.gs)
	}

	if !extend {
		e.anchor = e.caret
	}
}

func (e *Editor) stepLeft(extend bool) int {
	if start, _ := e.Selection(); e.HasSelection() && !extend {
		return start
	}
	return max(e.caret-1, 0)
}

func (e *Editor) stepRight(extend bool) int {
	if _, end := e.Selection(); e.HasSelection() && !extend {
		return end
	}
	return min(e.caret+1, len(e.gs))
}

// EndOfRow reports whether the insertion point is drawn at the end of the row
// above a soft wrap, rather than at the start of the row below it.
func (e *Editor) EndOfRow() bool {
	return e.endOfRow
}

// CursorPos returns the row and the column in it where the insertion point is drawn.
func (e *Editor) CursorPos() (row, col int) {
	row = e.caretRow()
	return row, e.colIn(e.layout()[row], e.caret)
}

// Rows returns the text laid out to the editor's width.
func (e *Editor) Rows() []Row {
	return e.layout()
}

func (e *Editor) layout() []Row {
	if !e.rowsOK {
		e.rows = Layout(e.gs, e.width, e.method)
		e.rowsOK = true
	}
	return e.rows
}

func (e *Editor) caretRow() int {
	rows := e.layout()
	r := sort.Search(len(rows), func(i int) bool { return rows[i].End >= e.caret })
	if r+1 < len(rows) && rows[r].End == e.caret && rows[r+1].Start == e.caret && !e.endOfRow {
		return r + 1
	}
	return r
}

func (e *Editor) colIn(row Row, index int) int {
	if index >= row.End {
		return row.EndCol
	}
	return row.Cols[index-row.Start]
}

func (e *Editor) moveVertically(from, dir int) {
	rows := e.layout()
	if !e.wantColSet {
		e.wantCol = e.colIn(rows[from], e.caret)
		e.wantColSet = true
	}
	to := from + dir
	if to < 0 {
		e.caret = 0
		return
	}
	if to >= len(rows) {
		e.caret = len(e.gs)
		return
	}

	row := rows[to]
	limit := row.End
	if to+1 < len(rows) && rows[to+1].Start == row.End {
		limit = row.End - 1
	}
	i := row.Start
	for i < limit && e.colIn(row, i+1) <= e.wantCol {
		i++
	}
	e.caret = i
}

type charClass int

const (
	gap charClass = iota
	letter
	symbol
)

func classOf(g string) charClass {
	r, _ := utf8.DecodeRuneInString(g)
	switch {
	case r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r):
		return letter
	case unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsControl(r):
		return gap
	default:
		return symbol
	}
}

func (e *Editor) wordLeft(i int) int {
	for i > 0 && classOf(e.gs[i-1]) == gap {
		i--
	}
	if i > 0 && classOf(e.gs[i-1]) == symbol {
		return i - 1
	}
	for i > 0 && classOf(e.gs[i-1]) == letter {
		i--
	}
	return i
}

func (e *Editor) wordRight(i int) int {
	for i < len(e.gs) && classOf(e.gs[i]) == gap {
		i++
	}
	if i < len(e.gs) && classOf(e.gs[i]) == symbol {
		return i + 1
	}
	for i < len(e.gs) && classOf(e.gs[i]) == letter {
		i++
	}
	return i
}
