package termtest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var drawonce, redrawer, noreader string

func TestMain(m *testing.M) {
	os.Exit(runWithFixtures(m))
}

func runWithFixtures(m *testing.M) int {
	dir, err := os.MkdirTemp("", "termtest-fixtures-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer func() { _ = os.RemoveAll(dir) }()
	drawonce = filepath.Join(dir, "drawonce")
	redrawer = filepath.Join(dir, "redrawer")
	noreader = filepath.Join(dir, "noreader")
	for pkg, out := range map[string]string{
		"./testdata/drawonce": drawonce,
		"./testdata/redrawer": redrawer,
		"./testdata/noreader": noreader,
	} {
		if err := BuildBinary(pkg, out, ""); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	return m.Run()
}

func TestBuildBinaryWritesAnExecutable(t *testing.T) {
	info, err := os.Stat(drawonce)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&0o111 == 0 {
		t.Fatalf("%s is not executable", drawonce)
	}
}

func TestStartUsesTheGivenSizeAndFixedEnvironment(t *testing.T) {
	t.Setenv("NO_COLOR", "1") // must not leak into the program
	term := Start(t, drawonce, Options{Width: 60, Height: 20})

	term.WaitForText("size 60x20")
	text := term.Screen().Text()
	for _, want := range []string{"TERM xterm-256color", "NO_COLOR unset"} {
		if !strings.Contains(text, want) {
			t.Errorf("screen doesn't show %q:\n%s", want, text)
		}
	}
}

func TestStartCanTurnColourOff(t *testing.T) {
	term := Start(t, drawonce, Options{Width: 60, Height: 20, NoColor: true})
	term.WaitForText("NO_COLOR 1")
}

func TestCursorPositionAfterDrawing(t *testing.T) {
	term := Start(t, drawonce, Options{Width: 60, Height: 20})
	term.WaitForText("NO_COLOR")
	term.WaitQuiet(200 * time.Millisecond)

	x, y, visible := term.Screen().Cursor()
	if x != 2 || y != 4 || !visible {
		t.Fatalf("cursor at %d, %d, visible %v; want 2, 4, true", x, y, visible)
	}
}

func TestSendAndExitStatus(t *testing.T) {
	term := Start(t, drawonce, Options{Width: 60, Height: 20})
	term.WaitForText("size")
	term.Send("q")
	if code := term.Wait(); code != 3 {
		t.Fatalf("exit status %d, want 3", code)
	}
}

func TestRawKeepsEveryByte(t *testing.T) {
	term := Start(t, drawonce, Options{Width: 60, Height: 20})
	term.WaitForText("size")
	term.Send("c")
	term.WaitForText("red")
	if !strings.Contains(string(term.Raw()), "\x1b[31mred\x1b[0m") {
		t.Fatalf("raw output is missing the coloured text: %q", term.Raw())
	}
}

func TestResizeTellsTheProgram(t *testing.T) {
	term := Start(t, drawonce, Options{Width: 60, Height: 20})
	term.WaitForText("size 60x20")
	term.Resize(40, 15)
	term.WaitForText("size 40x15")
}

func TestIdleDetectorCatchesARedrawOnATimer(t *testing.T) {
	term := Start(t, redrawer, Options{Width: 60, Height: 20})
	term.WaitForText("the same text")
	if n := term.IdleBytes(); n == 0 {
		t.Fatal("the idle detector saw no output from a program that redraws every 700 ms")
	}
}

func TestIdleDetectorPassesAStillProgram(t *testing.T) {
	term := Start(t, drawonce, Options{Width: 60, Height: 20})
	term.WaitForText("NO_COLOR")
	if n := term.IdleBytes(); n != 0 {
		t.Fatalf("the idle detector saw %d bytes from a program that draws once", n)
	}
}

func TestAProgramThatNeverReadsItsInputCantHangTheScreen(t *testing.T) {
	term := Start(t, noreader, Options{Width: 60, Height: 20})
	term.WaitForText("done")
}
