package termtest

import (
	"io"
	"strings"
	"sync"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/vt"
)

// Mode is a terminal mode that a program can turn on and must turn off
// again before it exits.
type Mode string

// The modes the screen tracks.
const (
	ModeAltScreen            Mode = "alternate screen"
	ModeCursorHidden         Mode = "cursor hidden"
	ModeMouse                Mode = "mouse reporting"
	ModeMouseSGR             Mode = "SGR mouse encoding"
	ModeBracketedPaste       Mode = "bracketed paste"
	ModeFocusEvents          Mode = "focus events"
	ModeKeyboardEnhancements Mode = "keyboard enhancements"
)

var allModes = []Mode{
	ModeAltScreen,
	ModeBracketedPaste,
	ModeCursorHidden,
	ModeFocusEvents,
	ModeKeyboardEnhancements,
	ModeMouse,
	ModeMouseSGR,
}

// Screen is a terminal emulator for tests. Write a program's output to it,
// then read back the text on screen, the cursor and the terminal modes.
type Screen struct {
	mu sync.Mutex

	emu    *vt.Emulator
	parser *ansi.Parser

	decModes map[int]bool
	everOn   map[Mode]bool

	kittyStack      []int
	kittyFlags      int
	modifyOtherKeys int

	replies    *replyQueue
	readerDone chan struct{}
	closeOnce  sync.Once
}

// NewScreen returns a blank screen of the given size.
func NewScreen(width, height int) *Screen {
	s := &Screen{
		emu:        vt.NewEmulator(width, height),
		parser:     ansi.NewParser(),
		decModes:   map[int]bool{},
		everOn:     map[Mode]bool{},
		replies:    newReplyQueue(),
		readerDone: make(chan struct{}),
	}
	s.emu.SetCallbacks(vt.Callbacks{
		EnableMode:  func(m ansi.Mode) { s.setDECMode(m, true) },
		DisableMode: func(m ansi.Mode) { s.setDECMode(m, false) },
	})
	s.parser.SetHandler(ansi.Handler{HandleCsi: s.handleKeyboardCsi})

	// The emulator blocks on its reply pipe unless something drains it.
	go s.readReplies()
	return s
}

// SetReplyWriter sends the emulator's answers to queries to w, which is
// normally the program's input.
func (s *Screen) SetReplyWriter(w io.Writer) {
	s.replies.setWriter(w)
}

func (s *Screen) readReplies() {
	defer close(s.readerDone)
	buf := make([]byte, 4096)
	for {
		n, err := s.emu.Read(buf)
		if n > 0 {
			s.replies.push(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

// Write feeds a program's output to the screen.
func (s *Screen) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, b := range p {
		s.parser.Advance(b)
	}
	n, err := s.emu.Write(p)
	s.noteModes()
	return n, err
}

// Resize changes the screen's size, as a terminal window would.
func (s *Screen) Resize(width, height int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.emu.Resize(width, height)
}

// Close stops the screen. Further writes fail, and closing again does nothing.
func (s *Screen) Close() error {
	var err error
	s.closeOnce.Do(func() {
		if pipe, ok := s.emu.InputPipe().(io.Closer); ok {
			_ = pipe.Close() // ends readReplies before the emulator is closed
		}
		<-s.readerDone
		s.replies.close()
		s.mu.Lock()
		defer s.mu.Unlock()
		err = s.emu.Close()
	})
	return err
}

// ModeOn reports whether mode is on now.
func (s *Screen) ModeOn(mode Mode) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.modeOn(mode)
}

// ModesTurnedOn lists every mode that was on at some point.
func (s *Screen) ModesTurnedOn() []Mode {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Mode
	for _, m := range allModes {
		if s.everOn[m] {
			out = append(out, m)
		}
	}
	return out
}

// ModesLeftOn lists every mode that is on now.
func (s *Screen) ModesLeftOn() []Mode {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []Mode
	for _, m := range allModes {
		if s.modeOn(m) {
			out = append(out, m)
		}
	}
	return out
}

// Lines returns the text on screen, one string per row, with trailing
// spaces removed. Wide characters appear once.
func (s *Screen) Lines() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	width, height := s.emu.Width(), s.emu.Height()
	lines := make([]string, height)
	for y := range height {
		var row strings.Builder
		for x := 0; x < width; x++ {
			c := s.emu.CellAt(x, y)
			switch {
			case c == nil || c.Content == "":
				if c != nil && c.Width == 0 {
					continue // the second half of a wide character
				}
				row.WriteByte(' ')
			default:
				row.WriteString(c.Content)
			}
		}
		lines[y] = strings.TrimRight(row.String(), " ")
	}
	return lines
}

// Text returns Lines joined with newlines.
func (s *Screen) Text() string {
	return strings.Join(s.Lines(), "\n")
}

// Cursor returns the cursor's column and row, counted from zero, and
// whether it is visible.
func (s *Screen) Cursor() (x, y int, visible bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos := s.emu.CursorPosition()
	return pos.X, pos.Y, !s.cursorHidden()
}

func (s *Screen) modeOn(mode Mode) bool {
	switch mode {
	case ModeAltScreen:
		return s.emu.IsAltScreen()
	case ModeCursorHidden:
		return s.cursorHidden()
	case ModeMouse:
		return s.decModes[9] || s.decModes[1000] || s.decModes[1001] ||
			s.decModes[1002] || s.decModes[1003]
	case ModeMouseSGR:
		return s.decModes[1006]
	case ModeBracketedPaste:
		return s.decModes[2004]
	case ModeFocusEvents:
		return s.decModes[1004]
	case ModeKeyboardEnhancements:
		return s.kittyFlags != 0 || s.modifyOtherKeys != 0
	}
	return false
}

// cursorHidden follows DEC mode 25 itself: x/vt keeps cursor visibility per
// screen, while real terminals keep it once.
func (s *Screen) cursorHidden() bool {
	shown, set := s.decModes[25]
	return set && !shown
}

func (s *Screen) noteModes() {
	for _, m := range allModes {
		if s.modeOn(m) {
			s.everOn[m] = true
		}
	}
}

func (s *Screen) setDECMode(m ansi.Mode, on bool) {
	if dec, ok := m.(ansi.DECMode); ok {
		s.decModes[int(dec)] = on
	}
}

// handleKeyboardCsi follows the Kitty keyboard protocol and modifyOtherKeys.
func (s *Screen) handleKeyboardCsi(cmd ansi.Cmd, params ansi.Params) {
	switch {
	case cmd.Final() == 'u' && cmd.Prefix() == '>': // push flags
		flags, _, _ := params.Param(0, 0)
		s.kittyStack = append(s.kittyStack, s.kittyFlags)
		s.kittyFlags = flags
	case cmd.Final() == 'u' && cmd.Prefix() == '<': // pop n entries
		n, _, _ := params.Param(0, 1)
		for range max(n, 1) {
			if len(s.kittyStack) == 0 {
				s.kittyFlags = 0
				break
			}
			s.kittyFlags = s.kittyStack[len(s.kittyStack)-1]
			s.kittyStack = s.kittyStack[:len(s.kittyStack)-1]
		}
	case cmd.Final() == 'u' && cmd.Prefix() == '=': // set flags
		flags, _, _ := params.Param(0, 0)
		how, _, _ := params.Param(1, 1)
		switch how {
		case 1:
			s.kittyFlags = flags
		case 2:
			s.kittyFlags |= flags
		case 3:
			s.kittyFlags &^= flags
		}
	case cmd.Final() == 'm' && cmd.Prefix() == '>': // modifyOtherKeys
		resource, _, _ := params.Param(0, -1)
		if resource == 4 {
			level, _, _ := params.Param(1, 0)
			s.modifyOtherKeys = level
		}
	}
}
