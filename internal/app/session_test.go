package app_test

import (
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
	seen := &observed{}
	mt := termtest.StartModel(t, probe{inner: app.New(), seen: seen}, opts)
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

// frame waits until the terminal's screen and cursor match the view the
// model drew after the last message it handled, so a check never reads an
// earlier frame, and returns the screen's rows.
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
			s.t.Fatalf("the screen never showed the model's latest frame, with the cursor at column %d, row %d, visible %v. It shows the cursor at column %d, row %d, visible %v:\n%s\nwant:\n%s",
				wantX+1, wantY+1, wantSet, x+1, y+1, visible, strings.Join(rows, "\n"), strings.Join(wantRows, "\n"))
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

// wantCursorColumn checks the terminal's cursor is on the given screen
// column, counted from 1, in the prompt box's text row.
func (s *session) wantCursorColumn(col int) {
	s.t.Helper()
	lines := s.frame()
	label := labelRow(s.t, lines)
	x, y, visible := s.mt.Screen().Cursor()
	if !visible || x+1 != col || y != label+1 {
		s.t.Fatalf("the cursor is at column %d, row %d, visible %v; want column %d, row %d, the prompt box's text row. The screen shows:\n%s",
			x+1, y+1, visible, col, label+2, strings.Join(lines, "\n"))
	}
}

// wantBoxRow checks the text row of the prompt box shows want, which is
// the row the insertion point is on.
func (s *session) wantBoxRow(want string) {
	s.t.Helper()
	lines := s.frame()
	got := strings.TrimPrefix(lines[labelRow(s.t, lines)+1], "  ")
	if got != strings.TrimRight(want, " ") {
		s.t.Errorf("the prompt box's text row shows %q, want %q", got, want)
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

func named(code rune, mod tea.KeyMod) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: code, Mod: mod}
}

func words(n int) string {
	return strings.TrimSuffix(strings.Repeat("word ", n), " ")
}
