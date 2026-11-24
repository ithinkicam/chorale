package editor

// Backspace removes the selection, or the character before the insertion point.
func (e *Editor) Backspace() {
	e.remove(max(e.caret-1, 0), e.caret)
}

// Delete removes the selection, or the character after the insertion point.
func (e *Editor) Delete() {
	e.remove(e.caret, min(e.caret+1, len(e.gs)))
}

// DeleteWordBack removes the selection, or the word before the insertion point.
func (e *Editor) DeleteWordBack() {
	e.remove(e.wordLeft(e.caret), e.caret)
}

// DeleteWordForward removes the selection, or the word after the insertion point.
func (e *Editor) DeleteWordForward() {
	e.remove(e.caret, e.wordRight(e.caret))
}

// DeleteToRowStart removes the selection, or everything from the start of the
// row to the insertion point.
func (e *Editor) DeleteToRowStart() {
	e.remove(e.layout()[e.caretRow()].Start, e.caret)
}

// DeleteToRowEnd removes the selection, or everything from the insertion
// point to the end of the row.
func (e *Editor) DeleteToRowEnd() {
	e.remove(e.caret, e.layout()[e.caretRow()].End)
}

// Clear removes the whole prompt.
func (e *Editor) Clear() {
	e.replace(0, len(e.gs), "")
}

func (e *Editor) remove(start, end int) {
	if e.HasSelection() {
		start, end = e.Selection()
	}
	e.replace(start, end, "")
}
