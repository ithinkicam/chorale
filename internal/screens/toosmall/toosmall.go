// Package toosmall is the screen shown when the window is smaller than
// Chorale's minimum size.
package toosmall

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"github.com/ithinkicam/chorale/internal/screens"
)

var keyStyle = lipgloss.NewStyle().Bold(true)

// Render draws the message for a window of the given size, saying the
// size Chorale needs. The message wraps to the window's width, and the
// key line stays on the last row.
func Render(width, height, minWidth, minHeight int) screens.Screen {
	message := fmt.Sprintf("Chorale needs at least %d by %d.", minWidth, minHeight)
	wrapped := strings.Split(lipgloss.NewStyle().Width(width).Render(message), "\n")
	for i, line := range wrapped {
		wrapped[i] = strings.TrimRight(line, " ")
	}
	keys := keyStyle.Render("q") + " quit"

	rows := make([]string, max(height, 1))
	last := len(rows) - 1
	// Keep the key line visible even when the message doesn't fit.
	copy(rows[:last], wrapped)
	rows[last] = keys
	if last == 0 {
		rows[0] = keys
	}
	return screens.Screen{Content: strings.Join(rows, "\n")}
}
