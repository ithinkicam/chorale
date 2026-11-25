// Package app holds the Bubble Tea model for the whole of Chorale.
package app

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/ithinkicam/chorale/internal/editor"
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
	prompt        *editor.Editor
}

const textInset = 3

// New returns the model Chorale starts with.
func New() Model {
	return Model{prompt: editor.New(MinWidth-textInset, ansi.WcWidth)}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		if !m.tooSmall() {
			m.prompt.SetWidth(m.width-textInset, ansi.WcWidth)
		}
	case tea.KeyPressMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) tooSmall() bool {
	return m.width < MinWidth || m.height < MinHeight
}

// View implements tea.Model.
func (m Model) View() tea.View {
	view := tea.NewView("")
	view.AltScreen = true
	if m.width == 0 || m.height == 0 {
		return view
	}

	var screen screens.Screen
	if m.tooSmall() {
		screen = toosmall.Render(m.width, m.height, MinWidth, MinHeight)
	} else {
		screen = sessionlist.Render(m.height, m.promptBox())
	}
	view.SetContent(screen.Content)
	cursor := tea.NewCursor(screen.CursorX, screen.CursorY)
	cursor.Blink = false // nothing animates
	view.Cursor = cursor
	return view
}

func (m Model) promptBox() sessionlist.Prompt {
	row, col := m.prompt.CursorPos()
	r := m.prompt.Rows()[row]
	return sessionlist.Prompt{
		Text:      strings.Join(m.prompt.Chars()[r.Start:r.End], ""),
		CursorCol: col,
	}
}
