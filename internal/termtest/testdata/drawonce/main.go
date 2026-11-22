// Command drawonce is a fixture for the termtest harness that draws one screen
// describing its terminal.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"
)

func main() {
	draw()

	resized := make(chan os.Signal, 1)
	signal.Notify(resized, syscall.SIGWINCH)
	keys := make(chan byte)
	go readKeys(keys)

	for {
		select {
		case <-resized:
			draw()
		case k := <-keys:
			switch k {
			case 'c':
				fmt.Print("\x1b[31mred\x1b[0m")
			case 'q':
				os.Exit(3)
			}
		}
	}
}

func draw() {
	ws, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		fmt.Println("no size:", err)
		return
	}
	noColor, set := os.LookupEnv("NO_COLOR")
	if !set {
		noColor = "unset"
	}
	fmt.Print("\x1b[2J\x1b[H")
	fmt.Printf("size %dx%d\r\n", ws.Col, ws.Row)
	fmt.Printf("TERM %s\r\n", os.Getenv("TERM"))
	fmt.Printf("NO_COLOR %s\r\n", noColor)
	fmt.Print("\x1b[5;3H") // leave the cursor at column 2, row 4
}

func readKeys(keys chan<- byte) {
	// Raw mode, so each key arrives on its own and isn't echoed.
	fd := int(os.Stdin.Fd())
	if t, err := unix.IoctlGetTermios(fd, ioctlGetTermios); err == nil {
		t.Lflag &^= unix.ICANON | unix.ECHO
		_ = unix.IoctlSetTermios(fd, ioctlSetTermios, t)
	}
	buf := make([]byte, 1)
	for {
		if _, err := os.Stdin.Read(buf); err != nil {
			return
		}
		keys <- buf[0]
	}
}
