// Package sessionlist is the screen that lists sessions.
package sessionlist

import (
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/ithinkicam/chorale/internal/screens"
)

// Title and Empty are the words on the empty session list.
const (
	Title = "♪ chorale"
	Empty = "No sessions yet."
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Magenta)
	keyStyle   = lipgloss.NewStyle().Bold(true)
	labelStyle = lipgloss.NewStyle().Bold(true)
)

const (
	emptyRow    = 2
	emptyColumn = 2
	indent      = 2
	label       = "▸ Prompt"
)

// Prompt is what the prompt box shows: one row of text, the column the
// cursor is in, and the message on the status row.
type Prompt struct {
	Text      string
	CursorCol int
	Status    string
}

// Render draws the empty session list for a window of the given height,
// which must be at least 20 rows, with the prompt box above the key line.
func Render(height int, prompt Prompt) screens.Screen {
	statusRow := height - 2
	textRow := height - 3
	rows := make([]string, height)
	rows[0] = titleStyle.Render(Title)
	rows[emptyRow] = strings.Repeat(" ", emptyColumn) + Empty
	rows[textRow-1] = labelStyle.Render(label)
	rows[textRow] = strings.Repeat(" ", indent) + prompt.Text
	if prompt.Status != "" {
		rows[statusRow] = strings.Repeat(" ", indent) + prompt.Status
	}
	rows[height-1] = keyStyle.Render("Enter") + " send · " + keyStyle.Render("Ctrl-C") + " quit"
	return screens.Screen{
		Content: strings.Join(rows, "\n"),
		CursorX: indent + prompt.CursorCol,
		CursorY: textRow,
	}
}
