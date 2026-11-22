// Package app holds the Bubble Tea model for the whole of Chorale.
package app

import (
	tea "charm.land/bubbletea/v2"

	"github.com/ithinkicam/chorale/internal/screens"
	"github.com/ithinkicam/chorale/internal/screens/sessionlist"
	"github.com/ithinkicam/chorale/internal/screens/toosmall"
)

// MinWidth and MinHeight are the smallest window Chorale draws its screens in.
const (
	MinWidth  = 60
	MinHeight = 20
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyPressMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m Model) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true
	if m.width == 0 || m.height == 0 {
		return view
	}

	var screen screens.Screen
	if m.width < MinWidth || m.height < MinHeight {
		screen = toosmall.Render(m.width, m.height, MinWidth, MinHeight)
	} else {
		screen = sessionlist.Render(m.width, m.height)
	}
	view.SetContent(screen.Content)
	cursor := tea.NewCursor(screen.CursorX, screen.CursorY)
	cursor.Blink = false // nothing animates
	view.Cursor = cursor
	return view
}
