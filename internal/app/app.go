// Package app holds the Bubble Tea model for the whole of Chorale.
package app

import (
	tea "charm.land/bubbletea/v2"

	"github.com/ithinkicam/chorale/internal/screens/sessionlist"
)

// Model is Chorale's top-level Bubble Tea model.
type Model struct {
	width, height int
}

// New returns the model Chorale starts with.
func New() Model { return Model{} }

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width, m.height = size.Width, size.Height
	}
	return m, nil
}

// View implements tea.Model. Nothing is drawn until the window size is
// known, so Chorale never shows a screen meant for another size.
func (m Model) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true
	if m.width == 0 || m.height == 0 {
		return view
	}

	screen := sessionlist.Render(m.width, m.height)
	view.SetContent(screen.Content)
	cursor := tea.NewCursor(screen.CursorX, screen.CursorY)
	cursor.Blink = false // nothing animates
	view.Cursor = cursor
	return view
}
