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

func BenchmarkKeystrokeAtTheEndOfALargePrompt(b *testing.B) {
	var paste strings.Builder
	for i := range 10000 {
		fmt.Fprintf(&paste, "line number %d of a long prompt\n", i)
	}
	e := newEditor()
	e.Paste(paste.String())
	e.Move(MoveEnd, false)

	for b.Loop() {
		e.Insert("x")
		e.CursorPos()
	}
}

func BenchmarkUndoInALargePrompt(b *testing.B) {
	var paste strings.Builder
	for i := range 10000 {
		fmt.Fprintf(&paste, "line number %d of a long prompt\n", i)
	}
	e := newEditor()
	e.Paste(paste.String())
	e.Insert("x")

	for b.Loop() {
		e.Undo()
		e.Redo()
	}
}
