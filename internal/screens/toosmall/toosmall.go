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

// Render draws the message for a window of the given size, saying the size
// Chorale needs.
func Render(width, height, minWidth, minHeight int) screens.Screen {
	message := fmt.Sprintf("Chorale needs at least %d by %d.", minWidth, minHeight)
	wrapped := strings.Split(lipgloss.NewStyle().Width(width).Render(message), "\n")
	for i, line := range wrapped {
		wrapped[i] = strings.TrimRight(line, " ")
	}
	rows := make([]string, max(height, 1))
	copy(rows, wrapped)
	if len(wrapped) < len(rows) {
		rows[len(rows)-1] = keyStyle.Render("q") + " quit"
	}
	return screens.Screen{Content: strings.Join(rows, "\n")}
}
