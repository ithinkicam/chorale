package app_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/ithinkicam/chorale/internal/app"
	"github.com/ithinkicam/chorale/internal/termtest"
)

const askWords = "Press Ctrl-C again to quit."

func ctrlC() tea.KeyPressMsg { return tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl} }

func promptText(t *testing.T, mt *termtest.ModelTerm) string {
	t.Helper()
	lines := mt.Lines()
	label := labelRow(t, lines)
	rows := lines[label+1 : len(lines)-2]
	for i, row := range rows {
		rows[i] = strings.TrimPrefix(row, "  ")
	}
	return strings.Join(rows, "\n")
}

func waitForPrompt(t *testing.T, mt *termtest.ModelTerm, want string) {
	t.Helper()
	mt.WaitFor("the prompt to be "+want, func(string) bool {
		return promptText(t, mt) == want
	})
}

func requireRunning(t *testing.T, mt *termtest.ModelTerm) {
	t.Helper()
	if !mt.Screen().ModeOn(termtest.ModeAltScreen) {
		t.Fatalf("the program has left the alternate screen, so it has ended. The screen shows:\n%s", mt.Text())
	}
}

func waitForNo(t testing.TB, mt *termtest.ModelTerm, words string) {
	t.Helper()
	mt.WaitFor("the screen to stop showing "+words, func(screen string) bool {
		return !strings.Contains(screen, words)
	})
}

func TestQTypesOnTheSessionList(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)

	mt.Type("q")

	waitForPrompt(t, mt, "q")
	requireRunning(t, mt)
}

func TestTheFirstCtrlCAsksBeforeQuitting(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)

	mt.Send(ctrlC())

	mt.WaitForText(askWords)
	requireRunning(t, mt)
}

func TestAnotherKeyCancelsQuitting(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)

	mt.Send(ctrlC())
	mt.WaitForText(askWords)
	mt.Send(tea.KeyPressMsg{Code: tea.KeyLeft})
	waitForNo(t, mt, askWords)
	mt.Send(ctrlC())

	mt.WaitForText(askWords)
	requireRunning(t, mt)
}

func TestClearingThenCtrlCDoesNotQuit(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)
	mt.Type("half a thought")
	waitForPrompt(t, mt, "half a thought")

	mt.Send(ctrlC())
	mt.Send(ctrlC())

	mt.WaitForText(askWords)
	requireRunning(t, mt)
	if got := promptText(t, mt); got != "" {
		t.Errorf("the prompt is %q, want it empty", got)
	}
	mt.Send(tea.KeyPressMsg{Code: 'z', Mod: tea.ModCtrl})
	waitForPrompt(t, mt, "half a thought")
}

func TestTwoCtrlCsEndTheProgramOnTheSessionList(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)

	mt.Send(ctrlC())
	mt.Send(ctrlC())

	mt.WaitFinished()
}

func TestTwoCtrlCsEndTheProgramOnTheTooSmallScreen(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 40, Height: 20})
	mt.WaitForText(tooSmallWords)

	mt.Send(ctrlC())
	mt.Send(ctrlC())

	mt.WaitFinished()
}

func TestAPasteCancelsQuitting(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)

	mt.Send(ctrlC())
	mt.WaitForText(askWords)
	mt.Send(tea.PasteMsg{Content: "pasted"})
	waitForPrompt(t, mt, "pasted")
	waitForNo(t, mt, askWords)
	requireRunning(t, mt)
}

func TestAResizeCancelsQuitting(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)

	mt.Send(ctrlC())
	mt.WaitForText(askWords)
	mt.Resize(70, 24)

	waitForNo(t, mt, askWords)
	requireRunning(t, mt)
}

func TestWidthMethodFollowsTheModeReport(t *testing.T) {
	tests := []struct {
		name      string
		report    *tea.ModeReportMsg
		wantWidth int
	}{
		{"no report", nil, 1},
		{"set", &tea.ModeReportMsg{Mode: ansi.ModeUnicodeCore, Value: ansi.ModeSet}, 2},
		{"reset", &tea.ModeReportMsg{Mode: ansi.ModeUnicodeCore, Value: ansi.ModeReset}, 2},
		{"permanently set", &tea.ModeReportMsg{Mode: ansi.ModeUnicodeCore, Value: ansi.ModePermanentlySet}, 2},
		{"unknown", &tea.ModeReportMsg{Mode: ansi.ModeUnicodeCore, Value: ansi.ModeNotRecognized}, 1},
		{"permanently reset", &tea.ModeReportMsg{Mode: ansi.ModeUnicodeCore, Value: ansi.ModePermanentlyReset}, 1},
		{"another mode set", &tea.ModeReportMsg{Mode: ansi.ModeSynchronizedOutput, Value: ansi.ModeSet}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := startPromptSession(t)
			if tt.report != nil {
				s.mt.Send(*tt.report)
			}

			s.typeText("a❤️")

			waitForPrompt(t, s.mt, "a❤️")
			s.wantCursorColumn(3 + 1 + tt.wantWidth)
		})
	}
}

func TestKeysThatDoNothingClearTheRememberedColumn(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyPressMsg
	}{
		{"tab", tea.KeyPressMsg{Code: tea.KeyTab}},
		{"escape", tea.KeyPressMsg{Code: tea.KeyEscape}},
		{"shift tab", tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
			mt.WaitForText(emptyWords)
			mt.Send(tea.PasteMsg{Content: "abcdef\nab\nabcdef"})
			waitForPrompt(t, mt, "abcdef\nab\nabcdef")
			mt.Send(tea.KeyPressMsg{Code: tea.KeyHome, Mod: tea.ModCtrl})
			for range 5 {
				mt.Send(tea.KeyPressMsg{Code: tea.KeyRight})
			}

			mt.Send(tea.KeyPressMsg{Code: tea.KeyDown})
			mt.Send(tt.key)
			mt.Send(tea.KeyPressMsg{Code: tea.KeyDown})
			mt.Send(tea.KeyPressMsg{Code: 'X', Text: "X"})

			waitForPrompt(t, mt, "abcdef\nab\nabXcdef")
		})
	}
}

func TestTypedTextIsCleanedLikeAPaste(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"control sequence", "a\u009b31mb", "ab!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pasted := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
			pasted.WaitForText(emptyWords)
			pasted.Send(tea.PasteMsg{Content: tt.text})
			pasted.Send(tea.PasteMsg{Content: "!"})
			waitForPrompt(t, pasted, tt.want)

			typed := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
			typed.WaitForText(emptyWords)
			typed.Send(tea.KeyPressMsg{Text: tt.text})
			typed.Send(tea.PasteMsg{Content: "!"})

			waitForPrompt(t, typed, tt.want)
		})
	}
}

func TestCancellingTheQuitOnTheTooSmallScreen(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 40, Height: 20})
	mt.WaitForText(tooSmallWords)

	mt.Send(ctrlC())
	mt.Send(tea.KeyPressMsg{Code: 'q', Text: "q"})
	mt.Send(ctrlC())
	mt.Resize(60, 20)

	mt.WaitForText(emptyWords)
	mt.Send(tea.KeyPressMsg{Code: 'q', Text: "q"})
	waitForPrompt(t, mt, "q")
	requireRunning(t, mt)
}

func TestCtrlCClearsTheRememberedColumnOnTheSessionList(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)
	mt.Send(tea.PasteMsg{Content: "abcdef\nab\nabcdef"})
	waitForPrompt(t, mt, "abcdef\nab\nabcdef")
	mt.Send(tea.KeyPressMsg{Code: tea.KeyHome, Mod: tea.ModCtrl})
	for range 5 {
		mt.Send(tea.KeyPressMsg{Code: tea.KeyRight})
	}

	mt.Send(tea.KeyPressMsg{Code: tea.KeyDown, Mod: tea.ModShift})
	mt.Send(ctrlC())
	mt.Send(tea.KeyPressMsg{Code: tea.KeyDown})
	mt.Send(tea.KeyPressMsg{Code: 'X', Text: "X"})

	waitForPrompt(t, mt, "abcdef\nab\nabXcdef")
}

func TestCtrlCOnTheTooSmallScreenKeepsTheRememberedColumn(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)
	mt.Send(tea.PasteMsg{Content: "abcdef\nab\nabcdef"})
	waitForPrompt(t, mt, "abcdef\nab\nabcdef")
	mt.Send(tea.KeyPressMsg{Code: tea.KeyHome, Mod: tea.ModCtrl})
	for range 5 {
		mt.Send(tea.KeyPressMsg{Code: tea.KeyRight})
	}
	mt.Send(tea.KeyPressMsg{Code: tea.KeyDown})

	mt.Resize(40, 20)
	mt.WaitForText(tooSmallWords)
	mt.Send(ctrlC())
	mt.Send(tea.KeyPressMsg{Code: 'q', Text: "q"})
	mt.Resize(60, 20)
	mt.WaitForText(emptyWords)
	mt.Send(tea.KeyPressMsg{Code: tea.KeyDown})
	mt.Send(tea.KeyPressMsg{Code: 'X', Text: "X"})

	waitForPrompt(t, mt, "abcdef\nab\nabcdeXf")
}

func TestEnterEndsTheUndoGroup(t *testing.T) {
	mt := termtest.StartModel(t, app.New(), termtest.Options{Width: 60, Height: 20})
	mt.WaitForText(emptyWords)

	mt.Type("ab")
	mt.Send(tea.KeyPressMsg{Code: tea.KeyEnter})
	mt.Type("c")
	waitForPrompt(t, mt, "abc")
	mt.Send(tea.KeyPressMsg{Code: 'z', Mod: tea.ModCtrl})

	waitForPrompt(t, mt, "ab")
}

func TestCancellingTheQuitOnTheTooSmallScreenWithoutATerminal(t *testing.T) {
	var m tea.Model = app.New()
	m, _ = m.Update(tea.WindowSizeMsg{Width: 40, Height: 20})

	m, _ = m.Update(ctrlC())
	m, _ = m.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	m, cmd := m.Update(ctrlC())

	if isQuit(cmd) {
		t.Error("the second Ctrl-C quit, after another key cancelled the first")
	}
	if !app.PendingQuit(m) {
		t.Error("the second Ctrl-C didn't ask before quitting")
	}
	_, cmd = m.Update(ctrlC())
	if !isQuit(cmd) {
		t.Error("the third Ctrl-C didn't quit")
	}
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, quit := cmd().(tea.QuitMsg)
	return quit
}
