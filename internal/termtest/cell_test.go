package termtest

import (
	"image/color"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestScreenCellReadsTextAndStyle(t *testing.T) {
	s := NewScreen(20, 3)
	write(t, s, "ab\x1b[7;34mcd\x1b[m\x1b[44mef\x1b[m合")

	cases := []struct {
		x    int
		want Cell
	}{
		{0, Cell{Text: "a"}},
		{2, Cell{Text: "c", Inverse: true, Fg: ansi.Blue}},
		{3, Cell{Text: "d", Inverse: true, Fg: ansi.Blue}},
		{4, Cell{Text: "e", Bg: ansi.Blue}},
		{6, Cell{Text: "合"}},
		{7, Cell{}},
	}
	for _, c := range cases {
		got := s.Cell(c.x, 0)
		if got.Text != c.want.Text || got.Inverse != c.want.Inverse ||
			!sameColour(got.Fg, c.want.Fg) || !sameColour(got.Bg, c.want.Bg) {
			t.Errorf("Cell(%d, 0) = %+v, want %+v", c.x, got, c.want)
		}
	}
}

func TestScreenCellOutsideTheScreenIsBlank(t *testing.T) {
	s := NewScreen(4, 2)
	if got := s.Cell(10, 10); got != (Cell{}) {
		t.Errorf("Cell(10, 10) = %+v, want a blank cell", got)
	}
}

func sameColour(a, b color.Color) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}
