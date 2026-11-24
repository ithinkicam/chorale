package editor

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkLargePaste(b *testing.B) {
	var paste strings.Builder
	for i := range 10000 {
		fmt.Fprintf(&paste, "line number %d of a long prompt\n", i)
	}
	text := paste.String()

	for b.Loop() {
		e := newEditor()
		e.Paste(text)
		e.Rows()
		e.Move(MoveStart, false)
		e.Move(MoveEnd, false)
		e.CursorPos()
	}
}
