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
)

const (
	emptyRow    = 2
	emptyColumn = 2
)

// Render draws the empty session list for a window of the given height,
// which must be at least 20 rows.
func Render(height int) screens.Screen {
	rows := make([]string, height)
	rows[0] = titleStyle.Render(Title)
	rows[emptyRow] = strings.Repeat(" ", emptyColumn) + Empty
	rows[height-1] = keyStyle.Render("q") + " quit"
	return screens.Screen{
		Content: strings.Join(rows, "\n"),
		CursorX: emptyColumn,
		CursorY: emptyRow,
	}
}
