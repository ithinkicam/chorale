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
	method        ansi.Method
	pendingQuit   bool
	status        string
	send          func(prompt string) (status string)
}

const (
	askToQuit      = "Press Ctrl-C again to quit."
	promptCleared  = "Prompt cleared. Ctrl-Z brings it back."
	noAgentMessage = "No agent is connected yet. Your prompt is kept."
)

const textInset = 3

// New returns the model Chorale starts with.
func New() Model {
	return Model{prompt: editor.New(MinWidth-textInset, ansi.WcWidth), method: ansi.WcWidth, send: noAgent}
}

func noAgent(string) string { return noAgentMessage }

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.pendingQuit = false
		m.layOutPrompt()
	case tea.ModeReportMsg:
		if msg.Mode == ansi.ModeUnicodeCore && supportsGraphemes(msg.Value) {
			m.method = ansi.GraphemeWidth
			m.layOutPrompt()
		}
	case tea.PasteMsg:
		m.pendingQuit = false
		m.status = ""
		if !m.tooSmall() {
			m.prompt.Paste(msg.Content)
		}
	case tea.KeyPressMsg:
		return m.press(msg)
	}
	return m, nil
}

func supportsGraphemes(setting ansi.ModeSetting) bool {
	return setting == ansi.ModeSet || setting == ansi.ModeReset || setting == ansi.ModePermanentlySet
}

func (m Model) layOutPrompt() {
	if !m.tooSmall() {
		m.prompt.SetWidth(m.width-textInset, m.method)
	}
}

func (m Model) press(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	action, bound := editor.Lookup(msg.String())
	m.status = ""
	if bound && action.Kind == editor.ActionInterrupt {
		if !m.tooSmall() {
			m.prompt.EndGroup()
		}
		return m.interrupt()
	}
	m.pendingQuit = false
	switch {
	case m.tooSmall():
	case bound:
		m.prompt.Apply(action)
		if action.Kind == editor.ActionSend && strings.Trim(m.prompt.Text(), " \n") != "" {
			m.status = m.send(m.prompt.Text())
		}
	case msg.Text != "":
		m.prompt.Insert(msg.Text)
	default:
		m.prompt.EndGroup()
	}
	return m, nil
}

func (m Model) interrupt() (tea.Model, tea.Cmd) {
	switch {
	case m.pendingQuit:
		return m, tea.Quit
	case m.tooSmall() || m.prompt.Text() == "":
		m.pendingQuit = true
	case !m.prompt.HasSelection():
		m.prompt.Clear()
		m.status = promptCleared
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
	prompt := sessionlist.Prompt{
		Text:      strings.Join(m.prompt.Chars()[r.Start:r.End], ""),
		CursorCol: col,
	}
	switch {
	case m.prompt.HasSelection():
		prompt.CtrlC = "copy"
	case m.prompt.Text() != "":
		prompt.CtrlC = "clear"
	default:
		prompt.CtrlC = "quit"
	}
	prompt.Status = m.status
	if m.pendingQuit {
		prompt.Status = askToQuit
	}
	return prompt
}
