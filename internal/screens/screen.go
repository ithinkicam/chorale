// Package screens holds what every screen renders to.
package screens

// Screen is a rendered screen: its content and where the terminal cursor goes.
type Screen struct {
	Content          string
	CursorX, CursorY int
}
