package app

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// PromptState is what a Model's prompt holds, as the tests read it.
type PromptState struct {
	Text        string
	Caret       int
	SelStart    int
	SelEnd      int
	Selected    string
	Rows        []string
	RowWidths   []int
	CaretRow    int
	PendingQuit bool
}

// StateOf reads the prompt of m, which must be a Model.
func StateOf(m tea.Model) PromptState {
	model := m.(Model)
	p := model.prompt
	start, end := p.Selection()
	row, _ := p.CursorPos()
	state := PromptState{
		Text:        p.Text(),
		Caret:       p.Caret(),
		SelStart:    start,
		SelEnd:      end,
		Selected:    p.SelectedText(),
		CaretRow:    row,
		PendingQuit: model.pendingQuit,
	}
	for _, r := range p.Rows() {
		state.Rows = append(state.Rows, strings.Join(p.Chars()[r.Start:r.End], ""))
		state.RowWidths = append(state.RowWidths, r.EndCol)
	}
	return state
}

// PendingQuit reports whether m, which must be a Model, is waiting for a second Ctrl-C.
func PendingQuit(m tea.Model) bool {
	return StateOf(m).PendingQuit
}

// WithSend returns m, which must be a Model, with send as its send hook.
func WithSend(m tea.Model, send func(prompt string) (status string)) tea.Model {
	model := m.(Model)
	model.send = send
	return model
}
