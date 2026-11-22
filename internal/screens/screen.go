// Package screens holds what every screen renders to. Each screen is a
// package under it.
package screens

// Screen is a rendered screen: its content, and where the terminal cursor
// goes so macOS Zoom follows the focus.
type Screen struct {
	Content          string
	CursorX, CursorY int
}
