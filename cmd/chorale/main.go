// Command chorale is a voice-first terminal workbench for coding agents.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/ithinkicam/chorale/internal/app"
	"github.com/ithinkicam/chorale/internal/version"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run parses the command line and does what it asks, returning the exit
// status.
func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("chorale", flag.ContinueOnError)
	flags.SetOutput(stderr)
	showVersion := flags.Bool("version", false, "print the version and exit")
	if err := flags.Parse(args); err != nil {
		return 2
	}

	if *showVersion {
		if _, err := fmt.Fprintf(stdout, "chorale %s\n", version.Version); err != nil {
			return 1
		}
		return 0
	}

	if _, err := tea.NewProgram(app.New()).Run(); err != nil {
		_, _ = fmt.Fprintf(stderr, "chorale: %v\n", err)
		return 1
	}
	return 0
}
