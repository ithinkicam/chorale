package termtest

import (
	"image/color"

	uv "github.com/charmbracelet/ultraviolet"
)

// Cell is what one cell of the screen shows. The second half of a wide
// character is a blank cell.
type Cell struct {
	Text    string
	Inverse bool
	Fg, Bg  color.Color
}

// Cell returns the cell at column x and row y, counted from zero. A cell
// outside the screen is blank.
func (s *Screen) Cell(x, y int) Cell {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := s.emu.CellAt(x, y)
	if c == nil {
		return Cell{}
	}
	return Cell{
		Text:    c.Content,
		Inverse: c.Style.Attrs&uv.AttrReverse != 0,
		Fg:      c.Style.Fg,
		Bg:      c.Style.Bg,
	}
}
