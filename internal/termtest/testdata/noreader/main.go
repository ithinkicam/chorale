// Command noreader is a fixture that asks the terminal thousands of
// questions and never reads the answers.
package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Print(strings.Repeat("\x1b[?2004$p", 50000))
	fmt.Print("done")
	select {}
}
