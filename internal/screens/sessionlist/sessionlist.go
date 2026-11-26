// Package sessionlist is the screen that lists sessions.
package sessionlist

import (
	"strconv"
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

// Prompt is what the prompt box shows: the text rows in view, which row and
// column the cursor is in, how many rows are out of sight above and below,
// the message on the status row, and what Ctrl-C does next.
type Prompt struct {
	Rows      []string
	CursorRow int
	CursorCol int
	Above     int
	Below     int
	Status    string
	CtrlC     string
}

// Render draws the empty session list for a window of the given height,
// which must be at least 20 rows, with the prompt box above the key line.
func Render(height int, prompt Prompt) screens.Screen {
	statusRow := height - 2
	textTop := statusRow - len(prompt.Rows)
	rows := make([]string, height)
	rows[0] = titleStyle.Render(Title)
	rows[emptyRow] = strings.Repeat(" ", emptyColumn) + Empty
	rows[textTop-1] = labelStyle.Render(label) + counts(prompt.Above, prompt.Below)
	for i, text := range prompt.Rows {
		rows[textTop+i] = strings.Repeat(" ", indent) + text
	}
	if prompt.Status != "" {
		rows[statusRow] = strings.Repeat(" ", indent) + prompt.Status
	}
	rows[height-1] = keyStyle.Render("Enter") + " send · " + keyStyle.Render("Ctrl-C") + " " + prompt.CtrlC
	return screens.Screen{
		Content: strings.Join(rows, "\n"),
		CursorX: indent + prompt.CursorCol,
		CursorY: textTop + prompt.CursorRow,
	}
}

func counts(above, below int) string {
	var out string
	if above > 0 {
		out += " · " + lines(above) + " above"
	}
	if below > 0 {
		out += " · " + lines(below) + " below"
	}
	return out
}

func lines(n int) string {
	if n == 1 {
		return "1 line"
	}
	return strconv.Itoa(n) + " lines"
}
