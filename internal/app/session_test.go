package app_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/ithinkicam/chorale/internal/app"
	"github.com/ithinkicam/chorale/internal/editor"
	"github.com/ithinkicam/chorale/internal/termtest"
)

type syncMsg int

type observed struct {
	mu        sync.Mutex
	state     app.PromptState
	view      string
	cursorX   int
	cursorY   int
	cursorSet bool
	synced    int
}

// probe wraps the app's model and records its state and view after every
// message, so a test reads what the model holds without racing the program.
type probe struct {
	inner tea.Model
	seen  *observed
}

func (p probe) Init() tea.Cmd { return p.inner.Init() }

func (p probe) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	next, cmd := p.inner.Update(msg)
	p.inner = next
	p.seen.mu.Lock()
	defer p.seen.mu.Unlock()
	p.seen.state = app.StateOf(next)
	view := next.View()
	p.seen.view = ansi.Strip(view.Content)
	p.seen.cursorSet = view.Cursor != nil
	if view.Cursor != nil {
		p.seen.cursorX, p.seen.cursorY = view.Cursor.X, view.Cursor.Y
	}
	if s, ok := msg.(syncMsg); ok {
		p.seen.synced = int(s)
	}
	return p, cmd
}

func (p probe) View() tea.View { return p.inner.View() }

// session is a running app whose prompt a test types into and reads back.
type session struct {
	t     *testing.T
	mt    *termtest.ModelTerm
	seen  *observed
	syncs int
}

func startSession(t *testing.T, opts termtest.Options) *session {
	t.Helper()
	return startSessionOf(t, app.New(), opts)
}

func startSessionOf(t *testing.T, model tea.Model, opts termtest.Options) *session {
	t.Helper()
	seen := &observed{}
	mt := termtest.StartModel(t, probe{inner: model, seen: seen}, opts)
	mt.WaitForText(emptyWords)
	return &session{t: t, mt: mt, seen: seen}
}

func startPromptSession(t *testing.T) *session {
	t.Helper()
	return startSession(t, termtest.Options{Width: 60, Height: 20})
}

// settle waits until the model has handled everything sent so far.
func (s *session) settle() {
	s.t.Helper()
	s.syncs++
	s.mt.Send(syncMsg(s.syncs))
	deadline := time.Now().Add(termtest.DefaultTimeout)
	for {
		s.seen.mu.Lock()
		done := s.seen.synced >= s.syncs
		s.seen.mu.Unlock()
		if done {
			return
		}
		if time.Now().After(deadline) {
			s.t.Fatalf("the model didn't handle what was sent. The screen shows:\n%s", s.mt.Text())
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (s *session) state() app.PromptState {
	s.t.Helper()
	s.settle()
	s.seen.mu.Lock()
	defer s.seen.mu.Unlock()
	return s.seen.state
}

// view returns the text of the view the model draws now, with no styling.
func (s *session) view() string {
	s.t.Helper()
	s.settle()
	s.seen.mu.Lock()
	defer s.seen.mu.Unlock()
	return s.seen.view
}

func (s *session) press(keys ...tea.KeyPressMsg) {
	s.t.Helper()
	for _, k := range keys {
		if _, bound := editor.Lookup(k.String()); !bound && k.Text == "" {
			s.t.Fatalf("%q is not a key in the key table", k.String())
		}
		s.mt.Send(k)
	}
	s.settle()
}

// pressUnbound sends a key that has no binding, which Chorale ignores.
func (s *session) pressUnbound(k tea.KeyPressMsg) {
	s.t.Helper()
	if _, bound := editor.Lookup(k.String()); bound {
		s.t.Fatalf("%q is a key in the key table", k.String())
	}
	s.mt.Send(k)
	s.settle()
}

func (s *session) typeText(text string) {
	s.t.Helper()
	s.mt.Type(text)
	s.settle()
}

func (s *session) paste(text string) {
	s.t.Helper()
	s.mt.Send(tea.PasteMsg{Content: text})
	s.settle()
}

// given puts text in the prompt with the insertion point at caret, a
// character index, and nothing selected.
func (s *session) given(text string, caret int) {
	s.t.Helper()
	s.paste(text)
	s.press(named(tea.KeyUp, tea.ModAlt))
	for range caret {
		s.press(named(tea.KeyRight, 0))
	}
	if got := s.state(); got.Text != text || got.Caret != caret {
		s.t.Fatalf("the set-up left the prompt %q with the insertion point at %d, want %q and %d", got.Text, got.Caret, text, caret)
	}
}

func (s *session) givenAtEnd(text string) {
	s.t.Helper()
	s.paste(text)
}

// givenSelected puts text in the prompt with its last count characters
// selected and the insertion point before them.
func (s *session) givenSelected(text string, count int) {
	s.t.Helper()
	s.givenAtEnd(text)
	for range count {
		s.press(named(tea.KeyLeft, tea.ModShift))
	}
}

func (s *session) wantText(want string) {
	s.t.Helper()
	if got := s.state().Text; got != want {
		s.t.Errorf("the prompt is %q, want %q", got, want)
	}
}

func (s *session) wantCaret(want int) {
	s.t.Helper()
	if got := s.state().Caret; got != want {
		s.t.Errorf("the insertion point is at %d, want %d", got, want)
	}
}

func (s *session) wantSelected(want string) {
	s.t.Helper()
	if got := s.state().Selected; got != want {
		s.t.Errorf("%q is selected, want %q", got, want)
	}
}

func (s *session) wantRows(want ...string) {
	s.t.Helper()
	got := s.state().Rows
	for i, row := range want {
		if i >= len(got) || strings.TrimRight(got[i], " ") != row {
			s.t.Errorf("text rows are %q, want them to start with %q", got, want)
			return
		}
	}
}

func (s *session) wantRowsNoWiderThan(width int) {
	s.t.Helper()
	for i, w := range s.state().RowWidths {
		if w > width {
			s.t.Errorf("text row %d is %d columns wide, wider than %d", i+1, w, width)
		}
	}
}

// frame waits until the terminal's screen and cursor match the model's latest
// view, and returns the screen's rows.
func (s *session) frame() []string {
	s.t.Helper()
	s.settle()
	s.seen.mu.Lock()
	wantRows := strings.Split(s.seen.view, "\n")
	wantX, wantY, wantSet := s.seen.cursorX, s.seen.cursorY, s.seen.cursorSet
	s.seen.mu.Unlock()
	deadline := time.Now().Add(termtest.DefaultTimeout)
	for {
		rows := s.mt.Lines()
		x, y, visible := s.mt.Screen().Cursor()
		if sameRows(rows, wantRows) && visible == wantSet && (!wantSet || x == wantX && y == wantY) {
			return rows
		}
		if time.Now().After(deadline) {
			wantCursor := "with the cursor hidden"
			if wantSet {
				wantCursor = fmt.Sprintf("with the cursor at column %d, row %d", wantX+1, wantY+1)
			}
			s.t.Fatalf("the screen shows the cursor at column %d, row %d, visible %v:\n%s\nwant the model's latest frame, %s:\n%s",
				x+1, y+1, visible, strings.Join(rows, "\n"), wantCursor, strings.Join(wantRows, "\n"))
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func sameRows(screen, view []string) bool {
	for i, row := range screen {
		want := ""
		if i < len(view) {
			want = strings.TrimRight(view[i], " ")
		}
		if row != want {
			return false
		}
	}
	return len(view) <= len(screen)
}

// boxRows returns the prompt box's text rows as the terminal's screen shows
// them, without the indent or the spaces after the text.
func (s *session) boxRows() []string {
	s.t.Helper()
	lines := s.frame()
	rows := lines[labelRow(s.t, lines)+1 : len(lines)-2]
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = strings.TrimRight(strings.TrimPrefix(row, "  "), " ")
	}
	return out
}

// wantLabel checks the prompt box's label row on the screen reads exactly want.
func (s *session) wantLabel(want string) {
	s.t.Helper()
	lines := s.frame()
	if got := strings.TrimRight(lines[labelRow(s.t, lines)], " "); got != want {
		s.t.Errorf("the label row reads %q, want %q", got, want)
	}
}

// wantCursorAt checks the terminal's cursor is on the given text row of the
// prompt box and the given screen column, both counted from 1.
func (s *session) wantCursorAt(row, col int) {
	s.t.Helper()
	lines := s.frame()
	label := labelRow(s.t, lines)
	x, y, visible := s.mt.Screen().Cursor()
	if !visible || x+1 != col || y != label+row {
		s.t.Fatalf("the cursor is at column %d, row %d, visible %v; want column %d on text row %d of the prompt box, screen row %d. The screen shows:\n%s",
			x+1, y+1, visible, col, row, label+row+1, strings.Join(lines, "\n"))
	}
}

// wantCursorColumn checks the terminal's cursor is on the given screen
// column, counted from 1, in the prompt box's first text row.
func (s *session) wantCursorColumn(col int) {
	s.t.Helper()
	s.wantCursorAt(1, col)
}

// wantBoxRows checks the first text rows of the prompt box on the screen
// show want, one string for each row.
func (s *session) wantBoxRows(want ...string) {
	s.t.Helper()
	got := s.boxRows()
	if len(got) < len(want) {
		s.t.Errorf("the prompt box shows %d text rows, %q, want at least %d, starting %q", len(got), got, len(want), want)
		return
	}
	for i, row := range want {
		if got[i] != strings.TrimRight(row, " ") {
			s.t.Errorf("the prompt box shows text rows %q, want them to start with %q", got, want)
			return
		}
	}
}

// wantBoxRowsNoWiderThan checks no text row of the prompt box, with its
// indent, is wider than width columns on the terminal's screen.
func (s *session) wantBoxRowsNoWiderThan(width int) {
	s.t.Helper()
	lines := s.frame()
	for i, row := range lines[labelRow(s.t, lines)+1 : len(lines)-2] {
		if w := ansi.StringWidth(row); w > width {
			s.t.Errorf("the prompt box's text row %d is %d columns wide, wider than %d: %q", i+1, w, width, row)
		}
	}
}

// wantViewFits checks the view the model draws has no more rows than
// height and none wider than width.
func (s *session) wantViewFits(width, height int) {
	s.t.Helper()
	rows := strings.Split(s.view(), "\n")
	if len(rows) > height {
		s.t.Errorf("the view has %d rows, more than the window's %d", len(rows), height)
	}
	for i, row := range rows {
		if w := ansi.StringWidth(row); w > width {
			s.t.Errorf("view row %d is %d columns wide, wider than the window's %d: %q", i+1, w, width, row)
		}
	}
}

func (s *session) wantScreenWith(words string) {
	s.t.Helper()
	if screen := strings.Join(s.frame(), "\n"); !strings.Contains(screen, words) {
		s.t.Errorf("the terminal's screen doesn't show %q:\n%s", words, screen)
	}
}

func (s *session) wantLastLine(want string) {
	s.t.Helper()
	lines := s.frame()
	if last := lines[len(lines)-1]; !strings.Contains(last, want) {
		s.t.Errorf("the last line is %q, want it to show %q", last, want)
	}
}

func (s *session) wantViewWithout(words string) {
	s.t.Helper()
	if view := s.view(); strings.Contains(view, words) {
		s.t.Errorf("the screen shows %q, want it not to:\n%s", words, view)
	}
}

func (s *session) wantScreenWithout(words string) {
	s.t.Helper()
	if screen := strings.Join(s.frame(), "\n"); strings.Contains(screen, words) {
		s.t.Errorf("the terminal's screen shows %q, want it not to:\n%s", words, screen)
	}
}

func named(code rune, mod tea.KeyMod) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code, Mod: mod}
}

func words(n int) string {
	return strings.TrimSuffix(strings.Repeat("word ", n), " ")
}
