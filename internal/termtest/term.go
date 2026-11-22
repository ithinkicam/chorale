package termtest

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/creack/pty"
)

// Timing used by the helpers. Waits poll rather than sleep a fixed time,
// and give up after DefaultTimeout.
const (
	DefaultTimeout = 5 * time.Second
	pollInterval   = 10 * time.Millisecond

	// IdleSettle is how long a program must write nothing before the idle
	// check starts counting, and IdleWindow is how long it then counts.
	IdleSettle = 500 * time.Millisecond
	IdleWindow = 2 * time.Second
)

// Options describe how to start a program in a pseudo-terminal.
type Options struct {
	Width, Height int
	// NoColor sets NO_COLOR=1 in the program's environment.
	NoColor bool
	// Args are the program's command-line arguments.
	Args []string
	// Env adds variables, as KEY=value, to the fixed environment.
	Env []string
}

// Term is a program running in a pseudo-terminal, with a Screen fed by
// everything it writes.
type Term struct {
	t      testing.TB
	cmd    *exec.Cmd
	ptmx   *os.File
	screen *Screen

	mu        sync.Mutex
	raw       []byte
	lastWrite time.Time

	readDone chan struct{}
	exited   chan struct{}
	exitCode int
}

// Start runs binary in a pseudo-terminal of the given size, with a fixed
// environment so results don't depend on the machine: TERM=xterm-256color,
// a UTF-8 locale, a temporary HOME, and NO_COLOR only if opts.NoColor is
// set. The program is stopped when the test ends.
func Start(t testing.TB, binary string, opts Options) *Term {
	t.Helper()
	cmd := exec.Command(binary, opts.Args...)
	cmd.Env = fixedEnv(t, opts)
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{
		Cols: uint16(opts.Width),
		Rows: uint16(opts.Height),
	})
	if err != nil {
		t.Fatalf("starting %s: %v", binary, err)
	}

	term := &Term{
		t:         t,
		cmd:       cmd,
		ptmx:      ptmx,
		screen:    NewScreen(opts.Width, opts.Height),
		lastWrite: time.Now(),
		readDone:  make(chan struct{}),
		exited:    make(chan struct{}),
	}
	term.screen.SetReplyWriter(ptmx)
	go term.read()
	go term.wait()

	t.Cleanup(term.stop)
	return term
}

func fixedEnv(t testing.TB, opts Options) []string {
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + t.TempDir(),
		"TERM=xterm-256color",
		"LANG=en_US.UTF-8",
		"LC_ALL=en_US.UTF-8",
	}
	if opts.NoColor {
		env = append(env, "NO_COLOR=1")
	}
	return append(env, opts.Env...)
}

func (term *Term) read() {
	defer close(term.readDone)
	buf := make([]byte, 32*1024)
	for {
		n, err := term.ptmx.Read(buf)
		if n > 0 {
			term.mu.Lock()
			term.raw = append(term.raw, buf[:n]...)
			term.lastWrite = time.Now()
			term.mu.Unlock()
			_, _ = term.screen.Write(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

func (term *Term) wait() {
	err := term.cmd.Wait()
	var exitErr *exec.ExitError
	switch {
	case err == nil:
		term.exitCode = 0
	case errors.As(err, &exitErr):
		term.exitCode = exitErr.ExitCode()
	default:
		term.exitCode = -1
	}
	close(term.exited)
}

func (term *Term) stop() {
	select {
	case <-term.exited:
	default:
		_ = term.cmd.Process.Kill()
		<-term.exited
	}
	_ = term.ptmx.Close()
	_ = term.screen.Close()
}

// Screen returns the screen the program draws on.
func (term *Term) Screen() *Screen { return term.screen }

// Raw returns a copy of every byte the program has written so far.
func (term *Term) Raw() []byte {
	term.mu.Lock()
	defer term.mu.Unlock()
	return append([]byte(nil), term.raw...)
}

// Send types text into the program, as raw bytes.
func (term *Term) Send(text string) {
	term.t.Helper()
	if _, err := term.ptmx.WriteString(text); err != nil {
		term.t.Fatalf("sending %s: %v", quote(text), err)
	}
}

// Resize changes the terminal's size. The program gets SIGWINCH.
func (term *Term) Resize(width, height int) {
	term.t.Helper()
	term.screen.Resize(width, height)
	err := pty.Setsize(term.ptmx, &pty.Winsize{
		Cols: uint16(width),
		Rows: uint16(height),
	})
	if err != nil {
		term.t.Fatalf("resizing to %dx%d: %v", width, height, err)
	}
}

// WaitForText waits until the screen shows text, and fails the test if it
// doesn't within DefaultTimeout.
func (term *Term) WaitForText(text string) {
	term.t.Helper()
	if !term.poll(func() bool { return strings.Contains(term.screen.Text(), text) }) {
		term.t.Fatalf("timed out waiting for %q. The screen shows:\n%s", text, term.screen.Text())
	}
}

// WaitForCondition waits until ok returns true, and fails the test with
// what if it doesn't within DefaultTimeout.
func (term *Term) WaitForCondition(what string, ok func() bool) {
	term.t.Helper()
	if !term.poll(ok) {
		term.t.Fatalf("timed out waiting for %s. The screen shows:\n%s", what, term.screen.Text())
	}
}

// WaitQuiet waits until the program has written nothing for quiet.
func (term *Term) WaitQuiet(quiet time.Duration) {
	term.t.Helper()
	ok := term.poll(func() bool {
		term.mu.Lock()
		defer term.mu.Unlock()
		return time.Since(term.lastWrite) >= quiet
	})
	if !ok {
		term.t.Fatalf("the program never went quiet for %v", quiet)
	}
}

// IdleBytes waits for the program to go quiet for IdleSettle, then counts
// the bytes it writes over IdleWindow with no input. A still program
// writes none. Comparing screens instead would miss a program that clears
// and redraws the same picture, which is flicker.
func (term *Term) IdleBytes() int {
	term.t.Helper()
	term.WaitQuiet(IdleSettle)
	before := len(term.Raw())
	time.Sleep(IdleWindow)
	return len(term.Raw()) - before
}

// Wait waits for the program to exit and for its output to be read, and
// returns its exit status. It fails the test after DefaultTimeout.
func (term *Term) Wait() int {
	term.t.Helper()
	select {
	case <-term.exited:
	case <-time.After(DefaultTimeout):
		term.t.Fatalf("the program didn't exit. The screen shows:\n%s", term.screen.Text())
	}
	select {
	case <-term.readDone:
	case <-time.After(DefaultTimeout):
		term.t.Fatal("the program exited, but its output never finished")
	}
	return term.exitCode
}

func (term *Term) poll(ok func() bool) bool {
	deadline := time.Now().Add(DefaultTimeout)
	for time.Now().Before(deadline) {
		if ok() {
			return true
		}
		time.Sleep(pollInterval)
	}
	return ok()
}
