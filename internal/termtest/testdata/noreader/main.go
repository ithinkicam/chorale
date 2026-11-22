// Command noreader is a fixture that asks the terminal thousands of
// questions and never reads the answers.
package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

func main() {
	fd := int(os.Stdin.Fd())
	if t, err := unix.IoctlGetTermios(fd, ioctlGetTermios); err == nil {
		t.Lflag &^= unix.ICANON | unix.ECHO
		_ = unix.IoctlSetTermios(fd, ioctlSetTermios, t)
	}
	fmt.Print(strings.Repeat("\x1b[?2004$p", 50000))
	fmt.Print("done")
	select {}
}
