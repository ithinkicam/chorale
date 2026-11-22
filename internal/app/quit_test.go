package app_test

import (
	"testing"

	"github.com/ithinkicam/chorale/internal/app"
	"github.com/ithinkicam/chorale/internal/termtest"
)

// Part of the scenario "Quitting restores the terminal": q on the session
// list ends the program. The terminal side is checked in e2e.
func TestQOnTheSessionListEndsTheProgram(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)
	mt.Type("q")
	mt.WaitFinished()
}
