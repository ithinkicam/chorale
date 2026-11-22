package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ithinkicam/chorale/internal/termtest"
)

const testVersion = "0.0.0-test"

var chorale, choraleDev string

func TestMain(m *testing.M) {
	os.Exit(run(m))
}

func run(m *testing.M) int {
	dir, err := os.MkdirTemp("", "chorale-e2e-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer func() { _ = os.RemoveAll(dir) }() // a leftover temporary folder is harmless
	chorale = filepath.Join(dir, "chorale")
	choraleDev = filepath.Join(dir, "chorale-dev")
	if err := termtest.BuildChorale(chorale, testVersion); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := termtest.BuildChorale(choraleDev, ""); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return m.Run()
}
