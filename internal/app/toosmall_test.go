package app_test

import (
	"strings"
	"testing"

	"github.com/ithinkicam/chorale/internal/app"
	"github.com/ithinkicam/chorale/internal/termtest"
)

const tooSmallWords = "Chorale needs at least 60 by 20."

func TestStartedAt40Columns(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 40, Height: 20})
	mt.WaitForText(tooSmallWords)
	mt.WaitForText(keyWords)
	assertViewFits(t, 40, 20)
	mt.RequireGolden()
}

func TestStartedAt15Rows(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 80, Height: 15})
	mt.WaitForText(tooSmallWords)
	mt.WaitForText(keyWords)
	assertViewFits(t, 80, 15)
	mt.RequireGolden()
}

func TestMessageWrapsInATinyWindow(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 20, Height: 10})
	mt.WaitFor("the wrapped message", func(screen string) bool {
		return strings.Contains(joinWithSpaces(screen), tooSmallWords)
	})
	mt.WaitForText(keyWords)
	assertViewFits(t, 20, 10)
	mt.RequireGolden()
}

func TestTooSmallSnapshotWithColourOff(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 40, Height: 20, NoColor: true})
	mt.WaitForText(tooSmallWords)
	mt.WaitForText(keyWords)
	mt.RequireGolden()
}

func TestWindowWidenedTo60Columns(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 40, Height: 20})
	mt.WaitForText(tooSmallWords)

	mt.Resize(60, 20)
	mt.WaitFor("the empty session list", func(screen string) bool {
		return strings.Contains(screen, titleWords) &&
			strings.Contains(screen, emptyWords) &&
			strings.Contains(screen, keyWords) &&
			!strings.Contains(screen, "Chorale needs at least")
	})
}

func TestWindowNarrowedWhileRunning(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)

	mt.Resize(40, 20)
	mt.WaitFor("the too-small message", func(screen string) bool {
		return strings.Contains(screen, tooSmallWords) &&
			strings.Contains(screen, keyWords) &&
			!strings.Contains(screen, emptyWords)
	})
	assertViewFits(t, 40, 20)
}

func joinWithSpaces(screen string) string {
	return strings.Join(strings.Fields(screen), " ")
}
