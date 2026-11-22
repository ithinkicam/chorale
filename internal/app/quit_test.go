package app_test

import (
	"testing"

	"github.com/ithinkicam/chorale/internal/app"
	"github.com/ithinkicam/chorale/internal/termtest"
)

func TestQOnTheSessionListEndsTheProgram(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)
	mt.Type("q")
	mt.WaitFinished()
}
