package termtest

import (
	"cmp"
	"fmt"
	"io"
	"iter"
	"maps"
	"slices"
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
	ModeCursorStyle          Mode = "cursor style"
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
	ModeCursorStyle,
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
	cursorStyle     int
	cursorStyleSet  bool

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
	s.parser.SetHandler(ansi.Handler{HandleCsi: s.handleCsi})

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
	return s.emu.Write(p)
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

// ModesTurnedOn lists every mode that differed from the terminal's default
// at some point, including any DEC private mode, as "DEC mode N on".
func (s *Screen) ModesTurnedOn() []Mode {
	s.mu.Lock()
	defer s.mu.Unlock()
	return sortModes(maps.Keys(s.everOn))
}

// ModesLeftOn lists every mode that differs from the terminal's default now.
func (s *Screen) ModesLeftOn() []Mode {
	s.mu.Lock()
	defer s.mu.Unlock()
	return sortModes(slices.Values(s.modesNow()))
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
	case ModeCursorStyle:
		return s.cursorStyle != 0
	}
	return false
}

func (s *Screen) cursorHidden() bool {
	// x/vt keeps cursor visibility per screen, but real terminals keep it once.
	shown, set := s.decModes[25]
	return set && !shown
}

// namedDECModes are reported under a Mode name rather than by number.
// Mode 1048 only saves and restores the cursor, so it is never "left on".
var namedDECModes = map[int]bool{
	9: true, 25: true, 47: true, 1000: true, 1001: true, 1002: true, 1003: true,
	1004: true, 1006: true, 1047: true, 1048: true, 1049: true, 2004: true,
}

// decModesOnByDefault are the DEC private modes a terminal starts with on.
var decModesOnByDefault = map[int]bool{7: true, 25: true}

func (s *Screen) modesNow() []Mode {
	var out []Mode
	for _, m := range allModes {
		if s.modeOn(m) {
			out = append(out, m)
		}
	}
	for n, on := range s.decModes {
		if namedDECModes[n] || on == decModesOnByDefault[n] {
			continue
		}
		state := "on"
		if !on {
			state = "off"
		}
		out = append(out, Mode(fmt.Sprintf("DEC mode %d %s", n, state)))
	}
	return out
}

func (s *Screen) noteModes() {
	for _, m := range s.modesNow() {
		s.everOn[m] = true
	}
}

func (s *Screen) setDECMode(m ansi.Mode, on bool) {
	if dec, ok := m.(ansi.DECMode); ok {
		s.decModes[int(dec)] = on
		s.noteModes()
	}
}

// sortModes puts named modes first, in allModes order, then DEC modes by name.
func sortModes(modes iter.Seq[Mode]) []Mode {
	out := slices.Collect(modes)
	rank := func(m Mode) int {
		if i := slices.Index(allModes, m); i >= 0 {
			return i
		}
		return len(allModes)
	}
	slices.SortFunc(out, func(a, b Mode) int {
		if c := cmp.Compare(rank(a), rank(b)); c != 0 {
			return c
		}
		return cmp.Compare(a, b)
	})
	return out
}

// CursorStyle returns the cursor style last set with DECSCUSR, where 0 is
// the terminal's default, and whether a program ever set one.
func (s *Screen) CursorStyle() (style int, set bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cursorStyle, s.cursorStyleSet
}

// CursorBlinks reports whether the current cursor style blinks. The default
// style, 0, counts as blinking.
func (s *Screen) CursorBlinks() bool {
	style, _ := s.CursorStyle()
	return style == 0 || style%2 == 1
}

// handleCsi follows what the emulator ignores: the cursor style, the Kitty
// keyboard protocol and modifyOtherKeys.
func (s *Screen) handleCsi(cmd ansi.Cmd, params ansi.Params) {
	defer s.noteModes()
	switch {
	case cmd.Final() == 'q' && cmd.Intermediate() == ' ':
		s.cursorStyle, _, _ = params.Param(0, 0)
		s.cursorStyleSet = true
	case cmd.Final() == 'u' && cmd.Prefix() == '>':
		flags, _, _ := params.Param(0, 0)
		s.kittyStack = append(s.kittyStack, s.kittyFlags)
		s.kittyFlags = flags
	case cmd.Final() == 'u' && cmd.Prefix() == '<':
		n, _, _ := params.Param(0, 1)
		for range max(n, 1) {
			if len(s.kittyStack) == 0 {
				s.kittyFlags = 0
				break
			}
			s.kittyFlags = s.kittyStack[len(s.kittyStack)-1]
			s.kittyStack = s.kittyStack[:len(s.kittyStack)-1]
		}
	case cmd.Final() == 'u' && cmd.Prefix() == '=':
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
	case cmd.Final() == 'm' && cmd.Prefix() == '>':
		resource, _, _ := params.Param(0, -1)
		if resource == 4 {
			level, _, _ := params.Param(1, 0)
			s.modifyOtherKeys = level
		}
	}
}
