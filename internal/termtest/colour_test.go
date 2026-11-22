package termtest

import (
	"slices"
	"testing"
)

func TestColourCodes(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want []ColourCode
	}{
		{"none", "plain \x1b[1mbold\x1b[0m \x1b[4munderline\x1b[22;24m", nil},
		{"palette foreground", "\x1b[31mx", []ColourCode{{Sequence: "\x1b[31m", Palette: true}}},
		{"bright background", "\x1b[105mx", []ColourCode{{Sequence: "\x1b[105m", Palette: true}}},
		{"default colours", "\x1b[39;49mx", []ColourCode{{Sequence: "\x1b[39;49m", Palette: true}}},
		{"colour among other attributes", "\x1b[1;32;4mx", []ColourCode{{Sequence: "\x1b[1;32;4m", Palette: true}}},
		{"256 colours", "\x1b[38;5;208mx", []ColourCode{{Sequence: "\x1b[38;5;208m", Palette: false}}},
		{"24-bit colon form", "\x1b[48:2::10:20:30mx", []ColourCode{{Sequence: "\x1b[48:2::10:20:30m", Palette: false}}},
		{"cursor moves are not colour", "\x1b[31;5H\x1b[2J", nil},
		{"private modes are not colour", "\x1b[?31h", nil},
		{"split across two", "\x1b[3", nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ColourCodes([]byte(c.raw)); !slices.Equal(got, c.want) {
				t.Fatalf("ColourCodes(%s) = %+v, want %+v", quote(c.raw), got, c.want)
			}
		})
	}
}
