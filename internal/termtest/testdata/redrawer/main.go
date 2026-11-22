// Command redrawer is a fixture for the termtest harness that redraws the
// same text on a timer.
package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Print("\x1b[2J\x1b[Hthe same text\r\n")
		time.Sleep(700 * time.Millisecond)
	}
}
