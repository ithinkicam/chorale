package termtest

import (
	"regexp"
	"strconv"
	"strings"
)

// ColourCode is one Select Graphic Rendition sequence that sets a
// foreground or background colour.
type ColourCode struct {
	Sequence string
	Palette  bool
}

var sgrPattern = regexp.MustCompile(`\x1b\[([0-9;:]*)m`)

// ColourCodes returns every colour code in raw, in order.
func ColourCodes(raw []byte) []ColourCode {
	var out []ColourCode
	for _, m := range sgrPattern.FindAllSubmatch(raw, -1) {
		hasColour, palette := classifySGR(string(m[1]))
		if hasColour {
			out = append(out, ColourCode{Sequence: string(m[0]), Palette: palette})
		}
	}
	return out
}

func classifySGR(params string) (hasColour, palette bool) {
	palette = true
	fields := strings.Split(params, ";")
	for i := 0; i < len(fields); i++ {
		field := fields[i]
		head, _, hasSub := strings.Cut(field, ":")
		n, err := strconv.Atoi(head)
		if err != nil {
			continue
		}
		switch {
		case n == 38 || n == 48:
			hasColour, palette = true, false
			if hasSub {
				continue // the colon form carries its arguments in this field
			}
			if i+1 < len(fields) {
				switch fields[i+1] {
				case "5":
					i += 2
				case "2":
					i += 4
				}
			}
		case n >= 30 && n <= 39, n >= 40 && n <= 49, n >= 90 && n <= 97, n >= 100 && n <= 107:
			hasColour = true
		}
	}
	return hasColour, palette && hasColour
}
