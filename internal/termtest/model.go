package termtest

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/charmbracelet/x/exp/teatest/v2"
)

// ModelTerm runs a Bubble Tea model through teatest and reads back its
// screen as plain text.
type ModelTerm struct {
	t      testing.TB
	tm     *teatest.TestModel
	screen *Screen

	finished bool
}

// StartModel runs m at the given size. Only opts.Width, Height and NoColor
// are used.
func StartModel(t testing.TB, m tea.Model, opts Options) *ModelTerm {
	t.Helper()
	profile := colorprofile.ANSI
	env := []string{"TERM=xterm-256color"}
	if opts.NoColor {
		profile = colorprofile.Ascii
		env = append(env, "NO_COLOR=1")
	}
	mt := &ModelTerm{
		t:      t,
		screen: NewScreen(opts.Width, opts.Height),
		tm: teatest.NewTestModel(t, m,
			teatest.WithInitialTermSize(opts.Width, opts.Height),
			teatest.WithProgramOptions(
				tea.WithColorProfile(profile),
				tea.WithEnvironment(env),
			),
		),
	}
	t.Cleanup(mt.Finish)
	return mt
}

func (mt *ModelTerm) pump() {
	data, err := io.ReadAll(mt.tm.Output())
	if err != nil {
		mt.t.Fatalf("reading the program's output: %v", err)
	}
	if len(data) > 0 {
		_, _ = mt.screen.Write(onlcr(data))
	}
}

// onlcr turns each line feed into a carriage return and line feed, because
// teatest writes to a buffer with no terminal driver.
func onlcr(data []byte) []byte {
	return bytes.ReplaceAll(data, []byte("\n"), []byte("\r\n"))
}

// Screen returns the screen, up to date with everything drawn so far.
func (mt *ModelTerm) Screen() *Screen {
	mt.pump()
	return mt.screen
}

// Lines returns the screen's rows as plain text.
func (mt *ModelTerm) Lines() []string { return mt.Screen().Lines() }

// Text returns the screen as plain text.
func (mt *ModelTerm) Text() string { return mt.Screen().Text() }

// WaitForText waits until the screen shows text, and fails the test if it
// doesn't within DefaultTimeout.
func (mt *ModelTerm) WaitForText(text string) {
	mt.t.Helper()
	mt.WaitFor("the screen to show "+text, func(screen string) bool {
		return strings.Contains(screen, text)
	})
}

// WaitFor waits until ok accepts the screen's text, and fails the test
// with what if it never does within DefaultTimeout.
func (mt *ModelTerm) WaitFor(what string, ok func(screen string) bool) {
	mt.t.Helper()
	deadline := time.Now().Add(DefaultTimeout)
	for {
		text := mt.Text()
		if ok(text) {
			return
		}
		if time.Now().After(deadline) {
			mt.t.Fatalf("timed out waiting for %s. The screen shows:\n%s", what, text)
		}
		time.Sleep(pollInterval)
	}
}

// Send sends a message to the program.
func (mt *ModelTerm) Send(msg tea.Msg) { mt.tm.Send(msg) }

// Type sends each character of s as a key press.
func (mt *ModelTerm) Type(s string) { mt.tm.Type(s) }

// Resize resizes the screen and tells the program, as a terminal would.
func (mt *ModelTerm) Resize(width, height int) {
	mt.pump()
	mt.screen.Resize(width, height)
	mt.tm.Send(tea.WindowSizeMsg{Width: width, Height: height})
}

// WaitFinished waits for the program to end by itself, and fails the test
// if it hasn't within DefaultTimeout.
func (mt *ModelTerm) WaitFinished() {
	mt.t.Helper()
	mt.tm.WaitFinished(mt.t, teatest.WithFinalTimeout(DefaultTimeout), teatest.WithTimeoutFn(func(tb testing.TB) {
		tb.Fatalf("the program didn't end. The screen shows:\n%s", mt.Text())
	}))
	mt.finished = true
}

// Finish ends the program from the test side. It runs automatically when
// the test ends.
func (mt *ModelTerm) Finish() {
	mt.end()
	_ = mt.screen.Close()
}

// FinalView ends the program and returns the content of the view the
// running model last drew. Read the screen before calling it, because
// ending the program leaves the alternate screen.
func (mt *ModelTerm) FinalView() string {
	mt.t.Helper()
	mt.end()
	return mt.tm.FinalModel(mt.t, teatest.WithFinalTimeout(DefaultTimeout)).View().Content
}

func (mt *ModelTerm) end() {
	if !mt.finished {
		mt.finished = true
		_ = mt.tm.Quit()
		mt.tm.WaitFinished(mt.t, teatest.WithFinalTimeout(DefaultTimeout))
	}
}

// RequireGolden compares the screen's plain text with the golden file
// testdata/<test name>.golden, which `make golden` rewrites on purpose only.
func (mt *ModelTerm) RequireGolden() {
	mt.t.Helper()
	golden.RequireEqual(mt.t, mt.Text()+"\n")
}
