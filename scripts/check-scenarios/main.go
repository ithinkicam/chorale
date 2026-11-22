// Command check-scenarios fails if any OpenSpec scenario has no matching
// test that passed. It reads the go test -json output that make test and
// make e2e write.
//
// Usage:
//
//	go run ./scripts/check-scenarios [-openspec dir] results.json...
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout))
}

func run(args []string, out io.Writer) int {
	flags := flag.NewFlagSet("check-scenarios", flag.ContinueOnError)
	dir := flags.String("openspec", "openspec", "the OpenSpec folder to read")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() == 0 {
		_, _ = fmt.Fprintln(out, "check-scenarios: give at least one go test -json file")
		return 2
	}

	var results []io.Reader
	for _, path := range flags.Args() {
		f, err := os.Open(path) //nolint:gosec // a path the developer passed
		if err != nil {
			_, _ = fmt.Fprintf(out, "check-scenarios: %v\n", err)
			return 2
		}
		defer func() { _ = f.Close() }() // read-only
		results = append(results, f)
	}

	report, err := Check(*dir, results...)
	if err != nil {
		_, _ = fmt.Fprintf(out, "check-scenarios: %v\n", err)
		return 2
	}
	report.Print(out)
	if report.Failed() {
		return 1
	}
	return 0
}
