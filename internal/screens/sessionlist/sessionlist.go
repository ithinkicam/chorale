// Package sessionlist is the screen that lists sessions. Until sessions
// exist, it shows the empty state.
package sessionlist

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// The words on the empty session list.
const (
	Title = "♪ chorale"
	Empty = "No sessions yet."
)

// Screen is the rendered screen: its content, and where the terminal
// cursor goes so macOS Zoom follows the focus.
type Screen struct {
	Content          string
	CursorX, CursorY int
}

// The default theme uses the terminal's own palette, so the user's
// terminal colours apply. Colour is never the only signal: every state is
// in words.
var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Magenta)
	keyStyle   = lipgloss.NewStyle().Bold(true)
)

// emptyRow and emptyColumn place the empty message, and the cursor on it.
const (
	emptyRow    = 2
	emptyColumn = 2
)

// Render draws the empty session list for a window of the given size,
// which must be at least 60 by 20.
func Render(width, height int) Screen {
	rows := make([]string, height)
	rows[0] = titleStyle.Render(Title)
	rows[emptyRow] = strings.Repeat(" ", emptyColumn) + Empty
	rows[height-1] = keyStyle.Render("q") + " quit"
	_ = width // every row is far narrower than the 60-column minimum
	return Screen{
		Content: strings.Join(rows, "\n"),
		CursorX: emptyColumn,
		CursorY: emptyRow,
	}
}
