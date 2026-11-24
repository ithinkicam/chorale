package editor

const (
	maxSnapshots    = 200
	maxHistoryBytes = 16 << 20
)

type snapshot struct {
	text   string
	caret  int
	anchor int
}

type history struct {
	snaps []snapshot
	cur   int
	size  int
}

func newHistory(text string) history {
	return history{snaps: []snapshot{{text: text}}, size: len(text)}
}

// HistoryLen returns how many snapshots the undo history holds.
func (e *Editor) HistoryLen() int {
	return len(e.hist.snaps)
}

// HistorySize returns the bytes of text the undo history holds.
func (e *Editor) HistorySize() int {
	return e.hist.size
}

// Undo reverts the last edit, and reports whether there was one.
func (e *Editor) Undo() bool {
	h := &e.hist
	if h.cur == 0 {
		return false
	}
	h.snaps[h.cur].caret, h.snaps[h.cur].anchor = e.caret, e.anchor
	h.cur--
	e.restore(h.snaps[h.cur])
	return true
}

// Redo repeats the last undone edit, and reports whether there was one.
func (e *Editor) Redo() bool {
	h := &e.hist
	if h.cur+1 >= len(h.snaps) {
		return false
	}
	h.snaps[h.cur].caret, h.snaps[h.cur].anchor = e.caret, e.anchor
	h.cur++
	e.restore(h.snaps[h.cur])
	return true
}

func (e *Editor) restore(s snapshot) {
	e.gs = Graphemes(s.text)
	e.caret, e.anchor = s.caret, s.anchor
	e.endOfRow = false
	e.wantColSet = false
	e.last = editOther
	e.rowsOK = false
}

// record adds the editor's current state to the history. The state before
// the edit keeps the insertion point and anchor it had then, so undo returns
// to them. A joining edit replaces the latest snapshot instead.
func (e *Editor) record(prevCaret, prevAnchor int, joining bool) {
	h := &e.hist
	now := snapshot{text: e.Text(), caret: e.caret, anchor: e.anchor}
	if joining && h.cur == len(h.snaps)-1 {
		h.size += len(now.text) - len(h.snaps[h.cur].text)
		h.snaps[h.cur] = now
		h.trim()
		return
	}

	h.snaps[h.cur].caret, h.snaps[h.cur].anchor = prevCaret, prevAnchor
	for _, dropped := range h.snaps[h.cur+1:] {
		h.size -= len(dropped.text)
	}
	h.snaps = append(h.snaps[:h.cur+1], now)
	h.cur++
	h.size += len(now.text)
	h.trim()
}

func (h *history) trim() {
	for h.cur > 0 && (len(h.snaps) > maxSnapshots || h.size > maxHistoryBytes) {
		h.size -= len(h.snaps[0].text)
		h.snaps = append(h.snaps[:0], h.snaps[1:]...)
		h.cur--
	}
}
